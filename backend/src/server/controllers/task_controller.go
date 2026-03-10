package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/server/services"
	"github.com/supakorn5039-boon/saas-task-backend/src/utils"
)

func GetTasks(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	taskService := services.NewTaskService()

	tasks, err := taskService.GetAllTasks(userID)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, tasks)
}

func CreateTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var body struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	taskService := services.NewTaskService()
	task, err := taskService.CreateTask(userID, body.Title, body.Description)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, task)
}

func ToggleTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, "invalid task id", http.StatusBadRequest)
		return
	}

	var body struct {
		Completed bool `json:"completed"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	taskService := services.NewTaskService()
	task, err := taskService.ToggleTask(userID, uint(taskID), body.Completed)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, task)
}

func DeleteTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, "invalid task id", http.StatusBadRequest)
		return
	}

	taskService := services.NewTaskService()
	if err := taskService.DeleteTask(userID, uint(taskID)); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Task deleted"})
}
