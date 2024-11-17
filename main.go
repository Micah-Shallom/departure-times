package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Micah-Shallom/departure-times/internal/config"
	"github.com/Micah-Shallom/departure-times/pkg/repository/storage/redis"
	"github.com/Micah-Shallom/departure-times/pkg/router"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	rds "github.com/redis/go-redis/v9"
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
	var (
		cacheLoadedKey = "cache:Loaded"
		ctx            = context.Background()
	)
	logger := utility.NewLogger()

	cfg := config.GetConfig(logger)

	//configurations
	config.LoadEnv(logger)
	validator := validator.New()

	//services
	redisClient := redis.ConnectToRedis(logger, cfg.Redis)

	isLoaded, err := redisClient.Get(ctx, cacheLoadedKey).Result()
	if err == rds.Nil || isLoaded == "" {
		log.Println("cache not loaded, starting the loading process........")
		go redis.LoadCache(logger, redisClient)
		redisClient.Set(ctx, cacheLoadedKey, "true", 24*time.Hour)
	} else if err != nil {
		log.Printf("Error Checking cache status: %v", err)
	} else {
		log.Println("Cache already loaded. Skipping cache loading. 🙂‍↔️🙂‍↔️🙂‍↔️🙂‍↔️🙂‍↔️")
	}

	app := NewRouter(cfg)

	r := router.Setup(app.engine, app.config, logger, validator)

	log.Fatal(r.Run(fmt.Sprintf(":%s", cfg.Server.Port)))
}
