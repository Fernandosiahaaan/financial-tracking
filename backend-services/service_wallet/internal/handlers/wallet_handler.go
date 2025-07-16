package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"service-wallet/internal/models"
	"service-wallet/internal/models/request"
	"service-wallet/internal/models/response"
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
	respOut, err := h.service.CreateNewWallet(&wallet)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, *respOut)
		return
	}

	response.CreateResponseHttp(c, http.StatusCreated, *respOut)
}

func (h *WalletHandler) GetWalletById(c *gin.Context) {
	walletID := c.Param("id")
	if walletID == "" {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: "Invalid User ID uri"})
		return
	}

	respOut, err := h.service.GetWalletById(walletID)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, *respOut)
		return
	}

	response.CreateResponseHttp(c, http.StatusOK, *respOut)
}

func (h *WalletHandler) WalletUpdate(c *gin.Context) {
	var req request.UpdateWallet

	walletID := c.Param("id")
	if walletID == "" {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: "Invalid User ID uri"})
		return
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: "failed parse body request", MessageErr: fmt.Sprintf("failed parse body request. err : %v", err)})
		return
	}

	req.WalletID = walletID
	errRedaksi, errSystem := validation.ValidationUpdateWallet(req)
	if errRedaksi != nil {
		response.CreateResponseHttp(c, http.StatusBadRequest, response.ResponseHttp{IsError: true, Message: errRedaksi.Error(), MessageErr: errSystem.Error()})
		return
	}

	balance, _ := strconv.Atoi(req.Balance)
	var wallet models.Wallet = models.Wallet{
		Name:    req.Name,
		Type:    req.Type,
		Balance: int64(balance),
		ID:      req.WalletID,
	}
	respOut, err := h.service.UpdateWallet(wallet)
	if err != nil {
		response.CreateResponseHttp(c, http.StatusInternalServerError, *respOut)
		return
	}

	response.CreateResponseHttp(c, http.StatusCreated, *respOut)
}

func (h *WalletHandler) Close() {
	h.cancel()
}
