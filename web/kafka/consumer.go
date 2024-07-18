package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"product_management/database"
	"product_management/models"

	"github.com/IBM/sarama"
)

func InitConsumer() {
	config := SaramaConfig()
	brokers := []string{"kafka:9092"}
	consumerGroup, err := sarama.NewConsumerGroup(brokers, "web-service-group", config)
	if err != nil {
		log.Fatalf("Failed to start Sarama consumer group: %v", err)
	}
	defer consumerGroup.Close()

	handler := &ConsumerGroupHandler{}
	for {
		if err := consumerGroup.Consume(context.Background(), []string{"register-user"}, handler); err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}

}

type ConsumerGroupHandler struct{}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var username string
		err := json.Unmarshal(msg.Value, &username)
		if err != nil {
			fmt.Printf("failed to unmarshal username: %v\n", err)
			continue
		}
		authUser := models.User{Username: username}
		if err := database.DB.DB.Create(&authUser).Error; err != nil {
			fmt.Printf("failed to save user to web database: %v\n", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
