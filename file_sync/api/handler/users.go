package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	middleware "github.com/GouthamGuna/file_sync/api/middleware"
	dto "github.com/GouthamGuna/file_sync/pkg/dto"
	response "github.com/GouthamGuna/file_sync/pkg/response"
	logger "github.com/GouthamGuna/file_sync/pkg/logger"
)

func GenerateUserToken(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling user token generation request")
	var credentials dto.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	username := credentials.Username
	password := credentials.Password

	if isValidUser(username, password) {
		token, err := middleware.CreateToken(username)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, "Failed to generate token", err.Error())
			return
		}

		response.WriteAccessTokenResponse(w, http.StatusOK, token)
	} else {
		response.WriteError(w, http.StatusUnauthorized, "Invalid credentials", "Unauthorized")
		return
	}
	logger.Info("User token generation request completed")
}

func isValidUser(username, password string) bool {
	var user_credentials = os.Getenv("USERS_CREDENTIALS")
	credentials := strings.Split(user_credentials, ",")

	for _, cred := range credentials {
		parts := strings.Split(cred, ":")
		if len(parts) == 2 && parts[0] == username && parts[1] == password {
			return true
		}
	}
	return false
}
