package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// ConnectDB - Provide connection to database with gorm
func ConnectDB() (db *gorm.DB) {
	db, err := gorm.Open(connectionVars())
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	db.LogMode(true)
	return
}

// MigrateDB - Do migration database with gorm.DB
func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Transaction{})
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
				"%s:%s@tcp(%s:3306)/%s?parseTime=true&charset=utf8",
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
