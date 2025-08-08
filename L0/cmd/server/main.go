package main

import (
	"L0/internal/config"
	"L0/internal/db/postgres"
	"L0/internal/kafka"
	"context"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

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
	kafkaConsumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID)
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
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("storage", storage)
		c.Next()
	})
	r.GET("/order/:id", func(c *gin.Context) {
		orderID := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"order_id": orderID,
			"message":  "", //TODO
		})
	})
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}
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
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}

	log.Println("Server stopped")
}
