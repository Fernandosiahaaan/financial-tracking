package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"service-user/internal/model"

	_ "github.com/lib/pq"
)

type UserStore struct {
	db     *sql.DB
	ctx    context.Context
	cancel context.CancelFunc
}

func NewUserStore(ctx context.Context) (*UserStore, error) {
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
	INSERT INTO users (id, username, password, full_name, email, phone_number, role, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id 
	`
	err := r.db.QueryRowContext(r.ctx, query, user.ID, user.Username, user.Password, user.FullName, user.Email, user.PhoneNumber, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("err query create user with username '%s' : %v", user.Username, err)
	}
	return id, nil
}

func (r *UserStore) UpdateUser(user model.User) (string, error) {
	var id string
	query := `
        UPDATE users 
        SET username = $1, password = $2, full_name = $3, email = $4, phone_number = $5, role = $6, updated_at = $7
        WHERE id = $8
        RETURNING id
    `
	err := r.db.QueryRowContext(r.ctx, query,
		user.Username,
		user.Password,
		user.FullName,
		user.Email,
		user.PhoneNumber,
		user.Role,
		user.UpdatedAt,
		user.ID,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("err query update user with id '%s' : %v", user.ID, err)
	}
	return id, nil
}

func (r *UserStore) GetUserByName(username string) (*model.User, error) {
	query := `
	SELECT id, username, password, full_name, email, phone_number, role, created_at, updated_at
	FROM users 
	WHERE username=$1
	`
	var existUser model.User
	err := r.db.QueryRowContext(r.ctx, query, username).Scan(
		&existUser.ID,
		&existUser.Username,
		&existUser.Password,
		&existUser.FullName,
		&existUser.Email,
		&existUser.PhoneNumber,
		&existUser.Role,
		&existUser.CreatedAt,
		&existUser.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return &existUser, fmt.Errorf("err query get user with username '%s' : %v", username, err)
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &existUser, nil
}

func (r *UserStore) GetUserById(userId string) (*model.User, error) {
	query := `
	SELECT id, username, password, full_name, email, phone_number, role, created_at, updated_at
	FROM users 
	WHERE id=$1
	`
	var existUser *model.User = &model.User{}
	err := r.db.QueryRowContext(r.ctx, query, userId).Scan(
		&existUser.ID,
		&existUser.Username,
		&existUser.Password,
		&existUser.FullName,
		&existUser.Email,
		&existUser.PhoneNumber,
		&existUser.Role,
		&existUser.CreatedAt,
		&existUser.UpdatedAt,
	)
	if err != nil {
		return existUser, fmt.Errorf("err query get user with id '%s' : %v", existUser.ID, err)
	}
	return existUser, nil
}

func (r *UserStore) GetAllUsers() ([]model.User, error) {
	query := `
	SELECT id, username, full_name, email, phone_number, role, created_at, updated_at
	FROM users
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
			&user.ID,
			&user.Username,
			&user.FullName,
			&user.Email,
			&user.PhoneNumber,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("err scan user : %v", user)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("err get all user rows : %v", err)
	}

	return users, nil
}

func (r *UserStore) Close() {
	r.db.Close()
	r.cancel()
}
