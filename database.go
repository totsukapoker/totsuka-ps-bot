package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yakkun/totsuka-ps-bot/models"
)

// ConnectDB - Provide connection to database with gorm
func ConnectDB() (db *gorm.DB) {
	db, err := gorm.Open(connectionVars())
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	db.LogMode(true)
	db.DB().SetMaxIdleConns(0) // To avoid an error "Invalid Connection" on Heroku
	return
}

// MigrateDB - Do migration database with gorm.DB
func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Game{})
	db.AutoMigrate(&models.Transaction{})
}

func connectionVars() (driver string, source string) {
	databaseURL := ""
	if os.Getenv("DATABASE_URL") != "" {
		databaseURL = os.Getenv("DATABASE_URL")
	} else if os.Getenv("CLEARDB_DATABASE_URL") != "" {
		databaseURL = os.Getenv("CLEARDB_DATABASE_URL")
	}

	if databaseURL != "" {
		re, _ := regexp.Compile("([^:]+)://([^:]+):([^@]+)@([^/]+)/([^?]+)")
		match := re.FindStringSubmatch(databaseURL)
		driver = match[1]
		if driver == "mysql" {
			source = fmt.Sprintf(
				"%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=true&loc=Asia%%2FTokyo",
				match[2],
				match[3],
				match[4],
				match[5],
			)
		} else {
			source = databaseURL
		}
	}

	fmt.Println("driver:", driver)
	fmt.Println("source:", source)

	return
}