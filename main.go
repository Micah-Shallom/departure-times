package main

import (
	"fmt"
	"log"

	"github.com/Micah-Shallom/departure-times/internal/config"
	"github.com/Micah-Shallom/departure-times/pkg/repository/storage/redis"
	"github.com/Micah-Shallom/departure-times/pkg/router"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/gin-gonic/gin"
)

type Router struct {
	config *config.Config
	engine *gin.Engine
}

func NewRouter(cfg *config.Config) *Router {
	gin.SetMode(cfg.Server.Mode)
	return &Router{
		config: cfg,
		engine: gin.New(),
	}
}

func main() {
	logger := utility.NewLogger()

	cfg := config.GetConfig(logger)

	//loading configurations
	config.LoadEnv(logger)

	//load services
	redis.ConnectToRedis(logger, cfg.Redis)

	//instantiate application
	app := NewRouter(cfg)

	r := router.Setup(app.engine, app.config ,logger)
	fmt.Println(cfg.Server.Port)

	log.Fatal(r.Run(fmt.Sprintf(":%s", cfg.Server.Port)))
}
