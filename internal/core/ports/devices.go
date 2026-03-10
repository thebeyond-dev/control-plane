package ports

import (
	"context"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
)

type DeviceRepository interface {
	GetByPublicKey(ctx context.Context, userID, pubkey string) (domain.Device, error)
	List(ctx context.Context, userID string) ([]domain.Device, error)
	Save(ctx context.Context, device domain.Device) error
	Delete(ctx context.Context, userID, pubkey string) error
}

type DeviceUseCase interface {
	List(ctx context.Context, userID string) ([]domain.Device, error)
	Create(ctx context.Context, userID, nodeID, name string, bandwidth int) (string, error)
	Delete(ctx context.Context, userID, pubkey string) error
}
