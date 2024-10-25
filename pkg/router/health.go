package router

import (
	"fmt"

	"github.com/Micah-Shallom/departure-times/internal/config"
	"github.com/Micah-Shallom/departure-times/pkg/controller/health"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/gin-gonic/gin"
)

func Health(r *gin.Engine, ApiVersion string, logger *utility.Logger, cfg *config.Config) {
	healthURL := r.Group(fmt.Sprintf("%s", ApiVersion))
	health := health.Controller{Logger: logger}

	healthURL.POST("/health", health.POST)
	healthURL.GET("/health", health.GET)
}
