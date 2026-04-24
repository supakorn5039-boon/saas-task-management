package service

import (
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

var allowedSortColumns = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
	"title":      "title",
	"status":     "status",
}

type ListTasksOptions struct {
	UserID  uint
	Page    int
	PerPage int
	Status  model.TaskStatus
	Search  string
	Sort    string
	Order   string
}

func (s *TaskService) ListTasks(opts ListTasksOptions) (*model.TaskListResponse, error) {
	q := s.db.Model(&model.Task{}).Where("user_id = ?", opts.UserID)

	if opts.Status != "" {
		q = q.Where("status = ?", opts.Status)
	}
	if opts.Search != "" {
		like := "%" + opts.Search + "%"
		q = q.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	sortCol, ok := allowedSortColumns[opts.Sort]
	if !ok {
		sortCol = "created_at"
	}
	order := "desc"
	if opts.Order == "asc" {
		order = "asc"
	}

	var tasks []model.Task
	err := q.Order(sortCol + " " + order).
		Limit(opts.PerPage).
		Offset((opts.Page - 1) * opts.PerPage).
		Find(&tasks).Error
	if err != nil {
		return nil, err
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
	err := s.db.Model(&model.Task{}).
		Select("status, count(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Find(&rows).Error
	if err != nil {
		return model.TaskListCounts{}, err
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

func (s *TaskService) CreateTask(userID uint, title, description string) (*model.TaskDto, error) {
	task := model.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      model.TaskStatusTodo,
	}
	if err := s.db.Create(&task).Error; err != nil {
		return nil, err
	}
	return task.ToDto(), nil
}

func (s *TaskService) UpdateStatus(userID uint, taskID uint, status model.TaskStatus) (*model.TaskDto, error) {
	var task model.Task
	if err := s.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		return nil, err
	}
	task.Status = status
	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}
	return task.ToDto(), nil
}

func (s *TaskService) DeleteTask(userID uint, taskID uint) error {
	return s.db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&model.Task{}).Error
}
