package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	reader   *kafka.Reader
	Messages chan kafka.Message
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			MaxBytes:    10e6,
			StartOffset: kafka.LastOffset,
			Logger:      kafka.LoggerFunc(log.Printf),
			ErrorLogger: kafka.LoggerFunc(log.Printf),
		}),
		Messages: make(chan kafka.Message, 10e2),
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return err
			}
			c.Messages <- m
		}
	}
}
