package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/totsukapoker/totsuka-ps-bot/handlers"
	"github.com/totsukapoker/totsuka-ps-bot/repositories"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"github.com/totsukapoker/totsuka-ps-bot/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file. But you could be ignore me.")
	}

	conf, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %+v", err)
	}

	// Prepare http router (gin)
	router := gin.New()
	router.Use(gin.Logger())

	// Config for statis root
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	// db connection (gorm)
	db := ConnectDB(conf.DbURL)
	defer db.Close()
	MigrateDB(db)

	// Repositories
	ur := repositories.NewUserRepository(db)
	gr := repositories.NewGameRepository(db)
	tr := repositories.NewTransactionRepository(db)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/callback", func(c *gin.Context) {
		handlers.NewCallbackHandler(c, conf, ur, gr, tr).Run()
	})

	router.GET("/result/:id", func(c *gin.Context) {
		handlers.NewResultHandler(c, conf, ur, gr, tr).Run()
	})

	router.Run(":" + strconv.Itoa(conf.Port))
}
