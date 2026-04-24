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

func (s *TaskService) GetAllTasks(userID uint) ([]*model.TaskDto, error) {
	var tasks []model.Task
	if err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&tasks).Error; err != nil {
		return nil, err
	}

	dtos := make([]*model.TaskDto, len(tasks))
	for i, task := range tasks {
		dtos[i] = task.ToDto()
	}
	return dtos, nil
}

func (s *TaskService) CreateTask(userID uint, title, description string) (*model.TaskDto, error) {
	task := model.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Completed:   false,
	}

	if err := s.db.Create(&task).Error; err != nil {
		return nil, err
	}

	return task.ToDto(), nil
}

func (s *TaskService) ToggleTask(userID uint, taskID uint, completed bool) (*model.TaskDto, error) {
	var task model.Task
	if err := s.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		return nil, err
	}

	task.Completed = completed
	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return task.ToDto(), nil
}

func (s *TaskService) DeleteTask(userID uint, taskID uint) error {
	return s.db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&model.Task{}).Error
}
