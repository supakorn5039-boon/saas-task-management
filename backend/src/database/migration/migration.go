package migration

import (
	"log"
	"sort"

	"gorm.io/gorm"
)

type Migration interface {
	// ID returns the migration identifier, matching the file name prefix.
	ID() string
	// Up runs the migration (raw SQL).
	Up(db *gorm.DB) error
	// Down rolls back the migration (raw SQL).
	Down(db *gorm.DB) error
}

type schemaMigration struct {
	ID string `gorm:"primaryKey"`
}

var registry []Migration

func Register(m Migration) {
	registry = append(registry, m)
}

func sorted() []Migration {
	sort.Slice(registry, func(i, j int) bool {
		return registry[i].ID() < registry[j].ID()
	})
	return registry
}

func Run(db *gorm.DB) error {
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		return err
	}

	for _, m := range sorted() {
		var record schemaMigration
		if db.First(&record, "id = ?", m.ID()).Error == nil {
			continue // already ran
		}

		log.Printf("[migration] running: %s", m.ID())
		if err := m.Up(db); err != nil {
			return err
		}

		if err := db.Create(&schemaMigration{ID: m.ID()}).Error; err != nil {
			return err
		}
		log.Printf("[migration] done: %s", m.ID())
	}

	return nil
}

func Rollback(db *gorm.DB) error {
	migrations := sorted()

	// Find the last executed migration
	for i := len(migrations) - 1; i >= 0; i-- {
		var record schemaMigration
		if db.First(&record, "id = ?", migrations[i].ID()).Error != nil {
			continue // not executed
		}

		log.Printf("[migration] rolling back: %s", migrations[i].ID())
		if err := migrations[i].Down(db); err != nil {
			return err
		}

		if err := db.Delete(&schemaMigration{}, "id = ?", migrations[i].ID()).Error; err != nil {
			return err
		}
		log.Printf("[migration] rolled back: %s", migrations[i].ID())
		return nil
	}

	log.Println("[migration] nothing to rollback")
	return nil
}

func Status(db *gorm.DB) {
	db.AutoMigrate(&schemaMigration{})

	log.Println("[migration] status:")
	for _, m := range sorted() {
		var record schemaMigration
		status := "pending"
		if db.First(&record, "id = ?", m.ID()).Error == nil {
			status = "migrated"
		}
		log.Printf("  %-50s %s", m.ID(), status)
	}
}
