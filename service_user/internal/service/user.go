package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"service-user/infrastructure/redis"
	"service-user/internal/model"
	"service-user/internal/store"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	repo   *store.UserStore
	ctx    context.Context
	cancel context.CancelFunc
	redis  *redis.RedisCln
}

func NewUserService(ctx context.Context, redis *redis.RedisCln, repo *store.UserStore) *UserService {
	serviceCtx, serviceCancel := context.WithCancel(ctx)
	return &UserService{
		repo:   repo,
		ctx:    serviceCtx,
		cancel: serviceCancel,
		redis:  redis,
	}
}

func (s *UserService) CreateNewUser(user model.User) (userID string, err error) {
	user.Password = strings.TrimSpace(user.Password)
	hashPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("failed hash password. err = %s", err.Error())
	}

	user.Password = hashPassword
	existUser, err := s.repo.GetUserByName(user.Username)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return "", err
	}

	if existUser.Username == user.Username {
		return "", errors.New("user already created")
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Id = uuid.New().String()

	userID, err = s.repo.CreateNewUser(user)
	if err != nil {
		return "", err
	}

	if err = s.redis.SaveUserInfo(user); err != nil {
		return "", err
	}

	return userID, err
}

func (s *UserService) GetUserByName(user model.User) (model.User, error) {
	existUser, err := s.repo.GetUserByName(user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return existUser, fmt.Errorf("username not found")
		}
		return existUser, fmt.Errorf("error sql. err = %s", err.Error())
	}

	// Verifikasi apakah password cocok dengan hash
	match := s.VerifyPassword(user.Password, existUser.Password)
	if !match {
		return existUser, fmt.Errorf("password not equal")
	}

	return existUser, nil
}

func (s *UserService) GetUserById(userId string) (*model.User, error) {
	exitUser, err := s.redis.GetUserInfo(userId)
	if (err == nil) && (exitUser != nil) {
		return exitUser, nil
	}

	existUser, err := s.repo.GetUserById(userId)
	if err != nil && err != sql.ErrNoRows {
		return existUser, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return existUser, nil
}

func (s *UserService) Close() {
	s.cancel()
}
