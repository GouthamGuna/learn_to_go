package video 

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type HLSProcessor struct {
	OutputDir    string
	SegmentTime  int    // Segment duration in seconds
	PlaylistName string // Master playlist name
}

func NewHLSProcessor(outputDir string) *HLSProcessor {
	os.MkdirAll(outputDir, 0755)
	return &HLSProcessor{
		OutputDir:    outputDir,
		SegmentTime:  10,
		PlaylistName: "stream.m3u8",
	}
}

// ConvertToHLS converts a video file to HLS format
func (h *HLSProcessor) ConvertToHLS(inputPath string, qualities []Quality) error {
	// Check if ffmpeg is available
	if !h.isFFmpegAvailable() {
		return fmt.Errorf("ffmpeg is required for HLS conversion")
	}

	// Generate quality-specific streams
	var streams []string
	for _, quality := range qualities {
		streamPath, err := h.createQualityStream(inputPath, quality)
		if err != nil {
			return err
		}
		streams = append(streams, streamPath)
	}

	// Create master playlist
	return h.createMasterPlaylist(streams, qualities)
}

type Quality struct {
	Name      string
	Width     int
	Height    int
	Bitrate   string
	Bandwidth string
}

// DefaultQualities returns common streaming qualities
func DefaultQualities() []Quality {
	return []Quality{
		{
			Name:      "1080p",
			Width:     1920,
			Height:    1080,
			Bitrate:   "5000k",
			Bandwidth: "5000000",
		},
		{
			Name:      "720p",
			Width:     1280,
			Height:    720,
			Bitrate:   "2800k",
			Bandwidth: "2800000",
		},
		{
			Name:      "480p",
			Width:     854,
			Height:    480,
			Bitrate:   "1400k",
			Bandwidth: "1400000",
		},
		{
			Name:      "360p",
			Width:     640,
			Height:    360,
			Bitrate:   "800k",
			Bandwidth: "800000",
		},
	}
}

func (h *HLSProcessor) createQualityStream(inputPath string, quality Quality) (string, error) {
	qualityDir := filepath.Join(h.OutputDir, quality.Name)
	os.MkdirAll(qualityDir, 0755)

	outputPattern := filepath.Join(qualityDir, "segment_%03d.ts")
	playlistPath := filepath.Join(qualityDir, "playlist.m3u8")

	args := []string{
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=%d:%d", quality.Width, quality.Height),
		"-c:v", "libx264",
		"-b:v", quality.Bitrate,
		"-c:a", "aac",
		"-b:a", "128k",
		"-f", "hls",
		"-hls_time", fmt.Sprintf("%d", h.SegmentTime),
		"-hls_list_size", "0", // Keep all segments in playlist
		"-hls_segment_filename", outputPattern,
		playlistPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stderr = os.Stderr // Capture ffmpeg output

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg conversion failed: %v", err)
	}

	return playlistPath, nil
}

func (h *HLSProcessor) createMasterPlaylist(streams []string, qualities []Quality) error {
	masterPlaylist := "#EXTM3U\n#EXT-X-VERSION:3\n"

	for i, stream := range streams {
		quality := qualities[i]
		relativePath := strings.TrimPrefix(stream, h.OutputDir+"/")
		masterPlaylist += fmt.Sprintf(
			"#EXT-X-STREAM-INF:BANDWIDTH=%s,RESOLUTION=%dx%d,NAME=\"%s\"\n%s\n",
			quality.Bandwidth, quality.Width, quality.Height, quality.Name, relativePath,
		)
	}

	masterPath := filepath.Join(h.OutputDir, h.PlaylistName)
	return os.WriteFile(masterPath, []byte(masterPlaylist), 0644)
}

func (h *HLSProcessor) isFFmpegAvailable() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// CleanupOldSegments removes segments older than specified duration
func (h *HLSProcessor) CleanupOldSegments(maxAge time.Duration) {
	for {
		time.Sleep(5 * time.Minute) // Run cleanup every 5 minutes
		h.cleanupSegments(maxAge)
	}
}

func (h *HLSProcessor) cleanupSegments(maxAge time.Duration) {
	cutoffTime := time.Now().Add(-maxAge)
	
	err := filepath.Walk(h.OutputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if strings.HasSuffix(path, ".ts") && info.ModTime().Before(cutoffTime) {
			os.Remove(path)
		}
		return nil
	})
	
	if err != nil {
		log.Printf("Error cleaning up segments: %v", err)
	}
}