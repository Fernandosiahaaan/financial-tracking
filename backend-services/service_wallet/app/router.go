package main

import (
	"net/http"
	"service-wallet/internal/handlers"

	"github.com/gin-gonic/gin"
)

func routing(handler *handlers.WalletHandler) (*gin.Engine, error) {
	rout := gin.Default()

	// with no middleware
	rout.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	rout.POST("/wallet", handler.WalletCreate)
	return rout, nil
}
