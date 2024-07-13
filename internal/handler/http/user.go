package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"hard/internal/domain/user"
	"hard/internal/service/tasker"
	"hard/pkg/server/response"
	"hard/pkg/store"
)

type UserHandler struct {
	taskerService *tasker.Service
}

func NewUserHandler(s *tasker.Service) *UserHandler {
	return &UserHandler{taskerService: s}
}

func (h *UserHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/users")
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

func (h *UserHandler) list(c *gin.Context) {
	res, err := h.taskerService.ListUsers(c)
	if err != nil {
		response.InternalServerError(c, err)
	}

	response.OK(c, res)
}

func (h *UserHandler) add(c *gin.Context) {
	req := user.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.taskerService.CreateUser(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *UserHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.taskerService.GetUser(c, id)
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

func (h *UserHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := user.Request{}

	if err := c.Bind(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	if req.Email == nil && req.Role == nil && req.FullName == nil {
		err := fmt.Errorf("bad request")
		response.BadRequest(c, err, req)
		return
	}

	if err := h.taskerService.UpdateUser(c, id, req); err != nil {
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

func (h *UserHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.taskerService.DeleteUser(c, id); err != nil {
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

func (h *UserHandler) search(c *gin.Context) {
	name := c.Query("name")
	email := c.Query("email")
	if name == "" && email == "" {
		response.BadRequest(c, errors.New("name or email query parameter required"), nil)
		return
	}

	res, err := h.taskerService.SearchUser(c, name, email)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *UserHandler) listTasks(c *gin.Context) {
	id := c.Param("id")

	res, err := h.taskerService.GetTasksByUser(c, id)
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
