package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"service-user/internal/model"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Midleware struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMidleware(ctx context.Context) *Midleware {
	midlewareCtx, midlewareCancel := context.WithCancel(ctx)
	var middleware *Midleware = &Midleware{
		ctx:    midlewareCtx,
		cancel: midlewareCancel,
	}
	return middleware
}

func (m *Midleware) CreateToken(username string, password string) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"password": password,
			"exp":      time.Now().Add(model.UserSessionTime).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func (m *Midleware) VerifyToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func (m *Midleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Header("Content-Type", "application/json")
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			model.CreateResponseHttp(c, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Authentication header null"})
			return
		}

		bearerToken := strings.Split(authToken, " ")
		if len(bearerToken) != 2 {
			model.CreateResponseHttp(c, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Invalid format token"})
			return
		}

		var jwtToken string = bearerToken[1]
		token, err := m.VerifyToken(jwtToken)
		if err != nil {
			model.CreateResponseHttp(c, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Failed token. err = %s", err)})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			model.CreateResponseHttp(c, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Fail claims token"})
			return
		}

		c.Set("jwtToken", jwtToken)
		c.Set("user", claims)
		c.Next()
	}

}

func (m *Midleware) Close() {
	m.cancel()
}
