package handlers

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/totsukapoker/totsuka-ps-bot/config"
	"github.com/totsukapoker/totsuka-ps-bot/models"
	"github.com/totsukapoker/totsuka-ps-bot/repositories"
)

// ResultHandler struct
type ResultHandler struct {
	c    *gin.Context
	conf *config.Config
	ur   *repositories.UserRepository
	gr   *repositories.GameRepository
	tr   *repositories.TransactionRepository
}

// NewResultHandler creates new ResultHandler.
func NewResultHandler(c *gin.Context, conf *config.Config, ur *repositories.UserRepository, gr *repositories.GameRepository, tr *repositories.TransactionRepository) *ResultHandler {
	return &ResultHandler{
		c:    c,
		conf: conf,
		ur:   ur,
		gr:   gr,
		tr:   tr,
	}
}

// Run executes handler.
func (h *ResultHandler) Run() {
	// id が "current" の場合は現在行われているゲームがあればその結果へリダイレクトする
	if h.c.Param("id") == "current" {
		game := h.gr.Current()
		if game.ID == 0 {
			h.showErrorHTML(http.StatusNotFound, "No game is running now.")
		}
		h.c.Redirect(http.StatusMovedPermanently, "/result/"+fmt.Sprint(game.ID))
		return
	}

	id, err := strconv.Atoi(h.c.Param("id"))
	if err != nil {
		h.showErrorHTML(http.StatusInternalServerError, "strconv error")
		return
	}
	if id <= 0 {
		h.showErrorHTML(http.StatusBadRequest, "Need valid id")
		return
	}

	game := h.gr.First(uint(id))
	if game.ID == 0 {
		h.showErrorHTML(http.StatusNotFound, "Not found")
		return
	}

	transactions := h.tr.FindByGame(game)

	var userIDs []uint
L:
	for _, t := range transactions {
		for _, i := range userIDs {
			if t.UserID == i {
				continue L
			}
		}
		userIDs = append(userIDs, t.UserID)
	}
	users := h.ur.FindByIDs(userIDs)

	type Stat struct {
		User          models.User
		CurrentAmount int
		BuyinAmount   int
		ROI           float64
		StartedAt     time.Time
		UpdatedAt     time.Time
	}
	type TotalStat struct {
		CurrentAmount       int
		BuyinAmount         int
		DifferenceAmount    int
		DifferenceAbsAmount int
	}
	var stats []Stat
	var totalstat TotalStat
	for _, u := range users {
		var s Stat
		s.User = u
		stats = append(stats, s)
	}
	for _, t := range transactions {
		var stat *Stat
		for i, s := range stats {
			if s.User.ID == t.UserID {
				stat = &stats[i]
				break
			}
		}
		stat.CurrentAmount += t.Amount
		if t.IsBuyin == true {
			stat.BuyinAmount += t.Amount
		}
		if stat.StartedAt.IsZero() == true || stat.StartedAt.After(t.CreatedAt) == true {
			stat.StartedAt = t.CreatedAt
		}
		if stat.UpdatedAt.Before(t.CreatedAt) == true {
			stat.UpdatedAt = t.CreatedAt
		}
	}
	for i, s := range stats {
		totalstat.CurrentAmount += s.CurrentAmount
		totalstat.BuyinAmount += s.BuyinAmount
		if s.BuyinAmount > 0 {
			stats[i].ROI = float64(s.CurrentAmount)/float64(s.BuyinAmount)*100 - 100
		}
	}
	totalstat.DifferenceAmount = totalstat.CurrentAmount - totalstat.BuyinAmount
	totalstat.DifferenceAbsAmount = int(math.Abs(float64(totalstat.DifferenceAmount)))
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].StartedAt.Before(stats[j].StartedAt)
	})

	type Log struct {
		ID        uint
		Amount    int
		IsBuyin   bool
		CreatedAt time.Time
		User      models.User
	}
	var logs []Log
	for _, t := range transactions {
		user := models.User{}
		for _, u := range users {
			if u.ID == t.UserID {
				user = u
			}
		}
		var l Log
		l.ID = t.ID
		l.Amount = t.Amount
		l.IsBuyin = t.IsBuyin
		l.CreatedAt = t.CreatedAt
		l.User = user
		logs = append(logs, l)
	}

	h.c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
		"currentTime": time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)),
		"game":        game,
		"stats":       stats,
		"totalstat":   totalstat,
		"logs":        logs,
	})
}

func (h *ResultHandler) showErrorHTML(code int, message string) {
	h.c.HTML(code, "error.tmpl.html", gin.H{
		"code":    code,
		"message": message,
	})
}
