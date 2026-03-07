package ports

import (
	"context"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
)

type UserRepository interface {
	GetByIdentity(ctx context.Context, provider, providerID string) (domain.User, error)
	UpdateLanguage(ctx context.Context, userID, languageCode string) error
	UpdateCurrency(ctx context.Context, userID, currencyCode string) error
	CreateWithIdentity(ctx context.Context, user domain.User, identity domain.Identity) error
}

type AuthUseCase interface {
	Login(ctx context.Context, provider, providerID string) (domain.User, error)
}

type UserSettingsUseCase interface {
	SetLanguage(ctx context.Context, userID, languageCode string) error
	SetCurrency(ctx context.Context, userID, currencyCode string) error
}
