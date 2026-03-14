package usersettings

import (
	"context"

	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type Interactor struct {
	userRepo ports.UserRepository
}

func NewInteractor(userRepo ports.UserRepository) ports.UserSettingsUseCase {
	return &Interactor{userRepo}
}

func (uc *Interactor) SetLanguage(ctx context.Context, userID, languageCode string) error {
	return uc.userRepo.UpdateLanguage(ctx, userID, languageCode)
}

func (uc *Interactor) SetCurrency(ctx context.Context, userID, currencyCode string) error {
	return uc.userRepo.UpdateCurrency(ctx, userID, currencyCode)
}
