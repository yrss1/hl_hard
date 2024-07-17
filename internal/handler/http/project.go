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

// Routes sets up the routes for project management
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

// listProjects godoc
//
//	@Summary		List all projects
//	@Description	Get a list of all projects
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		project.Response
//	@Failure		500	{object}	response.Object
//	@Router			/projects [get]
func (h *ProjectHandler) list(c *gin.Context) {
	res, err := h.taskerService.ListProjects(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// addProject godoc
//
//	@Summary		Add a new project
//	@Description	Create a new project
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			project	body		project.Request	true	"Project Request"
//	@Success		200		{object}	project.Response
//	@Failure		400		{object}	response.Object
//	@Failure		500		{object}	response.Object
//	@Router			/projects [post]
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

// getProject godoc
//
//	@Summary		Get project by ID
//	@Description	Get details of a specific project by ID
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Success		200	{object}	project.Response
//	@Failure		404	{object}	response.Object
//	@Failure		500	{object}	response.Object
//	@Router			/projects/{id} [get]
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

// updateProject godoc
//
//	@Summary		Update a project
//	@Description	Update an existing project by ID
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Project ID"
//	@Param			project	body		project.Request	true	"Project Request"
//	@Success		200		{string}	string			"ok"
//	@Failure		400		{object}	response.Object
//	@Failure		404		{object}	response.Object
//	@Failure		500		{object}	response.Object
//	@Router			/projects/{id} [put]
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

// deleteProject godoc
//
//	@Summary		Delete a project
//	@Description	Delete a project by ID
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Success		200	{string}	string	"Deleted Project ID"
//	@Failure		404	{object}	response.Object
//	@Failure		500	{object}	response.Object
//	@Router			/projects/{id} [delete]
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

// searchProjects godoc
//
//	@Summary		Search projects
//	@Description	Search projects by title or manager_id
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			title		query		string	false	"Project Title"
//	@Param			manager_id	query		string	false	"Manager ID"
//	@Success		200			{array}		project.Response
//	@Failure		400			{object}	response.Object
//	@Failure		500			{object}	response.Object
//	@Router			/projects/search [get]
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

// listTasks godoc
//
//	@Summary		List tasks by project
//	@Description	Get a list of all tasks for a specific project
//	@Tags			projects
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Project ID"
//	@Success		200	{array}		task.Response
//	@Failure		404	{object}	response.Object
//	@Failure		500	{object}	response.Object
//	@Router			/projects/{id}/tasks [get]
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
