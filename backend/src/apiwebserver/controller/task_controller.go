package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
)

type TaskController struct {
	svc *service.TaskService
}

func NewTaskController() *TaskController {
	return &TaskController{svc: service.NewTaskService()}
}

func (ctrl *TaskController) RegisterRoutes(r *gin.RouterGroup) {
	tasks := r.Group("/tasks")
	tasks.Use(middleware.Protected())
	{
		tasks.GET("", ctrl.getTasks)
		tasks.POST("", ctrl.createTask)
		tasks.PUT("/:id", ctrl.updateTask)
		tasks.DELETE("/:id", ctrl.deleteTask)
	}
}

const (
	defaultPerPage = 10
	maxPerPage     = 100
)

func (ctrl *TaskController) getTasks(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	opts := service.ListTasksOptions{
		UserID:  userID,
		Page:    parsePositiveInt(c.Query("page"), 1),
		PerPage: clampInt(parsePositiveInt(c.Query("per_page"), defaultPerPage), 1, maxPerPage),
		Status:  model.TaskStatus(c.Query("status")),
		Search:  c.Query("search"),
		Sort:    c.Query("sort"),
		Order:   c.Query("order"),
	}
	if opts.Status != "" && !opts.Status.Valid() {
		badRequest(c, "invalid status filter")
		return
	}

	result, err := ctrl.svc.ListTasks(opts)
	if err != nil {
		errorResponse(c, err)
		return
	}
	successResponse(c, result)
}

func parsePositiveInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func (ctrl *TaskController) createTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var body struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		badRequest(c, err.Error())
		return
	}

	task, err := ctrl.svc.CreateTask(userID, body.Title, body.Description)
	if err != nil {
		errorResponse(c, err)
		return
	}

	successResponse(c, task)
}

func (ctrl *TaskController) updateTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		badRequest(c, "invalid task id")
		return
	}

	// All optional — caller may patch any subset of {title, description, status}.
	var body struct {
		Title       *string           `json:"title"`
		Description *string           `json:"description"`
		Status      *model.TaskStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		badRequest(c, err.Error())
		return
	}
	if body.Title == nil && body.Description == nil && body.Status == nil {
		badRequest(c, "no fields to update")
		return
	}
	if body.Title != nil && *body.Title == "" {
		badRequest(c, "title cannot be empty")
		return
	}
	if body.Status != nil && !body.Status.Valid() {
		badRequest(c, "invalid status")
		return
	}

	task, err := ctrl.svc.UpdateTask(userID, uint(taskID), service.UpdateTaskInput{
		Title:       body.Title,
		Description: body.Description,
		Status:      body.Status,
	})
	if err != nil {
		errorResponse(c, err)
		return
	}

	successResponse(c, task)
}

func (ctrl *TaskController) deleteTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		badRequest(c, "invalid task id")
		return
	}

	if err := ctrl.svc.DeleteTask(userID, uint(taskID)); err != nil {
		errorResponse(c, err)
		return
	}

	successResponse(c, gin.H{"message": "Task deleted"})
}
