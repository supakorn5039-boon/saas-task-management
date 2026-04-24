package migration

import "gorm.io/gorm"

func init() {
	Register(&createRoles{})
}

type createRoles struct{}

func (m *createRoles) ID() string { return "2026_04_10_00_00_01_create_roles" }

func (m *createRoles) Up(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id          BIGSERIAL    PRIMARY KEY,
			name        VARCHAR(255) NOT NULL UNIQUE,
			description VARCHAR(255),
			created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			deleted_at  TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at);
	`).Error
}

func (m *createRoles) Down(db *gorm.DB) error {
	return db.Exec(`DROP TABLE IF EXISTS roles;`).Error
}
