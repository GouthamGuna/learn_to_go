package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error(), "Bad Request")
		return nil
	}
	defer file.Close()

	dst, err := os.Create(filepath.Join(syncDir, filepath.Base(fileHeader.Filename)))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}

	signedURL, err := generateSignedURL(dst)
	if err != nil {
		return err
	}

	response.WriteSuccess(w, http.StatusOK, "File uploaded successfully", linkGenerated(signedURL, dst))
	return nil
}

func generateSignedURL(file *os.File) (string, error) {

	baseURL := os.Getenv("BASE_URL")

	expires := time.Now().Add(15 * time.Minute).Unix()

	// Create the data to be signed: the file path and expiration time.
	dataToSign := fmt.Sprintf("%s:%d", getFileName(file), expires)

	// Create a new HMAC with SHA256.
	h := hmac.New(sha256.New, []byte(os.Getenv("SIGNING_KEY")))
	h.Write([]byte(dataToSign))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Create the full URL with query parameters.
	u, err := url.Parse(baseURL + common.APIFilePreviewPath)
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("file", getFileName(file))
	query.Set("expires", strconv.FormatInt(expires, 10))
	query.Set("signature", signature)
	u.RawQuery = query.Encode()

	return u.String(), nil
}

func linkGenerated(preSignedURL string, file *os.File) dto.DateURL {
	data, downloadURL := dataAccessURLGenerator(file)
	return dto.DateURL{
		PreviewURL:  preSignedURL,
		URL:         data,
		DownloadURL: downloadURL,
	}
}

func dataAccessURLGenerator(file *os.File) (string, string) {
	data := os.Getenv("BASE_URL") + common.APIDataPath + getFileName(file)
	downloadURL := os.Getenv("BASE_URL") + common.APIDownloadPath + getFileName(file)
	return data, downloadURL
}

func getFileName(file *os.File) string {
	return url.PathEscape(filepath.Base(file.Name()))
}
