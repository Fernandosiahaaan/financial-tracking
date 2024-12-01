package service

import (
	"context"
	"service-user/infrastructure/redis"
	"service-user/internal/store"
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

func (s *UserService) Close() {
	s.cancel()
}
