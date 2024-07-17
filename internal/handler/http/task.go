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
