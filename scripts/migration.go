package scripts

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(dsn string) {
	log.Println("📦 Running migrations via golang-migrate...")

	m, err := migrate.New("file://../migrate/migrations", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to initialize migrate: %v", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Migrations complete")
}
