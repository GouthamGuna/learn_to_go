package pkg

type FileMetadata struct {
	Name         string                  `json:"name"`
	Path         string                  `json:"path"` // Relative path from base directory
	Size         string                  `json:"size"`
	CreateDate   string                  `json:"create_date"`
	ModifiedDate string                  `json:"modified_date"`
	Type         string                  `json:"type"` // "file" or "directory"
	Links        *FileLinks              `json:"links,omitempty"` // Omit for directories
	Files        map[string]FileMetadata `json:"files,omitempty"` // Child files and directories
}