package alert

import (
	"context"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

type UseCase struct {
	notificationSender ports.NotificationSender
	currencies         domain.Items
}

func New(
	notificationSender ports.NotificationSender,
	currencies domain.Items,
) ports.AlertUseCase {
	return &UseCase{notificationSender, currencies}
}

func (uc *UseCase) NotifyNewPayment(ctx context.Context, user domain.User, amount, days int) error {
	currency, _ := uc.currencies.Get(user.CurrencyCode)
	text := i18n.Get(user.LanguageCode, "NewPayment", map[string]any{
		"Amount":         amount,
		"CurrencySymbol": currency.Symbol,
		"Devices":        user.Devices,
		"Days":           days,
	}, nil)

	if err := uc.notificationSender.Send(ctx, user.Identity.ProviderID, text); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) NotifySubscriptionExpiring(ctx context.Context, user domain.User, daysRemaining int) error {
	panic("not implemented")
}
