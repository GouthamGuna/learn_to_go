package pkg

import (
	"encoding/json"
	"net/http"

	dto "github.com/GouthamGuna/file_sync/pkg/dto"
)

type link struct {
	Method      string `json:"method,omitempty"`
	Description string `json:"description,omitempty"`
	Preview     any    `json:"preview"`
	Download    any    `json:"download"`
	URL         any    `json:"url"`
}

type returnResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Links   link   `json:"links"`
}

func WriteSuccess(w http.ResponseWriter, status int, message string, data dto.DateURL) {
	w.WriteHeader(status)
	response := returnResponse{
		Status:  status,
		Message: message,
		Links: link{
			Method:      "GET",
			Description: "This preview URL is temporary and will expire soon.",
			Preview:     data.PreviewURL,
			Download:    data.DownloadURL,
			URL:         data.URL,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func FileUploadedResponse(w http.ResponseWriter, message string, encryptedFilename string) {
	WriteSuccess(w, http.StatusOK, message, dto.DateURL{
		PreviewURL:  nil,
		DownloadURL: nil,
		URL:        encryptedFilename,
	})
}
