package billing

import (
	"context"

	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type Interactor struct {
	subscriptionRepo ports.SubscriptionRepository
}

func NewInteractor(subscriptionRepo ports.SubscriptionRepository) ports.BillingUseCase {
	return &Interactor{subscriptionRepo}
}

func (uc *Interactor) RenewSubscription(ctx context.Context, userID string, devices, days int) error {
	return uc.subscriptionRepo.Update(ctx, userID, devices, days)
}

func (uc *Interactor) CancelSubscription(ctx context.Context, userID string) error {
	return uc.subscriptionRepo.Deactivate(ctx, userID)
}
