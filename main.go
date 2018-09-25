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
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

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
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						fmt.Println("ERROR ReplyMessage:", err)
						c.AbortWithStatus(500)
					}
				}
			}
		}
	})

	router.Run(":" + port)
}
