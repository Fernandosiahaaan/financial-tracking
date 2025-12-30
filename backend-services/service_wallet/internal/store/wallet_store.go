package store

import (
	"context"
	"fmt"
	"service-wallet/infrastructure"
	"service-wallet/internal/models"
	"service-wallet/internal/models/request"
	"service-wallet/internal/store/postgree"
	"service-wallet/internal/store/redis"
	"service-wallet/utils"

	_ "github.com/lib/pq"
)

type WalletStore struct {
	dbPq   postgree.StoreWalletPq
	cache  redis.RedisCln
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWalletStore(ctx context.Context) (error, *WalletStore) {
	repoPostgree, err := infrastructure.NewConnectionPostgree(ctx)
	if err != nil {
		errMsg := utils.MessageError("infrastructure::NewConnectionPostgree", err)
		return errMsg, nil
	}

	cacheRedis, err := infrastructure.NewConnectionRedis(ctx)
	if err != nil {
		errMsg := utils.MessageError("infrastructure::NewConnectionRedis", err)
		return errMsg, nil
	}

	serviceCtx, serviceCancel := context.WithCancel(ctx)
	return nil, &WalletStore{
		ctx:    serviceCtx,
		cancel: serviceCancel,
		dbPq:   *repoPostgree,
		cache:  *cacheRedis,
	}
}

func (s *WalletStore) CreateNewWallet(wallet models.Wallet) (id string, err error) {
	id, err = s.dbPq.CreateNewWallet(wallet)
	if err != nil {
		errMsg := utils.MessageError("postgree::CreateNewWallet", err)
		return id, errMsg
	}

	// Get data by Id and save to cache
	item, err := s.dbPq.GetWalletById(id)
	if err != nil {
		errMsg := utils.MessageError("postgree::GetWalletById", err)
		return id, errMsg
	}

	s.cache.SaveWallet(*item)

	return id, nil
}

func (s *WalletStore) GetWalletByName(walletName string) (*models.Wallet, error) {
	wallet, err := s.dbPq.GetWalletByName(walletName)
	if err != nil {
		errMsg := utils.MessageError("postgree::GetWalletByName", err)
		return nil, errMsg
	}
	return wallet, err
}

func (s *WalletStore) GetWalletById(walletId string) (*models.Wallet, error) {

	// get from cache
	fmt.Println("walletId2 = ", walletId)

	wallet, err := s.cache.GetWallet(walletId)
	if (err == nil) && (wallet != nil) {
		return wallet, err
	}

	fmt.Println("wallet2 = ", wallet)

	// get from db
	walletData, err := s.dbPq.GetWalletById(walletId)
	if err != nil {
		errMsg := utils.MessageError("postgree::GetWalletById", err)
		return nil, errMsg
	}

	fmt.Println("walletData2 = ", walletData)

	if walletData != nil {
		s.cache.SaveWallet(*walletData)
	}

	return walletData, err
}

func (r *WalletStore) UpdateWalletById(wallet models.Wallet) error {
	err := r.dbPq.UpdateWalletById(wallet)
	if err != nil {
		errMsg := utils.MessageError("postgree::UpdateWalletById", err)
		return errMsg
	}

	// update wallet in cache
	if err = r.cache.DeleteWallet(wallet.ID); err != nil {
		newWallet, err := r.dbPq.GetWalletById(wallet.ID)
		if err == nil {
			r.cache.SaveWallet(*newWallet)
		}
	}

	return nil
}

func (r *WalletStore) DeleteWalletById(id string) error {
	err := r.dbPq.DeleteWalletById(id)
	if err != nil {
		errMsg := utils.MessageError("postgree::DeleteWalletById", err)
		return errMsg
	}
	r.cache.DeleteWallet(id)

	return nil
}

func (r *WalletStore) GetListWallets(params request.GetListWalletRequest) ([]models.Wallet, int, error) {
	wallets, total, err := r.dbPq.GetListWallets(params)
	if err != nil {
		errMsg := utils.MessageError("postgree::GetListWallets", err)
		return wallets, total, errMsg
	}
	return wallets, total, err
}

func (s *WalletStore) Close() {
	s.dbPq.Close()
	s.cache.Close()
}
