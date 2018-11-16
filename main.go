package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/totsukapoker/totsuka-ps-bot/repositories"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"github.com/totsukapoker/totsuka-ps-bot/config"
	"github.com/totsukapoker/totsuka-ps-bot/models"

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

	// GET: /
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	// POST: /callback
	router.POST("/callback", func(c *gin.Context) {
		callback(c, conf, ur, gr, tr)
	})

	// GET: /result/:id
	// ゲーム(id=game_id)の現在の状況及び結果を表示する
	router.GET("/result/:id", func(c *gin.Context) {
		// id が "current" の場合は現在行われているゲームがあればその結果へリダイレクトする
		if c.Param("id") == "current" {
			game := gr.Current()
			if game.ID == 0 {
				showErrorHTML(c, http.StatusNotFound, "No game is running now.")
			}
			c.Redirect(http.StatusMovedPermanently, "/result/"+fmt.Sprint(game.ID))
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			showErrorHTML(c, http.StatusInternalServerError, "strconv error")
			return
		}
		if id <= 0 {
			showErrorHTML(c, http.StatusBadRequest, "Need valid id")
			return
		}

		game := gr.First(uint(id))
		if game.ID == 0 {
			showErrorHTML(c, http.StatusNotFound, "Not found")
			return
		}

		transactions := tr.FindByGame(game)

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
		users := ur.FindByIDs(userIDs)

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

		c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
			"currentTime": time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)),
			"game":        game,
			"stats":       stats,
			"totalstat":   totalstat,
			"logs":        logs,
		})
	})

	router.Run(":" + strconv.Itoa(conf.Port))
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

func showErrorJSON(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
