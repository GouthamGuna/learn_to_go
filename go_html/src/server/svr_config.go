package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func SetResHTMLContent(htmlPath string, response http.ResponseWriter) {
	setResponseHeader(response)

	// Read the contents of the src/data_source/md/home.md file
	content, err := os.ReadFile(htmlPath)
	if err != nil {
		http.Error(response, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Set the content type to text/markdown or text/html
	response.Header().Set("Content-Type", "text/html")

	// Write the file contents to the response
	fmt.Fprintf(response, "%s", content)
}

func setResponseHeader(response http.ResponseWriter) {
	response.Header().Set("Server", "Cerpsoft - GoLang")
	response.Header().Set("Author", "Gowtham Sankar Gunasekaran")
}

func ReceiveHandler(w http.ResponseWriter, r *http.Request) {
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
    fmt.Printf("\nReceived byte array: %v\n", data)
    // Process the received byte array here
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Read the byte array data
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Receive the path from the client
	path := r.FormValue("path")

	// Read the file contents into a byte array
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Combine the byte array data and file/directory data
	finalData := append(data, fileData...)

	// Create the downloadable file
	downloadFile := "E:\\download.zip"
	err = ioutil.WriteFile(downloadFile, finalData, 0644)
	if err != nil {
		http.Error(w, "Error creating download file", http.StatusInternalServerError)
		return
	}

	// Send the downloadable file to the client
	sendDownloadToClient(w, r, downloadFile)
}

func sendDownloadToClient(w http.ResponseWriter, r *http.Request, filename string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, filename)
}
