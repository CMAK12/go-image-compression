package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP  HTTPConfig  `envPrefix:"HTTP_CONFIG"`
		GRPC  GRPCConfig  `envPrefix:"GRPC_CONFIG"`
		Minio MinioConfig `envPrefix:"MINIO_CONFIG"`
		NATS  NATSConfig  `envPrefix:"NATS_CONFIG"`
		Topic Topic       `envPrefix:"TOPIC_CONFIG"`
	}

	HTTPConfig struct {
		Port string `env:"HTTP_PORT"`
	}

	GRPCConfig struct {
		Port string `env:"GRPC_PORT"`
	}

	MinioConfig struct {
		Endpoint        string `env:"MINIO_ENDPOINT"`
		AccessKeyID     string `env:"MINIO_ACCESS_KEY_ID"`
		SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY"`
		SSLMode         bool   `env:"MINIO_SSL_MODE"`
	}

	NATSConfig struct {
		URL string `env:"NATS_URL"`
	}

	Topic struct {
		ImageCreated string `env:"IMAGE_CREATED_TOPIC"`
		ImageStream  string `env:"IMAGE_STREAM"`
	}
)

func MustLoad() (*Config, error) {
	config := &Config{}
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("config.MustLoad: %w", err)
	}

	if err := cleanenv.ReadEnv(config); err != nil {
		return nil, fmt.Errorf("config.MustLoad: %w", err)
	}

	return config, nil
}
