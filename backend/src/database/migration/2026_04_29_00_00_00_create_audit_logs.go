package migration

import "gorm.io/gorm"

func init() {
	Register(&createAuditLogs{})
}

type createAuditLogs struct{}

func (m *createAuditLogs) ID() string { return "2026_04_29_00_00_00_create_audit_logs" }

// Audit log is append-only: no updated_at, no soft delete. The actor's email
// is denormalized at write time so the log stays readable even if the user
// row is later deleted. Metadata is JSONB for flexible, queryable context
// (target name, before/after values, etc).
func (m *createAuditLogs) Up(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS audit_logs (
			id            BIGSERIAL    PRIMARY KEY,
			actor_user_id BIGINT,
			actor_email   VARCHAR(255),
			action        VARCHAR(64)  NOT NULL,
			target_type   VARCHAR(32),
			target_id     BIGINT,
			status        VARCHAR(16)  NOT NULL DEFAULT 'success',
			ip            VARCHAR(64),
			user_agent    VARCHAR(512),
			metadata      JSONB,
			created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_created ON audit_logs(actor_user_id, created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_action_created ON audit_logs(action, created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at DESC);
	`).Error
}

func (m *createAuditLogs) Down(db *gorm.DB) error {
	return db.Exec(`DROP TABLE IF EXISTS audit_logs;`).Error
}
