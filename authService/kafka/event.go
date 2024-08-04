package kafka

import (
	"auth-service/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type requestPayload struct {
	NewUsername string `json:"new_username"`
	OldUsername string `json:"new_username"`
}

func ProduceEvent(data interface{}, topic string) {
	var msg *sarama.ProducerMessage
	producer := createProducer()
	defer producer.Close()

	if topic == "register-user" {
		msg = &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(data.(string)),
		}

	}
	if topic == "delete-user" {
		msg = &sarama.ProducerMessage{
			Topic: "delete-user",
			Value: sarama.StringEncoder(data.(string)),
		}
	}
	if topic == "update-user" {
		jsonData, err := json.Marshal(data.(models.RequestPayload))
		if err != nil {
			log.Printf("Error in marshal operation: %v", err)
			return
		}
		msg = &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(jsonData),
		}
	}
	if topic == "change-username" {
		log.Printf("Data: %+v", data)
		jsonData, err := json.Marshal(data.(models.RequestPayload))
		if err != nil {
			log.Printf("Error in marshal operation: %v", err)
			return
		}

		log.Printf("JSON Data: %s", string(jsonData))

		msg = &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(jsonData),
		}
	}

	fmt.Print("PRODUCER CREATED AND WAITING TO SEND MESSAGE...")
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Printf("failed to send message: %v\n", err)
	}
	fmt.Print("MESSAGE SENT")

}
