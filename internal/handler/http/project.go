package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"hard/internal/domain/project"
	"hard/internal/service/tasker"
	"hard/pkg/helpers"
	"hard/pkg/server/response"
	"hard/pkg/store"
	"strings"
)

type ProjectHandler struct {
	taskerService *tasker.Service
}

func NewProjectHandler(s *tasker.Service) *ProjectHandler {
	return &ProjectHandler{taskerService: s}
}

func (h *ProjectHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/projects")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)

		api.GET("/:id", h.get)
		api.GET("/:id/tasks", h.listTasks)
		api.PUT("/:id", h.update)
		api.DELETE("/:id", h.delete)

		api.GET("/search", h.search)

	}
}

func (h *ProjectHandler) list(c *gin.Context) {
	res, err := h.taskerService.ListProjects(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
func (h *ProjectHandler) add(c *gin.Context) {
	req := project.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.taskerService.CreateProject(c, req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "failed to parse:"):
			response.BadRequest(c, err, req)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, res)
}

func (h *ProjectHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.taskerService.GetProject(c, id)
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

func (h *ProjectHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := project.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if req.IsEmpty() {
		err := fmt.Errorf("bad request")
		response.BadRequest(c, err, req)
		return
	}

	if err := h.taskerService.UpdateProject(c, id, req); err != nil {
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

func (h *ProjectHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.taskerService.DeleteProject(c, id); err != nil {
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

func (h *ProjectHandler) search(c *gin.Context) {
	req := project.Request{
		Title:     helpers.GetStringPtr(c.Query("title")),
		ManagerID: helpers.GetStringPtr(c.Query("manager_id")),
	}
	if req.IsEmpty() {
		response.BadRequest(c, errors.New("query parameters required"), nil)
		return
	}

	res, err := h.taskerService.SearchProjects(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ProjectHandler) listTasks(c *gin.Context) {
	id := c.Param("id")

	res, err := h.taskerService.GetTasksByProject(c, id)
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
