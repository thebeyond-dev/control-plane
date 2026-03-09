package ports

import (
	"context"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
)

type CommandHandler interface {
	Execute(ctx context.Context, msg input.Message) error
}
