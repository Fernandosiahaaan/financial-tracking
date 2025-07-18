package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-user/infrastructure/redis"
	"service-user/internal/model"
	"service-user/internal/model/response"
	"service-user/internal/service"
	"service-user/internal/store"
	"service-user/middlewares"
	"service-user/utils"

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
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: "failed parse body request", MessageErr: fmt.Sprintf("failed parse body request. err : %v", err)})
		return
	}
	if user.Role == "" {
		user.Role = model.RoleUser
	}

	var statusCode int = 201
	bodyResp, err := h.service.CreateNewUser(user)
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	response.CreateResponseHttp(c, statusCode, bodyResp)
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var user model.User
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: "failed parse body request"})
		return
	}

	user, err := h.service.GetUserByName(user)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: err.Error()})
		return
	}

	tokenString, err := h.Middleware.CreateToken(user.Username, user.Password)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "failed created token"})
		return
	}

	h.redis.SaveUserInfo(user)
	err = h.redis.SetLoginInfo(h.ctx, tokenString, model.LoginCacheData{Id: user.ID, Username: user.Username, Role: user.Role})
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: err.Error()})
		return
	}

	dataResponse := model.LoginData{Token: tokenString, Id: user.ID}
	response.CreateResponseHttp(c, http.StatusOK, response.ResponseHttp{IsError: false, Message: "Success login", Data: dataResponse})
}

func (h *UserHandler) UserLogout(c *gin.Context) {
	tokenStr, flag := c.Get("jwtToken")
	if !flag {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "failed read token value"})
		return
	}

	_, err := h.redis.GetLoginInfo(tokenStr.(string))
	if err != nil {
		var msgErr error = utils.MessageError("redis::GetLoginInfo", err)
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "failed get session user", MessageErr: msgErr.Error()})
		return
	}

	err = h.redis.DeleteLoginInfo(tokenStr.(string))
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "Failed logout session"})
		return
	}

	response.CreateResponseHttp(c, http.StatusOK, response.ResponseHttp{IsError: false, Message: "Success logout session"})

}

func (h *UserHandler) UserGetByID(c *gin.Context) {
	tokenStr, flag := c.Get("jwtToken")
	if !flag {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "failed read token value"})
		return
	}

	loginInfo, err := h.redis.GetLoginInfo(tokenStr.(string))
	if err != nil {
		var msgErr error = utils.MessageError("redis::GetLoginInfo", err)
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "failed get session user", MessageErr: msgErr.Error()})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: "Invalid User ID uri"})
		return
	}

	user, err := h.service.GetUserById(userID)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: fmt.Sprintf("Invalid username and password. err := %v", err)})
		return
	} else if user == nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: fmt.Sprintf("username not found with id '%s'", userID)})
		return
	} else if (loginInfo.Role == model.RoleUser) && (userID != loginInfo.Id) {
		response.CreateResponseHttp(c, http.StatusForbidden, response.ResponseHttp{IsError: true, Message: fmt.Sprintf("user %s doesn't have access get user %s info", loginInfo.Username, user.Username)})
		return
	}

	response.CreateResponseHttp(c, http.StatusOK, response.ResponseHttp{IsError: false, Message: "success get info me", Data: user})
}

func (h *UserHandler) UserGetAll(c *gin.Context) {
	tokenStr, flag := c.Get("jwtToken")
	if !flag {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "failed read token value"})
		return
	}

	loginInfo, err := h.redis.GetLoginInfo(tokenStr.(string))
	if err != nil {
		var msgErr error = utils.MessageError("redis::GetLoginInfo", err)
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "failed get session user", MessageErr: msgErr.Error()})
		return
	}

	if loginInfo.Role == model.RoleUser {
		response.CreateResponseHttp(c, http.StatusForbidden, response.ResponseHttp{IsError: true, Message: fmt.Sprintf("user %s doesn't have access to get all users", loginInfo.Username)})
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: fmt.Sprintf("Invalid get users. err = %s", err.Error())})
		return
	}

	response.CreateResponseHttp(c, http.StatusOK, response.ResponseHttp{IsError: true, Message: "success get all users", Data: users})
}

func (h *UserHandler) UserUpdateByID(c *gin.Context) {
	var userUpdateReq model.User

	tokenStr, flag := c.Get("jwtToken")
	if !flag {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "failed read token value"})
		return
	}

	loginInfo, err := h.redis.GetLoginInfo(tokenStr.(string))
	if err != nil {
		var msgErr error = utils.MessageError("redis::GetLoginInfo", err)
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "failed get session user", MessageErr: msgErr.Error()})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: "Invalid User ID uri"})
		return
	}

	user, err := h.service.GetUserById(userID)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: fmt.Sprintf("Invalid username and password. err := %v", err)})
		return
	} else if user == nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "username not found"})
		return
	} else if (loginInfo.Role != "superadmin") && (userID != loginInfo.Id) {
		response.CreateResponseHttp(c, http.StatusForbidden, response.ResponseHttp{IsError: true, Message: fmt.Sprintf("user %s with role %s doesn't have access update user %s info", loginInfo.Username, loginInfo.Role, user.Username)})
		return
	}

	userUpdateReq.ID = user.ID
	if err = json.NewDecoder(c.Request.Body).Decode(&userUpdateReq); err != nil {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed to decode body request, err = %s", err.Error())})
		return
	}

	var statusCode int = 200
	bodyResp, err := h.service.UpdateUser(userUpdateReq)
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	response.CreateResponseHttp(c, statusCode, bodyResp)
}

func (h *UserHandler) Close() {
	h.cancel()
}
