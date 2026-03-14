package node

import (
	"context"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type Interactor struct {
	repo     ports.NodeRepository
	nodePool ports.NodePool
}

func NewInteractor(
	repo ports.NodeRepository,
	nodePool ports.NodePool,
) ports.NodeUseCase {
	return &Interactor{repo, nodePool}
}

func (uc *Interactor) ListLocations(ctx context.Context) ([]domain.Node, error) {
	return uc.repo.List(ctx)
}
