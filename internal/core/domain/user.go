package domain

import (
	"context"
	"time"
)

type User struct {
	ID string
	Subscription
	CreatedAt time.Time
}

type UserRepository interface {
	GetByIdentity(ctx context.Context, provider, providerID string) (User, error)
	CreateWithIdentity(ctx context.Context, user User, identity Identity) error
}

type AuthService interface {
	Login(ctx context.Context, provider, providerID string) (User, error)
}
