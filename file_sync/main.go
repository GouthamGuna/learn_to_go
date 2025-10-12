package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	config "github.com/GouthamGuna/file_sync/configs"
	logger "github.com/GouthamGuna/file_sync/pkg/logger"
)

func main() {
	// Define command line flags
	var ipAddress string
	flag.StringVar(&ipAddress, "ip", "", "Server IP address with port (e.g. 127.0.0.1:2427)")
	flag.Parse()

	// Check if IP address was provided
	if ipAddress == "" {
		fmt.Println("Error: IPv4 address is required")
		fmt.Println("Usage: go run . --ip <ip_address>")
		fmt.Println("Example: go run . --ip 127.0.0.1")
		os.Exit(1)
	}

	config.LoadEnv()

	if !strings.HasPrefix(ipAddress, "http://") && !strings.HasPrefix(ipAddress, "https://") {
		serverPort := strings.TrimPrefix(os.Getenv("SERVER_PORT"), ":")
		ipAddress = fmt.Sprintf("http://%s:%s", ipAddress, serverPort)
	} 
	os.Setenv("BASE_URL", ipAddress)
	logger.Infof("BASE_URL set to %s", ipAddress)

	var server_port = os.Getenv("SERVER_PORT")

	router := RegisterRoutes()

	server := http.Server{
		Addr:    server_port,
		Handler: router,
	}

	logger.Infof("Starting server on %s", server_port)
	logger.Fatal(server.ListenAndServe())

	config.RunAsService(&server)
}
