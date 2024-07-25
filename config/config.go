package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DatabaseURL  string   `yaml:"database_url"`
	KafkaBrokers []string `yaml:"kafka_brokers"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
