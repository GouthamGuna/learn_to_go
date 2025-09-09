package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GouthamGuna/file_sync/api/video"
	common "github.com/GouthamGuna/file_sync/pkg/common"
)

type VideoHandler struct {
	HLSProcessor       *video.HLSProcessor
	EventStreamHandler *video.EventStreamHandler
}

func NewVideoHandler(hlsProcessor *video.HLSProcessor, eventHandler *video.EventStreamHandler) *VideoHandler {
	return &VideoHandler{
		HLSProcessor:       hlsProcessor,
		EventStreamHandler: eventHandler,
	}
}

// Add this helper function at the top level
func checkFFmpegAvailable() error {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg is not available: %v", err)
	}
	return nil
}

// UploadVideoHandler handles video upload and HLS conversion
func (vh *VideoHandler) UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(100 << 20) // 100 MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving video", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var homeDir, er = common.GetBaseDirectory()
	if er != nil {
		log.Fatalf("Error getting base directory: %v", er)
	}

	// Save uploaded video
	tempPath := filepath.Join(homeDir, os.Getenv("FILE_SYNC_PATH"), header.Filename)
	//outFile, err := os.Create(tempPath)
	log.Printf("Saving video to %s\n", tempPath)
	//log.Printf("outFile: %v\n", outFile)
	if err != nil {
		log.Printf("Error: %v\n", err)
		http.Error(w, "Error saving video", http.StatusInternalServerError)
		return
	}
	//defer outFile.Close()
	//defer os.Remove(tempPath) // Cleanup temp file

	// _, err = io.Copy(tempPath, file)
	// if err != nil {
	// 	http.Error(w, "Error copying video", http.StatusInternalServerError)
	// 	return
	// }

	// Check FFmpeg availability before proceeding
	if err := checkFFmpegAvailable(); err != nil {
		log.Printf("FFmpeg check failed: %v", err)
		http.Error(w, "Server configuration error: FFmpeg is not available", http.StatusInternalServerError)
		return
	}

	videoID := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	outputDir := filepath.Join(vh.HLSProcessor.OutputDir, videoID)

	hlsProcessor := video.NewHLSProcessor(outputDir)
	qualities := video.DefaultQualities()

	// Convert to HLS
	err = hlsProcessor.ConvertToHLS(tempPath, qualities)
	if err != nil {
		log.Printf("HLS conversion failed: %v", err)
		// Cleanup any partial files that may have been created
		os.RemoveAll(outputDir)
		http.Error(w, fmt.Sprintf("Video conversion failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"status":"success","video_id":"%s","message":"Video converted to HLS"}`, videoID)
}

// StreamVideoHandler handles video streaming
func (vh *VideoHandler) StreamVideoHandler(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("id")
	if videoID == "" {
		http.Error(w, "Video ID required", http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")

	switch format {
	case "hls":
		vh.serveHLSPlaylist(w, r, videoID)
	case "direct":
		videoPath := r.URL.Query().Get("path")
		if videoPath == "" {
			http.Error(w, "Video path required", http.StatusBadRequest)
			return
		}
		vh.EventStreamHandler.StreamVideo(w, r, videoPath)
	case "live":
		vh.EventStreamHandler.LiveStream(w, r, videoID)
	default:
		// Default to HLS
		vh.serveHLSPlaylist(w, r, videoID)
	}
}

func (vh *VideoHandler) serveHLSPlaylist(w http.ResponseWriter, r *http.Request, videoID string) {
	playlistPath := filepath.Join(vh.HLSProcessor.OutputDir, videoID, "stream.m3u8")

	// Check if playlist exists
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	// Serve HLS playlist
	http.ServeFile(w, r, playlistPath)
}

// ServeHLSSegment serves individual HLS segments
func (vh *VideoHandler) ServeHLSSegment(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("video_id")
	quality := r.URL.Query().Get("quality")
	segment := r.URL.Query().Get("segment")

	if videoID == "" || quality == "" || segment == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	segmentPath := filepath.Join(vh.HLSProcessor.OutputDir, videoID, quality, segment)
	http.ServeFile(w, r, segmentPath)
}
