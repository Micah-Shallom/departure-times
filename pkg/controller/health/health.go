package health

import (
	"github.com/Micah-Shallom/departure-times/internal/config"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Logger *utility.Logger
	Config *config.Config
}

func (base *Controller) GET(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (base *Controller) POST(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
