package domain

import "context"

type Node struct {
	ID          string
	Address     string
	LoadPercent int
}

type NodeRepository interface {
	GetActive(ctx context.Context) ([]Node, error)
	UpdateLoad(ctx context.Context, id string, load int) error
}

type NodeService interface {
	GetOptimizedNode(ctx context.Context) (Node, error)
}
