package consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/major1ink/simple-notification-telegram/internal/converter/kafka"
	def "github.com/major1ink/simple-notification-telegram/internal/service"
	"github.com/major1ink/simple-notification-telegram/pkg/kafka"
)

type service struct {
	consumer        kafka.Consumer
	decoder         kafkaConverter.OrderAssembledDecoder
	telegramService def.TelegramService
	logger          *zap.Logger
}

func NewService(
	consumer kafka.Consumer,
	decoder kafkaConverter.OrderAssembledDecoder,
	telegramService def.TelegramService,
	logger *zap.Logger,
) *service {
	return &service{
		consumer:        consumer,
		decoder:         decoder,
		telegramService: telegramService,
		logger:          logger,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	s.logger.Info("Starting consumer service")

	err := s.consumer.Consume(ctx, s.Handler)
	if err != nil {
		s.logger.Error("Consume topic error", zap.Error(err))
		return err
	}

	return nil
}
