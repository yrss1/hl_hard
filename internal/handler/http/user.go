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

// Routes sets up the routes for user management
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

// listUsers godoc
//	@Summary		List all users
//	@Description	Get a list of all users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		user.Response
//	@Failure		500	{object}	response.Object
//	@Router			/users [get]
func (h *UserHandler) list(c *gin.Context) {
	res, err := h.taskerService.ListUsers(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// addUser godoc
//	@Summary		Add a new user
//	@Description	Create a new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		user.Request	true	"User Request"
//	@Success		200		{object}	user.Response
//	@Failure		400		{object}	response.Object
//	@Failure		500		{object}	response.Object
//	@Router			/users [post]
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

// getUser godoc
//	@Summary		Get user by ID
//	@Description	Get details of a specific user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	user.Response
//	@Failure		404	{object}	response.Object
//	@Failure		500	{object}	response.Object
//	@Router			/users/{id} [get]
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

// updateUser godoc
//	@Summary		Update a user
//	@Description	Update an existing user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"User ID"
//	@Param			user	body		user.Request	true	"User Request"
//	@Success		200		{string}	string			"ok"
//	@Failure		400		{object}	response.Object
//	@Failure		404		{object}	response.Object
//	@Failure		500		{object}	response.Object
//	@Router			/users/{id} [put]
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

// deleteUser godoc
//	@Summary		Delete a user
//	@Description	Delete a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{string}	string	"Deleted User ID"
//	@Failure		404	{object}	response.Object
//	@Failure		500	{object}	response.Object
//	@Router			/users/{id} [delete]
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

// searchUsers godoc
//	@Summary		Search users
//	@Description	Search users by name or email
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"User Name"
//	@Param			email	query		string	false	"User Email"
//	@Success		200		{array}		user.Response
//	@Failure		400		{object}	response.Object
//	@Failure		500		{object}	response.Object
//	@Router			/users/search [get]
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

// listUserTasks godoc
//	@Summary		List tasks for a user
//	@Description	Get a list of tasks for a specific user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{array}		task.Response
//	@Failure		404	{object}	response.Object
//	@Failure		500	{object}	response.Object
//	@Router			/users/{id}/tasks [get]
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
