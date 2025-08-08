package main

import (
	"L0/internal/kafka"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "test-topic"
	groupID := "test-group"
	kafkaConsumer := kafka.NewConsumer(brokers, topic, groupID)
	defer func(kafkaConsumer *kafka.Consumer) {
		err := kafkaConsumer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(kafkaConsumer)
	err := kafkaConsumer.Start()
	if err != nil {
		log.Println(err)
	}
	r := gin.Default()
	r.Run(":8080")
}
