package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-user/infrastructure/redis"
	"service-user/internal/model"
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

func (h *UserHandler) UserCreate(c *gin.Context) {
	var user model.User
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		model.CreateResponseHttp(c, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed parse body request"})
		return
	}
	if user.Role == "" {
		user.Role = "user"
	}

	userID, err := h.service.CreateNewUser(user)
	if err != nil {
		model.CreateResponseHttp(c, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	user.Id = userID
	model.CreateResponseHttp(c, http.StatusOK, model.ResponseHttp{Error: false, Message: fmt.Sprintf("Success created user %s", user.Id), Data: user})
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var user model.User
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		model.CreateResponseHttp(c, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed parse body request"})
		return
	}

	user, err := h.service.GetUserByName(user)
	if err != nil {
		model.CreateResponseHttp(c, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	tokenString, err := h.Middleware.CreateToken(user.Username, user.Password)
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed created token"})
		return
	}

	err = h.redis.SaveUserInfo(user)
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	dataResponse := model.LoginData{Token: tokenString, Id: user.Id}
	model.CreateResponseHttp(c, http.StatusOK, model.ResponseHttp{Error: false, Message: "Success login", Data: dataResponse})
}

func (h *UserHandler) UserGetByID(c *gin.Context) {
	// token, flag := c.Get("jwtToken")
}

func (h *UserHandler) UserLogout(c *gin.Context)     {}
func (h *UserHandler) UserGetAll(c *gin.Context)     {}
func (h *UserHandler) UserUpdateByID(c *gin.Context) {}

func (h *UserHandler) Close() {
	h.cancel()
}
