package model

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

func (s TaskStatus) Valid() bool {
	switch s {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	}
	return false
}

// TaskPriority — Jira-style priority levels. Stored as a string so adding a
// new level (e.g. "blocker") is a non-breaking change for callers that already
// know how to render an unknown value.
type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityUrgent TaskPriority = "urgent"
)

func (p TaskPriority) Valid() bool {
	switch p {
	case TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh, TaskPriorityUrgent:
		return true
	}
	return false
}

type Task struct {
	gorm.Model
	Title       string       `gorm:"not null"`
	Description string       `gorm:"null"`
	Status      TaskStatus   `gorm:"type:varchar(20);not null;default:'todo'"`
	Priority    TaskPriority `gorm:"type:varchar(16);not null;default:'medium'"`
	StartDate   *time.Time   ``
	DueDate     *time.Time   ``
	UserID      uint         `gorm:"not null"`
	User        User         `gorm:"foreignKey:UserID"`
	// AssigneeUserID — who currently owns the task. Defaults to UserID
	// (creator) for backwards compatibility but can be reassigned at any
	// time by the creator or an admin.
	AssigneeUserID *uint ``
	Assignee       *User `gorm:"foreignKey:AssigneeUserID"`
}

type TaskDto struct {
	ID             uint         `json:"id"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Status         TaskStatus   `json:"status"`
	Priority       TaskPriority `json:"priority"`
	StartDate      *string      `json:"startDate,omitempty"`
	DueDate        *string      `json:"dueDate,omitempty"`
	AssigneeID     *uint        `json:"assigneeId,omitempty"`
	AssigneeEmail  string       `json:"assigneeEmail,omitempty"`
	UserID         uint         `json:"userId"`
	CreatedAt      string       `json:"createdAt"`
	UpdatedAt      string       `json:"updatedAt"`
}

type TaskListMeta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"perPage"`
	Total   int64 `json:"total"`
}

type TaskListCounts struct {
	All        int64 `json:"all"`
	Todo       int64 `json:"todo"`
	InProgress int64 `json:"in_progress"`
	Done       int64 `json:"done"`
}

type TaskListResponse struct {
	Data   []*TaskDto     `json:"data"`
	Meta   TaskListMeta   `json:"meta"`
	Counts TaskListCounts `json:"counts"`
}

func (t *Task) ToDto() *TaskDto {
	dto := &TaskDto{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    t.Priority,
		UserID:      t.UserID,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if t.StartDate != nil {
		s := t.StartDate.Format("2006-01-02T15:04:05Z07:00")
		dto.StartDate = &s
	}
	if t.DueDate != nil {
		d := t.DueDate.Format("2006-01-02T15:04:05Z07:00")
		dto.DueDate = &d
	}
	if t.AssigneeUserID != nil {
		dto.AssigneeID = t.AssigneeUserID
	}
	if t.Assignee != nil && t.Assignee.ID != 0 {
		dto.AssigneeEmail = t.Assignee.Email
	}
	return dto
}
