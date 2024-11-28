package router

import (
	"net/http"

	"github.com/Micah-Shallom/departure-times/internal/config"
	"github.com/Micah-Shallom/departure-times/pkg/middleware"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Setup(r *gin.Engine, cfg *config.Config, logger *utility.Logger, validator *validator.Validate) *gin.Engine {
	ApiVersion := "/api/v1"

	//middlewares
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"192.168.0.1", "192.168.0.2"})
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORS())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.MaxMultipartMemory = 1 << 20 // 1 MiB

	//routers
	Health(r, ApiVersion, logger, cfg)
	NextBus(r, ApiVersion, logger, validator ,cfg)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":    "Not Found",
			"message": "The requested URL was not found on this server.",
			"code":    http.StatusNotFound,
			"status":  "error",
		})
	})

	return r
}
