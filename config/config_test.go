package config

import (
	"os"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	config, err := Load()
	if err != nil {
		t.Errorf("%#v", err)
	}
	if config.Port != 8000 {
		t.Fatalf("WRONG: default port, current: \"%d\"", config.Port)
	}
	if config.DbURL != "mysql://root:@localhost/totsuka_ps_bot" {
		t.Fatalf("WRONG: default DbURL, current: \"%v\"", config.DbURL)
	}
	if config.ProxyURL != "" {
		t.Fatalf("WRONG: default ProxyURL, current: \"%v\"", config.ProxyURL)
	}
	if config.LineChannelSecret != "" {
		t.Fatalf("WRONG: default LineChannelSecret, current: \"%v\"", config.LineChannelSecret)
	}
	if config.LineChannelAccessToken != "" {
		t.Fatalf("WRONG: default LineChannelAccessToken, current: \"%v\"", config.LineChannelAccessToken)
	}
}

func TestLoadPort(t *testing.T) {
	if value := os.Getenv("PORT"); value != "" {
		os.Unsetenv("PORT")
		defer func() { os.Setenv("PORT", value) }()
	}

	os.Setenv("PORT", "8989")
	config, _ := Load()
	if config.Port != 8989 {
		t.Fatalf("WRONG: Port, current: \"%d\"", config.Port)
	}
}

func TestLoadDbURL(t *testing.T) {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		os.Unsetenv("DATABASE_URL")
		defer func() { os.Setenv("DATABASE_URL", value) }()
	}
	if value := os.Getenv("CLEARDB_DATABASE_URL"); value != "" {
		os.Unsetenv("CLEARDB_DATABASE_URL")
		defer func() { os.Setenv("CLEARDB_DATABASE_URL", value) }()
	}

	var config *Config

	os.Setenv("DATABASE_URL", "somedatabaseurl1")
	config, _ = Load()
	if config.DbURL != "somedatabaseurl1" {
		t.Fatalf("WRONG: DbURL, current: \"%v\"", config.DbURL)
	}

	os.Setenv("DATABASE_URL", "")
	os.Setenv("CLEARDB_DATABASE_URL", "somecleardburl1")
	config, _ = Load()
	if config.DbURL != "somecleardburl1" {
		t.Fatalf("WRONG: DbURL, current: \"%v\"", config.DbURL)
	}

	os.Setenv("DATABASE_URL", "somedatabaseurl2")
	os.Setenv("CLEARDB_DATABASE_URL", "somecleardburl2")
	config, _ = Load()
	if config.DbURL != "somedatabaseurl2" {
		t.Fatalf("WRONG: DbURL, current: \"%v\"", config.DbURL)
	}
}

func TestLoadProxyURL(t *testing.T) {
	if value := os.Getenv("PROXY_URL"); value != "" {
		os.Unsetenv("PROXY_URL")
		defer func() { os.Setenv("PROXY_URL", value) }()
	}
	if value := os.Getenv("FIXIE_URL"); value != "" {
		os.Unsetenv("FIXIE_URL")
		defer func() { os.Setenv("FIXIE_URL", value) }()
	}

	var config *Config

	os.Setenv("PROXY_URL", "someproxyurl1")
	config, _ = Load()
	if config.ProxyURL != "someproxyurl1" {
		t.Fatalf("WRONG: ProxyURL, current: \"%v\"", config.ProxyURL)
	}

	os.Setenv("PROXY_URL", "")
	os.Setenv("FIXIE_URL", "somefixieurl1")
	config, _ = Load()
	if config.ProxyURL != "somefixieurl1" {
		t.Fatalf("WRONG: ProxyURL, current: \"%v\"", config.ProxyURL)
	}

	os.Setenv("PROXY_URL", "someproxyurl2")
	os.Setenv("FIXIE_URL", "somefixieurl2")
	config, _ = Load()
	if config.ProxyURL != "someproxyurl2" {
		t.Fatalf("WRONG: ProxyURL, current: \"%v\"", config.ProxyURL)
	}
}

func TestLineChannelSecret(t *testing.T) {
	if value := os.Getenv("LINE_CHANNEL_SECRET"); value != "" {
		os.Unsetenv("LINE_CHANNEL_SECRET")
		defer func() { os.Setenv("LINE_CHANNEL_SECRET", value) }()
	}

	os.Setenv("LINE_CHANNEL_SECRET", "somelinesecret")
	config, _ := Load()
	if config.LineChannelSecret != "somelinesecret" {
		t.Fatalf("WRONG: LineChannelSecret, current: \"%v\"", config.LineChannelSecret)
	}
}

func TestLineChannelAccessToken(t *testing.T) {
	if value := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"); value != "" {
		os.Unsetenv("LINE_CHANNEL_ACCESS_TOKEN")
		defer func() { os.Setenv("LINE_CHANNEL_ACCESS_TOKEN", value) }()
	}

	os.Setenv("LINE_CHANNEL_ACCESS_TOKEN", "somelineaccesstoken")
	config, _ := Load()
	if config.LineChannelAccessToken != "somelineaccesstoken" {
		t.Fatalf("WRONG: LineChannelAccessToken, current: \"%v\"", config.LineChannelAccessToken)
	}
}
