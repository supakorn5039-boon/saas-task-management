package migration

import "gorm.io/gorm"

func init() {
	Register(&taskStatus{})
}

type taskStatus struct{}

func (m *taskStatus) ID() string { return "2026_04_24_00_00_00_task_status" }

func (m *taskStatus) Up(db *gorm.DB) error {
	return db.Exec(`
		ALTER TABLE tasks ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'todo';
		UPDATE tasks SET status = 'done' WHERE completed = TRUE;
		ALTER TABLE tasks DROP COLUMN completed;
	`).Error
}

func (m *taskStatus) Down(db *gorm.DB) error {
	return db.Exec(`
		ALTER TABLE tasks ADD COLUMN completed BOOLEAN NOT NULL DEFAULT FALSE;
		UPDATE tasks SET completed = TRUE WHERE status = 'done';
		ALTER TABLE tasks DROP COLUMN status;
	`).Error
}
