package model

import "gorm.io/gorm"

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

type Task struct {
	gorm.Model
	Title       string     `gorm:"not null"`
	Description string     `gorm:"null"`
	Status      TaskStatus `gorm:"type:varchar(20);not null;default:'todo'"`
	UserID      uint       `gorm:"not null"`
	User        User       `gorm:"foreignKey:UserID"`
}

type TaskDto struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   string     `json:"createdAt"`
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
	return &TaskDto{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
