package main

import (
	"net/http"

	"github.com/GouthamGuna/file_sync/configs"
	"github.com/gorilla/mux"

	handlers "github.com/GouthamGuna/file_sync/api/handler"
	middleware "github.com/GouthamGuna/file_sync/api/middleware"
	common "github.com/GouthamGuna/file_sync/pkg/common"
)

func RegisterRoutes() *mux.Router {

	router := mux.NewRouter()
	router.Use(configs.Cors)

	// Unauthenticated routes
	publicRouter := router.PathPrefix(common.APIBasePath_V1).Subrouter()
	publicRouter.HandleFunc("/users/generate_token", handlers.GenerateUserToken).Methods("POST")
	publicRouter.HandleFunc("/file_preview", handlers.FilePreviewHandler).Methods("GET")

	// Authenticated routes
	protectedRouter := router.PathPrefix(common.APIBasePath_V1).Subrouter()
	protectedRouter.Use(middleware.AuthenticateMiddleware)
	protectedRouter.Use(configs.CSRF)
	protectedRouter.HandleFunc("/csrf-token", configs.CSRFTokenHandler).Methods("GET")
	protectedRouter.HandleFunc("/uploads", handlers.Uploads).Methods("POST")
	protectedRouter.HandleFunc("/data", handlers.StreamingFile).Methods("GET")
	protectedRouter.HandleFunc("/download", handlers.DownloadFile).Methods("GET")
	protectedRouter.HandleFunc("/files", handlers.FilesHandler).Methods("GET")
	protectedRouter.HandleFunc("/delete", handlers.DeleteFile).Methods("DELETE")

	// Handle 404 - Route not found
	router.NotFoundHandler = http.HandlerFunc(common.NotFoundHandler)

	// Handle 405 - Method not allowed
	router.MethodNotAllowedHandler = http.HandlerFunc(common.MethodNotAllowedHandler)

	// Serve index.html at root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	}).Methods("GET")

	// Handle favicon.ico
	router.HandleFunc("/favicon.ico", common.FaviconHandler).Methods("GET")

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	middleware.NewHeaderMiddleware(middleware.DefaultHeaders()).SetHeaders(router)

	return router
}
