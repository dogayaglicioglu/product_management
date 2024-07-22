package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

func ProduceEvent(username string) {
	producer := createProducer()
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: "register-user",
		Value: sarama.StringEncoder(username),
	}
	fmt.Print("pRODUCER CREATED AND WAITING TO SEND MESSAGE...")
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Printf("failed to send message: %v\n", err)
	}
	fmt.Print("MESSAGE SENT")

}
