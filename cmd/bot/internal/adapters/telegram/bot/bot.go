package bot

import (
	"github.com/go-telegram/bot"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"go.uber.org/zap"
)

type Bot struct {
	client *bot.Bot
}

func New(bot *bot.Bot, appLogger *zap.Logger) ports.Bot {
	return &Bot{bot}
}
