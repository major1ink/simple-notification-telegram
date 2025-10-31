package yaml

type TelegramConfig struct {
	TelegramBotToken string `yaml:"telegram_bot_token"`
	TelegramChatID   int64  `yaml:"telegram_chat_id"`
}

func (t *TelegramConfig) GetTelegramBotToken() string {
	return t.TelegramBotToken
}

func (t *TelegramConfig) GetTelegramChatID() int64 {
	return t.TelegramChatID
}
