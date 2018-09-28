package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	"github.com/line/line-bot-sdk-go/linebot"

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

	// db connection (gorm)
	db := ConnectDB()
	defer db.Close()
	// db migration
	MigrateDB(db)

	// Config for statis root
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	// GET: /
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	// POST: /callback
	router.POST("/callback", func(c *gin.Context) {
		proxyURL, _ := url.Parse(os.Getenv("FIXIE_URL"))
		client := &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		}
		bot, err := linebot.New(
			os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"), linebot.WithHTTPClient(client),
		)
		if err != nil {
			fmt.Println("ERROR linebot.New:", err)
			c.AbortWithStatus(500)
			return
		}

		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			fmt.Println("ERROR ParseRequest:", err)
			if err == linebot.ErrInvalidSignature {
				c.AbortWithStatus(400)
			} else {
				c.AbortWithStatus(500)
			}
			return
		}

		for _, event := range events {
			// User loading
			user := User{}
			if event.Source.UserID != "" {
				profile, err := bot.GetProfile(event.Source.UserID).Do()
				if err != nil {
					fmt.Println("ERROR GetProfile:", err, "UserID:", event.Source.UserID)
					c.AbortWithStatus(400)
				}
				db.Where(User{UserID: event.Source.UserID}).Assign(User{DisplayName: profile.DisplayName, PictureURL: profile.PictureURL, StatusMessage: profile.StatusMessage}).FirstOrCreate(&user)
			}

			// Support types: EventTypeMessage, EventTypeFollow, EventTypeUnfollow, EventTypePostback
			// Unsupport types: EventTypeJoin, EventTypeLeave, EventTypeBeacon
			// -> do nothing (ignore it)
			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(user.DisplayName+"さんこんにちは！\n"+message.Text),
					).Do(); err != nil {
						fmt.Println("ERROR TypeMessage(Text) ReplyMessage:", err)
						c.AbortWithStatus(500)
					}
				}
			case linebot.EventTypeFollow:
				if _, err = bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(
						"『怪物と闘う者は、その過程で自らが怪物と化さぬよう心せよ。"+user.DisplayName+"が長く深淵を覗くならば、深淵もまた等しく"+user.DisplayName+"を見返すのだ』",
					),
				).Do(); err != nil {
					fmt.Println("ERROR TypeFollow ReplyMessage:", err)
					c.AbortWithStatus(500)
				}
			case linebot.EventTypeUnfollow:
				// do something
			case linebot.EventTypePostback:
				// do something
			}
		}
	})

	router.Run(":" + port)
}
