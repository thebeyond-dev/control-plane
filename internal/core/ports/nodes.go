package ports

import (
	"context"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
)

type NodeRepository interface {
	Get(ctx context.Context, id string) (domain.Node, error)
	List(ctx context.Context) ([]domain.Node, error)
}

type VPNClient interface {
	CreatePeer(ctx context.Context, bandwidth int) (pubkey, config string, err error)
	DeletePeer(ctx context.Context, pubkey string) error
}

type NodePool interface {
	GetVPNClient(ctx context.Context, nodeID string) (VPNClient, error)
}

type NodeUseCase interface {
	ListLocations(ctx context.Context) ([]domain.Node, error)
}
