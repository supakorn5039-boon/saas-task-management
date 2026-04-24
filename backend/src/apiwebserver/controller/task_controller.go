package controller

import (
	"net/http"
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
		tasks.PUT("/:id", ctrl.updateStatus)
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
		errorResponse(c, "invalid status filter", http.StatusBadRequest)
		return
	}

	result, err := ctrl.svc.ListTasks(opts)
	if err != nil {
		errorResponse(c, err.Error(), http.StatusInternalServerError)
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
		errorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := ctrl.svc.CreateTask(userID, body.Title, body.Description)
	if err != nil {
		errorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	successResponse(c, task)
}

func (ctrl *TaskController) updateStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorResponse(c, "invalid task id", http.StatusBadRequest)
		return
	}

	var body struct {
		Status model.TaskStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	if !body.Status.Valid() {
		errorResponse(c, "invalid status", http.StatusBadRequest)
		return
	}

	task, err := ctrl.svc.UpdateStatus(userID, uint(taskID), body.Status)
	if err != nil {
		errorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	successResponse(c, task)
}

func (ctrl *TaskController) deleteTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorResponse(c, "invalid task id", http.StatusBadRequest)
		return
	}

	if err := ctrl.svc.DeleteTask(userID, uint(taskID)); err != nil {
		errorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	successResponse(c, gin.H{"message": "Task deleted"})
}
