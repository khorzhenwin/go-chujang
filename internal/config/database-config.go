package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	DSN      string
}

func LoadDBConfig() (*DBConfig, error) {
	_ = godotenv.Load() // Loads from .env file

	cfg := &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL"),
		DSN:      os.Getenv("DB_DSN"),
	}

	// Validate required fields
	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.Name == "" {
		return nil, fmt.Errorf("incomplete DB config")
	}

	return cfg, nil
}

func (c *DBConfig) GetFormattedDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode,
	)
}
