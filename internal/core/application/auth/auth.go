package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/thebeyond-net/control-plane/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type Interactor struct {
	userRepo     ports.UserRepository
	newbieConfig input.NewbieConfig
}

func NewInteractor(
	userRepo ports.UserRepository,
	newbieConfig input.NewbieConfig,
) ports.AuthUseCase {
	return &Interactor{userRepo, newbieConfig}
}

func (uc *Interactor) Login(ctx context.Context, provider, providerID string, input input.Login) (domain.User, error) {
	user, err := uc.userRepo.GetByIdentity(ctx, provider, providerID)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, err
	}

	if err := uc.userRepo.CreateWithIdentity(ctx, domain.User{
		Devices:   uc.newbieConfig.Devices,
		Bandwidth: uc.newbieConfig.Bandwidth,
		Subscription: domain.Subscription{
			ExpiresAt: time.Now().Add(uc.newbieConfig.TrialDuration),
		},
		Referral: domain.Referral{
			ReferrerID: input.ReferrerID,
		},
		LanguageCode: uc.newbieConfig.LanguageCode,
		CurrencyCode: uc.newbieConfig.CurrencyCode,
	}, domain.Identity{
		Provider:   provider,
		ProviderID: providerID,
	}); err != nil {
		return domain.User{}, err
	}

	return uc.userRepo.GetByIdentity(ctx, provider, providerID)
}
