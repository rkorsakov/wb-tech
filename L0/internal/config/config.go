package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Kafka    KafkaConfig    `yaml:"kafka"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	Topic   string   `yaml:"topic" env:"KAFKA_TOPIC" env-default:"test-topic"`
	GroupID string   `yaml:"group_id" env:"KAFKA_GROUP_ID" env-default:"test-group"`
}

type ServerConfig struct {
	Port string `yaml:"port" env:"SERVER_PORT" env-default:":8080"`
}

type DatabaseConfig struct {
	URL string `yaml:"url" env:"DATABASE_URL" env-default:"postgres://postgres:postgres@localhost/postgres?sslmode=disable"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
