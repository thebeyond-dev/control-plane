package ports

import (
	"context"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
)

type NotificationSender interface {
	Send(ctx context.Context, recipientID, text string) error
}

type AlertUseCase interface {
	NotifyNewPayment(ctx context.Context, user domain.User, amount, days int) error
	NotifySubscriptionExpiring(ctx context.Context, user domain.User, daysRemaining int) error
}
