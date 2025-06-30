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
	"service-user/middlewares"

	"github.com/gin-gonic/gin"
)

type ParamHandler struct {
	Ctx        context.Context
	Middleware *middlewares.Midleware
	Service    *service.UserService
	Redis      *redis.RedisCln
	Store      *store.UserStore
}

type UserHandler struct {
	ctx        context.Context
	cancel     context.CancelFunc
	Middleware *middlewares.Midleware
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
		user.Role = model.RoleUser
	}

	userID, err := h.service.CreateNewUser(user)
	if err != nil {
		model.CreateResponseHttp(c, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	user.ID = userID
	model.CreateResponseHttp(c, http.StatusCreated, model.ResponseHttp{Error: false, Message: fmt.Sprintf("Success created user %s", user.ID), Data: user})
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

	h.redis.SaveUserInfo(user)
	err = h.redis.SetLoginInfo(h.ctx, tokenString, model.LoginCacheData{Id: user.ID, Username: user.Username, Role: user.Role})
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	dataResponse := model.LoginData{Token: tokenString, Id: user.ID}
	model.CreateResponseHttp(c, http.StatusOK, model.ResponseHttp{Error: false, Message: "Success login", Data: dataResponse})
}

func (h *UserHandler) UserLogout(c *gin.Context) {
	tokenStr, flag := c.Get("jwtToken")
	if !flag {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed read token value"})
		return
	}

	_, err := h.redis.GetLoginInfo(tokenStr.(string))
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	err = h.redis.DeleteLoginInfo(tokenStr.(string))
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "Failed logout session"})
		return
	}

	model.CreateResponseHttp(c, http.StatusOK, model.ResponseHttp{Error: false, Message: "Success logout session"})

}

func (h *UserHandler) UserGetByID(c *gin.Context) {
	tokenStr, flag := c.Get("jwtToken")
	if !flag {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed read token value"})
		return
	}

	loginInfo, err := h.redis.GetLoginInfo(tokenStr.(string))
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		model.CreateResponseHttp(c, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "Invalid User ID uri"})
		return
	}

	user, err := h.service.GetUserById(userID)
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Invalid username and password. err := %v", err)})
		return
	} else if user == nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: fmt.Sprintf("username not found with id '%s'", userID)})
		return
	} else if (loginInfo.Role == model.RoleUser) && (userID != loginInfo.Id) {
		model.CreateResponseHttp(c, http.StatusForbidden, model.ResponseHttp{Error: true, Message: fmt.Sprintf("user %s doesn't have access get user %s info", loginInfo.Username, user.Username)})
		return
	}

	model.CreateResponseHttp(c, http.StatusOK, model.ResponseHttp{Error: false, Message: "success get info me", Data: user})
}

func (h *UserHandler) UserGetAll(c *gin.Context) {
	tokenStr, flag := c.Get("jwtToken")
	if !flag {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed read token value"})
		return
	}

	loginInfo, err := h.redis.GetLoginInfo(tokenStr.(string))
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	if loginInfo.Role == model.RoleUser {
		model.CreateResponseHttp(c, http.StatusForbidden, model.ResponseHttp{Error: true, Message: fmt.Sprintf("user %s doesn't have access to get all users", loginInfo.Username)})
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Invalid get users. err = %s", err.Error())})
		return
	}

	model.CreateResponseHttp(c, http.StatusOK, model.ResponseHttp{Error: true, Message: "success get all users", Data: users})
}

func (h *UserHandler) UserUpdateByID(c *gin.Context) {
	var userUpdateReq model.User

	tokenStr, flag := c.Get("jwtToken")
	if !flag {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed read token value"})
		return
	}

	loginInfo, err := h.redis.GetLoginInfo(tokenStr.(string))
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		model.CreateResponseHttp(c, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "Invalid User ID uri"})
		return
	}

	user, err := h.service.GetUserById(userID)
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Invalid username and password. err := %v", err)})
		return
	} else if user == nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "username not found"})
		return
	} else if (loginInfo.Role != "superadmin") && (userID != loginInfo.Id) {
		model.CreateResponseHttp(c, http.StatusForbidden, model.ResponseHttp{Error: true, Message: fmt.Sprintf("user %s with role %s doesn't have access update user %s info", loginInfo.Username, loginInfo.Role, user.Username)})
		return
	}

	userUpdateReq.ID = user.ID
	if err = json.NewDecoder(c.Request.Body).Decode(&userUpdateReq); err != nil {
		model.CreateResponseHttp(c, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("failed to decode body request, err = %s", err.Error())})
		return
	}
	fmt.Println("user update = ", userUpdateReq)

	userUpdateReq, err = h.service.UpdateUser(userUpdateReq)
	if err != nil {
		model.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	model.CreateResponseHttp(c, http.StatusOK, model.ResponseHttp{Error: false, Message: fmt.Sprintf("success update user %s", user.Username), Data: userUpdateReq})
}

func (h *UserHandler) Close() {
	h.cancel()
}
