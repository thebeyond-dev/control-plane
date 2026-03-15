package tgnotifier

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type BotProvider interface {
	GetBot() *bot.Bot
}

type NotificationSender struct {
	provider BotProvider
}

func New(provider BotProvider) ports.NotificationSender {
	return &NotificationSender{provider}
}

func (s *NotificationSender) Send(ctx context.Context, recipientID string, text string) error {
	b := s.provider.GetBot()
	if b == nil {
		return fmt.Errorf("bot is not initialized yet")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: recipientID,
		Text:   text,
	})
	return err
}
