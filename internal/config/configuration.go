package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/joho/godotenv"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	App    Application
	Server ServerMode
	Redis  Redis
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       string
}

type ServerMode struct {
	Port string
	Mode string
}

type Application struct {
	NextBusURL string
}

func LoadEnv(logger *utility.Logger) *Config {
	once.Do(func() {
		// Load .env file
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error loading .env file")
			logger.Error("Warning: .env file not found: %v", err)
			return
		}

		config = &Config{
			App: Application{
				NextBusURL: getEnv("NextBus_URL", "https://retro.umoiq.com/service/publicXMLFeed"),
			},
			Server: ServerMode{
				Port: getEnv("PORT", "8080"),
				Mode: getEnv("GIN_MODE", "debug"),
			},
			Redis: Redis{
				Host:     getEnv("REDIS_HOST", "localhost"),
				Port:     getEnv("REDIS_PORT", "6379"),
				Password: getEnv("REDIS_PASSWORD", ""),
				DB:       "0",
			},
		}
	})
	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetConfig(logger *utility.Logger) *Config {
	if config == nil {
		return LoadEnv(logger)
	}
	return config
}
