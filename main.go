package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/yakkun/totsuka-ps-bot/models"

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
			user := models.User{}
			if event.Source.UserID != "" {
				profile, err := bot.GetProfile(event.Source.UserID).Do()
				if err != nil {
					fmt.Println("ERROR GetProfile:", err, "UserID:", event.Source.UserID)
					c.AbortWithStatus(400)
					return
				}
				db.Where(models.User{UserID: event.Source.UserID}).Assign(models.User{DisplayName: profile.DisplayName, PictureURL: profile.PictureURL, StatusMessage: profile.StatusMessage}).FirstOrCreate(&user)
			}

			if user.ID == 0 {
				fmt.Println("ERROR User is not specified, event.Source.UserID:", event.Source.UserID)
				c.AbortWithStatus(400)
				return
			}

			// Game loading
			game := models.Game{}
			db.Where("? BETWEEN started_at AND ended_at", time.Now()).First(&game)
			if game.ID == 0 {
				fmt.Println("ERROR Game is not exist")
				if _, err = bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage(
						"現在開催されているゲームがないので利用できません。",
					),
				).Do(); err != nil {
					fmt.Println("ERROR Game-is-not-exist ReplyMessage:", err)
					c.AbortWithStatus(500)
				}
				return
			}

			// Support types: EventTypeMessage, EventTypeFollow, EventTypeUnfollow, EventTypePostback
			// Unsupport types: EventTypeJoin, EventTypeLeave, EventTypeBeacon
			// -> do nothing (ignore it)
			switch event.Type {
			case linebot.EventTypeMessage:
				replyMessage := ""
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					m := normalizeMessage(message.Text)
					switch {
					case checkRegexp(`^\+[0-9]+$`, m): // バイイン時
						m, _ := strconv.Atoi(m)
						transaction := models.Transaction{UserID: user.ID, GameID: game.ID, Amount: m, IsBuyin: true}
						db.Create(&transaction)
						replyMessage = "バイインの入力をしました"
					case checkRegexp(`^[0-9]+$`, m): // 現在額入力時
						m, _ := strconv.Atoi(m)
						type Result struct {
							Total int
						}
						var result Result
						db.Table("transactions").Select("SUM(amount) AS total").Where("user_id = ? AND game_id = ?", user.ID, game.ID).Scan(&result)
						transaction := models.Transaction{UserID: user.ID, GameID: game.ID, Amount: m - result.Total, IsBuyin: false}
						db.Create(&transaction)
						replyMessage = "現在額の入力をしました"
					case checkRegexp(`^(今|いま)いく(つ|ら)(？|\?)$`, m): // 自分の状態質問時
						type Result struct {
							Total string
						}
						var all Result
						db.Table("transactions").Select("IFNULL(SUM(amount), 0) AS total").Where("user_id = ? AND game_id = ?", user.ID, game.ID).Scan(&all)
						var buyin Result
						db.Table("transactions").Select("IFNULL(SUM(amount), 0) AS total").Where("user_id = ? AND game_id = ? AND is_buyin = ?", user.ID, game.ID, true).Scan(&buyin)
						replyMessage = "現在額:" + all.Total + "\nバイイン:" + buyin.Total
					default:
						replyMessage = usageMessage()
					}
				default:
					replyMessage = usageMessage()
				}
				if replyMessage != "" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						fmt.Println("ERROR TypeMessage(Text) ReplyMessage:", err)
						c.AbortWithStatus(500)
						return
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
					return
				}
			case linebot.EventTypeUnfollow:
				// do something
			case linebot.EventTypePostback:
				// do something
			}
		}
	})

	// GET: /result/:id
	// ゲーム(id=game_id)の現在の状況及び結果を表示する
	router.GET("/result/:id", func(c *gin.Context) {
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

func usageMessage() string {
	return "使い方:\n・現在の持ち点をそのまま入力(例:12340)\n・バイインで増やした点を入力(例:+5000)"
}

func normalizeMessage(m string) (msg string) {
	msg = m
	return
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
