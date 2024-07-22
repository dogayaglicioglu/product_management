package kafka

import (
	"context"
	"fmt"
	"log"
	"product_management/database"
	"product_management/models"

	"github.com/IBM/sarama"
)

func InitConsumer() {
	config := SaramaConfig()
	brokers := []string{"kafka:9092"}
	consumerGroup, err := sarama.NewConsumerGroup(brokers, "web-service-group-new", config)
	if err != nil {
		log.Fatalf("Failed to start Sarama consumer group: %v", err)
	}
	fmt.Print("CONSUMER CREATED...")
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
		username := string(msg.Value)
		fmt.Printf("USERNAME: %s\n", username)

		authUser := models.User{Username: username}
		fmt.Printf("AUTH USER: %v\n", authUser)

		if err := database.DB.DB.Create(&authUser).Error; err != nil {
			fmt.Printf("Failed to save user to web database: %v\n", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
