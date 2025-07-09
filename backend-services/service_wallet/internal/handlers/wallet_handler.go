package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-wallet/internal/models"
	"service-wallet/internal/models/response"
	model "service-wallet/internal/models/response"
	"service-wallet/internal/services"

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
	var wallet models.Wallet
	if err := json.NewDecoder(c.Request.Body).Decode(&wallet); err != nil {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: "failed parse body request"})
		return
	}

	userID, err := h.service.CreateNewWallet(wallet)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusBadRequest, model.ResponseHttp{IsError: true, Message: err.Error()})
		return
	}

	wallet.ID = userID
	model.CreateResponseHttp(c, http.StatusCreated, model.ResponseHttp{IsError: false, Message: fmt.Sprintf("Success created wallet %s", wallet.ID), Data: wallet})
}

func (h *WalletHandler) Close() {
	h.cancel()
}
