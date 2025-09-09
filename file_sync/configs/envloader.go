package configs

import (
	"github.com/joho/godotenv"
	"os"

	logger "github.com/GouthamGuna/file_sync/pkg/logger"
)

func LoadEnv() {
	logger.Info("Loading environment variables")
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file")
	}
	info()
	logger.Info("Environment variables loaded successfully")
}

func info() {
	logger.Infof("Current environment: %s", os.Getenv("ENV"))
	logger.Infof("File sync path: %s", os.Getenv("FILE_SYNC_PATH"))
}
