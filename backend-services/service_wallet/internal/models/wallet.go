package models

import "time"

type Wallet struct {
	ID        string     `json:"id"`
	UserId    string     `json:"user_id"`
	Name      string     `json:"name"`
	Type      WalletType `json:"type"`
	Balance   string     `json:"balance"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type WalletType string

const (
	WalletCash    WalletType = "CASH"
	WalletBank    WalletType = "BANK"
	WalletEwallet WalletType = "EWALLET"
)
