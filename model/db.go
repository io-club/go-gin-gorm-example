package model

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var db *gorm.DB
	var err error
	switch MODE {
	case "debug":
		// connect to sqlite
		db, err = gorm.Open(sqlite.Open("./cache/test.db"), &gorm.Config{})
	case "release":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", DBUser, DBPassword, DBHost, DBPort, DBName, DBCharset)
		// set timeout
		dsn += "&timeout=10s&readTimeout=30s&writeTimeout=30s&parseTime=true"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		panic("mode not correct")
	}
	if err != nil {
		panic("failed to connect database")
	}
	DB = db

	// Create tables
	DB.AutoMigrate(&Fabric{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Image{})
	DB.AutoMigrate(&Brand{})
	DB.AutoMigrate(&Trend{})
	DB.AutoMigrate(&Cloth{})
	DB.AutoMigrate(&Dress{})
	DB.AutoMigrate(&News{})
}

var tableNames *[]string = nil

func GetTableNames() []string {
	if tableNames == nil {
		tableNames = &[]string{}
		DB.Raw("SHOW TABLES").Scan(tableNames)
	}
	return *tableNames
}

func TableExists(tableName string) bool {
	tables := GetTableNames()
	for _, table := range tables {
		if table == tableName {
			return true
		}
	}
	return false
}
