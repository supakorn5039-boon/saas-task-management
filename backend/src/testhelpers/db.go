// Package testhelpers wires up a fresh Postgres-backed *gorm.DB for tests.
//
// Usage:
//
//	func TestSomething(t *testing.T) {
//	    db := testhelpers.SetupTestDB(t)
//	    // ... services that read database.DB now use the test database
//	}
//
// Connection: reads DATABASE_URL_TEST env var (DSN string, e.g.
// "host=localhost port=5432 user=... password=... dbname=saas_test sslmode=disable").
// Falls back to a sensible local-docker-compose default if unset.
//
// If the database is unreachable, t.Skip is called rather than t.Fatal —
// this lets developers run other test suites without a Postgres instance.
package testhelpers

import (
	"os"
	"testing"

	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultTestDSN = "host=localhost port=5432 user=username password=password dbname=saas_test sslmode=disable TimeZone=UTC"

// SetupTestDB connects to the test database, ensures the schema exists,
// truncates all tables, and points database.DB at the test connection so
// services constructed with NewXService() pick it up.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL_TEST")
	if dsn == "" {
		dsn = defaultTestDSN
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("test database unreachable (set DATABASE_URL_TEST or start docker-compose): %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Task{}, &model.Role{}); err != nil {
		t.Fatalf("auto-migrate test schema: %v", err)
	}

	// CASCADE so foreign-key relations don't block the truncate.
	if err := db.Exec(`TRUNCATE TABLE users, tasks, roles RESTART IDENTITY CASCADE`).Error; err != nil {
		t.Fatalf("truncate test tables: %v", err)
	}

	database.DB = db
	return db
}
