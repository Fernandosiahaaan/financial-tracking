package main

import (
	"net/http"
	"service-user/internal/handler"

	"github.com/gin-gonic/gin"
)

func routing(handler *handler.UserHandler) (*gin.Engine, error) {
	rout := gin.Default()

	// with no middleware
	rout.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	rout.POST("/login", handler.UserLogin)
	rout.POST("/user", handler.UserCreate)

	// with midleware
	auth := rout.Group("/", handler.Middleware.AuthMiddleware())
	auth.POST("user/logout", handler.UserLogout)
	auth.GET("user", handler.UserGetAll)
	auth.GET("user/:id", handler.UserGetByID)
	auth.POST("user/:id", handler.UserUpdateByID)

	return rout, nil
}
