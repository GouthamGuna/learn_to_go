package pkg

import (
	"net/http"
	"os"

	logger "github.com/GouthamGuna/file_sync/pkg/logger"
	response "github.com/GouthamGuna/file_sync/pkg/response"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Initiating 404 handler")
	response.WriteError(w, http.StatusNotFound, "404 - Not Found", "Not Found")
	logger.Infof("End Requested URL: %s", r.URL.Path)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Initiating 405 handler")
	response.WriteError(w, http.StatusMethodNotAllowed, "405 - Method Not Allowed", "Method Not Allowed")
	logger.Infof("End Requested URL: %s", r.URL.Path)
}

func GETRequestValidator(w http.ResponseWriter, r *http.Request) {
	logger.Info("Validating GET request")
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "Only GET requests are allowed", "Method not allowed")
		return
	}
	logger.Info("GET request is valid")
}

func GetBaseDirectory() (string, error) {
	return os.UserHomeDir()
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Serving favicon")
	http.ServeFile(w, r, "web/static/favicon.ico")
}
