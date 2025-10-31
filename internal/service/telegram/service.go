package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/major1ink/simple-notification-telegram/internal/client/http"
	"github.com/major1ink/simple-notification-telegram/internal/model"
)

//go:embed templates/assembled_notification.tmpl
var templateFS embed.FS

type assembledTemplateData struct {
	EventUuid string
	TypeEvent string
	App       string
	Message   string
}

var (
	assembledTemplate = template.Must(template.ParseFS(templateFS, "templates/assembled_notification.tmpl"))
)

type service struct {
	telegramClient http.TelegramClient
	logger         *zap.Logger
	chatID         int64
}

func NewService(telegramClient http.TelegramClient, logger *zap.Logger, chatID int64) *service {
	return &service{
		telegramClient: telegramClient,
		logger:         logger,
		chatID:         chatID,
	}
}

func (s *service) SendAssembledNotification(ctx context.Context, assembledEvent model.AssembledEvent) error {
	message, err := s.buildAssembledMessage(assembledEvent)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, s.chatID, message)
	if err != nil {
		return err
	}

	s.logger.Debug("Telegram message sent to chat", zap.Int64("chat_id", s.chatID), zap.String("message", message))
	return nil
}

func (s *service) buildAssembledMessage(assembledEvent model.AssembledEvent) (string, error) {
	data := assembledTemplateData{
		EventUuid: assembledEvent.EventUuid,
		TypeEvent: assembledEvent.TypeEvent,
		App:       assembledEvent.App,
		Message:   assembledEvent.Message,
	}

	var buf bytes.Buffer
	err := assembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
