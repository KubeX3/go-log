package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds environment variables
type Config struct {
	Environment string
	LogEnabled  bool
	LogFilePath string
}

var DOTENV Config

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	environment := getEnv("Environment", "development")

	now := time.Now()
	fmt.Printf(
		"\x1b[90m[%d/%02d/%02d - %02d:%02d:%02d] \x1b[32m[  INFO   ] \x1b[37m- \x1b[90m[config.go] \x1b[37m- \x1b[37m%s\n",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		"Dotenv inject",
	)

	// Initialize the Global Config
	DOTENV = Config{
		Environment:     environment,
		LogEnabled:  getEnvBool("LOG_ENABLED", true),
		LogFilePath: getEnv("LOG_FILE_PATH", "./logs/system.log"),
	}
}

// Helper to get string env with default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Helper to get boolean env with default
func getEnvBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val == "true"
}
