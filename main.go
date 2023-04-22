package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/madeinatria/love-all-backend/internal/app"
	"github.com/madeinatria/love-all-backend/internal/config"
	"github.com/madeinatria/love-all-backend/internal/logger"
	"go.uber.org/zap"
)

var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {

	// Parse command line flags
	configFile := flag.String("config", "", "Path to configuration file")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version information and exit
	if *showVersion {
		fmt.Printf("Service version: %s\n", version)
		fmt.Printf("Built on: %s\n", buildDate)
		os.Exit(0)
	}

	// Check if the configuration file path is provided
	log.Println(configFile)
	if *configFile == "" {
		fmt.Println("error: configuration file path is required")
		os.Exit(1)
	}

	// Load the configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("error loading configuration file: %v\n", err)
		os.Exit(1)
	}

	// Create the logger
	logger := logger.New(cfg.LogLevel)
	defer logger.Sync()

	// Add version information to the logger
	logger = logger.With(
		zap.String("version", version),
		zap.String("buildDate", buildDate),
		zap.String("configFile", *configFile),
	)

	// Create the service
	svc := app.New(cfg, logger)

	// Run the service
	if err := svc.Run(); err != nil {
		logger.Error("service stopped with error", zap.Error(err))
		os.Exit(1)
	}
}
