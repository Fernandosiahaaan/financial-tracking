package main

import (
	"fmt"
	"log"
	"net/http"
	"service-wallet/internal/handlers"
	"service-wallet/internal/models/response"

	"github.com/gin-gonic/gin"
)

func routing(handler *handlers.WalletHandler) (*gin.Engine, error) {
	rout := gin.New()
	rout.Use(CustomRecovery())

	// with no middleware
	rout.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	rout.POST("/wallet", handler.WalletCreate)
	rout.GET("/wallet/:id", handler.GetWalletById)
	rout.PUT("/wallet/:id", handler.WalletUpdate)
	return rout, nil
}

func CustomRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ResponseHttp{IsError: true, Message: "internal server error", MessageErr: fmt.Sprintf("Error Panic recovered: %v", recovered)})
	})
}
