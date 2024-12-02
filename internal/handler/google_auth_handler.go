package handler

import (
	"net/http"

	"github.com/aldotp/golang-login-with-google/internal/service"
	"github.com/aldotp/golang-login-with-google/pkg/util"
	"github.com/gin-gonic/gin"
)

type GoogleAuthHandlerInterface interface {
	GoogleAuthHandler() gin.HandlerFunc
	GoogleTokenHandler() gin.HandlerFunc
}

type googleAuthHandler struct {
	AuthService service.AuthService
}

func NewGoogleAuthHandler(authService service.AuthService) GoogleAuthHandlerInterface {
	return &googleAuthHandler{
		AuthService: authService,
	}
}
func (g *googleAuthHandler) GoogleAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := g.AuthService.GenerateGoogleAuthURL()
		response := util.ApiResponse(http.StatusOK, "Get Google Auth URL Success", "success", url)
		c.JSON(http.StatusOK, response)
	}
}

func (g *googleAuthHandler) GoogleTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Code string `json:"code"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			response := util.ApiResponse(http.StatusBadRequest, err.Error(), "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		token, err := g.AuthService.ExchangeGoogleCodeForToken(payload.Code)
		if err != nil {
			response := util.ApiResponse(http.StatusBadRequest, err.Error(), "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		// Create a struct to hold the token
		tokenResponse := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}

		response := util.ApiResponse(http.StatusOK, "Get Google Token Success", "success", tokenResponse)
		c.JSON(http.StatusOK, response)
	}
}
