package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	middleware "github.com/GouthamGuna/file_sync/api/middleware"
	common "github.com/GouthamGuna/file_sync/pkg/common"
	logger "github.com/GouthamGuna/file_sync/pkg/logger"
	response "github.com/GouthamGuna/file_sync/pkg/response"
)

func FilePreviewHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling file preview request")
	if !isValidSignedURL(r) {
		response.WriteError(w, http.StatusUnauthorized, "URL is invalid or expired", "Unauthorized")
		return
	}

	// Link is valid; retrieve the file path from the URL.
	encodedFileName := r.URL.Query().Get("file")

	// Decode the Base64 encoded filename
	fileName, err := decodeFileName(encodedFileName)
	if err != nil {
		logger.Infof("Error decoding filename: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid file name encoding", "Bad Request")
		return
	}

	// Get the user's home directory
	var homeDir, err2 = common.GetBaseDirectory()
	if err2 != nil {
		response.WriteError(w, http.StatusInternalServerError, err2.Error(), "Internal Server Error")
		return
	}

	// Construct the full path to the file
	filePath := filepath.Join(homeDir, os.Getenv("FILE_SYNC_PATH"), fileName)
	logger.Infof("Serving file: %s", filePath)

	// Serve the file
	http.ServeFile(w, r, filePath)
	logger.Info("File preview request completed")
}

// isValidSignedURL validates the signature and expiration time of a signed URL.
func isValidSignedURL(r *http.Request) bool {
	// Extract query parameters.
	encodedFilePath := r.URL.Query().Get("file")
	expiresStr := r.URL.Query().Get("expires")
	signature := r.URL.Query().Get("signature")

	if encodedFilePath == "" || expiresStr == "" || signature == "" {
		return false
	}

	// Check for expiration time.
	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil || time.Now().Unix() > expires {
		return false
	}

	// Recreate the signature on the server side and compare.
	// Note: We use the encoded filename in the signature verification since that's what was signed
	dataToSign := fmt.Sprintf("%s:%d", encodedFilePath, expires)
	h := hmac.New(sha256.New, []byte(os.Getenv("SIGNING_KEY")))
	h.Write([]byte(dataToSign))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// Common file serving function used by both streaming and download handlers
func serveFile(w http.ResponseWriter, r *http.Request, disposition string) {

	common.GETRequestValidator(w, r)

	// Get the Base64 encoded filename
	encodedFileName := r.URL.Query().Get("file")
	if encodedFileName == "" {
		response.WriteError(w, http.StatusBadRequest, "Please specify a file name", "File parameter is required")
		return
	}

	// Decode the Base64 encoded filename
	fileName, err := decodeFileName(encodedFileName)
	if err != nil {
		logger.Infof("Error decoding filename: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid file name encoding", "Bad Request")
		return
	}

	// Validate file name to prevent directory traversal attacks
	if strings.Contains(fileName, "..") || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		response.WriteError(w, http.StatusBadRequest, "File name contains invalid characters", "Invalid file name")
		return
	}

	// Get the user's home directory
	homeDir, err := common.GetBaseDirectory()
	if err != nil {
		logger.Infof("Error getting user home directory: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "Could not determine user directory", "Internal server error")
		return
	}

	// Construct the full path to the file
	fileSyncPath := os.Getenv("FILE_SYNC_PATH")
	if fileSyncPath == "" {
		logger.Infof("FILE_SYNC_PATH environment variable is not set")
		response.WriteError(w, http.StatusInternalServerError, "File sync path not configured", "Configuration error")
		return
	}

	filePath := filepath.Join(homeDir, fileSyncPath, fileName)

	// Verify the file exists and is a regular file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			response.WriteError(w, http.StatusNotFound, "The requested file does not exist", "File not found")
			return
		}
		logger.Infof("Error accessing file %s: %v", filePath, err)
		response.WriteError(w, http.StatusInternalServerError, "Could not access the file", "Internal server error")
		return
	}

	// Check if it's a directory
	if fileInfo.IsDir() {
		response.WriteError(w, http.StatusBadRequest, "The requested path is a directory, not a file", "Invalid file")
		return
	}

	// Check if it's an MP4 file that needs chunked streaming
	if strings.ToLower(filepath.Ext(filePath)) == ".mp4" && disposition == "stream" {
		serveMP4Chunked(w, r, filePath, fileInfo)
		return
	}

	// Set appropriate headers based on disposition type
	if disposition == "download" {
		logger.Infof("Downloading file: %s (size: %d bytes)", filePath, fileInfo.Size())
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", sanitizeFileName(fileName)))
		w.Header().Set("Cache-Control", "no-store") // Prevent caching of sensitive files
	} else {
		logger.Infof("Serving file: %s (size: %d bytes)", filePath, fileInfo.Size())
		w.Header().Set("Content-Type", getContentType(filePath))
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	http.ServeFile(w, r, filePath)
}

func StreamingFile(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling file streaming request")
	serveFile(w, r, "stream")
	logger.Info("File streaming request completed")
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling file download request")
	serveFile(w, r, "download")
	logger.Info("File download request completed")
}

// getContentType determines the Content-Type based on file extension
func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	default:
		return "application/octet-stream"
	}
}

// sanitizeFileName removes potentially dangerous characters from file names
func sanitizeFileName(fileName string) string {
	// Remove path components and replace problematic characters
	cleanName := filepath.Base(fileName)
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	cleanName = strings.ReplaceAll(cleanName, "\"", "")
	cleanName = strings.ReplaceAll(cleanName, "'", "")
	cleanName = strings.ReplaceAll(cleanName, "`", "")
	return cleanName
}

// decodeFileName decodes a Base64 encoded filename back to its original form
func decodeFileName(encodedName string) (string, error) {
	return common.DecodeFileName(encodedName)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling file deletion request")
	role, ok := middleware.GetUserRoleFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "Unable to determine user role", "Unauthorized")
		return
	}

	if role == "admin" {
		deleteRequestValidator(w, r)
	} else {
		response.WriteError(w, http.StatusForbidden, "You do not have permission to access this resource", "Forbidden")
	}
	logger.Info("File deletion request completed")
}

func deleteRequestValidator(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method
	if r.Method != http.MethodDelete {
		response.WriteError(w, http.StatusMethodNotAllowed, "Only DELETE method is allowed", "Method not allowed")
		return
	}

	// Get the Base64 encoded filename
	encodedFileName := r.URL.Query().Get("file")
	if encodedFileName == "" {
		response.WriteError(w, http.StatusBadRequest, "File name is required", "Invalid file name")
		return
	}

	// Decode the Base64 encoded filename
	fileName, err := decodeFileName(encodedFileName)
	if err != nil {
		logger.Infof("Error decoding filename: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid file name encoding", "Bad Request")
		return
	}

	// Validate file name to prevent directory traversal attacks
	if strings.Contains(fileName, "..") || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		response.WriteError(w, http.StatusBadRequest, "File name contains invalid characters", "Invalid file name")
		return
	}

	// Get the user's home directory
	homeDir, err := common.GetBaseDirectory()
	if err != nil {
		logger.Infof("Error getting user home directory: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "Could not determine user directory", "Internal server error")
		return
	}

	// Construct the full path to the file
	fileSyncPath := os.Getenv("FILE_SYNC_PATH")
	if fileSyncPath == "" {
		logger.Infof("FILE_SYNC_PATH environment variable is not set")
		response.WriteError(w, http.StatusInternalServerError, "File sync path not configured", "Configuration error")
		return
	}

	filePath := filepath.Join(homeDir, fileSyncPath, fileName)

	// Delete the file
	err = os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			response.WriteError(w, http.StatusNotFound, "The requested file does not exist", "File not found")
			return
		}
		logger.Infof("Error deleting file %s: %v", filePath, err)
		response.WriteError(w, http.StatusInternalServerError, "Could not delete the file", "Internal server error")
		return
	}

	// Respond with success
	//w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "File deleted successfully",
	})
}

func serveMP4Chunked(w http.ResponseWriter, r *http.Request, filePath string, fileInfo os.FileInfo) {
	const chunkSize = 1024 * 1024 // 1MB chunks

	file, err := os.Open(filePath)
	if err != nil {
		logger.Infof("Error opening file: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "Could not open file", "Internal server error")
		return
	}
	defer file.Close()

	// Handle range requests for seeking
	rangeHeader := r.Header.Get("Range")
	var start, end int64
	if rangeHeader != "" {
		if strings.HasPrefix(rangeHeader, "bytes=") {
			rangeHeader = strings.TrimPrefix(rangeHeader, "bytes=")
			parts := strings.Split(rangeHeader, "-")
			if len(parts) == 2 {
				start, _ = strconv.ParseInt(parts[0], 10, 64)
				if parts[1] != "" {
					end, _ = strconv.ParseInt(parts[1], 10, 64)
				} else {
					end = fileInfo.Size() - 1
				}
			}
		}
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileInfo.Size()))
		w.WriteHeader(http.StatusPartialContent)
	} else {
		end = fileInfo.Size() - 1
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Transfer-Encoding", "chunked")

	buffer := make([]byte, chunkSize)
	file.Seek(start, 0)
	bytesRemaining := end - start + 1

	for bytesRemaining > 0 {
		n := chunkSize
		if bytesRemaining < int64(chunkSize) {
			n = int(bytesRemaining)
		}

		bytesRead, err := file.Read(buffer[:n])
		if err != nil {
			logger.Infof("Error reading file chunk: %v", err)
			return
		}

		_, err = w.Write(buffer[:bytesRead])
		if err != nil {
			logger.Infof("Error writing chunk to response: %v", err)
			return
		}

		bytesRemaining -= int64(bytesRead)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}
