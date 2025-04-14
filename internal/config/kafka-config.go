package config

import (
	"fmt"
	"os"
)

type KafkaConfig struct {
	Broker           string
	Username         string
	Password         string
	TickerPriceTopic string
}

func LoadKafkaConfig() (*KafkaConfig, error) {
	cfg := &KafkaConfig{
		Broker:           os.Getenv("KAFKA_BROKER"),
		Username:         os.Getenv("KAFKA_USERNAME"),
		Password:         os.Getenv("KAFKA_PASSWORD"),
		TickerPriceTopic: os.Getenv("KAFKA_TICKER_PRICE_TOPIC"),
	}

	if cfg.Username == "" || cfg.Password == "" {
		return nil, fmt.Errorf("incomplete Kafka config")
	}

	return cfg, nil
}
