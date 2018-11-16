package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/totsukapoker/totsuka-ps-bot/config"
	"github.com/totsukapoker/totsuka-ps-bot/repositories"
)

func TestNewCallbackHandler(t *testing.T) {
	c := gin.Context{}
	conf := config.Config{}
	ur := repositories.UserRepository{}
	gr := repositories.GameRepository{}
	tr := repositories.TransactionRepository{}

	h := NewCallbackHandler(&c, &conf, &ur, &gr, &tr)
	if h.c != &c {
		t.Errorf("CallbackHandler.c = %#v; want: %#v", h.c, c)
	}
	if h.conf != &conf {
		t.Errorf("CallbackHandler.conf = %#v; want: %#v", h.conf, conf)
	}
	if h.ur != &ur {
		t.Errorf("CallbackHandler.ur = %#v; want: %#v", h.ur, ur)
	}
	if h.gr != &gr {
		t.Errorf("CallbackHandler.gr = %#v; want: %#v", h.gr, gr)
	}
	if h.tr != &tr {
		t.Errorf("CallbackHandler.tr =  %#v; want: %#v", h.tr, tr)
	}
	expectedUsage := "こう使え:\n・現在額をそのまま入力(例:12340)\n・バイインした額を入力(例:+5000)\n・｢取消｣で1つ前の入力を取消\n・｢名前をxxxにして｣\n・｢名前をリセット｣"
	if h.usage != expectedUsage {
		t.Errorf("CallbackHandler.usage = %#v; want: %#v", h.usage, expectedUsage)
	}
}

func TestCallbackHandler_Run(t *testing.T) {
	t.Skip("implement me")
}

func TestCallbackHandler_checkRegexp(t *testing.T) {
	c := gin.Context{}
	conf := config.Config{}
	ur := repositories.UserRepository{}
	gr := repositories.GameRepository{}
	tr := repositories.TransactionRepository{}
	h := NewCallbackHandler(&c, &conf, &ur, &gr, &tr)

	tests := []struct {
		reg, str string
		want     bool
	}{
		{"", "", true},
		{".", "A", true},
		{".", "", false},
		{"^やっくん$", "やっくん", true},
		{"^もっくん$", "やっくん", false},
		{"^やっ", "やっくん", true},
		{"^やっ", "もっくん", false},
		{"っくん$", "やっくん", true},
		{"っくん$", "やっくそ", false},
		{"^(や|も)っくん$", "もっくん", true},
		{"^(や|も)っくん$", "とっくん", false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := h.checkRegexp(tt.reg, tt.str); got != tt.want {
				t.Errorf("CallbackHandler.checkRegexp(%#v, %#v) = %#v; want: %#v", tt.reg, tt.str, got, tt.want)
			}
		})
	}
}

func TestCallbackHandler_normalizeMessage(t *testing.T) {
	c := gin.Context{}
	conf := config.Config{}
	ur := repositories.UserRepository{}
	gr := repositories.GameRepository{}
	tr := repositories.TransactionRepository{}
	h := NewCallbackHandler(&c, &conf, &ur, &gr, &tr)

	tests := []struct {
		str, want string
	}{
		{"", ""},
		{"やっくん", "やっくん"},
		{"や　く ん", "や　く ん"},
		{"AbCdE", "AbCdE"},
		{"８９3", "８９3"},
		{"ＹＡＫ", "ＹＡＫ"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := h.normalizeMessage(tt.str); got != tt.want {
				t.Errorf("CallbackHandler.normalizeMessage(%#v) = %#v; want: %#v", tt.str, got, tt.want)
			}
		})
	}
}
