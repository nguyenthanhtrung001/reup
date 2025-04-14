package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Telegram interface {
	SendMessage(chatID int64, text, mode string) (tgbotapi.Message, error)
}
type ChatIDs struct {
	ReportBug     int64
	ReportPayment int64
}
type implTelegram struct {
	bot     *tgbotapi.BotAPI
	chatIDs ChatIDs
}

func New(bot_key string, chatIDS ChatIDs) Telegram {
	bot, err := tgbotapi.NewBotAPI(bot_key)
	if err != nil {
		panic(err)
	}
	return &implTelegram{
		bot:     bot,
		chatIDs: chatIDS,
	}
}
