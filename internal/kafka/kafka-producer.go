package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/khorzhenwin/go-chujang/internal/config"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"log"
	"os"
	"sync"
	"time"
)

var client *kgo.Client
var clientError error

func InitKafkaProducer(kafkaConfig *config.KafkaConfig) {
	broker := kafkaConfig.Broker
	username := kafkaConfig.Username
	password := kafkaConfig.Password

	client, clientError = kgo.NewClient(
		kgo.SeedBrokers(broker),
		kgo.DialTLSConfig(&tls.Config{}),
		kgo.SASL(scram.Auth{
			User: username,
			Pass: password,
		}.AsSha256Mechanism()),
		kgo.WithLogger(kgo.BasicLogger(os.Stdout, kgo.LogLevelDebug, nil)),
	)

	if clientError != nil {
		log.Fatalf("❌ Failed to initialize franz Kafka client: %v", clientError)
	}

	log.Println("🚀 Franz Kafka producer initialized")
}

func CloseKafkaProducer() {
	if client != nil {
		client.Close()
		log.Println("👋 Closed Franz Kafka client")
	}
}

func PushToKafkaTopic[T any](topic string, data T, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	value, err := json.Marshal(data)
	if err != nil {
		log.Println("❌ JSON encode error: %w", err)
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	client.Produce(ctx, record, func(r *kgo.Record, err error) {
		defer wg.Done()
		if err != nil {
			log.Printf("❌ Failed to produce message: %v", err)
		} else {
			log.Printf("✅ Kafka message sent (offset=%d): %s", r.Offset, key)
		}
	})

	wg.Wait()
}
