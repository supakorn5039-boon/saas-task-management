package migration

import "gorm.io/gorm"

func init() {
	Register(&createUsers{})
}

type createUsers struct{}

func (m *createUsers) ID() string { return "2026_04_10_00_00_00_create_users" }

func (m *createUsers) Up(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id         BIGSERIAL    PRIMARY KEY,
			email      VARCHAR(255) NOT NULL UNIQUE,
			password   VARCHAR(255) NOT NULL,
			status     INT          NOT NULL DEFAULT 0,
			role       VARCHAR(50)  NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
	`).Error
}

func (m *createUsers) Down(db *gorm.DB) error {
	return db.Exec(`DROP TABLE IF EXISTS users;`).Error
}
