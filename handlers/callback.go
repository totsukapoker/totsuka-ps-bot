package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/totsukapoker/totsuka-ps-bot/config"
	"github.com/totsukapoker/totsuka-ps-bot/models"
	"github.com/totsukapoker/totsuka-ps-bot/repositories"

	"github.com/line/line-bot-sdk-go/linebot"
)

// CallbackHandler struct
type CallbackHandler struct {
	c     *gin.Context
	conf  *config.Config
	ur    *repositories.UserRepository
	gr    *repositories.GameRepository
	tr    *repositories.TransactionRepository
	usage string
}

// NewCallbackHandler creates new CallbackHandler.
func NewCallbackHandler(c *gin.Context, conf *config.Config, ur *repositories.UserRepository, gr *repositories.GameRepository, tr *repositories.TransactionRepository) *CallbackHandler {
	return &CallbackHandler{
		c:     c,
		conf:  conf,
		ur:    ur,
		gr:    gr,
		tr:    tr,
		usage: "こう使え:\n・現在額をそのまま入力(例:12340)\n・バイインした額を入力(例:+5000)\n・｢取消｣で1つ前の入力を取消\n・｢名前をxxxにして｣\n・｢名前をリセット｣",
	}
}

// Run executes handler.
func (h *CallbackHandler) Run() {
	proxyURL, _ := url.Parse(h.conf.ProxyURL)
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}
	bot, err := linebot.New(
		h.conf.LineChannelSecret, h.conf.LineChannelAccessToken, linebot.WithHTTPClient(client),
	)
	if err != nil {
		fmt.Println("ERROR linebot.New:", err)
		h.c.AbortWithStatus(500)
		return
	}

	events, err := bot.ParseRequest(h.c.Request)
	if err != nil {
		fmt.Println("ERROR ParseRequest:", err)
		if err == linebot.ErrInvalidSignature {
			h.c.AbortWithStatus(400)
		} else {
			h.c.AbortWithStatus(500)
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
				h.c.AbortWithStatus(400)
				return
			}
			user = h.ur.FirstOrCreate(event.Source.UserID, profile.DisplayName, profile.PictureURL, profile.StatusMessage)
		}

		if user.ID == 0 {
			fmt.Println("ERROR User is not specified, event.Source.UserID:", event.Source.UserID)
			h.c.AbortWithStatus(400)
			return
		}

		// Game loading
		game := h.gr.Current()
		if game.ID == 0 {
			fmt.Println("ERROR Game is not exist")
			if _, err = bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(
					"今お前と遊んでいる暇はない。",
				),
			).Do(); err != nil {
				fmt.Println("ERROR Game-is-not-exist ReplyMessage:", err)
				h.c.AbortWithStatus(500)
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
				m := h.normalizeMessage(message.Text)
				switch {
				case h.checkRegexp(`^\+[0-9]+$`, m): // バイイン時
					m, _ := strconv.Atoi(m)
					h.tr.Create(user, game, m, true)
					replyMessage = "バイインの入力をしたぞ！"
				case h.checkRegexp(`^[0-9]+$`, m): // 現在額入力時
					m, _ := strconv.Atoi(m)
					all := h.tr.CurrentAmountBy(user, game)
					h.tr.Create(user, game, m-all, false)
					replyMessage = "現在額の入力をしたぞ！"
				case h.checkRegexp(`^(今|いま)いく(つ|ら)(？|\?)$`, m): // 自分の状態質問時
					all := h.tr.CurrentAmountBy(user, game)
					buyin := h.tr.CurrentAmountBuyinBy(user, game)
					replyMessage = "現在額:" + strconv.Itoa(all) + "\nバイイン:" + strconv.Itoa(buyin)
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
				case h.checkRegexp(`^(取消|取り消し|取消し|とりけし|トリケシ|undo|UNDO|Undo)$`, m): // 1つ前のアクションを取り消し
					replyMessage = "お前に使う時間はない"
					t := h.tr.LastBy(user, game)
					if t.ID > 0 {
						h.tr.Delete(&t)
						replyMessage = "次はないぞ？心しろ。"
					}
				case h.checkRegexp(`^名前を.+にして$`, m): // 名前を設定
					// FIXME: 効率悪いけどもう一回正規表現使って名前部分だけを抜き出す
					r := regexp.MustCompile("^名前を(.+)にして$")
					g := r.FindStringSubmatch(m) // g[1] が名前になる
					h.ur.SetMyName(&user, g[1])
					replyMessage = g[1] + "にしたぞ。"
				case h.checkRegexp(`^名前を((消|け)して|リセット)$`, m): // 設定した名前をリセット
					replyMessage = "お前に使う時間はない"
					if user.MyName != "" {
						h.ur.SetMyName(&user, "")
						replyMessage = "今の名前(" + user.DisplayName + ")に戻したぞ。"
					}
				default:
					replyMessage = h.usage
				}
			default:
				replyMessage = h.usage
			}
			if replyMessage != "" {
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					fmt.Println("ERROR TypeMessage(Text) ReplyMessage:", err)
					h.c.AbortWithStatus(500)
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
				h.c.AbortWithStatus(500)
				return
			}
		case linebot.EventTypeUnfollow:
			// do something
		case linebot.EventTypePostback:
			// do something
		}
	}
}

func (h *CallbackHandler) checkRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

func (h *CallbackHandler) normalizeMessage(m string) (msg string) {
	msg = m
	return
}
