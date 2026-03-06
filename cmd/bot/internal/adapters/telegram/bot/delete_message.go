package bot

import (
	"context"

	"github.com/go-telegram/bot"
)

func (a *Bot) DeleteMessage(
	ctx context.Context,
	chatID, messageID int,
) error {
	_, err := a.client.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: messageID,
	})
	return err
}
