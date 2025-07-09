package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"service-wallet/internal/models"

	_ "github.com/lib/pq"
)

type WalletStore struct {
	db     *sql.DB
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWalletStore(ctx context.Context) (*WalletStore, error) {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database. err : %v", err)
	}

	dbCtx, dbCancel := context.WithCancel(ctx)
	return &WalletStore{
		db:     db,
		ctx:    dbCtx,
		cancel: dbCancel,
	}, nil
}

func (s *WalletStore) CreateNewWallet(wallet models.Wallet) (string, error) {
	var id string
	query := `
	INSERT INTO wallets (user_id, name, type, balance, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id 
	`
	err := s.db.QueryRowContext(s.ctx, query,
		wallet.UserId,
		wallet.Name,
		wallet.Type,
		wallet.Balance,
		wallet.CreatedAt,
		wallet.UpdatedAt,
	).Scan(id)

	return id, err
}

func (s *WalletStore) GetWalletByName(walletName string) (models.Wallet, error) {
	query := `
	SELECT id, user_id, name, type, balance, created_at, updated_at
	FROM wallets
	WWHERE name = $1  
	`
	var existWallet models.Wallet
	err := s.db.QueryRowContext(s.ctx, query, walletName).Scan(
		&existWallet.ID,
		&existWallet.UserId,
		&existWallet.Name,
		&existWallet.Type,
		&existWallet.Balance,
		&existWallet.CreatedAt,
		&existWallet.UpdatedAt,
	)
	if err != nil {
		return existWallet, err
	}

	return existWallet, nil
}

func (s *WalletStore) Close() {
	s.db.Close()
	s.cancel()
}
