//go:build windows
// +build windows

package configs

import (
	"context"
	"net/http"
	"os"

	logger "github.com/GouthamGuna/file_sync/pkg/logger"
	"github.com/kardianos/service"
)

var (
	httpServer  *http.Server
	autoService = os.Getenv("AUTO_SERVICE") == "true"
)

type program struct {
	cancel context.CancelFunc
}

func (p *program) Start(s service.Service) error {
	// Start the HTTP server in a goroutine
	go p.run()
	return nil
}

func (p *program) run() {
	_, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	logger.Infof("Starting HTTP server on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("HTTP server error: %v", err)
	}
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Shutting down Windows service...")

	// Gracefully shutdown HTTP server
	if httpServer != nil {
		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Errorf("Error during server shutdown: %v", err)
		}
	}

	// Cancel context
	if p.cancel != nil {
		p.cancel()
	}

	logger.Info("Service shutdown complete")
	return nil
}

func RunAsService(server *http.Server) {
	// Store the server reference for the service
	httpServer = server

	svcConfig := &service.Config{
		Name:        "FileSyncService",
		DisplayName: "File Sync Service",
		Description: "File synchronization service built in Go",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Fatalf("Failed to create service: %v", err)
	}

	// If autoService is true and running interactively, auto-install and start service
	if autoService && service.Interactive() {
		logger.Info("Auto-service mode enabled. Installing and starting service...")
		if err := s.Install(); err != nil {
			if err.Error() != "The specified service already exists." {
				logger.Fatalf("Failed to auto-install service: %v", err)
			} else {
				logger.Info("Service already installed.")
			}
		} else {
			logger.Info("Service installed successfully.")
		}
		if err := s.Start(); err != nil {
			logger.Fatalf("Failed to auto-start service: %v", err)
		}
		logger.Info("Service started successfully. Exiting main process.")
		os.Exit(0)
	}

	// Check if running as Windows service or normal console
	if isWindowsService() {
		// Handle service commands if provided
		if len(os.Args) > 1 {
			handleServiceCommand(s)
			return
		}
		// Run as service
		logger.Info("Running as Windows service")
		if err := s.Run(); err != nil {
			logger.Fatalf("Service failed: %v", err)
		}
	} else {
		// Running as normal console application
		logger.Infof("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("Server failed: %v", err)
		}
	}
}

func isWindowsService() bool {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		return cmd == "install" || cmd == "uninstall" || cmd == "start" ||
			cmd == "stop" || cmd == "restart" || cmd == "status"
	}
	return !service.Interactive()
}

func handleServiceCommand(s service.Service) {
	cmd := os.Args[1]
	switch cmd {
	case "install":
		if err := s.Install(); err != nil {
			logger.Fatalf("Failed to install service: %v", err)
		}
		logger.Info("Service installed successfully")
		os.Exit(0)
	case "uninstall":
		if err := s.Uninstall(); err != nil {
			logger.Fatalf("Failed to uninstall service: %v", err)
		}
		logger.Info("Service uninstalled successfully")
		os.Exit(0)
	case "start":
		if err := s.Start(); err != nil {
			logger.Fatalf("Failed to start service: %v", err)
		}
		logger.Info("Service started successfully")
		os.Exit(0)
	case "stop":
		if err := s.Stop(); err != nil {
			logger.Fatalf("Failed to stop service: %v", err)
		}
		logger.Info("Service stopped successfully")
		os.Exit(0)
	case "restart":
		if err := s.Restart(); err != nil {
			logger.Fatalf("Failed to restart service: %v", err)
		}
		logger.Info("Service restarted successfully")
		os.Exit(0)
	default:
		// Unknown command, run as console app
		logger.Infof("Starting server on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatalf("Server failed: %v", err)
		}
	}
}
