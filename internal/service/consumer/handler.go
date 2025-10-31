package consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/major1ink/simple-notification-telegram/pkg/kafka/consumer"
)

func (s *service) Handler(ctx context.Context, msg consumer.Message) error {
	event, err := s.decoder.DecodeAssembled(msg.Value)
	if err != nil {
		s.logger.Error("Failed to decode assembled event", zap.Error(err))
		return err
	}

	return s.telegramService.SendAssembledNotification(ctx, event)
}
