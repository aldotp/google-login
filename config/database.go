package config

import (
	"database/sql"
	"log"
	"sync"
)

var (
	once sync.Once
	db   *sql.DB
)

func NewPostgresDatabase() {
	once.Do(func() {
		var err error
		host := GetEnv("DB_HOST")
		port := GetEnv("DB_PORT")
		user := GetEnv("DB_USER")
		password := GetEnv("DB_PASSWORD")
		dbname := GetEnv("DB_NAME")
		dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal(err)
		}

		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}

		log.Println("Connected to database")

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)

	})
}

func GetConnection() *sql.DB {
	return db
}
