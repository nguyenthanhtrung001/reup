package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t implTelegram) SendMessage(chatID int64, text, mode string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	if mode != "" {
		msg.ParseMode = mode
	}

	return t.bot.Send(msg)
}
