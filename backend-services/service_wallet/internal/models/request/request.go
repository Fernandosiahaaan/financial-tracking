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
	ReqId       string `json:"reqId"`
	Page        string `json:"page"`
	PageItem    string `json:"pageItem"`
	FilterBy    string `json:"filterBy"`
	FilterValue string `json:"filterValue"`
}
