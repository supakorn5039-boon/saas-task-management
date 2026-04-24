package seeder

import (
	"log"
	"sort"

	"gorm.io/gorm"
)

type Seeder interface {
	ID() string

	Description() string

	Seed(db *gorm.DB) error
}

var registry []Seeder

func Register(s Seeder) {
	registry = append(registry, s)
}

func sorted() []Seeder {
	sort.Slice(registry, func(i, j int) bool {
		return registry[i].ID() < registry[j].ID()
	})
	return registry
}

func Seed(db *gorm.DB) error {
	for _, s := range sorted() {
		log.Printf("[seeder] running: %s — %s", s.ID(), s.Description())
		if err := s.Seed(db); err != nil {
			return err
		}
		log.Printf("[seeder] done: %s", s.ID())
	}
	return nil
}
