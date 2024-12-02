package router

import (
	"github.com/aldotp/golang-login-with-google/config"
	"github.com/aldotp/golang-login-with-google/internal/handler"
	"github.com/aldotp/golang-login-with-google/internal/repository"
	"github.com/aldotp/golang-login-with-google/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	db := config.GetConnection()

	repo := repository.NewUserRepository(db)
	authService := service.NewAuthService(repo)

	r := gin.Default()

	authHandler := handler.NewAuthHandler(authService)
	callbackHandler := handler.NewCallbackHandler(authService)
	googleAuthHandler := handler.NewGoogleAuthHandler(authService)

	apiV1 := r.Group("/api/v1")

	authRouter(apiV1.Group("/auth"), authHandler)
	googleAuthRouter(apiV1.Group("/auth/google"), googleAuthHandler)
	callbackRouter(apiV1.Group("/auth/google"), callbackHandler)

	return r
}

func authRouter(r *gin.RouterGroup, h handler.AuthHandlerInterface) {
	r.POST("/login", h.LoginHandler())
	r.POST("/register", h.RegisterHandler())
}

func googleAuthRouter(r *gin.RouterGroup, h handler.GoogleAuthHandlerInterface) {
	r.GET("", h.GoogleAuthHandler())
	r.POST("/token", h.GoogleTokenHandler())
}

func callbackRouter(r *gin.RouterGroup, h handler.CallbackHandlerInterface) {
	r.GET("/callback", h.GoogleCallbackHandler)
}
