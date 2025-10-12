package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	middleware "github.com/GouthamGuna/file_sync/api/middleware"
	common "github.com/GouthamGuna/file_sync/pkg/common"
	dto "github.com/GouthamGuna/file_sync/pkg/dto"
	logger "github.com/GouthamGuna/file_sync/pkg/logger"
	response "github.com/GouthamGuna/file_sync/pkg/response"
)

func Uploads(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling file uploads")
	if role, ok := middleware.GetUserRoleFromContext(r.Context()); ok {
		if !ok {
			response.WriteError(w, http.StatusInternalServerError, "Unable to determine user role", "Unauthorized")
			return
		}

		if role == "admin" {
			if err := handleFileUpload(w, r); err != nil {
				response.WriteError(w, http.StatusInternalServerError, err.Error(), "Internal Server Error")
				return
			}
		} else {
			response.WriteError(w, http.StatusForbidden, "You do not have permission to upload files", "Forbidden")
		}
	}
	logger.Info("File upload handler completed")
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) error {
	// Get the current user's home directory
	var homeDir, err = common.GetBaseDirectory()
	if err != nil {
		return err
	}

	// Create the file_sync directory
	syncDir := filepath.Join(homeDir, os.Getenv("FILE_SYNC_PATH"))
	if err := os.MkdirAll(syncDir, 0755); err != nil {
		return err
	}

	// Increase the max memory for multipart form to handle larger files
	// 5GB = 5 * 1024 * 1024 * 1024 bytes
	maxMemory := int64(5 * 1024 * 1024 * 1024)

	// Parse the multipart form with 5GB memory limit
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error(), "Bad Request")
		return nil
	}

	// Get all files from the form
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		response.WriteError(w, http.StatusBadRequest, "No files were uploaded", "Bad Request")
		return nil
	}

	logger.Infof("Processing %d files concurrently", len(files))

	// Set up concurrency
	type fileResult struct {
		index   int
		fileURL dto.DateURL
		err     error
	}

	// Create a channel for results
	resultChan := make(chan fileResult, len(files))

	// Process files concurrently
	for i, fileHeader := range files {
		// Start a goroutine for each file
		go func(index int, header *multipart.FileHeader) {
			result := fileResult{index: index}

			// Open the uploaded file
			file, err := header.Open()
			if err != nil {
				result.err = err
				resultChan <- result
				return
			}
			defer file.Close()

			// Replace spaces with underscores in filename
			originalFilename := filepath.Base(header.Filename)
			safeFilename := strings.ReplaceAll(originalFilename, " ", "_")

			// Create the destination file
			dstPath := filepath.Join(syncDir, safeFilename)
			dst, err := os.Create(dstPath)
			if err != nil {
				result.err = err
				resultChan <- result
				return
			}
			defer dst.Close()

			// Copy the uploaded file to the destination
			if _, err := io.Copy(dst, file); err != nil {
				result.err = err
				resultChan <- result
				return
			}

			// Generate signed URL for this file
			signedURL, err := generateSignedURL(dst)
			if err != nil {
				result.err = err
				resultChan <- result
				return
			}

			// Create the file URL object
			result.fileURL = linkGenerated(signedURL, dst, safeFilename)

			// Send the result
			resultChan <- result

		}(i, fileHeader)
	}

	// Collect results
	uploadedFiles := make([]dto.DateURL, len(files))
	for i := 0; i < len(files); i++ {
		result := <-resultChan
		if result.err != nil {
			return result.err
		}
		uploadedFiles[result.index] = result.fileURL
	}

	// Log completion
	logger.Infof("Successfully processed %d files", len(files))

	// If only one file was uploaded, use the existing response structure
	if len(uploadedFiles) == 1 {
		response.WriteSuccess(w, http.StatusOK, "File uploaded successfully", uploadedFiles[0])
	} else {
		// For multiple files, create a custom response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Create a custom response structure for multiple files
		type multiFileResponse struct {
			Status  int           `json:"status"`
			Message string        `json:"message"`
			Files   []dto.DateURL `json:"files"`
		}

		resp := multiFileResponse{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("%d files uploaded successfully", len(files)),
			Files:   uploadedFiles,
		}

		json.NewEncoder(w).Encode(resp)
	}
	return nil
}

func generateSignedURL(file *os.File) (string, error) {
	baseURL := os.Getenv("BASE_URL")

	// Get the base filename and encode it
	fileName := filepath.Base(file.Name())
	encodedFileName := encodeFileName(fileName)

	expires := time.Now().Add(15 * time.Minute).Unix()

	// Create the data to be signed: the encoded file path and expiration time
	dataToSign := fmt.Sprintf("%s:%d", encodedFileName, expires)

	// Create a new HMAC with SHA256
	h := hmac.New(sha256.New, []byte(os.Getenv("SIGNING_KEY")))
	h.Write([]byte(dataToSign))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Create the full URL with query parameters
	u, err := url.Parse(baseURL + common.APIFilePreviewPath)
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("file", encodedFileName)
	query.Set("expires", strconv.FormatInt(expires, 10))
	query.Set("signature", signature)
	u.RawQuery = query.Encode()

	return u.String(), nil
}

func linkGenerated(preSignedURL string, file *os.File, fileName string) dto.DateURL {
	data, downloadURL := dataAccessURLGenerator(file)

	return dto.DateURL{
		FileName:    fileName, // Keep original filename for display
		PreviewURL:  preSignedURL,
		URL:         data,
		DownloadURL: downloadURL,
	}
}

func dataAccessURLGenerator(file *os.File) (string, string) {
	encodedFileName := encodeFileName(filepath.Base(file.Name()))
	data := os.Getenv("BASE_URL") + common.APIDataPath + encodedFileName
	downloadURL := os.Getenv("BASE_URL") + common.APIDownloadPath + encodedFileName
	return data, downloadURL
}

/*func getFileName(file *os.File) string {
	return encodeFileName(filepath.Base(file.Name()))
}*/

// encodeFileName encodes a filename to base64 to hide actual filename in URLs
func encodeFileName(fileName string) string {
	return common.EncodeFileName(fileName)
}
