package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ApiKey    string `env:"API_KEY"`
	RangersId string `env:"RANGERS_ID"`
}

func New() (*Config, error) {
	var config Config
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}

	if config.ApiKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable is not set")
	}

	if config.RangersId == "" {
		return nil, fmt.Errorf("RANGERS_ID environment variable is not set")
	}

	return &config, nil
}
