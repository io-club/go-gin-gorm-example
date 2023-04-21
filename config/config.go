package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	MODE       string
	DB_MODE    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBCharset  string
)

var (
    MaxImageNum int64
)

var SessionSecret = []byte("my-secret")

func init() {
	// read config from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MODE = os.Getenv("MODE")
	DB_MODE = os.Getenv("DB_MODE")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
	DBCharset = os.Getenv("DB_CHARSET")
	SessionSecret = []byte(os.Getenv("SESSION_SECRET"))
    MaxImageNumStr := os.Getenv("MAX_IMAGE_NUM")
    MaxImageNumInt, err := strconv.ParseInt(MaxImageNumStr, 10, 64)
    if err != nil {
        MaxImageNum = 20
    } else {
        MaxImageNum = MaxImageNumInt
    }

}
