package handler

import (
	"github.com/gin-gonic/gin"
	"hard/internal/config"
	"hard/internal/handler/http"
	"hard/internal/service/tasker"
	"hard/pkg/server/router"
)

type Dependencies struct {
	Configs       config.Configs
	TaskerService *tasker.Service
}
type Handler struct {
	dependencies Dependencies
	HTTP         *gin.Engine
}
type Configuration func(h *Handler) error

func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	h = &Handler{
		dependencies: d,
		HTTP:         router.New(),
	}

	for _, cfg := range configs {
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		h.HTTP = router.New()
		userHandler := http.NewUserHandler(h.dependencies.TaskerService)
		taskHandler := http.NewTaskHandler(h.dependencies.TaskerService)
		projectHandler := http.NewProjectHandler(h.dependencies.TaskerService)
		api := h.HTTP.Group("/api/v1/")
		{
			userHandler.Routes(api)
			taskHandler.Routes(api)
			projectHandler.Routes(api)
		}
		return
	}
}
