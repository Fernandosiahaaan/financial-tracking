package services

import (
	"context"
	"database/sql"
	"fmt"
	"service-wallet/internal/models"
	"service-wallet/internal/store"
	"time"
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

func (s *WalletService) CreateNewWallet(wallet *models.Wallet) (walletId string, err error) {
	existWallet, err := s.repo.GetWalletByName(wallet.Name)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("[repository] %v", err)
	} else if existWallet.Name == wallet.Name {
		return "", fmt.Errorf("[repository] wallet already created")
	}

	// TODO: Add request to user service for check user_id is valid

	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	walletId, err = s.repo.CreateNewWallet(*wallet)
	if err != nil {
		return "", fmt.Errorf("[repository] %v", err)
	}

	return walletId, nil
}

func (s *WalletService) Close() {
	s.cancel()
}
