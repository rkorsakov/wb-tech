package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	reader *kafka.Reader
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
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) Start() error {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
