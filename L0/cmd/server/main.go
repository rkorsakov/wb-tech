package main

import (
	"L0/internal/config"
	"L0/internal/kafka"
	"context"
	"github.com/gin-gonic/gin"
	"log"
)

type Config struct {
}

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kafkaConsumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID)
	go func() {
		if err := kafkaConsumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
			cancel()
		}
	}()
	r := gin.Default()
	r.Run(cfg.Server.Port)
}
