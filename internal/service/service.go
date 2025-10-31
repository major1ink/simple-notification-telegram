package service

import (
	"context"

	"github.com/major1ink/simple-notification-telegram/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendAssembledNotification(ctx context.Context, assembledEvent model.AssembledEvent) error
}
