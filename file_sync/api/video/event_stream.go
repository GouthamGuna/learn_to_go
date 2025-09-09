package video

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type EventStreamHandler struct {
	VideoDir    string
	ChunkSize   int64
	ContentType string
}

func NewEventStreamHandler(videoDir string) *EventStreamHandler {
	os.MkdirAll(videoDir, 0755)
	return &EventStreamHandler{
		VideoDir:    videoDir,
		ChunkSize:   1024 * 1024, // 1MB chunks
		ContentType: "video/mp4",
	}
}

// StreamVideo streams video using HTTP range requests
func (e *EventStreamHandler) StreamVideo(w http.ResponseWriter, r *http.Request, videoPath string) {
	filePath := filepath.Join(e.VideoDir, videoPath)
	
	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	fileSize := fileInfo.Size()
	
	// Handle range requests
	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		// Serve full file if no range requested
		e.serveFullFile(w, filePath, fileSize)
		return
	}

	// Parse range header
	start, end, err := e.parseRangeHeader(rangeHeader, fileSize)
	if err != nil {
		http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	e.servePartialContent(w, r, filePath, start, end, fileSize)
}

func (e *EventStreamHandler) serveFullFile(w http.ResponseWriter, filePath string, fileSize int64) {
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", e.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	w.Header().Set("Accept-Ranges", "bytes")

	io.Copy(w, file)
}

func (e *EventStreamHandler) servePartialContent(w http.ResponseWriter, r *http.Request, filePath string, start, end, totalSize int64) {
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Seek to start position
	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		http.Error(w, "Error seeking file", http.StatusInternalServerError)
		return
	}

	contentLength := end - start + 1
	
	w.Header().Set("Content-Type", e.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusPartialContent)

	// Stream the content in chunks
	buffer := make([]byte, e.ChunkSize)
	bytesRemaining := contentLength

	for bytesRemaining > 0 {
		bytesToRead := min(e.ChunkSize, bytesRemaining)
		n, err := file.Read(buffer[:bytesToRead])
		if err != nil && err != io.EOF {
			log.Printf("Error reading chunk: %v", err)
			break
		}

		if n == 0 {
			break
		}

		w.Write(buffer[:n])
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}

		bytesRemaining -= int64(n)

		// Add small delay to simulate streaming (optional)
		time.Sleep(10 * time.Millisecond)
	}
}

func (e *EventStreamHandler) parseRangeHeader(rangeHeader string, fileSize int64) (int64, int64, error) {
	// Range header format: "bytes=start-end"
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, fmt.Errorf("invalid range format")
	}

	rangeStr := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeStr, "-")
	
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format")
	}

	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start position")
	}

	var end int64
	if parts[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid end position")
		}
	}

	if start < 0 || end >= fileSize || start > end {
		return 0, 0, fmt.Errorf("range out of bounds")
	}

	return start, end, nil
}

// LiveStream handles Server-Sent Events for live streaming metadata
func (e *EventStreamHandler) LiveStream(w http.ResponseWriter, r *http.Request, videoID string) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create context for handling client disconnects
	ctx := r.Context()
	
	// Send initial event
	fmt.Fprintf(w, "event: connected\ndata: {\"video_id\": \"%s\", \"timestamp\": \"%s\"}\n\n", videoID, time.Now().Format(time.RFC3339))
	flusher.Flush()

	// Simulate live events (replace with actual video processing)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Client disconnected
			log.Printf("Client disconnected from live stream: %s", videoID)
			return
		case <-ticker.C:
			// Send periodic events
			event := fmt.Sprintf(
				"event: update\ndata: {\"video_id\": \"%s\", \"viewers\": %d, \"timestamp\": \"%s\"}\n\n",
				videoID,
				time.Now().Unix()%100, // Simulated viewer count
				time.Now().Format(time.RFC3339),
			)
			
			if _, err := fmt.Fprint(w, event); err != nil {
				log.Printf("Error sending event: %v", err)
				return
			}
			flusher.Flush()
		}
	}
}