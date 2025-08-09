package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "test-topic",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()
	files, err := os.ReadDir("testdata")
	if err != nil {
		log.Fatalf("Failed to read testdata directory: %v", err)
	}
	for _, file := range files {
		filePath := filepath.Join("testdata", file.Name())
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v", file.Name(), err)
			continue
		}
		var order map[string]interface{}
		if err := json.Unmarshal(jsonData, &order); err != nil {
			log.Printf("Invalid JSON in file %s: %v", file.Name(), err)
			continue
		}
		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: jsonData,
			},
		)
		if err != nil {
			log.Printf("Failed to send message from file %s: %v", file.Name(), err)
			continue
		}
		log.Printf("Successfully sent message from file: %s", file.Name())
	}
	log.Println("All files processed")
}
