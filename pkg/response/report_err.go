package response

import (
	"log"

	"github.com/nguyenthanhtrung001/reup/pkg/telegram"

	"github.com/gin-gonic/gin"
)

func sendServerTelegramMessageAsync(message string, c *gin.Context, t telegram.Telegram, chatID int64) {
	go func() {
		splitMessages := splitMessageForTelegram(message)
		for _, message := range splitMessages {
			_, err := t.SendMessage(chatID, message, "")
			if err != nil {
				log.Printf("Error sending Telegram message: %v\n", err)
			}
		}
	}()
}

func splitMessageForTelegram(message string) []string {
	const maxMessageLength = 4096
	var messages []string
	for len(message) > maxMessageLength {
		messages = append(messages, message[:maxMessageLength])
		message = message[maxMessageLength:]
	}
	messages = append(messages, message)
	return messages
}
