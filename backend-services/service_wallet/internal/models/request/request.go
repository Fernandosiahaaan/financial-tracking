package request

import (
	"service-wallet/internal/models"
)

type CreateWallet struct {
	ReqID   string            `json:"req_id"`
	UserId  string            `json:"user_id"`
	Name    string            `json:"name"`
	Type    models.WalletType `json:"type"`
	Balance string            `json:"balance"`
}
