package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"
	"go.uber.org/zap"

	httpClient "github.com/major1ink/simple-notification-telegram/internal/client/http"
	telegramClient "github.com/major1ink/simple-notification-telegram/internal/client/http/telegram"
	"github.com/major1ink/simple-notification-telegram/internal/config"
	kafkaConverter "github.com/major1ink/simple-notification-telegram/internal/converter/kafka"
	"github.com/major1ink/simple-notification-telegram/internal/converter/kafka/decoder"
	"github.com/major1ink/simple-notification-telegram/internal/service"
	assembledConsumer "github.com/major1ink/simple-notification-telegram/internal/service/consumer"
	telegramService "github.com/major1ink/simple-notification-telegram/internal/service/telegram"
	"github.com/major1ink/simple-notification-telegram/pkg/closer"
	wrappedKafka "github.com/major1ink/simple-notification-telegram/pkg/kafka"
	wrappedKafkaConsumer "github.com/major1ink/simple-notification-telegram/pkg/kafka/consumer"
)

type diContainer struct {
	assembleConsumerService service.ConsumerService
	telegramService         service.TelegramService

	assembledConsumerGroup sarama.ConsumerGroup

	assembledConsumer wrappedKafka.Consumer

	assembledDecoder kafkaConverter.OrderAssembledDecoder

	telegramClient httpClient.TelegramClient
	telegramBot    *bot.Bot

	logger *zap.Logger
	closer *closer.Closer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) SetLogger(l *zap.Logger) {
	d.logger = l
}

func (d *diContainer) SetCloser(c *closer.Closer) {
	d.closer = c
}

func (d *diContainer) AssembleConsumerService(ctx context.Context) service.ConsumerService {
	if d.assembleConsumerService == nil {
		d.assembleConsumerService = assembledConsumer.NewService(d.AssembledConsumer(), d.AssembledDecoder(), d.TelegramService(ctx), d.logger)
	}

	return d.assembleConsumerService
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(
			d.TelegramClient(ctx),
			d.logger,
			config.AppConfig().TelegramBot.GetTelegramChatID(),
		)
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.GetTelegramBotToken())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}

func (d *diContainer) AssembledConsumerGroup() sarama.ConsumerGroup {
	if d.assembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.GetBrokers(),
			config.AppConfig().Consumer.GetGroupId(),
			config.AppConfig().Consumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create assembled consumer group: %s\n", err.Error()))
		}
		d.closer.AddNamed("Kafka assembled consumer group", func(ctx context.Context) error {
			return d.assembledConsumerGroup.Close()
		})

		d.assembledConsumerGroup = consumerGroup
	}

	return d.assembledConsumerGroup
}

func (d *diContainer) AssembledConsumer() wrappedKafka.Consumer {
	if d.assembledConsumer == nil {
		d.assembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.AssembledConsumerGroup(),
			[]string{
				config.AppConfig().Consumer.GetTopic(),
			},
			d.logger,
		)
	}

	return d.assembledConsumer
}

func (d *diContainer) AssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.assembledDecoder == nil {
		d.assembledDecoder = decoder.NewOrderDecoderAssembled()
	}

	return d.assembledDecoder
}
