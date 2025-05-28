package main

import (
	"Advertising/configs"
	"Advertising/pkg/repository"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load environment variables
	configs.LoadEnvConfig()

	// Load YAML config and apply overrides
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Initialize database connection
	db, err := repository.NewDb(cfg)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()
	//todo handle error

	// Initialize web server
	e := echo.New()

	// Start HTTP server
	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s...", address)
	if err := e.Start(address); err != nil {
		log.Fatalf("error starting server: %v", err)
		os.Exit(1)
	}
}
