package router

import (
	"strconv"
	"strings"

	"github.com/go-telegram/bot/models"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
)

func toInputMessage(update *models.Update) (input.Message, string) {
	var (
		msg     input.Message
		userID  int64
		rawText string
	)

	switch {
	case update.Message != nil:
		message := update.Message
		msg.ChatID = int(message.Chat.ID)
		msg.ID = message.ID
		userID = message.From.ID
		rawText = strings.TrimPrefix(message.Text, "/")
	case update.CallbackQuery != nil:
		callbackQuery := update.CallbackQuery
		userID = callbackQuery.From.ID
		msg.InteractionID = callbackQuery.ID
		rawText = callbackQuery.Data

		callbackQueryMessage := callbackQuery.Message.Message
		if callbackQueryMessage != nil {
			msg.ChatID = int(callbackQueryMessage.Chat.ID)
			msg.ID = callbackQueryMessage.ID
		}
	default:
		return input.Message{}, ""
	}

	msg.Text, msg.Args = parseInput(rawText)

	return msg, strconv.FormatInt(userID, 10)
}

func parseInput(raw string) (string, []string) {
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}
