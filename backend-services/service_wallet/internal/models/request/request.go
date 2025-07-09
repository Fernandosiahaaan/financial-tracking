package request

import (
	"service-wallet/internal/models"
	"time"
)

type Wallet struct {
	ReqID     string            `json:"req_id"`
	ID        string            `json:"id"`
	UserId    string            `json:"user_id"`
	Name      string            `json:"name"`
	Type      models.WalletType `json:"type"`
	Balance   string            `json:"balance"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
