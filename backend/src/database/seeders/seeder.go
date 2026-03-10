package seeders

import (
	"log"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	log.Println("Seeding database...")

	allSeeders := []struct {
		name   string
		seeder func(*gorm.DB)
	}{
		{"RoleSeeder", RoleSeeder},
		{"UserSeeder", UserSeeder},
	}

	for _, s := range allSeeders {
		log.Printf("Running seeder: %s", s.name)
		s.seeder(db)
	}

	log.Println("Seeding completed successfully.")
}
