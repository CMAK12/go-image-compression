package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/nordew/go-errx"
)

type (
	Config struct {
		HTTP  HTTPConfig  `envPrefix:"HTTP_CONFIG"`
		Minio MinioConfig `envPrefix:"MINIO_CONFIG"`
		NATS  NATSConfig  `envPrefix:"NATS_CONFIG"`
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

	NATSConfig struct {
		URL string `env:"NATS_URL"`
	}
)

const codepath = "config/config.go"

func MustLoad() (*Config, error) {
	config := &Config{}

	if err := godotenv.Load(); err != nil {
		return nil, errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}

	if err := cleanenv.ReadEnv(config); err != nil {
		return nil, errx.NewInternal().WithDescriptionAndCause(codepath, err)
	}

	return config, nil
}
