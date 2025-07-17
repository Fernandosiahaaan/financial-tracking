package services

import (
	"context"
	"fmt"
	"math"
	"service-wallet/internal/models"
	"service-wallet/internal/models/request"
	"service-wallet/internal/models/response"
	"service-wallet/internal/store"
	"service-wallet/utils"
	"strconv"
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

func (s *WalletService) UpdateWalletById(id string, walletUpdate models.Wallet) (resp *response.ResponseHttp, err error) {
	err = nil
	var msgErr error = nil

	// Check wallet if exist
	_, err = s.repo.GetWalletById(id)
	if err != nil {
		msgErr = utils.MessageError("Repository::GetWalletById", err)
		return &response.ResponseHttp{IsError: true, Message: fmt.Sprintf("Unknown wallet with id '%s' [E001]", id), MessageErr: msgErr.Error()}, msgErr
	}

	// Update Wallet
	walletUpdate.ID = id
	walletUpdate.UpdatedAt = time.Now()
	err = s.repo.UpdateWalletById(walletUpdate)
	if err != nil {
		msgErr = utils.MessageError("Repository::UpdateWalletById", err)
		return &response.ResponseHttp{IsError: true, Message: fmt.Sprintf("Failed update wallet with id '%s' [E002]", walletUpdate.ID), MessageErr: msgErr.Error()}, msgErr
	}

	return &response.ResponseHttp{IsError: false, Message: fmt.Sprintf("Success update wallet with id '%s'", walletUpdate.ID)}, nil
}

func (s *WalletService) DeleteWalletById(id string) (resp *response.ResponseHttp, err error) {
	err = nil
	var msgErr error = nil

	// Check wallet if exist
	_, err = s.repo.GetWalletById(id)
	if err != nil {
		msgErr = utils.MessageError("Repository::GetWalletById", err)
		return &response.ResponseHttp{IsError: true, Message: fmt.Sprintf("Unknown wallet with id '%s' [E001]", id), MessageErr: msgErr.Error()}, msgErr
	}

	// Delete Wallet
	err = s.repo.DeleteWalletById(id)
	if err != nil {
		msgErr = utils.MessageError("Repository::DeleteWalletById", err)
		return &response.ResponseHttp{IsError: true, Message: fmt.Sprintf("Failed delete wallet with id '%s' [E002]", id), MessageErr: msgErr.Error()}, msgErr
	}

	return &response.ResponseHttp{IsError: false, Message: fmt.Sprintf("Success delete wallet with id '%s'", id)}, nil
}

func (s *WalletService) GetListWallets(params request.GetListWalletRequest) (resp *response.GetListWalletResponse, err error) {
	err = nil
	var msgErr error = nil

	// Check wallet if exist
	datas, total, err := s.repo.GetListWallets(params)
	if err != nil {
		msgErr = utils.MessageError("Repository::GetListWallets", err)
		return &response.GetListWalletResponse{IsError: true, Message: fmt.Sprintf("Failed Get List Wallet [E001]"), MessageErr: msgErr.Error()}, msgErr
	}

	// Hitung pagination
	pageNumber, _ := strconv.Atoi(params.Page)
	if pageNumber <= 0 {
		pageNumber = 1
	}
	pageSize, _ := strconv.Atoi(params.PageItem)
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNumber - 1) * pageSize

	pages := int(math.Ceil(float64(total) / float64(pageSize)))

	startItem := offset + 1
	endItem := offset + pageSize
	if endItem > total {
		endItem = total
	}

	return &response.GetListWalletResponse{
		IsError:      false,
		Message:      "Success",
		Start:        startItem,
		End:          endItem,
		Page:         pageNumber,
		Pages:        pages,
		RecordsTotal: total,
		Data:         datas,
	}, nil
}

func (s *WalletService) Close() {
	s.cancel()
}
