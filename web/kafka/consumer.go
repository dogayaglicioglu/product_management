package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"product_management/database"
	"product_management/models"
	"time"

	"github.com/IBM/sarama"
)

const (
	consumerGroupName = "web-service-group-new"
)

type KafkaConsumer struct {
	brokers         []string
	groupName       string
	topics          []string
	consumerGroup   sarama.ConsumerGroup
	consumerHandler ConsumerGroupHandler
}

func NewKafkaConsumer(brokers []string, groupID string, topics []string, config *sarama.Config, kafkaCreated chan bool) (*KafkaConsumer, error) {
	if err := waitForKafka(brokers, 1*time.Minute, config, kafkaCreated); err != nil {
		return nil, fmt.Errorf("Failed to connect to kafka: %w", err)
	}
	consumerGroup, err := sarama.NewConsumerGroup(brokers, consumerGroupName, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start Sarama consumer group: %w", err)
	}
	return &KafkaConsumer{
		brokers:         brokers,
		groupName:       "web-service-group-new",
		topics:          topics,
		consumerGroup:   consumerGroup,
		consumerHandler: ConsumerGroupHandler{},
	}, nil
}

func (kc *KafkaConsumer) Start() error {
	defer kc.consumerGroup.Close()
	for {
		if err := kc.consumerGroup.Consume(context.Background(), kc.topics, &kc.consumerHandler); err != nil {
			return fmt.Errorf("failed to consume messages: %w", err)
		}
	}
}
func InitConsumer(kafkaCreated chan bool) {
	config := SaramaConfig()
	brokers := []string{"kafka:9092"}
	topics := []string{"register-user", "change-username"}
	kafkaConsumer, err := NewKafkaConsumer(brokers, "web-service-group-new", topics, config, kafkaCreated)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	kafkaConsumer.Start()

}

type ConsumerGroupHandler struct{}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		switch msg.Topic {
		case "register-user":
			handlerRegisterUser(msg)
		case "change-username", "update-user":
			handlerUpdateUser(msg)
		case "delete-user":
			handlerDeleteUser(msg)
		default:
			log.Printf("Unhandled topic: %s\n", msg.Topic)

		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func waitForKafka(brokers []string, timeout time.Duration, config *sarama.Config, kafkaCreated chan bool) error { //BCS Kafka consumer group has not started
	startTime := time.Now()

	for time.Since(startTime) < timeout {
		client, err := sarama.NewClient(brokers, config)
		if err == nil {
			fmt.Print("Kafka connection success..")
			kafkaCreated <- true
			client.Close()
			return nil
		}
		fmt.Printf("Kafka connection failed, trying again.. %v\n", err)
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("Kafka connection failed")
}

func handlerRegisterUser(msg *sarama.ConsumerMessage) {
	var user models.User
	username := string(msg.Value)
	fmt.Printf("USERNAME: %s\n", user.Username)

	authUser := models.User{Username: username}
	fmt.Printf("AUTH USER: %v\n", authUser)

	if err := database.DB.DB.Create(&authUser).Error; err != nil {
		fmt.Printf("Failed to save user to web database: %v\n", err)
	}
}

func handlerUpdateUser(msg *sarama.ConsumerMessage) {
	var msgPayload map[string]interface{}
	if err := json.Unmarshal(msg.Value, &msgPayload); err != nil {
		fmt.Println("Error in unsmarshal operation. %v", err)
	}
	oldUsername, ok1 := msgPayload["old_username"].(string)
	newUsername, ok2 := msgPayload["new_username"].(string)

	if !ok1 || !ok2 {
		fmt.Println("Invalid data format: 'oldUsername' or 'newUsername' is missing or not a string")
		return
	}
	updatedUser := models.User{
		Username: newUsername,
	}
	result := database.DB.DB.Model(&models.User{}).Where("username = ?", oldUsername).Updates(updatedUser)

	if result.Error != nil {
		log.Printf("Failed to delete user in web database: %v \n", result.Error)
	} else if result.RowsAffected == 0 {
		fmt.Println("No user found with the given username.")
		return
	} else {
		fmt.Printf("User with username '%s' updated as a '%s'successfully.\n", oldUsername, newUsername)
	}
}

func handlerDeleteUser(msg *sarama.ConsumerMessage) {
	var msgPayload map[string]interface{}
	if err := json.Unmarshal(msg.Value, msgPayload); err != nil {
		log.Print("Error in unmarshal op. %v", err)
	}
	username, ok := msgPayload["oldUsername"].(string)
	if !ok {
		fmt.Println("Invalid data format: 'oldUsername' is missing or not a string")
		return
	}

	result := database.DB.DB.Where("username = ?", username).Delete(&models.User{})
	if result.Error != nil {
		log.Printf("Failed to delete user in web database: %v \n", result.Error)
	} else if result.RowsAffected == 0 { //bcs. if no user found the rowsAffectedValue is 0, but no error returns
		fmt.Println("No user found with the given username.")
		return
	} else {
		fmt.Printf("User with username '%s' deleted successfully.\n", username)
	}

}
