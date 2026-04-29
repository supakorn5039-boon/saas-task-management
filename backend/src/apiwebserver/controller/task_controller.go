package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
)

type TaskController struct {
	svc   *service.TaskService
	audit *service.AuditService
}

func NewTaskController() *TaskController {
	return &TaskController{
		svc:   service.NewTaskService(),
		audit: service.NewAuditService(),
	}
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
		UserID:   userID,
		Page:     parsePositiveInt(c.Query("page"), 1),
		PerPage:  clampInt(parsePositiveInt(c.Query("per_page"), defaultPerPage), 1, maxPerPage),
		Status:   model.TaskStatus(c.Query("status")),
		Priority: model.TaskPriority(c.Query("priority")),
		Search:   c.Query("search"),
		Sort:     c.Query("sort"),
		Order:    c.Query("order"),
	}
	if opts.Status != "" && !opts.Status.Valid() {
		badRequest(c, "invalid status filter")
		return
	}
	if opts.Priority != "" && !opts.Priority.Valid() {
		badRequest(c, "invalid priority filter")
		return
	}
	// "?assignee=me" is a tiny convenience — saves the frontend from injecting
	// its own user id and lets us hardcode the URL in nav links.
	if a := c.Query("assignee"); a == "me" {
		uid := userID
		opts.AssigneeID = &uid
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

// parseISO is a small helper that turns an optional date-time string into
// *time.Time. Empty / nil input → nil. Malformed input is treated as a 400
// at the controller layer.
func parseISO(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (ctrl *TaskController) createTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var body struct {
		Title       string             `json:"title" binding:"required"`
		Description string             `json:"description"`
		Priority    model.TaskPriority `json:"priority"`
		StartDate   string             `json:"startDate"`
		DueDate     string             `json:"dueDate"`
		AssigneeID  *uint              `json:"assigneeId"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		badRequest(c, err.Error())
		return
	}

	start, err := parseISO(body.StartDate)
	if err != nil {
		badRequest(c, "invalid startDate")
		return
	}
	due, err := parseISO(body.DueDate)
	if err != nil {
		badRequest(c, "invalid dueDate")
		return
	}

	task, err := ctrl.svc.CreateTask(userID, service.CreateTaskInput{
		Title:       body.Title,
		Description: body.Description,
		Priority:    body.Priority,
		StartDate:   start,
		DueDate:     due,
		AssigneeID:  body.AssigneeID,
	})
	if err != nil {
		errorResponse(c, err)
		return
	}

	ctrl.audit.Record(c, model.AuditActionTaskCreated, model.AuditStatusSuccess, service.RecordOpts{
		TargetType: "task",
		TargetID:   &task.ID,
		Metadata:   model.JSONB{"title": task.Title, "priority": string(task.Priority)},
	})
	successResponse(c, task)
}

// updateBody uses *string for "string fields that may be cleared". For the
// date and assignee fields we accept the JSON value and decide later — the
// service uses double-pointers so we can distinguish missing vs explicit null.
type updateTaskBody struct {
	Title       *string             `json:"title"`
	Description *string             `json:"description"`
	Status      *model.TaskStatus   `json:"status"`
	Priority    *model.TaskPriority `json:"priority"`
	// json.Unmarshal into a pointer-to-pointer so we can tell missing
	// (outer nil) from explicit null (inner nil). The frontend uses this to
	// clear a due date or unassign a task.
	StartDate  *string `json:"startDate"`
	DueDate    *string `json:"dueDate"`
	AssigneeID *uint   `json:"assigneeId"`
	// Set to true to clear the corresponding field — distinct from "not set".
	ClearStartDate bool `json:"clearStartDate"`
	ClearDueDate   bool `json:"clearDueDate"`
	ClearAssignee  bool `json:"clearAssignee"`
}

func (ctrl *TaskController) updateTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		badRequest(c, "invalid task id")
		return
	}

	var body updateTaskBody
	if err := c.ShouldBindJSON(&body); err != nil {
		badRequest(c, err.Error())
		return
	}

	in, err := buildUpdateTaskInput(body)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	task, err := ctrl.svc.UpdateTask(userID, uint(taskID), in)
	if err != nil {
		errorResponse(c, err)
		return
	}

	meta := model.JSONB{"title": task.Title}
	if body.Status != nil {
		meta["status"] = string(*body.Status)
	}
	if body.Priority != nil {
		meta["priority"] = string(*body.Priority)
	}
	ctrl.audit.Record(c, model.AuditActionTaskUpdated, model.AuditStatusSuccess, service.RecordOpts{
		TargetType: "task",
		TargetID:   &task.ID,
		Metadata:   meta,
	})
	successResponse(c, task)
}

// buildUpdateTaskInput translates the controller-side body shape into the
// service-side UpdateTaskInput, validating along the way. Pulled out so the
// handler stays focused on HTTP concerns.
func buildUpdateTaskInput(body updateTaskBody) (service.UpdateTaskInput, error) {
	if body.Title != nil && *body.Title == "" {
		return service.UpdateTaskInput{}, errEmptyTitle
	}
	if body.Status != nil && !body.Status.Valid() {
		return service.UpdateTaskInput{}, errInvalidStatus
	}
	if body.Priority != nil && !body.Priority.Valid() {
		return service.UpdateTaskInput{}, errInvalidPriority
	}

	in := service.UpdateTaskInput{
		Title:       body.Title,
		Description: body.Description,
		Status:      body.Status,
		Priority:    body.Priority,
	}

	// Date fields — set or clear. "Set" comes in as a non-empty string;
	// "clear" comes in as the explicit clear* flag. Anything else means "no
	// change" so the service leaves the column alone.
	if body.ClearStartDate {
		var nilTime *time.Time
		in.StartDate = &nilTime
	} else if body.StartDate != nil {
		t, err := parseISO(*body.StartDate)
		if err != nil {
			return service.UpdateTaskInput{}, errInvalidStartDate
		}
		in.StartDate = &t
	}
	if body.ClearDueDate {
		var nilTime *time.Time
		in.DueDate = &nilTime
	} else if body.DueDate != nil {
		t, err := parseISO(*body.DueDate)
		if err != nil {
			return service.UpdateTaskInput{}, errInvalidDueDate
		}
		in.DueDate = &t
	}
	if body.ClearAssignee {
		var nilID *uint
		in.AssigneeID = &nilID
	} else if body.AssigneeID != nil {
		uid := body.AssigneeID
		in.AssigneeID = &uid
	}

	if in.Title == nil && in.Description == nil && in.Status == nil &&
		in.Priority == nil && in.StartDate == nil && in.DueDate == nil &&
		in.AssigneeID == nil {
		return service.UpdateTaskInput{}, errNoFieldsToUpdate
	}
	return in, nil
}

// Sentinel errors used by buildUpdateTaskInput so the controller can render
// a stable 400 message without coupling tests to literal English text. They
// implement error so they can flow through err returns naturally.
var (
	errEmptyTitle       = stringError("title cannot be empty")
	errInvalidStatus    = stringError("invalid status")
	errInvalidPriority  = stringError("invalid priority")
	errInvalidStartDate = stringError("invalid startDate")
	errInvalidDueDate   = stringError("invalid dueDate")
	errNoFieldsToUpdate = stringError("no fields to update")
)

type stringError string

func (s stringError) Error() string { return string(s) }

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

	tid := uint(taskID)
	ctrl.audit.Record(c, model.AuditActionTaskDeleted, model.AuditStatusSuccess, service.RecordOpts{
		TargetType: "task",
		TargetID:   &tid,
	})
	successResponse(c, gin.H{"message": "Task deleted"})
}
