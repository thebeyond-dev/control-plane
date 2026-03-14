package devices

import (
	"context"
	"fmt"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type Interactor struct {
	repo     ports.DeviceRepository
	nodePool ports.NodePool
}

func NewInteractor(
	repo ports.DeviceRepository,
	nodePool ports.NodePool,
) ports.DeviceUseCase {
	return &Interactor{repo, nodePool}
}

func (uc *Interactor) List(ctx context.Context, userID string) ([]domain.Device, error) {
	return uc.repo.List(ctx, userID)
}

func (uc *Interactor) Create(ctx context.Context, userID, nodeID, name string, bandwidth int) (string, error) {
	vpnClient, err := uc.nodePool.GetVPNClient(ctx, nodeID)
	if err != nil {
		return "", fmt.Errorf("failed to get VPN client for country %s: %w", nodeID, err)
	}

	pubkey, config, err := vpnClient.CreatePeer(ctx, bandwidth)
	if err != nil {
		return "", err
	}

	return config, uc.repo.Save(ctx, domain.Device{
		UserID:    userID,
		PublicKey: pubkey,
		NodeID:    nodeID,
		Name:      name,
	})
}

func (uc *Interactor) Delete(ctx context.Context, userID, pubkey string) error {
	device, err := uc.repo.GetByPublicKey(ctx, userID, pubkey)
	if err != nil {
		return err
	}

	vpnClient, err := uc.nodePool.GetVPNClient(ctx, device.NodeID)
	if err != nil {
		return fmt.Errorf("failed to get VPN client for country %s: %w", device.NodeID, err)
	}

	if err := vpnClient.DeletePeer(ctx, pubkey); err != nil {
		return err
	}

	return uc.repo.Delete(ctx, userID, pubkey)
}
