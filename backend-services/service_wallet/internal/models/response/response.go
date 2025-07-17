package response

import (
	"fmt"
	"service-wallet/internal/models"

	"github.com/gin-gonic/gin"
)

type ResponseHttp struct {
	IsError    bool        `json:"is_error"`
	Message    string      `json:"message"`
	MessageErr string      `json:"message_error"`
	Data       interface{} `json:"data"`
}

type GetListWalletResponse struct {
	IsError      bool            `json:"is_error"`
	Message      string          `json:"message"`
	MessageErr   string          `json:"message_error"`
	Start        int             `json:"start"`
	End          int             `json:"end"`
	Page         int             `json:"page"`
	Pages        int             `json:"pages"`
	RecordsTotal int             `json:"records_total"`
	Data         []models.Wallet `json:"data"`
}

func CreateResponseHttp(c *gin.Context, statusCode int, response ResponseHttp) {
	c.JSON(statusCode, response)
	if response.IsError {
		fmt.Printf("❌  [%s] uri = '%s'; status code = %d; message = %s\n", c.Request.Method, c.Request.URL, statusCode, response.Message)
		return
	}
	fmt.Printf("✅  [%s] uri = '%s'; status code = %d; message = %s\n", c.Request.Method, c.Request.URL, statusCode, response.Message)
}
