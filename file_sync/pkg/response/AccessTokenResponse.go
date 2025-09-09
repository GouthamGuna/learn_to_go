package pkg

import (
	"encoding/json"
	"net/http"
	"time"
)

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	GenerateTime string `json:"generate_time"`
}

func WriteAccessTokenResponse(w http.ResponseWriter, statusCode int, accessToken string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(AccessTokenResponse{
		AccessToken:  accessToken,
		GenerateTime: time.Now().Format(time.RFC3339),
	})
}
