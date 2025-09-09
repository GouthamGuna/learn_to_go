package api

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"

	middleware "github.com/GouthamGuna/file_sync/api/middleware"
	common "github.com/GouthamGuna/file_sync/pkg/common"
	dto "github.com/GouthamGuna/file_sync/pkg/dto"
	response "github.com/GouthamGuna/file_sync/pkg/response"
	logger "github.com/GouthamGuna/file_sync/pkg/logger"
)

// formatFileSize converts bytes to human readable format
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	if exp >= len(sizes) {
		exp = len(sizes) - 1
	}

	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), sizes[exp])
}

// formatDate formats time to a readable string
func formatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// buildHierarchicalStructure builds a hierarchical file structure
func buildHierarchicalStructure(flatMap map[string]dto.FileMetadata) map[string]dto.FileMetadata {
	hierarchy := make(map[string]dto.FileMetadata)

	// First, create all directory entries
	for path, metadata := range flatMap {
		if metadata.Type == "directory" {
			hierarchy[path] = dto.FileMetadata{
				Name:         metadata.Name,
				Path:         metadata.Path,
				Size:         metadata.Size,
				CreateDate:   metadata.CreateDate,
				ModifiedDate: metadata.ModifiedDate,
				Type:         "directory",
				Files:        make(map[string]dto.FileMetadata),
			}
		}
	}

	// Then, add files to their parent directories
	for path, metadata := range flatMap {
		if metadata.Type == "file" {
			dirPath := filepath.Dir(path)
			if dirPath == "." {
				dirPath = "" // Root directory
			}

			if dirPath == "" {
				// File is in root directory
				if _, exists := hierarchy[path]; !exists {
					hierarchy[path] = metadata
				}
			} else {
				// File is in subdirectory
				if dir, exists := hierarchy[dirPath]; exists {
					dir.Files[metadata.Name] = metadata
					hierarchy[dirPath] = dir
				}
			}
		}
	}

	// Finally, add subdirectories to their parent directories
	for path, metadata := range hierarchy {
		if metadata.Type == "directory" && path != "" {
			dirPath := filepath.Dir(path)
			if dirPath == "." {
				dirPath = "" // Root directory
			}

			if dirPath != "" {
				if parentDir, exists := hierarchy[dirPath]; exists {
					parentDir.Files[metadata.Name] = metadata
					hierarchy[dirPath] = parentDir
				}
			}
		}
	}

	return hierarchy
}

// getFileMetadata retrieves metadata for a single file or directory
func getFileMetadata(filePath, basePath, baseURL string, fileInfo fs.FileInfo) (dto.FileMetadata, error) {
	// Get relative path from base directory
	relPath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		relPath = filePath // Fallback to full path if relative path fails
	}

	// Get file stats for creation and modification times
	stat, err := os.Stat(filePath)
	if err != nil {
		return dto.FileMetadata{}, err
	}

	metadata := dto.FileMetadata{
		Name:         fileInfo.Name(),
		Path:         relPath,
		Size:         formatFileSize(fileInfo.Size()),
		CreateDate:   formatDate(stat.ModTime()),
		ModifiedDate: formatDate(stat.ModTime()),
		Type:         "file",
		Files:        make(map[string]dto.FileMetadata),
	}

	if fileInfo.IsDir() {
		metadata.Type = "directory"
		metadata.Size = formatFileSize(0)
	} else {
		// Only add links for files, not directories
		encodedPath := url.QueryEscape(relPath)
		metadata.Links = &dto.FileLinks{
			URL:      fmt.Sprintf("%s"+common.APIDataPath+"%s", baseURL, encodedPath),
			Download: fmt.Sprintf("%s"+common.APIDownloadPath+"%s", baseURL, encodedPath),
			Delete:   fmt.Sprintf("%s"+common.APIDeletePath+"%s", baseURL, encodedPath),
		}
		metadata.Files = nil // Files don't have child files
	}

	return metadata, nil
}

// walkDirectory recursively walks through directory and collects metadata
func walkDirectory(rootPath, baseURL string) (map[string]dto.FileMetadata, error) {
	flatMap := make(map[string]dto.FileMetadata)

	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			logger.Infof("Warning: Error accessing path %s: %v", path, err)
			return nil // Continue walking despite errors
		}

		metadata, err := getFileMetadata(path, rootPath, baseURL, info)
		if err != nil {
			logger.Infof("Warning: Could not get metadata for %s: %v", path, err)
			return nil
		}

		// Use relative path as key
		flatMap[metadata.Path] = metadata

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	// Build hierarchical structure from flat map
	hierarchy := buildHierarchicalStructure(flatMap)

	// Return only root level items
	rootItems := make(map[string]dto.FileMetadata)
	for path, item := range hierarchy {
		if filepath.Dir(path) == "." || filepath.Dir(path) == "" {
			rootItems[path] = item
		}
	}

	return rootItems, nil
}

// sortFiles alphabetically for consistent output
func sortFiles(files map[string]dto.FileMetadata) []dto.FileMetadata {
	var sorted []dto.FileMetadata
	for _, file := range files {
		sorted = append(sorted, file)
	}

	sort.Slice(sorted, func(i, j int) bool {
		// Directories first, then files
		if sorted[i].Type != sorted[j].Type {
			return sorted[i].Type == "directory"
		}
		return sorted[i].Name < sorted[j].Name
	})

	return sorted
}

// GetFilesMetadataHandler is the main function that returns the JSON response
func getFilesMetadataHandler() ([]byte, error) {
	var userHome, err = common.GetBaseDirectory();
	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %v", err)
	}

	homeDir := filepath.Join(userHome, os.Getenv("FILE_SYNC_PATH"))

	filesMetadata, err := walkDirectory(homeDir, os.Getenv("BASE_URL"))
	if err != nil {
		return nil, err
	}

	// Sort files for consistent output
	sortedFiles := sortFiles(filesMetadata)

	// Convert to JSON with proper indentation
	jsonData, err := json.MarshalIndent(map[string]interface{}{
		"files": sortedFiles,
	}, "", "    ")

	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	return jsonData, nil
}

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling files request")
	role, ok := middleware.GetUserRoleFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusInternalServerError, "Unable to determine user role", "Unauthorized")
		return
	}

	if role == "admin" {
		data, err := getFilesMetadataHandler()
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error(), "Internal Server Error")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		response.WriteError(w, http.StatusForbidden, "You do not have permission to access this resource", "Forbidden")
	}
	logger.Info("Files request completed")
}
