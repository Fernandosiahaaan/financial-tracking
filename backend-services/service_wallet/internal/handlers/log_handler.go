package handlers

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *WalletHandler) ReadRequestBody(c *gin.Context, maxSize int64) (string, error) {
	if c.Request.Body == nil {
		return "", nil
	}

	bodyBytes, err := io.ReadAll(
		io.LimitReader(c.Request.Body, maxSize),
	)
	if err != nil {
		return "", err
	}

	// restore body supaya handler tetap bisa baca
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes), nil
}

func (h *WalletHandler) MaskSensitive(body string) string {
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return body // bukan JSON → return apa adanya
	}

	sensitiveKeys := []string{
		"password",
		"pin",
		"token",
		"authorization",
	}

	for _, key := range sensitiveKeys {
		if _, ok := data[key]; ok {
			data[key] = "***"
		}
	}

	masked, _ := json.Marshal(data)
	return string(masked)
}

func (h *WalletHandler) LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := h.ReadRequestBody(c, 10*1024) // 10KB max
		body = h.MaskSensitive(body)
		if h.logStash != nil {
			h.logStash.Info("http_request",
				zap.String("type", "request"),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("body", body),
				zap.String("request_id", c.GetString("request_id")),
			)
		}

		c.Next()
	}
}
