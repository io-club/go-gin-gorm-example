package model

import (
	"fmt"

	"fibric/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var db *gorm.DB
	var err error
	switch config.DB_MODE {
	case "debug":
		// connect to sqlite
		db, err = gorm.Open(sqlite.Open("./cache/test.db"), &gorm.Config{})
	case "release":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName, config.DBCharset)
		// set timeout
		dsn += "&timeout=10s&readTimeout=30s&writeTimeout=30s&parseTime=true"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		panic("mode not correct")
	}
	if err != nil {
		panic("failed to connect database")
	}

	// Create tables
	DB = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
	DB.AutoMigrate(&Fabric{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Image{})
	DB.AutoMigrate(&Brand{})
	DB.AutoMigrate(&Trend{})
	DB.AutoMigrate(&Cloth{})
	DB.AutoMigrate(&Dress{})
	DB.AutoMigrate(&News{})
	DB = db
}

var tableNames *[]string = nil

func GetTableNames() []string {
	if tableNames == nil {
		tableNames = &[]string{
			Fabric{}.TableName(),
			Brand{}.TableName(),
			Trend{}.TableName(),
			Cloth{}.TableName(),
			Dress{}.TableName(),
			News{}.TableName(),
		}
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
