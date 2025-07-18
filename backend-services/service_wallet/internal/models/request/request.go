package request

import (
	"service-wallet/internal/models"
)

type CreateWalletRequest struct {
	ReqID   string            `json:"req_id"`
	UserId  string            `json:"user_id"`
	Name    string            `json:"name"`
	Type    models.WalletType `json:"type"`
	Balance string            `json:"balance"`
}

type UpdateWalletRequest struct {
	ReqID    string            `json:"req_id"`
	WalletID string            `json:"wallet_id"`
	Name     string            `json:"name"`
	Type     models.WalletType `json:"type"`
	Balance  string            `json:"balance"`
}

type GetListWalletRequest struct {
	ReqID       string `json:"req_id"`
	Page        string `json:"page"`
	PageItem    string `json:"page_item"`
	FilterBy    string `json:"filter_by"`
	FilterValue string `json:"filter_value"`
}
