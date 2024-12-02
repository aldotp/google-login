package main

import (
	"log"

	"github.com/aldotp/golang-login-with-google/config"
	"github.com/aldotp/golang-login-with-google/internal/router"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func main() {
	config.LoadEnv()
	config.NewPostgresDatabase()
	r := router.SetupRouter()
	log.Println("Starting server on :8080")
	r.Run(":8080")
}
