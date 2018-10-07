package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/yakkun/totsuka-ps-bot/models"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Prepare http router (gin)
	router := gin.New()
	router.Use(gin.Logger())

	// Config for statis root
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	// GET: /
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	// POST: /callback
	router.POST("/callback", callback)

	// GET: /result/:id
	// ゲーム(id=game_id)の現在の状況及び結果を表示する
	router.GET("/result/:id", func(c *gin.Context) {
		// db connection (gorm)
		db := ConnectDB()
		defer db.Close()
		// db migration
		MigrateDB(db)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			showErrorHTML(c, 500, "strconv error")
			return
		}
		if id <= 0 {
			showErrorHTML(c, 400, "Need id")
			return
		}
		game := models.Game{}
		db.First(&game, id)
		if game.ID == 0 {
			showErrorHTML(c, 404, "Not found")
			return
		}
		c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
			"game": game,
		})
	})

	router.Run(":" + port)
}

func checkRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

func showErrorHTML(c *gin.Context, code int, message string) {
	c.HTML(code, "error.tmpl.html", gin.H{
		"code":    code,
		"message": message,
	})
}
