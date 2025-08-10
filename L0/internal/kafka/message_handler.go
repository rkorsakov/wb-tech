package kafka

import (
	"L0/internal/db/postgres"
	"L0/internal/models"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type MessageHandler struct {
	storage *postgres.Storage
}

func NewMessageHandler(storage *postgres.Storage) *MessageHandler {
	return &MessageHandler{storage: storage}
}

func (handler *MessageHandler) HandleMessage(m *kafka.Message) error {
	var order models.Order
	err := json.Unmarshal(m.Value, &order)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = handler.storage.SaveOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}
