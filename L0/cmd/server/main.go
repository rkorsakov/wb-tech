package main

import (
	"L0/internal/config"
	"L0/internal/db/postgres"
	kfk "L0/internal/kafka"
	"L0/internal/models"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
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
	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*")
	r.Use(func(c *gin.Context) {
		c.Set("storage", storage)
		c.Next()
	})
	r.GET("/order/:id", func(c *gin.Context) {
		orderID := c.Param("id")
		ctx := c.Request.Context()
		order, err := storage.GetOrder(ctx, orderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Просмотр заказов",
		})
	})

	r.POST("/search", func(c *gin.Context) {
		orderID := c.PostForm("order_id")
		ctx := c.Request.Context()
		order, err := storage.GetOrder(ctx, orderID)
		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"Title":   "Просмотр заказов",
				"Error":   "Не удалось найти заказ: " + err.Error(),
				"OrderID": orderID,
			})
			return
		}
		prettyJSON, _ := json.MarshalIndent(order, "", "    ")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":   "Просмотр заказов",
			"Order":   string(prettyJSON),
			"OrderID": orderID,
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
