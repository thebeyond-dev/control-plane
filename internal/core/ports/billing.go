package ports

import (
	"context"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
)

type Invoice interface {
	NewPayment(ctx context.Context, user domain.User, currency string, devices, days int, price float64) (string, error)
}

type SubscriptionRepository interface {
	Update(ctx context.Context, userID string, devices, days int) error
	Deactivate(ctx context.Context, userID string) error
}

type BillingUseCase interface {
	RenewSubscription(ctx context.Context, userID string, devices, days int) error
	CancelSubscription(ctx context.Context, userID string) error
}
