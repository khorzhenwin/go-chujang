package db

import (
	"fmt"
	"github.com/khorzhenwin/go-chujang/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(config *config.DBConfig) (*gorm.DB, error) {
	cfg := postgres.Config{
		DSN:                  config.GetFormattedDSN(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}

	db, err := gorm.Open(postgres.New(cfg), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("DB connection failed: %w", err)
	}
	return db, nil
}
