package config

import (
	"fmt"
	"os"

	structYaml "github.com/major1ink/simple-notification-telegram/internal/config/yaml"
	"gopkg.in/yaml.v3"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	Kafka       KafkaConfig
	Consumer    ConsumerConfig
	TelegramBot TelegramConfig
}

func Load(path ...string) error {
	var configPath string
	if len(path) > 0 && path[0] != "" {
		configPath = path[0]
	}

	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var yamlConfig struct {
		Logger   *structYaml.LoggerConfig   `yaml:"logger"`
		Kafka    *structYaml.KafkaConfig    `yaml:"kafkaConfig"`
		Consumer *structYaml.ConsumerConfig `yaml:"consumerConfig"`
		Telegram *structYaml.TelegramConfig `yaml:"telegramConfig"`
	}

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&yamlConfig); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}

	appConfig = &config{
		Logger:      yamlConfig.Logger,
		Kafka:       yamlConfig.Kafka,
		Consumer:    yamlConfig.Consumer,
		TelegramBot: yamlConfig.Telegram,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
