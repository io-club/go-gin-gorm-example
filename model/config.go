package model

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	MODE       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBCharset  string
)

var SessionSecret = []byte("my-secret")

func init() {
	// read config from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MODE = os.Getenv("MODE")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
	DBCharset = os.Getenv("DB_CHARSET")
	SessionSecret = []byte(os.Getenv("SESSION_SECRET"))
}
