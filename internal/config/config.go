package config

import (
	"log"
	"os"

	env "github.com/joho/godotenv"
)

type ConfigVars struct {
	Port string
	Env  string
}

// Must pattern in go: Must refers that it must load, other panic out from it
// it is just an convention, using Must is not mandatory, but it is a good practice to use it for functions that must succeed and panic if they fail. It is a common pattern in Go for functions that are expected to always succeed, such as loading configuration or initializing resources. The idea is that if the function fails, it indicates a serious problem that should not be ignored, and the program should terminate immediately.

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " environment variable is not set")
	}

	return value
}

func MustLoad() ConfigVars {
	log.Println("Loading environment variables...")

	if err := env.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	port := mustGetEnv("PORT")
	envVar := mustGetEnv("ENV")

	return ConfigVars{
		Port: port,
		Env:  envVar,
	}
}
