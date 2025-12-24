package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"service-wallet/internal/models"
	"service-wallet/utils"

	"github.com/go-redis/redis"
)

const (
	PrefixKeyWallet = "financial-tracking:wallet-service"
)

type RedisCln struct {
	Redis  *redis.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

func (r *RedisCln) SaveWallet(wallet models.Wallet) (key string, err error) {
	var msgErr error = nil
	walletJson, err := json.Marshal(wallet)
	if err != nil {
		msgErr = utils.MessageError("Redis::SaveWallet", fmt.Errorf("failed convert wallet '%s' to json", wallet.ID))
		return "", msgErr
	}

	// send wallet info to reddis data
	keyWalletInfo := fmt.Sprintf("%s:%s", PrefixKeyWallet, wallet.ID)
	err = r.Redis.Set(keyWalletInfo, walletJson, 0).Err()
	if err != nil {
		msgErr = utils.MessageError("Redis::SaveWallet", fmt.Errorf("error saving wallet data to redis. err = %s", err.Error()))
		return "", msgErr
	}

	return keyWalletInfo, nil
}

func (r *RedisCln) GetWallet(walletId string) (wallet *models.Wallet, err error) {
	var msgErr error = nil

	walletInfo := fmt.Sprintf("%s:%s", PrefixKeyWallet, walletId)
	walletJson, err := r.Redis.Get(walletInfo).Result()
	if err != nil {
		msgErr = utils.MessageError("Redis::GetWallet", fmt.Errorf("failed get wallet data from redis. err : %v", err))
		return nil, msgErr
	}
	err = json.Unmarshal([]byte(walletJson), &wallet)
	if err != nil {
		msgErr = utils.MessageError("Redis::GetWallet", fmt.Errorf("failed convert data wallet info from json. err : %v", err))
		return nil, msgErr
	}
	return wallet, nil
}

func (r *RedisCln) GetAllWallets() ([]models.Wallet, error) {
	var msgErr error = nil
	keysPattern := fmt.Sprintf("%s:*", PrefixKeyWallet)

	keys, err := r.Redis.Keys(keysPattern).Result()
	if err != nil {
		msgErr = utils.MessageError("Redis::GetAllWallets", fmt.Errorf("failed to get keys from redis with pattern %s. err: %v", keysPattern, err))
		return nil, msgErr
	}

	// MGET all of data in keys
	walletsJson, err := r.Redis.MGet(keys...).Result()
	if err != nil {
		msgErr = utils.MessageError("Redis::GetAllWallets", fmt.Errorf("failed to get multiple employments from redis. err: %v", err))
		return nil, msgErr
	}

	var wallets []models.Wallet
	for _, walletJson := range walletsJson {
		if walletJson == nil {
			continue
		}
		var wallet models.Wallet
		err := json.Unmarshal([]byte(walletJson.(string)), &wallet)
		if err != nil {
			msgErr = utils.MessageError("Redis::GetAllWallets", fmt.Errorf("failed to unmarshal employment json: %v", err))
			return nil, msgErr
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func (r *RedisCln) DeleteWallet(walletId string) error {
	key := fmt.Sprintf("%s:%s", PrefixKeyWallet, walletId)

	deleted, err := r.Redis.Del(key).Result()
	if err != nil {
		return utils.MessageError("Redis::DeleteWallet", fmt.Errorf("failed to delete wallet from redis: %w", err))
	}

	if deleted == 0 {
		return utils.MessageError("Redis::DeleteWallet", fmt.Errorf("wallet key %s not found", walletId))
	}

	return nil
}

func (r *RedisCln) Close() {
	r.Redis.Close()
	r.Cancel()
}
