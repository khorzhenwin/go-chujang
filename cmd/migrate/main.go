package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_DSN") // full DSN for migrations
	m, err := migrate.New(
		"file://migrations",
		"postgres://"+dsn)
	if err != nil {
		log.Fatal("migrate.New failed: ", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Migration applied successfully ✅")
}
