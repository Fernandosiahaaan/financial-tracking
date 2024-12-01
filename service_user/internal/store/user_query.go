package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"service-user/internal/model"

	"github.com/lib/pq"
)

type UserStore struct {
	db     *sql.DB
	ctx    context.Context
	cancel context.CancelFunc
}

func NewUserStore(ctx context.Context) (*UserStore, error) {
	sql.Register("postgres", &pq.Driver{})

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database. err : %v", err)
	}

	dbCtx, dbCancel := context.WithCancel(ctx)
	return &UserStore{
		db:     db,
		ctx:    dbCtx,
		cancel: dbCancel,
	}, nil
}

func (r *UserStore) CreateNewUser(user model.User) (string, error) {
	var id string
	query := `
	INSERT INTO users (id, username, password, email, role, created_at, updated_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)
	RETURNING id 
	`
	err := r.db.QueryRowContext(r.ctx, query, user.Id, user.Username, user.Password, user.Email, user.Role, user.CreatedAt, user.UpdatedAt, user.DeletedAt).Scan(&id)
	return id, err
}

func (r *UserStore) UpdateUser(user model.User) (string, error) {
	var id string
	query := `
        UPDATE users 
        SET username = $1, password = $2, email = $3, role = $4, updated_at = $5
        WHERE id = $6 and deleted_at = NULL
        RETURNING id
    `
	err := r.db.QueryRowContext(r.ctx, query, user.Username, user.Password, user.Email, user.Role, user.UpdatedAt, user.Id).Scan(&id)

	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *UserStore) GetUser(user model.User) (model.User, error) {
	query := `
	SELECT id, username, password, email, role, created_at, updated_at 
	FROM users 
	WHERE username=$1 and deleted_at=NULL
	`
	var existUser model.User
	err := r.db.QueryRowContext(r.ctx, query, user.Username).Scan(
		&existUser.Id,
		&existUser.Username,
		&existUser.Password,
		&existUser.Email,
		&existUser.Role,
		&existUser.CreatedAt,
		&existUser.UpdatedAt,
	)
	if err != nil {
		return existUser, err
	}
	return existUser, nil
}

func (r *UserStore) GetUserById(userId string) (*model.User, error) {
	query := `
	SELECT id, username, password, email, role, created_at, updated_at
	FROM users 
	WHERE id=$1 and deleted_at=NULL
	`
	var existUser *model.User = &model.User{}
	err := r.db.QueryRowContext(r.ctx, query, userId).Scan(
		&existUser.Id,
		&existUser.Username,
		&existUser.Password,
		&existUser.Email,
		&existUser.Role,
		&existUser.CreatedAt,
		&existUser.UpdatedAt,
	)
	if err != nil {
		return existUser, err
	}
	return existUser, nil
}

func (r *UserStore) GetAllUsers() ([]model.User, error) {
	query := `
	SELECT id, username, password, email, role, created_at, updated_at
	FROM users
	WHERE deleted_at=NULL
	`
	rows, err := r.db.QueryContext(r.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserStore) Close() {
	r.db.Close()
	r.cancel()
}
