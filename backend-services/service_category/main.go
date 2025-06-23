package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/orders", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Order list"})
	})
	r.Run(":8082")
}
