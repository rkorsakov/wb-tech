package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Kafka  KafkaConfig  `yaml:"kafka"`
	Server ServerConfig `yaml:"server"`
	Cache  CacheConfig  `yaml:"cache"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	Topic   string   `yaml:"topic" env:"KAFKA_TOPIC" env-default:"test-topic"`
	GroupID string   `yaml:"group_id" env:"KAFKA_GROUP_ID" env-default:"test-group"`
}

type ServerConfig struct {
	Port string `yaml:"port" env:"SERVER_PORT" env-default:":8080"`
}

type CacheConfig struct {
	Capacity int `yaml:"capacity" env:"CACHE_CAPACITY" env-default:"1000"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
