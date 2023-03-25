package model

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	DB = db
}

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", DBUser, DBPassword, DBHost, DBPort, DBName, DBCharset)
	// set timeout
	dsn += "&timeout=10s&readTimeout=30s&writeTimeout=30s&parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create tables
	db.AutoMigrate(&Fabric{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Image{})
	db.AutoMigrate(&Brand{})

	return db, nil
}
