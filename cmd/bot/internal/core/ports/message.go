package ports

import (
	"context"
	"io"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
)

type MessageBuilder interface {
	WithReplyMarkup(markup interaction.InlineKeyboardMarkup) MessageBuilder
	WithFile(name string, content io.Reader) MessageBuilder
	Send(ctx context.Context) error
	Edit(ctx context.Context, messageID int) error
}
