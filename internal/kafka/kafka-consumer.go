package kafka

import (
	"context"
	"encoding/json"
	"github.com/khorzhenwin/go-chujang/internal/models"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

// StartKafkaConsumer consumes messages from a Kafka topic and pushes them into a channel
func StartKafkaConsumer(brokers []string, topic string, groupID string, out chan<- models.TickerPrice) {
	log.Println("📥 Starting Kafka consumer...")

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		log.Fatalf("❌ Failed to create Kafka consumer: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	for {
		fetches := client.PollFetches(ctx)
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, record := range p.Records {
				var msg models.TickerPrice
				if err := json.Unmarshal(record.Value, &msg); err != nil {
					log.Printf("⚠️ Failed to parse message: %v", err)
					continue
				}
				out <- msg
			}
		})

		if err := fetches.Err(); err != nil {
			log.Printf("⚠️ Fetch error: %v", err)
		}
	}
}
