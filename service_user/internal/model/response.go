package model

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type ResponseHttp struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateResponseHttp(c *gin.Context, statusCode int, response ResponseHttp) {
	c.JSON(statusCode, response)
	if response.Error {
		fmt.Printf("❌  [%s] uri = '%s'; status code = %d; message = %s\n", c.Request.Method, c.Request.URL, statusCode, response.Message)
		return
	}
	fmt.Printf("✅  [%s] uri = '%s'; status code = %d; message = %s\n", c.Request.Method, c.Request.URL, statusCode, response.Message)
}
