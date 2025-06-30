package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"service-wallet/internal/models"
	"service-wallet/internal/store"
	"time"

	"github.com/google/uuid"
)

type WalletService struct {
	repo   *store.WalletStore
	ctx    context.Context
	cancel context.CancelFunc
	// red
}

func NewWalletService(ctx context.Context, repo *store.WalletStore) *WalletService {
	serviceCtx, serviceCancel := context.WithCancel(ctx)
	return &WalletService{
		repo:   repo,
		ctx:    serviceCtx,
		cancel: serviceCancel,
	}
}

func (s *WalletService) CreateNewWallet(wallet models.Wallet) (walletId string, err error) {
	existWallet, err := s.repo.GetWalletByName(wallet.Name)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return "", fmt.Errorf("[repository] %v", err)
	} else if existWallet.Name == wallet.Name {
		return "", fmt.Errorf("[repository] wallet already created", err)
	}

	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	wallet.ID = uuid.New().String()

	walletId, err = s.repo.CreateNewWallet(wallet)
	if err != nil {
		return "", fmt.Errorf("[repository] %v", err)
	}

	return walletId, err
}
