package redis

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Micah-Shallom/departure-times/internal/config"
	"github.com/Micah-Shallom/departure-times/pkg/repository/storage"
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
)

func ConnectToRedis(logger *utility.Logger, cfg config.Redis) *redis.Client {
	utility.LogAndPrint(logger, "connecting to redis server")
	connectedServer := connectToDb(cfg.Host, cfg.Port, cfg.DB, logger)
	fmt.Println("-------------------------------------", connectedServer)

	utility.LogAndPrint(logger, "connected to redis server")

	storage.DB.Redis = connectedServer

	return connectedServer
}

func connectToDb(host, port, db string, logger *utility.Logger) *redis.Client {
	if _, err := strconv.Atoi(port); err != nil {
		u, err := url.Parse(port)
		if err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("parsing url %v to get port failed with: %v", port, err))
			panic(err)
		}

		detectedPort := u.Port()
		if detectedPort == "" {
			utility.LogAndPrint(logger, fmt.Sprintf("detecting port from url %v failed with: %v", port, err))
			panic(err)
		}
		port = detectedPort
	}
	dbInst, err := strconv.Atoi(db)
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("parsing url %v to get port failed with: %v", port, err))
		panic(err)
	}

	addr := fmt.Sprintf("%v:%v", host, port)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       dbInst,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		utility.LogAndPrint(logger, fmt.Sprintln(addr))
		utility.LogAndPrint(logger, fmt.Sprintln("Redis db error: ", err))
	}

	pong, _ := redisClient.Ping(ctx).Result()
	utility.LogAndPrint(logger, fmt.Sprintln("Redis says: ", pong))

	return redisClient
}
