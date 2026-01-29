package domain

import (
	"context"
	"time"
)

type Subscription struct {
	NodeID      string
	TariffID    int
	TrafficUsed int64
	StartAt     time.Time
	ExpiresAt   time.Time
	IsActive    bool
}

type SubscriptionRepository interface {
	GetByUserID(ctx context.Context, userID string) (Subscription, error)
	UpdateTraffic(ctx context.Context, userID string, delta int64) error
	ChangeTariff(ctx context.Context, userID string, newTariffID int) error
	Activate(ctx context.Context, userID string, expiresAt time.Time) error
	Deactivate(ctx context.Context, userID string) error
}

type BillingService interface {
	ProcessTraffic(ctx context.Context, userID string, bytes int64) error
	PurchaseSubscription(ctx context.Context, userID string, tariffID int) error
	RenewSubscription(ctx context.Context, userID string) error
	CancelSubscription(ctx context.Context, userID string) error
}
