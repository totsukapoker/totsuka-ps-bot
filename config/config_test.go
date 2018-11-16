package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	got := New()
	want := &Config{
		Port:                   8000,
		DbURL:                  "mysql://root:@localhost/totsuka_ps_bot",
		ProxyURL:               "",
		LineChannelSecret:      "",
		LineChannelAccessToken: "",
	}
	if got.Port != want.Port {
		t.Errorf("config.Port() = %#v; want: %#v", got.Port, want.Port)
	}
	if got.DbURL != want.DbURL {
		t.Errorf("config.DbURL() = %#v; want: %#v", got.DbURL, want.DbURL)
	}
	if got.ProxyURL != want.ProxyURL {
		t.Errorf("config.ProxyURL() = %#v; want: %#v", got.ProxyURL, want.ProxyURL)
	}
	if got.LineChannelSecret != want.LineChannelSecret {
		t.Errorf("config.LineChannelSecret() = %#v; want: %#v", got.LineChannelSecret, want.LineChannelSecret)
	}
	if got.LineChannelAccessToken != want.LineChannelAccessToken {
		t.Errorf("config.LineChannelAccessToken() = %#v; want: %#v", got.LineChannelAccessToken, want.LineChannelAccessToken)
	}
}

func TestConfig_Load(t *testing.T) {
	config := New()
	err := config.Load()
	if err != nil {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func TestConfig_loadPort(t *testing.T) {
	if value := os.Getenv("PORT"); value != "" {
		os.Unsetenv("PORT")
		defer func() { os.Setenv("PORT", value) }()
	}

	tests := []struct {
		port string
		want int
	}{
		{"", 8000},
		{"8989", 8989},
		{"0080", 80},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			os.Setenv("PORT", tt.port)
			config := New()
			if err := config.loadPort(); err != nil {
				t.Fatalf("%#v", err)
			}
			if config.Port != tt.want {
				t.Errorf("config.Port = %d; want: %d", config.Port, tt.want)
			}
		})
	}
}

func TestConfig_loadDbURL(t *testing.T) {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		os.Unsetenv("DATABASE_URL")
		defer func() { os.Setenv("DATABASE_URL", value) }()
	}
	if value := os.Getenv("CLEARDB_DATABASE_URL"); value != "" {
		os.Unsetenv("CLEARDB_DATABASE_URL")
		defer func() { os.Setenv("CLEARDB_DATABASE_URL", value) }()
	}

	tests := []struct {
		dbURL, cleardbURL, want string
	}{
		{"", "", "mysql://root:@localhost/totsuka_ps_bot"},
		{"somedburl1", "", "somedburl1"},
		{"", "somecleardburl", "somecleardburl"},
		{"somedburl2", "somecleardbur2", "somecleardbur2"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			os.Setenv("DATABASE_URL", tt.dbURL)
			os.Setenv("CLEARDB_DATABASE_URL", tt.cleardbURL)
			config := New()
			if err := config.loadDbURL(); err != nil {
				t.Fatalf("%#v", err)
			}
			if config.DbURL != tt.want {
				t.Errorf("config.DbURL = %#v; want: %#v", config.DbURL, tt.want)
			}
		})
	}
}

func TestConfig_loadProxyURL(t *testing.T) {
	if value := os.Getenv("PROXY_URL"); value != "" {
		os.Unsetenv("PROXY_URL")
		defer func() { os.Setenv("PROXY_URL", value) }()
	}
	if value := os.Getenv("FIXIE_URL"); value != "" {
		os.Unsetenv("FIXIE_URL")
		defer func() { os.Setenv("FIXIE_URL", value) }()
	}

	tests := []struct {
		proxyURL, fixieURL, want string
	}{
		{"", "", ""},
		{"someproxyurl1", "", "someproxyurl1"},
		{"", "somefixieurl1", "somefixieurl1"},
		{"someproxyurl2", "somefixieurl2", "somefixieurl2"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			os.Setenv("PROXY_URL", tt.proxyURL)
			os.Setenv("FIXIE_URL", tt.fixieURL)
			config := New()
			if err := config.loadProxyURL(); err != nil {
				t.Fatalf("%#v", err)
			}
			if config.ProxyURL != tt.want {
				t.Errorf("config.ProxyURL = %#v; want: %#v", config.ProxyURL, tt.want)
			}
		})
	}
}

func TestConfig_loadLineChannelSecret(t *testing.T) {
	if value := os.Getenv("LINE_CHANNEL_SECRET"); value != "" {
		os.Unsetenv("LINE_CHANNEL_SECRET")
		defer func() { os.Setenv("LINE_CHANNEL_SECRET", value) }()
	}

	want := "somelinesecret"
	os.Setenv("LINE_CHANNEL_SECRET", want)
	config := New()
	if err := config.loadLineChannelSecret(); err != nil {
		t.Fatalf("%#v", err)
	}
	if config.LineChannelSecret != want {
		t.Errorf("config.LineChannelSecret = %#v; want: %#v", config.LineChannelSecret, want)
	}
}

func TestConfig_loadLineChannelAccessToken(t *testing.T) {
	if value := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"); value != "" {
		os.Unsetenv("LINE_CHANNEL_ACCESS_TOKEN")
		defer func() { os.Setenv("LINE_CHANNEL_ACCESS_TOKEN", value) }()
	}

	want := "somelineaccesstoken"
	os.Setenv("LINE_CHANNEL_ACCESS_TOKEN", want)
	config := New()
	if err := config.loadLineChannelAccessToken(); err != nil {
		t.Fatalf("%#v", err)
	}
	if config.LineChannelAccessToken != want {
		t.Errorf("config.LineChannelAccessToken = %#v; want: %#v", config.LineChannelAccessToken, want)
	}
}
