package main

import (
	"L0/internal/config"
	"L0/internal/db/postgres"
	kfk "L0/internal/kafka"
	"L0/internal/models"
	"L0/internal/server"
	"context"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleMessage(m *kafka.Message, storage *postgres.Storage) error {
	var order models.Order
	err := json.Unmarshal(m.Value, &order)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = storage.SaveOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	storage, err := postgres.NewPostgresStorage(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()
	kafkaConsumer := kfk.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID)
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			log.Printf("Failed to close Kafka consumer: %v", err)
		}
	}()
	go func() {
		if err := kafkaConsumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
			cancel()
		}
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-kafkaConsumer.Messages:
				if !ok {
					return
				}
				if err := handleMessage(&msg, storage); err != nil {
					log.Printf("Failed to handle message: %v", err)
				}
			}
		}
	}()
	srv := server.NewServer(cfg, storage)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		cancel()
	}()
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}
	log.Println("Server stopped")
}
