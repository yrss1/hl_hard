package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"hard/internal/domain/task"
	"hard/internal/service/tasker"
	"hard/pkg/helpers"
	"hard/pkg/server/response"
	"hard/pkg/store"
	"strings"
)

type TaskHandler struct {
	taskerService *tasker.Service
}

func NewTaskHandler(s *tasker.Service) *TaskHandler {
	return &TaskHandler{taskerService: s}
}

// Routes sets up the routes for task management
func (h *TaskHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/tasks")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)

		api.GET("/:id", h.get)
		api.PUT("/:id", h.update)
		api.DELETE("/:id", h.delete)

		api.GET("/search", h.search)
	}
}

// listTasks godoc
//	@Summary		List all tasks
//	@Description	Get a list of all tasks
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		task.Response
//	@Failure		500	{object}	response.Object
//	@Router			/tasks [get]
func (h *TaskHandler) list(c *gin.Context) {
	res, err := h.taskerService.ListTasks(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	//err := errors.New("repository error")
	//response.InternalServerError(c, err)
	response.OK(c, res)
}

// addTask godoc
//	@Summary		Add a new task
//	@Description	Create a new task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body		task.Request	true	"Task Request"
//	@Success		201		{object}	task.Response
//	@Failure		400		{object}	response.Object
//	@Failure		500		{object}	response.Object
//	@Router			/tasks [post]
func (h *TaskHandler) add(c *gin.Context) {
	req := task.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.taskerService.CreateTask(c, req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "failed to parse:"):
			response.BadRequest(c, err, req)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.Created(c, res)
}

// getTask godoc
//	@Summary		Get task by ID
//	@Description	Get details of a specific task by ID
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{object}	task.Response
//	@Failure		404	{object}	response.Object
//	@Failure		500	{object}	response.Object
//	@Router			/tasks/{id} [get]
func (h *TaskHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.taskerService.GetTask(c, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, res)
}

// updateTask godoc
//	@Summary		Update a task
//	@Description	Update an existing task by ID
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Task ID"
//	@Param			task	body		task.Request	true	"Task Request"
//	@Success		200		{string}	string			"ok"
//	@Failure		400		{object}	response.Object
//	@Failure		404		{object}	response.Object
//	@Failure		500		{object}	response.Object
//	@Router			/tasks/{id} [put]
func (h *TaskHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := task.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if req.Title == nil && req.Description == nil && req.Priority == nil &&
		req.Status == nil && req.AssigneeID == nil && req.ProjectID == nil {
		err := fmt.Errorf("bad request")
		response.BadRequest(c, err, req)
		return
	}

	if err := h.taskerService.UpdateTask(c, id, req); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, "ok")
}

// deleteTask godoc
//	@Summary		Delete a task
//	@Description	Delete a task by ID
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{string}	string	"Deleted Task ID"
//	@Failure		404	{object}	response.Object
//	@Failure		500	{object}	response.Object
//	@Router			/tasks/{id} [delete]
func (h *TaskHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.taskerService.DeleteTask(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, id)
}

// searchTasks godoc
//	@Summary		Search tasks
//	@Description	Search tasks by title, priority, status, assignee_id, or project_id
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			title		query		string	false	"Task Title"
//	@Param			priority	query		string	false	"Task Priority"
//	@Param			status		query		string	false	"Task Status"
//	@Param			assignee_id	query		string	false	"Assignee ID"
//	@Param			project_id	query		string	false	"Project ID"
//	@Success		200			{array}		task.Response
//	@Failure		400			{object}	response.Object
//	@Failure		500			{object}	response.Object
//	@Router			/tasks/search [get]
func (h *TaskHandler) search(c *gin.Context) {
	data := task.Request{
		Title:      helpers.GetStringPtr(c.Query("title")),
		Priority:   helpers.GetStringPtr(c.Query("priority")),
		Status:     helpers.GetStringPtr(c.Query("status")),
		AssigneeID: helpers.GetStringPtr(c.Query("assignee_id")),
		ProjectID:  helpers.GetStringPtr(c.Query("project_id")),
	}
	if task.IsEmpty(data) {
		response.BadRequest(c, errors.New("query parameters required"), nil)
		return
	}

	res, err := h.taskerService.SearchTasks(c, data)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
