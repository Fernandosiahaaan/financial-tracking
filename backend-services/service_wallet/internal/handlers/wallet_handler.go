package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-wallet/internal/models"
	"service-wallet/internal/models/request"
	"service-wallet/internal/models/response"
	model "service-wallet/internal/models/response"
	"service-wallet/internal/services"
	validation "service-wallet/internal/validations"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	ctx     context.Context
	cancel  context.CancelFunc
	service *services.WalletService
}

func NewUserHandler(ctx context.Context, service services.WalletService) *WalletHandler {
	handlerCtx, handlerCancel := context.WithCancel(ctx)
	return &WalletHandler{
		ctx:     handlerCtx,
		cancel:  handlerCancel,
		service: &service,
	}
}

func (h *WalletHandler) WalletCreate(c *gin.Context) {
	var req request.CreateWallet
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: "failed parse body request", MessageErr: fmt.Sprintf("failed parse body request. err : %v", err)})
		return
	}

	msg, err := validation.ValidationCreateWallet(req)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: msg.Error(), MessageErr: err.Error()})
		return
	}

	balance, _ := strconv.Atoi(req.Balance)
	var wallet models.Wallet = models.Wallet{
		Name:    req.Name,
		Type:    req.Type,
		Balance: int64(balance),
		UserId:  req.UserId,
	}
	walletID, err := h.service.CreateNewWallet(&wallet)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, model.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed created wallet '%s'", wallet.Name), MessageErr: err.Error()})
		return
	}

	wallet.ID = walletID
	model.CreateResponseHttp(c, http.StatusCreated, model.ResponseHttp{IsError: false, Message: fmt.Sprintf("Success created wallet '%s'", wallet.ID), Data: wallet})
}

func (h *WalletHandler) Close() {
	h.cancel()
}
