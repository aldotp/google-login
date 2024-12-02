package handler

import (
	"net/http"

	"github.com/aldotp/golang-login-with-google/internal/service"
	"github.com/aldotp/golang-login-with-google/pkg/util"
	"github.com/gin-gonic/gin"
)

type AuthHandlerInterface interface {
	LoginHandler() gin.HandlerFunc
	RegisterHandler() gin.HandlerFunc
}

type handler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandlerInterface {
	return &handler{
		authService: authService,
	}
}

func (h *handler) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Token string `json:"token"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			response := util.ApiResponse(http.StatusBadRequest, err.Error(), "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		user, err := h.authService.AuthenticateGoogle(payload.Token)
		if err != nil {
			response := util.ApiResponse(http.StatusUnauthorized, err.Error(), "error", nil)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		response := util.ApiResponse(http.StatusOK, "Login Success", "success", user)
		c.JSON(http.StatusOK, response)
	}
}

func (h *handler) RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			Token string `json:"token"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			response := util.ApiResponse(http.StatusBadRequest, err.Error(), "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		user, err := h.authService.RegisterGoogle(payload.Token)
		if err != nil {
			response := util.ApiResponse(http.StatusBadRequest, err.Error(), "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := util.ApiResponse(http.StatusOK, "Register Success", "success", user)
		c.JSON(http.StatusOK, response)
	}
}
