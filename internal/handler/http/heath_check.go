package http

import (
	"github.com/gin-gonic/gin"
	"hard/pkg/server/response"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/health")
	{
		api.GET("/", h.ok)
	}
}

// ok godoc
//
//	@Summary		HealthCheck
//	@Description	HealthСheck
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}	string			"ok"
//	@Router			/health [get]
func (h *HealthHandler) ok(c *gin.Context) {
	text := "Im okay, dont worry Morty"
	response.OK(c, text)
}
