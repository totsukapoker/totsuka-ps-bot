package config

import (
	"os"
	"strconv"
)

// Config struct for config values
type Config struct {
	Port                   int
	DbURL                  string
	ProxyURL               string
	LineChannelSecret      string
	LineChannelAccessToken string
}

// New returns config struct with default values.
func New() *Config {
	return &Config{
		Port:  8000,
		DbURL: "mysql://root:@localhost/totsuka_ps_bot",
	}
}

// Load is loading config values.
func (c *Config) Load() error {
	if err := c.loadPort(); err != nil {
		return err
	}
	if err := c.loadDbURL(); err != nil {
		return err
	}
	if err := c.loadProxyURL(); err != nil {
		return err
	}
	if err := c.loadLineChannelSecret(); err != nil {
		return err
	}
	if err := c.loadLineChannelAccessToken(); err != nil {
		return err
	}
	return nil

}

func (c *Config) loadPort() error {
	port := os.Getenv("PORT")
	if port != "" {
		p, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return err
		}
		c.Port = p
	}
	return nil
}

func (c *Config) loadDbURL() error {
	if os.Getenv("DATABASE_URL") != "" {
		c.DbURL = os.Getenv("DATABASE_URL")
	}
	if os.Getenv("CLEARDB_DATABASE_URL") != "" {
		c.DbURL = os.Getenv("CLEARDB_DATABASE_URL")
	}
	return nil
}

func (c *Config) loadProxyURL() error {
	if os.Getenv("PROXY_URL") != "" {
		c.ProxyURL = os.Getenv("PROXY_URL")
	}
	if os.Getenv("FIXIE_URL") != "" {
		c.ProxyURL = os.Getenv("FIXIE_URL")
	}
	return nil
}

func (c *Config) loadLineChannelSecret() error {
	c.LineChannelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	return nil
}

func (c *Config) loadLineChannelAccessToken() error {
	c.LineChannelAccessToken = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	return nil
}
