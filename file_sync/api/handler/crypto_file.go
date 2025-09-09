package api

import (
	"io"
	"net/http"
	"path/filepath"

	response "github.com/GouthamGuna/file_sync/pkg/response"
	crypto "github.com/GouthamGuna/file_sync/pkg/service/crypto"
)

type FileHandler struct {
	Storage *crypto.FileStorage
}

func NewFileHandler(storage *crypto.FileStorage) *FileHandler {
	return &FileHandler{Storage: storage}
}

// UploadHandler handles file upload and encryption
func (fh *FileHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.WriteError(w, http.StatusMethodNotAllowed, "Only POST method is allowed for file upload.", "Method not allowed")
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 30) // 10 GB max memory
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Error parsing form the request", err.Error())
		return
	}

	// Get the file from the form data
	file, header, err := r.FormFile("file")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Error retrieving file", err.Error())
		return
	}
	defer file.Close()

	// Save and encrypt the file
	encryptedFilename, err := fh.Storage.SaveEncryptedFile(file, header.Filename)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Error saving file", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	response.FileUploadedResponse(w, "File uploaded and encrypted successfully.", encryptedFilename)
}

// DownloadHandler handles file download with decryption
func (fh *FileHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	encryptedFilename := r.URL.Query().Get("file")
	if encryptedFilename == "" {
		http.Error(w, "File parameter is required", http.StatusBadRequest)
		response.WriteError(w, http.StatusBadRequest, "File parameter is required", "Missing file parameter")
		return
	}

	if !fh.Storage.FileExists(encryptedFilename) {
		http.Error(w, "File not found", http.StatusNotFound)
		response.WriteError(w, http.StatusNotFound, "File not found", "Inaccessible file")
		return
	}

	// Get original filename from header or query parameter
	originalFilename := r.URL.Query().Get("filename")
	if originalFilename == "" {
		originalFilename = encryptedFilename // Fallback to encrypted filename
	}

	// Get decrypted file reader
	file, err := fh.Storage.GetDecryptedFile(encryptedFilename)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Error decrypting file", err.Error())
		return
	}
	defer file.Close()

	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=\""+originalFilename+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")

	// Stream the decrypted file to the client
	_, err = io.Copy(w, file)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Error streaming file", err.Error())
		return
	}
}

// ServeHandler serves decrypted file content (for viewing)
func (fh *FileHandler) ServeHandler(w http.ResponseWriter, r *http.Request) {
	encryptedFilename := r.URL.Query().Get("file")
	if encryptedFilename == "" {
		response.WriteError(w, http.StatusBadRequest, "File parameter is required", "Missing file parameter")
		return
	}

	if !fh.Storage.FileExists(encryptedFilename) {
		http.Error(w, "File not found", http.StatusNotFound)
		response.WriteError(w, http.StatusNotFound, "File not found", "Inaccessible file")
		return
	}

	// Get decrypted content as bytes
	decryptedData, err := fh.Storage.GetDecryptedFileBytes(encryptedFilename)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Error decrypting file", err.Error())
		return
	}

	// Set appropriate content type based on file extension
	ext := filepath.Ext(encryptedFilename)
	switch ext {
	case ".txt", ".log":
		w.Header().Set("Content-Type", "text/plain")
	case ".html", ".htm":
		w.Header().Set("Content-Type", "text/html")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	case ".pdf":
		w.Header().Set("Content-Type", "application/pdf")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Write(decryptedData)
}

func InitializeFileStorage() *FileHandler {
	// Get encryption key from environment variable or use default
	encryptionKey := crypto.GetEncryptionKey()

	// Initialize file storage
	fileStorage := crypto.NewFileStorage("./uploads", encryptionKey)

	// Initialize handlers
	return NewFileHandler(fileStorage)
}
