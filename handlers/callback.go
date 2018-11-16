package handlers

import (
	"errors"
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
	c    *gin.Context
	conf *config.Config
	ur   *repositories.UserRepository
	gr   *repositories.GameRepository
	tr   *repositories.TransactionRepository

	Messages *CallbackMessages

	linebot *linebot.Client
	user    models.User
	game    models.Game
}

// CallbackMessages struct
type CallbackMessages struct {
	Usage             string
	NoGame            string
	Follow            string
	Gorilla           []string
	BuyinDone         string
	CurrentAmountDone string
	CurrentInfo       string
	NoUndo            string
	UndoDone          string
	NamedDone         string
	NoNamed           string
	ResetNamedDone    string
}

// NewCallbackHandler creates new CallbackHandler.
func NewCallbackHandler(c *gin.Context, conf *config.Config, ur *repositories.UserRepository, gr *repositories.GameRepository, tr *repositories.TransactionRepository) (*CallbackHandler, error) {
	h := &CallbackHandler{
		c:    c,
		conf: conf,
		ur:   ur,
		gr:   gr,
		tr:   tr,
	}
	if err := h.loadMessages(); err != nil {
		return h, err
	}
	return h, nil
}

// Run executes handler.
func (h *CallbackHandler) Run() error {
	if err := h.loadLinebot(); err != nil {
		h.c.AbortWithStatus(500)
		return err
	}

	events, err := h.linebot.ParseRequest(h.c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			h.c.AbortWithStatus(400)
		} else {
			h.c.AbortWithStatus(500)
		}
		return err
	}

	for _, event := range events {
		if err := h.loadUser(event.Source.UserID); err != nil {
			h.c.AbortWithStatus(400)
			return err
		}

		if err := h.loadGame(); err != nil {
			h.c.AbortWithStatus(500)
			return err
		}

		if h.game.ID == 0 {
			if err := h.replyNoGame(event.ReplyToken); err != nil {
				h.c.AbortWithStatus(500)
				return err
			}
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
					h.tr.Create(h.user, h.game, m, true)
					replyMessage = h.Messages.BuyinDone
				case h.checkRegexp(`^[0-9]+$`, m): // 現在額入力時
					m, _ := strconv.Atoi(m)
					all := h.tr.CurrentAmountBy(h.user, h.game)
					h.tr.Create(h.user, h.game, m-all, false)
					replyMessage = h.Messages.CurrentAmountDone
				case h.checkRegexp(`^(今|いま)いく(つ|ら)(？|\?)$`, m): // 自分の状態質問時
					replyMessage = fmt.Sprintf(h.Messages.CurrentInfo, h.tr.CurrentAmountBy(h.user, h.game), h.tr.CurrentAmountBuyinBy(h.user, h.game))
				case m == "ウホウホ": // ゴリラボタン
					rand.Seed(time.Now().UnixNano())
					replyMessage = h.Messages.Gorilla[rand.Intn(len(h.Messages.Gorilla))]
				case h.checkRegexp(`^(取消|取り消し|取消し|とりけし|トリケシ|undo|UNDO|Undo)$`, m): // 1つ前のアクションを取り消し
					replyMessage = h.Messages.NoUndo
					t := h.tr.LastBy(h.user, h.game)
					if t.ID > 0 {
						h.tr.Delete(&t)
						replyMessage = h.Messages.UndoDone
					}
				case h.checkRegexp(`^名前を.+にして$`, m): // 名前を設定
					// FIXME: 効率悪いけどもう一回正規表現使って名前部分だけを抜き出す
					r := regexp.MustCompile("^名前を(.+)にして$")
					g := r.FindStringSubmatch(m) // g[1] が名前になる
					h.ur.SetMyName(&h.user, g[1])
					replyMessage = fmt.Sprintf(h.Messages.NamedDone, g[1])
				case h.checkRegexp(`^名前を((消|け)して|リセット)$`, m): // 設定した名前をリセット
					replyMessage = h.Messages.NoNamed
					if h.user.MyName != "" {
						h.ur.SetMyName(&h.user, "")
						replyMessage = fmt.Sprintf(h.Messages.ResetNamedDone, h.user.DisplayName)
					}
				default:
					replyMessage = h.Messages.Usage
				}
			default:
				replyMessage = h.Messages.Usage
			}
			if replyMessage != "" {
				if _, err = h.linebot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					h.c.AbortWithStatus(500)
					return err
				}
			}
		case linebot.EventTypeFollow:
			if _, err = h.linebot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(fmt.Sprintf(h.Messages.Follow, h.user.DisplayName, h.user.DisplayName)),
			).Do(); err != nil {
				h.c.AbortWithStatus(500)
				return err
			}
		case linebot.EventTypeUnfollow:
			// do something
		case linebot.EventTypePostback:
			// do something
		}
	}
	return nil
}

func (h *CallbackHandler) loadMessages() error {
	h.Messages = &CallbackMessages{
		Usage:  "こう使え:\n・現在額をそのまま入力(例:12340)\n・バイインした額を入力(例:+5000)\n・｢取消｣で1つ前の入力を取消\n・｢名前をxxxにして｣\n・｢名前をリセット｣",
		NoGame: "今お前と遊んでいる暇はない。",
		Follow: "『怪物と闘う者は、その過程で自らが怪物と化さぬよう心せよ。%sが長く深淵を覗くならば、深淵もまた等しく%sを見返すのだ』",
		Gorilla: []string{
			"殺すぞ",
			"ブチ殺すぞ",
			"アア？",
			"俺とじゃれるか？",
			"お前は俺の玩具だ",
			"俺の握力は600kgだ",
		},
		BuyinDone:         "バイインの入力をしたぞ！",
		CurrentAmountDone: "現在額の入力をしたぞ！",
		CurrentInfo:       "現在額:%d\nバイイン:%d",
		NoUndo:            "お前に使う時間はない",
		UndoDone:          "次はないぞ？心しろ。",
		NamedDone:         "%sにしたぞ。",
		NoNamed:           "お前に使う時間はない",
		ResetNamedDone:    "今の名前(%s)に戻したぞ。",
	}
	return nil
}

func (h *CallbackHandler) loadLinebot() error {
	proxyURL, _ := url.Parse(h.conf.ProxyURL)
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}
	bot, err := linebot.New(
		h.conf.LineChannelSecret, h.conf.LineChannelAccessToken, linebot.WithHTTPClient(client),
	)
	if err != nil {
		return err
	}
	h.linebot = bot
	return nil
}

func (h *CallbackHandler) loadUser(userID string) error {
	if userID == "" {
		return errors.New("userID is empty")
	}
	profile, err := h.linebot.GetProfile(userID).Do()
	if err != nil {
		return err
	}
	h.user = h.ur.FirstOrCreate(userID, profile.DisplayName, profile.PictureURL, profile.StatusMessage)
	return nil
}

func (h *CallbackHandler) loadGame() error {
	h.game = h.gr.Current()
	return nil
}

func (h *CallbackHandler) replyNoGame(replyToken string) error {
	if h.game.ID > 0 {
		return errors.New("game is available")
	}
	if _, err := h.linebot.ReplyMessage(replyToken, linebot.NewTextMessage(h.Messages.NoGame)).Do(); err != nil {
		return err
	}
	return nil
}

func (h *CallbackHandler) checkRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

func (h *CallbackHandler) normalizeMessage(m string) (msg string) {
	msg = m
	return
}
