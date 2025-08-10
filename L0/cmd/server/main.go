package main

import (
	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/db/postgres"
	kfk "L0/internal/kafka"
	"L0/internal/server"
	"context"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	storage, err := postgres.NewPostgresStorage(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()
	orderCache := cache.New()
	orders, err := storage.GetAllOrders(ctx)
	if err != nil {
		log.Fatalf("Failed to get all orders: %v", err)
	}
	for _, order := range orders {
		orderCache.Set(order.OrderUID, order)
	}
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
	messageHandler := kfk.NewMessageHandler(storage)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-kafkaConsumer.Messages:
				if !ok {
					return
				}
				if err := messageHandler.HandleMessage(orderCache, &msg); err != nil {
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
