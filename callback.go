package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/yakkun/totsuka-ps-bot/config"
	"github.com/yakkun/totsuka-ps-bot/models"

	"github.com/line/line-bot-sdk-go/linebot"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func callback(c *gin.Context, db *gorm.DB, conf *config.Config) {
	proxyURL, _ := url.Parse(conf.ProxyURL)
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}
	bot, err := linebot.New(
		conf.LineChannelSecret, conf.LineChannelAccessToken, linebot.WithHTTPClient(client),
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
					"今お前と遊んでいる暇はない。",
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
					replyMessage = "バイインの入力をしたぞ！"
				case checkRegexp(`^[0-9]+$`, m): // 現在額入力時
					m, _ := strconv.Atoi(m)
					type Result struct {
						Total int
					}
					var result Result
					db.Table("transactions").Select("SUM(amount) AS total").Where("user_id = ? AND game_id = ?", user.ID, game.ID).Scan(&result)
					transaction := models.Transaction{UserID: user.ID, GameID: game.ID, Amount: m - result.Total, IsBuyin: false}
					db.Create(&transaction)
					replyMessage = "現在額の入力をしたぞ！"
				case checkRegexp(`^(今|いま)いく(つ|ら)(？|\?)$`, m): // 自分の状態質問時
					type Result struct {
						Total string
					}
					var all Result
					db.Table("transactions").Select("IFNULL(SUM(amount), 0) AS total").Where("user_id = ? AND game_id = ?", user.ID, game.ID).Scan(&all)
					var buyin Result
					db.Table("transactions").Select("IFNULL(SUM(amount), 0) AS total").Where("user_id = ? AND game_id = ? AND is_buyin = ?", user.ID, game.ID, true).Scan(&buyin)
					replyMessage = "現在額:" + all.Total + "\nバイイン:" + buyin.Total
				case m == "ウホウホ": // ゴリラボタン
					replyMessages := []string{
						"殺すぞ",
						"ブチ殺すぞ",
						"アア？",
						"俺とじゃれるか？",
						"お前は俺の玩具だ",
						"俺の握力は600kgだ",
					}
					rand.Seed(time.Now().UnixNano())
					replyMessage = replyMessages[rand.Intn(len(replyMessages))]
				case checkRegexp(`^(取消|取り消し|取消し|とりけし|トリケシ|undo|UNDO|Undo)$`, m): // 1つ前のアクションを取り消し
					replyMessage = "お前に使う時間はない"
					var t models.Transaction
					db.Where("user_id = ? AND game_id = ?", user.ID, game.ID).Order("id desc").First(&t)
					if t.ID > 0 {
						db.Delete(&t)
						replyMessage = "次はないぞ？心しろ。"
					}
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
}

func usageMessage() string {
	return "こう使え:\n・現在額をそのまま入力(例:12340)\n・バイインした額を入力(例:+5000)\n・｢取消｣で1つ前の入力を取消"
}

func normalizeMessage(m string) (msg string) {
	msg = m
	return
}
