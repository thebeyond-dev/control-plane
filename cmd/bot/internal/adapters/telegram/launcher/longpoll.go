package launcher

import (
	"context"

	"github.com/go-telegram/bot"
)

type LongPoll struct {
	bot *bot.Bot
}

func NewLongPoll(bot *bot.Bot) *LongPoll {
	return &LongPoll{bot}
}

func (lp *LongPoll) Launch(ctx context.Context) {
	lp.bot.Start(ctx)
}
