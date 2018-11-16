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

	h, err := NewCallbackHandler(&c, &conf, &ur, &gr, &tr)
	if err != nil {
		t.Fatalf("%#v", err)
	}
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
	if h.Messages.Usage == "" {
		t.Errorf("CallbackHandler.Messages.Usage is empty")
	}
	if h.Messages.NoGame == "" {
		t.Errorf("CallbackHandler.Messages.NoGame is empty")
	}
	if h.Messages.Follow == "" {
		t.Errorf("CallbackHandler.Messages.Follow is empty")
	}
	if len(h.Messages.Gorilla) == 0 {
		t.Errorf("CallbackHandler.Messages.Gorilla is empty")
	}
	if h.Messages.BuyinDone == "" {
		t.Errorf("CallbackHandler.Messages.BuyinDone is empty")
	}
	if h.Messages.CurrentAmountDone == "" {
		t.Errorf("CallbackHandler.Messages.CurrentAmountDone is empty")
	}
	if h.Messages.CurrentInfo == "" {
		t.Errorf("CallbackHandler.Messages.CurrentInfo is empty")
	}
	if h.Messages.NoUndo == "" {
		t.Errorf("CallbackHandler.Messages.NoUndo is empty")
	}
	if h.Messages.UndoDone == "" {
		t.Errorf("CallbackHandler.Messages.UndoDone is empty")
	}
	if h.Messages.NamedDone == "" {
		t.Errorf("CallbackHandler.Messages.NamedDone is empty")
	}
	if h.Messages.NoNamed == "" {
		t.Errorf("CallbackHandler.Messages.NoNamed is empty")
	}
	if h.Messages.ResetNamedDone == "" {
		t.Errorf("CallbackHandler.Messages.ResetNamedDone is empty")
	}
}

func TestCallbackHandler_Run(t *testing.T) {
	t.Skip("implement me")
}

func TestCallbackHandler_loadMessages(t *testing.T) {
	t.Skip("implement me")
}

func TestCallbackHandler_loadUser(t *testing.T) {
	t.Skip("implement me")
}

func TestCallbackHandler_loadGame(t *testing.T) {
	t.Skip("implement me")
}

func TestCallbackHandler_loadLinebot(t *testing.T) {
	t.Run("If secret and token was not set", func(t *testing.T) {
		h := &CallbackHandler{conf: &config.Config{
			LineChannelSecret: "",
		}}
		err := h.loadLinebot()
		if err == nil {
			t.Fatalf("want: some error")
		}
		want := "missing channel secret"
		if err.Error() != want {
			t.Errorf("error = %#v; want: %#v", err.Error(), want)
		}
	})

	t.Run("All clean", func(t *testing.T) {
		conf := config.Config{
			LineChannelSecret:      "SOMESECRET",
			LineChannelAccessToken: "SOMETOKEN",
			ProxyURL:               "SOMEPROXY",
		}
		h := &CallbackHandler{conf: &conf}
		if err := h.loadLinebot(); err != nil {
			t.Fatalf("%#v", err)
		}
	})
}

func TestCallbackHandler_replyNoGame(t *testing.T) {
	t.Skip("implement me")
}

func TestCallbackHandler_checkRegexp(t *testing.T) {
	c := gin.Context{}
	conf := config.Config{}
	ur := repositories.UserRepository{}
	gr := repositories.GameRepository{}
	tr := repositories.TransactionRepository{}
	h, err := NewCallbackHandler(&c, &conf, &ur, &gr, &tr)
	if err != nil {
		t.Fatalf("%#v", err)
	}

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
	h, err := NewCallbackHandler(&c, &conf, &ur, &gr, &tr)
	if err != nil {
		t.Fatalf("%#v", err)
	}

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
