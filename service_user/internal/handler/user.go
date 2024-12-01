package handler

import (
	"context"
	"service-user/infrastructure/redis"
	"service-user/internal/service"
	"service-user/internal/store"
	"service-user/middleware"

	"github.com/gin-gonic/gin"
)

type ParamHandler struct {
	Ctx        context.Context
	Middleware *middleware.Midleware
	Service    *service.UserService
	Redis      *redis.RedisCln
	Store      *store.UserStore
}

type UserHandler struct {
	ctx        context.Context
	cancel     context.CancelFunc
	Middleware *middleware.Midleware
	service    *service.UserService
	redis      *redis.RedisCln
	store      *store.UserStore
}

func NewUserHandler(param ParamHandler) *UserHandler {
	handlerCtx, handlerCancel := context.WithCancel(param.Ctx)
	return &UserHandler{
		ctx:        handlerCtx,
		cancel:     handlerCancel,
		service:    param.Service,
		redis:      param.Redis,
		store:      param.Store,
		Middleware: param.Middleware,
	}
}

func (h *UserHandler) UserLogin(c *gin.Context)      {}
func (h *UserHandler) UserCreate(c *gin.Context)     {}
func (h *UserHandler) UserLogout(c *gin.Context)     {}
func (h *UserHandler) UserGetAll(c *gin.Context)     {}
func (h *UserHandler) UserGetByID(c *gin.Context)    {}
func (h *UserHandler) UserUpdateByID(c *gin.Context) {}

func (h *UserHandler) Close() {
	h.cancel()
}
