package main

import (
	"net/http"
	"os"

	config "github.com/GouthamGuna/file_sync/configs"
	logger "github.com/GouthamGuna/file_sync/pkg/logger"
)

func main() {

	config.LoadEnv()
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
