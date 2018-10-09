package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config struct for config values
type Config struct {
	Port                   int
	DbURL                  string
	FixieURL               string
	LineChannelSecret      string
	LineChannelAccessToken string
}

// Load config.Load() to load config values from .env and os.Getenv
func Load() (*Config, error) {
	// default values
	config := &Config{
		Port:  8000,
		DbURL: "root@/totsuka_ps_bot",
	}

	loadDotenv()

	err := loadAll(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func loadDotenv() error {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file. But you could be ignore me.")
	}
	return err
}

func loadAll(config *Config) error {
	var err error
	err = loadPort(config)
	if err != nil {
		return err
	}
	err = loadDbURL(config)
	if err != nil {
		return err
	}
	err = loadFixieURL(config)
	if err != nil {
		return err
	}
	err = loadLineChannelSecret(config)
	if err != nil {
		return err
	}
	err = loadLineChannelAccessToken(config)
	if err != nil {
		return err
	}
	return err

}

func loadPort(config *Config) error {
	port := os.Getenv("PORT")
	if port != "" {
		p, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return fmt.Errorf("Invalid PORT: %v", err)
		}
		config.Port = p
	}
	return nil
}

func loadDbURL(config *Config) error {
	var dbURL string
	if os.Getenv("DATABASE_URL") != "" {
		dbURL = os.Getenv("DATABASE_URL")
	} else if os.Getenv("CLEARDB_DATABASE_URL") != "" {
		dbURL = os.Getenv("CLEARDB_DATABASE_URL")
	}
	if dbURL == "" {
		return fmt.Errorf("Require DATABASE_URL")
	}
	config.DbURL = dbURL
	return nil
}

func loadFixieURL(config *Config) error {
	config.FixieURL = os.Getenv("FIXIE_URL")
	return nil
}

func loadLineChannelSecret(config *Config) error {
	config.LineChannelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	return nil
}

func loadLineChannelAccessToken(config *Config) error {
	config.LineChannelAccessToken = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	return nil
}