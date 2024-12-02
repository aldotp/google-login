package handler

import (
	"fmt"
	"net/http"

	"github.com/aldotp/golang-login-with-google/internal/service"
	"github.com/gin-gonic/gin"
)

type CallbackHandlerInterface interface {
	GoogleCallbackHandler(c *gin.Context)
}

type callbackHandler struct {
}

func NewCallbackHandler(authService service.AuthService) CallbackHandlerInterface {
	return &callbackHandler{}
}

func (h *callbackHandler) GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code received"})
		return
	}

	fmt.Println("Authorization Code:", code)
	c.JSON(http.StatusOK, gin.H{"code": code})
}
