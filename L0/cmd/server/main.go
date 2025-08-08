package main

import (
	"L0/internal/kafka"
	"context"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "test-topic"
	groupID := "test-group"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kafkaConsumer := kafka.NewConsumer(brokers, topic, groupID)
	defer func(kafkaConsumer *kafka.Consumer) {
		err := kafkaConsumer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(kafkaConsumer)
	go func() {
		if err := kafkaConsumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
			cancel()
		}
	}()
	r := gin.Default()
	r.Run(":8080")
}
