package bot

import (
	"github.com/go-telegram/bot"
	"github.com/thebeyond-net/control-plane/cmd/bot/config"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/adapters/telegram/webhook"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"go.uber.org/zap"
)

type Bot struct {
	client *bot.Bot
}

func New(
	bot *bot.Bot,
	appLogger *zap.Logger,
	handler *webhook.Handler,
	cfg *config.Config,
) ports.Bot {
	return &Bot{bot}
}
