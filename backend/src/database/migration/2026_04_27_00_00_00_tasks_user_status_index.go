package migration

import "gorm.io/gorm"

func init() {
	Register(&tasksUserStatusIndex{})
}

type tasksUserStatusIndex struct{}

func (m *tasksUserStatusIndex) ID() string {
	return "2026_04_27_00_00_00_tasks_user_status_index"
}

func (m *tasksUserStatusIndex) Up(db *gorm.DB) error {
	return db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tasks_user_status ON tasks(user_id, status);
	`).Error
}

func (m *tasksUserStatusIndex) Down(db *gorm.DB) error {
	return db.Exec(`
		DROP INDEX IF EXISTS idx_tasks_user_status;
	`).Error
}
