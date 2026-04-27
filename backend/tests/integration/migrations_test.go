package integration

import (
	"os"
	"testing"

	"github.com/supakorn5039-boon/saas-task-backend/src/database/migration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// connectFreshDB opens a brand-new connection (does NOT call SetupTestDB) so we
// can wipe the schema and run migrations from a known clean state. Tests in
// this file aren't concerned with services — they prove the migration runner
// itself works.
func connectFreshDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL_TEST")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=username password=password dbname=saas_test sslmode=disable TimeZone=UTC"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("test database unreachable: %v", err)
	}
	return db
}

// dropAll wipes every table the migrations might create. We use the public
// schema's table list rather than enumerating tables, so the test still works
// if migrations add new tables in the future.
func dropAll(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`).Error; err != nil {
		t.Fatalf("drop schema: %v", err)
	}
}

func TestMigration_RunFromCleanState(t *testing.T) {
	db := connectFreshDB(t)
	dropAll(t, db)

	if err := migration.Run(db); err != nil {
		t.Fatalf("first run: %v", err)
	}

	// Idempotent: running again with everything already migrated must succeed.
	if err := migration.Run(db); err != nil {
		t.Fatalf("second run (should be idempotent): %v", err)
	}

	// Schema_migrations should exist and be populated.
	var count int64
	if err := db.Raw(`SELECT COUNT(*) FROM schema_migrations`).Scan(&count).Error; err != nil {
		t.Fatalf("count migrations: %v", err)
	}
	if count == 0 {
		t.Error("schema_migrations is empty — migrations did not record themselves")
	}
}

func TestMigration_RollbackThenForward(t *testing.T) {
	db := connectFreshDB(t)
	dropAll(t, db)

	if err := migration.Run(db); err != nil {
		t.Fatalf("initial run: %v", err)
	}

	// Roll back every applied migration one at a time.
	var applied int64
	db.Raw(`SELECT COUNT(*) FROM schema_migrations`).Scan(&applied)
	if applied == 0 {
		t.Fatal("nothing to roll back")
	}

	for i := int64(0); i < applied; i++ {
		if err := migration.Rollback(db); err != nil {
			t.Fatalf("rollback iteration %d: %v", i, err)
		}
	}

	// schema_migrations should now be empty.
	var remaining int64
	db.Raw(`SELECT COUNT(*) FROM schema_migrations`).Scan(&remaining)
	if remaining != 0 {
		t.Errorf("after full rollback: %d entries remain, want 0", remaining)
	}

	// Roll back one more time — should be a no-op, not an error.
	if err := migration.Rollback(db); err != nil {
		t.Errorf("rollback on empty schema: %v (should be a graceful no-op)", err)
	}

	// And a fresh forward run should succeed.
	if err := migration.Run(db); err != nil {
		t.Fatalf("forward after full rollback: %v", err)
	}
}
