package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func ProduceEvent(data interface{}, topic string) {
	var msg *sarama.ProducerMessage
	producer := createProducer()
	defer producer.Close()

	if topic == "register-user" {
		msg = &sarama.ProducerMessage{
			Topic: "register-user",
			Value: sarama.StringEncoder(data.(string)),
		}

	}
	if topic == "change-username" {
		jsonData, err := json.Marshal(data) //struct yollandı
		if err != nil {
			log.Print("Error in marshal operation %v", err)
		}

		msg = &sarama.ProducerMessage{
			Topic: "change-username",
			Value: sarama.ByteEncoder(jsonData),
		}
	}
	if topic == "delete-user" {
		jsonData, err := json.Marshal(data)
		if err != err {
			log.Print("Error in marshal operation %v", err)
		}

		msg = &sarama.ProducerMessage{
			Topic: "delete-user",
			Value: sarama.ByteEncoder(jsonData),
		}
	}

	fmt.Print("pRODUCER CREATED AND WAITING TO SEND MESSAGE...")
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Printf("failed to send message: %v\n", err)
	}
	fmt.Print("MESSAGE SENT")

}
