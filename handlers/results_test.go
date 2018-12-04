package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/totsukapoker/totsuka-ps-bot/config"
	"github.com/totsukapoker/totsuka-ps-bot/repositories"
)

func TestNewResultsHandler(t *testing.T) {
	c := gin.Context{}
	conf := config.Config{}
	gr := repositories.GameRepository{}

	h := NewResultsHandler(&c, &conf, &gr)
	if h.c != &c {
		t.Errorf("got: %v, expected: %v", h.c, c)
	}
	if h.conf != &conf {
		t.Errorf("got: %v, expected: %v", h.conf, conf)
	}
	if h.gr != &gr {
		t.Errorf("got: %v, expected: %v", h.gr, gr)
	}
}

func TestResultsHandler_Run(t *testing.T) {
	t.Skip("implement me")
}

func TestResultsHandler_showErrorHTML(t *testing.T) {
	t.Skip("implement me")
}
