package config

import (
	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	GetLevel() string
	GetLogDir() string
	GetLogMode() string
	GetRewriteLog() bool
}

type KafkaConfig interface {
	GetBrokers() []string
}

type ConsumerConfig interface {
	GetTopic() string
	GetGroupId() string
	Config() *sarama.Config
}

type TelegramConfig interface {
	GetTelegramBotToken() string
	GetTelegramChatID() int64
}
