package postgree

import (
	"context"
	"database/sql"
	"fmt"
	"service-wallet/internal/models"
	"service-wallet/internal/models/request"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var defaultValueTimeout time.Duration = 5 * time.Second

type StoreWalletPq struct {
	Db     *sql.DB
	Ctx    context.Context
	Cancel context.CancelFunc
}

func (s *StoreWalletPq) CreateNewWallet(wallet models.Wallet) (string, error) {
	var id string
	query := `
	INSERT INTO wallets (user_id, name, type, balance, description, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id 
	`
	err := s.Db.QueryRowContext(s.Ctx, query,
		wallet.UserId,
		wallet.Name,
		wallet.Type,
		wallet.Balance,
		wallet.Description,
		wallet.CreatedAt,
		wallet.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed create wallet to db. err : %v", err)
	}

	return id, nil

}

func (s *StoreWalletPq) GetWalletByName(walletName string) (*models.Wallet, error) {
	ctx, cancel := context.WithTimeout(s.Ctx, defaultValueTimeout)
	defer cancel()

	query := `
	SELECT id, user_id, name, type, balance, description, created_at, updated_at
	FROM wallets
	WHERE name = $1  
	`
	var existWallet *models.Wallet = &models.Wallet{}
	err := s.Db.QueryRowContext(ctx, query, walletName).Scan(
		&existWallet.ID,
		&existWallet.UserId,
		&existWallet.Name,
		&existWallet.Type,
		&existWallet.Balance,
		&existWallet.Description,
		&existWallet.CreatedAt,
		&existWallet.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed get wallet with name '%s' to db. err : %v", walletName, err)
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return existWallet, nil
}

func (s *StoreWalletPq) GetWalletById(walletId string) (*models.Wallet, error) {
	ctx, cancel := context.WithTimeout(s.Ctx, defaultValueTimeout)
	defer cancel()

	query := `
	SELECT id, user_id, name, type, balance, description, created_at, updated_at
	FROM wallets
	WHERE id = $1  
	`
	var existWallet *models.Wallet = &models.Wallet{}
	err := s.Db.QueryRowContext(ctx, query, walletId).Scan(
		&existWallet.ID,
		&existWallet.UserId,
		&existWallet.Name,
		&existWallet.Type,
		&existWallet.Balance,
		&existWallet.Description,
		&existWallet.CreatedAt,
		&existWallet.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed get wallet with id '%s' to db. err : %v", walletId, err)
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return existWallet, nil
}

func (r *StoreWalletPq) UpdateWalletById(wallet models.Wallet) error {
	ctx, cancel := context.WithTimeout(r.Ctx, defaultValueTimeout)
	defer cancel()

	query := `
        UPDATE wallets 
        SET name = $1, type = $2, balance = $3, description = $4, updated_at = $5 
        WHERE id = $6
    `
	result, err := r.Db.ExecContext(ctx, query,
		wallet.Name,
		wallet.Type,
		wallet.Balance,
		wallet.Description,
		wallet.UpdatedAt,
		wallet.ID,
	)
	if err != nil {
		return fmt.Errorf("failed exec update id. err : %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed check affected update id. err : %v", err)
	} else if affected <= 0 {
		return fmt.Errorf("not affected update in table wallets")
	}

	return nil
}

func (r *StoreWalletPq) DeleteWalletById(id string) error {
	ctx, cancel := context.WithTimeout(r.Ctx, defaultValueTimeout)
	defer cancel()

	query := `
        DELETE FROM wallets 
        WHERE id = $1
    `
	result, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed exec delete id. err : %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed check affected delete id. err : %v", err)
	} else if affected <= 0 {
		return fmt.Errorf("not affected delete in table wallets")
	}

	return nil
}

func (r *StoreWalletPq) GetListWallets(params request.GetListWalletRequest) ([]models.Wallet, int, error) {
	ctx, cancel := context.WithTimeout(r.Ctx, 5*time.Second)
	defer cancel()

	// Default pagination
	pageNumber, _ := strconv.Atoi(params.Page)
	if pageNumber <= 0 {
		pageNumber = 1
	}
	pageSize, _ := strconv.Atoi(params.PageItem)
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNumber - 1) * pageSize

	// Filter query
	filterQuery := `WHERE 1=1`
	if params.FilterBy != "" && params.FilterValue != "" {
		filterQuery += fmt.Sprintf(" AND %s ILIKE '%%%s%%'", params.FilterBy, params.FilterValue)
	}

	// Hitung total data
	var total int
	countQuery := `SELECT COUNT(*) FROM wallets ` + filterQuery
	if err := r.Db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count wallets: %v", err)
	}

	// Ambil data paginated
	query := fmt.Sprintf(`
        SELECT id, name, type, balance, description, updated_at
        FROM wallets
        %s
        ORDER BY updated_at DESC
        LIMIT $1 OFFSET $2
    `, filterQuery)

	rows, err := r.Db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query wallets: %v", err)
	}
	defer rows.Close()

	var wallets []models.Wallet
	for rows.Next() {
		var w models.Wallet
		if err := rows.Scan(&w.ID, &w.Name, &w.Type, &w.Balance, &w.Description, &w.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan wallet: %v", err)
		}
		wallets = append(wallets, w)
	}

	return wallets, total, nil
}

func (s *StoreWalletPq) Close() {
	s.Db.Close()
	s.Cancel()
}
