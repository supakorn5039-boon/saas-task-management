package migration

import "gorm.io/gorm"

func init() {
	Register(&tasksJiraFields{})
}

type tasksJiraFields struct{}

func (m *tasksJiraFields) ID() string {
	return "2026_04_29_00_00_01_tasks_jira_fields"
}

// Adds the Jira/Asana-style task fields:
//
//   - priority: one of low/medium/high/urgent (default 'medium').
//   - start_date / due_date: optional, used for scheduling.
//   - assignee_user_id: who owns the task. Defaults to user_id (creator) for
//     existing rows so the column is NOT NULL going forward without breaking
//     the historical data.
//
// Indexes pick the access patterns the UI actually uses: filtering/sorting
// by due date, filtering by priority, and grouping by assignee.
func (m *tasksJiraFields) Up(db *gorm.DB) error {
	return db.Exec(`
		ALTER TABLE tasks
			ADD COLUMN IF NOT EXISTS priority VARCHAR(16) NOT NULL DEFAULT 'medium',
			ADD COLUMN IF NOT EXISTS start_date TIMESTAMPTZ,
			ADD COLUMN IF NOT EXISTS due_date TIMESTAMPTZ,
			ADD COLUMN IF NOT EXISTS assignee_user_id BIGINT REFERENCES users(id);

		UPDATE tasks SET assignee_user_id = user_id WHERE assignee_user_id IS NULL;

		CREATE INDEX IF NOT EXISTS idx_tasks_assignee ON tasks(assignee_user_id);
		CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
		CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	`).Error
}

func (m *tasksJiraFields) Down(db *gorm.DB) error {
	return db.Exec(`
		DROP INDEX IF EXISTS idx_tasks_priority;
		DROP INDEX IF EXISTS idx_tasks_due_date;
		DROP INDEX IF EXISTS idx_tasks_assignee;
		ALTER TABLE tasks
			DROP COLUMN IF EXISTS assignee_user_id,
			DROP COLUMN IF EXISTS due_date,
			DROP COLUMN IF EXISTS start_date,
			DROP COLUMN IF EXISTS priority;
	`).Error
}
