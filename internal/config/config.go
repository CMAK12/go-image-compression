package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP HTTPConfig `envPrefix:"HTTP_CONFIG"`
}

type HTTPConfig struct {
	Port string `env:"HTTP_PORT"`
}

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
