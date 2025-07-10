package services

import (
	"context"
	"fmt"
	"service-wallet/internal/models"
	"service-wallet/internal/models/response"
	"service-wallet/internal/store"
	"service-wallet/utils"
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

func (s *WalletService) CreateNewWallet(wallet *models.Wallet) (resp *response.ResponseHttp, err error) {
	var msgErr error = nil
	existWallet, err := s.repo.GetWalletByName(wallet.Name)
	if err != nil {
		msgErr = utils.MessageError("Repository::GetWalletByName", err)
		return &response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed created wallet '%s' [E001]", wallet.Name), MessageErr: msgErr.Error()}, msgErr
	} else if existWallet != nil {
		msgErr = fmt.Errorf("wallet with name '%s' already created", wallet.Name)
		return &response.ResponseHttp{IsError: true, Message: msgErr.Error(), MessageErr: msgErr.Error()}, msgErr
	}

	// TODO: Add request to user service for check user_id is valid

	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	walletId, err := s.repo.CreateNewWallet(*wallet)
	if err != nil {
		msgErr = utils.MessageError("Repository::CreateNewWallet", err)
		return &response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed created wallet '%s' [E002]", wallet.Name), MessageErr: msgErr.Error()}, msgErr
	}
	wallet.ID = walletId

	return &response.ResponseHttp{IsError: false, Message: fmt.Sprintf("Success created wallet '%s'", wallet.ID), Data: wallet}, nil
}

func (s *WalletService) GetWalletById(id string) (resp *response.ResponseHttp, err error) {
	var msgErr error = nil
	existWallet, err := s.repo.GetWalletById(id)
	if err != nil {
		msgErr = utils.MessageError("Repository::GetWalletById", err)
		return &response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed get wallet with id '%s' [E001]", id), MessageErr: msgErr.Error()}, msgErr
	} else if existWallet == nil {
		msgErr = fmt.Errorf("wallet with id '%s' not  found in table", id)
		return &response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed get wallet with id '%s' [E002]", id), MessageErr: msgErr.Error()}, msgErr
	}

	return &response.ResponseHttp{IsError: false, Message: fmt.Sprintf("Success created wallet '%s'", id), Data: existWallet}, nil
}

func (s *WalletService) Close() {
	s.cancel()
}
