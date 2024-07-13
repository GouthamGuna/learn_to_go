package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func sendResJsonContent(filePath string, response http.ResponseWriter) {
	setResponseHeader(response)

	// Read the contents of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(response, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Set the content type to text/markdown or text/html
	response.Header().Set("Content-Type", "application/json")

	// Write the file contents to the response
	fmt.Fprintf(response, "%s", content)
}

/*func handleToStream(w http.ResponseWriter, r *http.Request) {
	setResponseHeader(w)

	filename := strings.TrimPrefix(r.URL.Path, "/stream/")
	fmt.Println("filename : ", filename)
	file, err := os.Open("datasource/videos/" + filename + ".mp4")
	fmt.Println("err : ", err)
	if err != nil {
		http.Error(w, "Video not found.", http.StatusNotFound)
		return
	}

	defer file.Close()
	w.Header().Set("Content-Type", "video/mp4")
	buffer := make([]byte, 1024*1024*1024) // 1 GB buffer size
	io.CopyBuffer(w, file, buffer)
}*/

func handleToStream(w http.ResponseWriter, r *http.Request) {
	setResponseHeader(w)

	filename := strings.TrimPrefix(r.URL.Path, "/stream/")
	fmt.Println("filename : ", filename)
	file, err := os.Open("datasource/videos/" + filename + ".mp4")
	fmt.Println("err : ", err)
	if err != nil {
		http.Error(w, "Video not found.", http.StatusNotFound)
		return
	}

	defer file.Close()
	w.Header().Set("Content-Type", "video/mp4")

	// Create a channel to send data chunks
	dataCh := make(chan []byte)

	// Start a goroutine to read data from the file
	go func() {
		defer close(dataCh)
		buffer := make([]byte, 1024*1024) // 1 MB buffer size
		for {
			n, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading file:", err)
				}
				return
			}
			dataCh <- buffer[:n]
		}
	}()

	// Check if the response writer implements the http.Flusher interface
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Write data from the channel to the response writer
	for data := range dataCh {
		_, err := w.Write(data)
		if err != nil {
			fmt.Println("Error writing to response:", err)
			return
		}
		flusher.Flush()
	}
}

func setResponseHeader(response http.ResponseWriter) {

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET") // "GET, POST, PUT, DELETE, OPTIONS"
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	response.Header().Set("Server", "Cerpsoft - GoLang")
	response.Header().Set("Author", "Gowtham Sankar Gunasekaran")
}
