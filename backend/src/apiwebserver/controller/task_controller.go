package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
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
		tasks.PATCH("/:id", ctrl.toggleTask)
		tasks.DELETE("/:id", ctrl.deleteTask)
	}
}

func (ctrl *TaskController) getTasks(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	tasks, err := ctrl.svc.GetAllTasks(userID)
	if err != nil {
		errorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	successResponse(c, tasks)
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

func (ctrl *TaskController) toggleTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorResponse(c, "invalid task id", http.StatusBadRequest)
		return
	}

	var body struct {
		Completed bool `json:"completed"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := ctrl.svc.ToggleTask(userID, uint(taskID), body.Completed)
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
