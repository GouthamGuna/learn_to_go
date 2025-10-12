package pkg

import (
	"encoding/base64"
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

// EncodeFileName encodes a filename to Base64 for use in URLs
func EncodeFileName(fileName string) string {
	return base64.URLEncoding.EncodeToString([]byte(fileName))
}

// DecodeFileName decodes a Base64-encoded filename back to its original form
func DecodeFileName(encodedName string) (string, error) {
	bytes, err := base64.URLEncoding.DecodeString(encodedName)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
