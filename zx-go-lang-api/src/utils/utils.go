package utils

import (
	"github.com/joho/godotenv"
	"os"
)

// LoadEnv loads environment variables from a .env file and returns an error if it fails.
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

// Initialize the environment variables once at package level.
var envLoaded = false

// Init initializes the environment variables if not already loaded.
func Init() error {
	if !envLoaded {
		if err := LoadEnv(); err != nil {
			return err
		}
		envLoaded = true
	}
	return nil
}

// GetServerPort retrieves the server port from the environment variables.
func GetServerPort() string {
	if err := Init(); err != nil {
		return ":8088"
	}
	return os.Getenv("SERVER_PORT")
}

// GetDBURL retrieves the database URL from the environment variables.
func GetDBURL() (string, error) {
	if err := Init(); err != nil {
		return "mongodb://localhost:27017/", err
	}
	return os.Getenv("DB_URL"), nil
}

func GETFavicon() string {
	if err := Init(); err != nil {
		return ""
	}
	return os.Getenv("FAVICON_PATH")
}

func GetRootPath() string {
	if err := Init(); err != nil {
		return ""
	}
	return os.Getenv("ROOT_PATH")
}

func GetJWTEncryptionKey() string {
	if err := Init(); err != nil {
        return ""
    }
    return os.Getenv("JWT_ENCRYPTION_KEY")
}
