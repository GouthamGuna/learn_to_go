package pkg

const (
	APIDataPath        = "/api/v1/data?file="
	APIDownloadPath    = "/api/v1/download?file="
	APIFilePreviewPath = "/api/v1/file_preview"
	APIDeletePath      = "/api/v1/delete?file="
	APIBasePath_V1     = "/api/v1"

	// Error messages
	ErrInvalidRequest      = "Invalid request"
	ErrFileNotFound        = "File not found"
	ErrMethodNotAllowed    = "Method not allowed"
	ErrInternalServerError = "Internal server error"
)
