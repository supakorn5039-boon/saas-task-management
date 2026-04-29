package service

import (
	"errors"
	"time"

	"github.com/supakorn5039-boon/saas-task-backend/src/apperror"
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService() *TaskService {
	return &TaskService{database.DB}
}

// allowedSortColumns whitelists what callers can sort by — protects against
// SQL injection from arbitrary `sort=` query strings and keeps the query
// shape predictable.
var allowedSortColumns = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
	"title":      "title",
	"status":     "status",
	"priority":   "priority",
	"due_date":   "due_date",
}

type ListTasksOptions struct {
	UserID   uint
	Page     int
	PerPage  int
	Status   model.TaskStatus
	Priority model.TaskPriority
	// AssigneeID filters tasks assigned to a specific user. When nil the
	// service uses the visibility rule below ("created by or assigned to me").
	AssigneeID *uint
	Search     string
	Sort       string
	Order      string
}

// scopeForUser returns a *gorm.DB scoped to tasks the calling user can see —
// either created by them OR assigned to them. This is the same query used by
// every read endpoint so the visibility rule is defined exactly once.
func (s *TaskService) scopeForUser(userID uint) *gorm.DB {
	return s.db.Model(&model.Task{}).
		Where("user_id = ? OR assignee_user_id = ?", userID, userID)
}

func (s *TaskService) ListTasks(opts ListTasksOptions) (*model.TaskListResponse, error) {
	q := s.scopeForUser(opts.UserID)

	if opts.Status != "" {
		q = q.Where("status = ?", opts.Status)
	}
	if opts.Priority != "" {
		q = q.Where("priority = ?", opts.Priority)
	}
	if opts.AssigneeID != nil {
		q = q.Where("assignee_user_id = ?", *opts.AssigneeID)
	}
	if opts.Search != "" {
		like := "%" + opts.Search + "%"
		q = q.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, apperror.Wrap(err, 500, "failed to count tasks")
	}

	sortCol, ok := allowedSortColumns[opts.Sort]
	if !ok {
		sortCol = "created_at"
	}
	order := "desc"
	if opts.Order == "asc" {
		order = "asc"
	}
	// Sorting by due_date with NULLs last — Postgres puts NULL first by
	// default which buries actual deadlines.
	orderExpr := sortCol + " " + order
	if sortCol == "due_date" {
		orderExpr = "due_date " + order + " NULLS LAST"
	}

	var tasks []model.Task
	err := q.Preload("Assignee").
		Order(orderExpr).
		Limit(opts.PerPage).
		Offset((opts.Page - 1) * opts.PerPage).
		Find(&tasks).Error
	if err != nil {
		return nil, apperror.Wrap(err, 500, "failed to list tasks")
	}

	dtos := make([]*model.TaskDto, len(tasks))
	for i, t := range tasks {
		dtos[i] = t.ToDto()
	}

	counts, err := s.statusCounts(opts.UserID)
	if err != nil {
		return nil, err
	}

	return &model.TaskListResponse{
		Data:   dtos,
		Meta:   model.TaskListMeta{Page: opts.Page, PerPage: opts.PerPage, Total: total},
		Counts: counts,
	}, nil
}

func (s *TaskService) statusCounts(userID uint) (model.TaskListCounts, error) {
	rows := []struct {
		Status model.TaskStatus
		Count  int64
	}{}
	err := s.scopeForUser(userID).
		Select("status, count(*) as count").
		Group("status").
		Find(&rows).Error
	if err != nil {
		return model.TaskListCounts{}, apperror.Wrap(err, 500, "failed to count tasks by status")
	}

	out := model.TaskListCounts{}
	for _, r := range rows {
		switch r.Status {
		case model.TaskStatusTodo:
			out.Todo = r.Count
		case model.TaskStatusInProgress:
			out.InProgress = r.Count
		case model.TaskStatusDone:
			out.Done = r.Count
		}
		out.All += r.Count
	}
	return out, nil
}

// CreateTaskInput is the full set of fields accepted at creation time. Only
// Title is mandatory — everything else has a sensible default (medium
// priority, todo status, assignee = creator).
type CreateTaskInput struct {
	Title       string
	Description string
	Priority    model.TaskPriority
	StartDate   *time.Time
	DueDate     *time.Time
	AssigneeID  *uint
}

func (s *TaskService) CreateTask(userID uint, in CreateTaskInput) (*model.TaskDto, error) {
	priority := in.Priority
	if priority == "" {
		priority = model.TaskPriorityMedium
	}
	if !priority.Valid() {
		return nil, apperror.BadRequest("invalid priority")
	}
	assignee := in.AssigneeID
	if assignee == nil {
		me := userID
		assignee = &me
	}

	task := model.Task{
		UserID:         userID,
		AssigneeUserID: assignee,
		Title:          in.Title,
		Description:    in.Description,
		Status:         model.TaskStatusTodo,
		Priority:       priority,
		StartDate:      in.StartDate,
		DueDate:        in.DueDate,
	}
	if err := s.db.Create(&task).Error; err != nil {
		return nil, apperror.Wrap(err, 500, "failed to create task")
	}
	// Reload with assignee preloaded so the response carries the email.
	if err := s.db.Preload("Assignee").First(&task, task.ID).Error; err != nil {
		return nil, apperror.Wrap(err, 500, "failed to reload task")
	}
	return task.ToDto(), nil
}

// UpdateTaskInput carries the fields that may be patched on a task. All
// pointer-typed so the service can tell "absent" from "set to empty/zero".
// DueDate / StartDate / AssigneeID use a double pointer so the caller can
// also explicitly clear them — e.g. dueDate=null in the JSON body.
type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *model.TaskStatus
	Priority    *model.TaskPriority
	// **time.Time / **uint encoding: outer-pointer-nil means "no change",
	// outer-pointer-non-nil with inner-pointer-nil means "clear to NULL".
	StartDate  **time.Time
	DueDate    **time.Time
	AssigneeID **uint
}

func (s *TaskService) UpdateTask(userID uint, taskID uint, in UpdateTaskInput) (*model.TaskDto, error) {
	var task model.Task
	if err := s.scopeForUser(userID).First(&task, "tasks.id = ?", taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("task not found")
		}
		return nil, apperror.Wrap(err, 500, "failed to load task")
	}

	if in.Title != nil {
		task.Title = *in.Title
	}
	if in.Description != nil {
		task.Description = *in.Description
	}
	if in.Status != nil {
		task.Status = *in.Status
	}
	if in.Priority != nil {
		if !in.Priority.Valid() {
			return nil, apperror.BadRequest("invalid priority")
		}
		task.Priority = *in.Priority
	}
	if in.StartDate != nil {
		task.StartDate = *in.StartDate
	}
	if in.DueDate != nil {
		task.DueDate = *in.DueDate
	}
	if in.AssigneeID != nil {
		task.AssigneeUserID = *in.AssigneeID
	}

	if err := s.db.Save(&task).Error; err != nil {
		return nil, apperror.Wrap(err, 500, "failed to update task")
	}
	if err := s.db.Preload("Assignee").First(&task, task.ID).Error; err != nil {
		return nil, apperror.Wrap(err, 500, "failed to reload task")
	}
	return task.ToDto(), nil
}

func (s *TaskService) DeleteTask(userID uint, taskID uint) error {
	// Only the creator can delete — assignees can update but not destroy.
	result := s.db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&model.Task{})
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, "failed to delete task")
	}
	if result.RowsAffected == 0 {
		return apperror.NotFound("task not found")
	}
	return nil
}
