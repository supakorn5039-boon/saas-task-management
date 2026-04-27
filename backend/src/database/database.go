package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/supakorn5039-boon/saas-task-backend/src/config"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/migration"
)

var DB *gorm.DB

func Connect() error {
	dsn := buildDSN(config.App.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return migration.Run(DB)
}

// buildDSN prefers DATABASE_URL when set (cloud Postgres providers like Neon,
// Supabase, Render hand out a DSN directly). Falls back to assembling one from
// discrete fields for local docker-compose dev.
func buildDSN(cfg config.DatabaseConfig) string {
	if cfg.DSN != "" {
		return cfg.DSN
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
	)
}
