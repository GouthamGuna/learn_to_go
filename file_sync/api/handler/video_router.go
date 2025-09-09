package api

import (
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/GouthamGuna/file_sync/api/video"
)

func InitializeVideoServices(router *mux.Router) {
	
	// Initialize video processors
	hlsProcessor := video.NewHLSProcessor("./hls_output")
	eventHandler := video.NewEventStreamHandler("./videos")

	// Start segment cleanup goroutine
	go hlsProcessor.CleanupOldSegments(24 * time.Hour) // Cleanup segments older than 24 hours

	// Initialize handlers
	videoHandler := NewVideoHandler(hlsProcessor, eventHandler)

	router.HandleFunc("/video/upload", videoHandler.UploadVideoHandler)
	router.HandleFunc("/video/stream", videoHandler.StreamVideoHandler)
	router.HandleFunc("/video/segment", videoHandler.ServeHLSSegment)

	// Serve HLS files statically
	router.Handle("/hls/", http.StripPrefix("/hls/", http.FileServer(http.Dir(hlsProcessor.OutputDir))))
}