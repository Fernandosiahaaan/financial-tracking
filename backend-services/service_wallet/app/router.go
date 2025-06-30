package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func routing() (*gin.Engine, error) {
	rout := gin.Default()

	// with no middleware
	rout.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return rout, nil
}
