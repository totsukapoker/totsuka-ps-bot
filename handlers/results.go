package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/totsukapoker/totsuka-ps-bot/config"
	"github.com/totsukapoker/totsuka-ps-bot/repositories"
)

// ResultsHandler struct
type ResultsHandler struct {
	c    *gin.Context
	conf *config.Config
	gr   *repositories.GameRepository
}

// NewResultsHandler creates new ResultHandler.
func NewResultsHandler(c *gin.Context, conf *config.Config, gr *repositories.GameRepository) *ResultsHandler {
	return &ResultsHandler{
		c:    c,
		conf: conf,
		gr:   gr,
	}
}

// Run executes handler.
func (h *ResultsHandler) Run() error {
	// id が "current" の場合は現在行われているゲームがあればその結果へリダイレクトする
	if h.c.Param("id") == "current" {
		game := h.gr.Current()
		if game.ID == 0 {
			h.showErrorHTML(http.StatusNotFound, "No game is running now.")
			return nil
		}
		h.c.Redirect(http.StatusFound, "/results/"+fmt.Sprint(game.ID))
		return nil
	}

	id, err := strconv.Atoi(h.c.Param("id"))
	if err != nil {
		h.showErrorHTML(http.StatusInternalServerError, "strconv error")
		return err
	}
	if id <= 0 {
		h.showErrorHTML(http.StatusBadRequest, "Need valid id")
		return nil
	}

	game := h.gr.First(uint(id))
	if game.ID == 0 {
		h.showErrorHTML(http.StatusNotFound, "Not found")
		return nil
	}

	h.c.HTML(http.StatusOK, "results.tmpl.html", gin.H{})
	return nil
}

func (h *ResultsHandler) showErrorHTML(code int, message string) {
	h.c.HTML(code, "error.tmpl.html", gin.H{
		"code":    code,
		"message": message,
	})
}
