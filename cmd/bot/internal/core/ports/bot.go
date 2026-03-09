package ports

import (
	"context"
)

type Bot interface {
	NewMessage(chatID int, text string) MessageBuilder
	DeleteMessage(ctx context.Context, chatID, messageID int) error
	ShowNotification(ctx context.Context, interactionID string, text string, alert bool) error
}
