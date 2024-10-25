package config

import (
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

// GetConfig returns the current config instance
func GetConfig(logger *utility.Logger) *Config {
	if config == nil {
		return LoadEnv(logger)
	}
	return config
}
