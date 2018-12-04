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

	conf := config.New()
	if err := conf.Load(); err != nil {
		log.Fatalf("Failed to load config: %#v", err)
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
		h, err := handlers.NewCallbackHandler(c, conf, ur, gr, tr)
		if err != nil {
			log.Printf("NewCallbackHandler error: %#v", err)
		}
		if err := h.Run(); err != nil {
			log.Printf("CallbackHandler.Run() error: %#v", err)
		}
	})

	router.GET("/result/:id", func(c *gin.Context) {
		if err := handlers.NewResultHandler(c, conf, ur, gr, tr).Run(); err != nil {
			log.Printf("ResultHandler error: %#v", err)
		}
	})

	router.GET("/results/:id", func(c *gin.Context) {
		if err := handlers.NewResultsHandler(c, conf, gr).Run(); err != nil {
			log.Printf("ResultsHandler error: %#v", err)
		}
	})

	api := router.Group("/api")
	{
		api.GET("/games/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				log.Fatalf("strconv error: %#v", err)
			}
			game := gr.First(uint(id))
			c.JSON(http.StatusOK, game)
		})
	}

	router.Run(":" + strconv.Itoa(conf.Port))
}
