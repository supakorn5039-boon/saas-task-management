package migration

import "gorm.io/gorm"

func init() {
	Register(&createTasks{})
}

type createTasks struct{}

func (m *createTasks) ID() string { return "2026_04_10_00_00_02_create_tasks" }

func (m *createTasks) Up(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id          BIGSERIAL    PRIMARY KEY,
			title       VARCHAR(255) NOT NULL,
			description TEXT,
			completed   BOOLEAN      NOT NULL DEFAULT FALSE,
			user_id     BIGINT       NOT NULL REFERENCES users(id),
			created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			deleted_at  TIMESTAMPTZ
		);
		CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);
		CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
	`).Error
}

func (m *createTasks) Down(db *gorm.DB) error {
	return db.Exec(`DROP TABLE IF EXISTS tasks;`).Error
}
