package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func createProducer() sarama.SyncProducer {
	config := SaramaConf()
	brokers := []string{"kafka:9092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	return producer
}
