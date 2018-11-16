package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/totsukapoker/totsuka-ps-bot/config"
	"github.com/totsukapoker/totsuka-ps-bot/repositories"
)

func TestNewResultHandler(t *testing.T) {
	c := gin.Context{}
	conf := config.Config{}
	ur := repositories.UserRepository{}
	gr := repositories.GameRepository{}
	tr := repositories.TransactionRepository{}

	h := NewResultHandler(&c, &conf, &ur, &gr, &tr)
	if h.c != &c {
		t.Errorf("got: %v, expected: %v", h.c, c)
	}
	if h.conf != &conf {
		t.Errorf("got: %v, expected: %v", h.conf, conf)
	}
	if h.ur != &ur {
		t.Errorf("got: %v, expected: %v", h.ur, ur)
	}
	if h.gr != &gr {
		t.Errorf("got: %v, expected: %v", h.gr, gr)
	}
	if h.tr != &tr {
		t.Errorf("got: %v, expected: %v", h.tr, tr)
	}
}

func TestResultHandler_Run(t *testing.T) {
	t.Skip("implement me")
}

func TestResultHandler_showErrorHTML(t *testing.T) {
	t.Skip("implement me")
}
