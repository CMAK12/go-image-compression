package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP  HTTPConfig  `envPrefix:"HTTP_CONFIG"`
		Minio MinioConfig `envPrefix:"MINIO_CONFIG"`
	}

	HTTPConfig struct {
		Port string `env:"HTTP_PORT"`
	}

	MinioConfig struct {
		Endpoint        string `env:"MINIO_ENDPOINT"`
		AccessKeyID     string `env:"MINIO_ACCESS_KEY_ID"`
		SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY"`
		SSLMode         bool   `env:"MINIO_SSL_MODE"`
	}
)

func MustLoad() (*Config, error) {
	config := &Config{}

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("config/config.go: %v", err)
	}

	if err := cleanenv.ReadEnv(config); err != nil {
		return nil, fmt.Errorf("config/config.go: %v", err)
	}

	return config, nil
}
