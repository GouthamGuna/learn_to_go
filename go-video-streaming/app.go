package main

import (
	"fmt"
	"strings"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	http.HandleFunc("/stream/{filename}", handleToStream)

	http.HandleFunc("/content_upload", handleContentUpload)

	http.HandleFunc("/upload", handleUpload)

	http.ListenAndServe(":2427", nil)
}

func handleToStream(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/stream/")
	fmt.Println("filename : ", filename)
	file, err := os.Open("videos/" + filename)
	if err != nil {
		http.Error(w, "Video not found.", http.StatusNotFound)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "video/mp4")
	buffer := make([]byte, 1024*1024*1024) // 1 GB buffer size
	io.CopyBuffer(w, file, buffer)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	// Create a new file on the filesystem
	// r.ParseMultipartForm(32 << 20) // 32 MB
	err := r.ParseMultipartForm(20 * 1024 * 1024 * 1024) // 20 GB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["video"]
	if len(files) == 0 {
		http.Error(w, "No file selected", http.StatusBadRequest)
		return
	}

	file, err := files[0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename := files[0].Filename
	dst, err := os.Create(filepath.Join("videos", filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully", filename)
}

func handleContentUpload(w http.ResponseWriter, r *http.Request) {
	// Read the contents of the src/data_source/md/home.md file
	content, err := os.ReadFile("public/upload.html")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Set the content type to text/markdown or text/html
	w.Header().Set("Content-Type", "text/html")

	// Write the file contents to the response
	fmt.Fprintf(w, "%s", content)
}
