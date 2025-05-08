package main

import (
	"Advertising/configs"
	"Advertising/pkg/repository"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	// Load environment variables from .env (if any)
	configs.LoadEnvConfig()

	// Load configuration from config.yaml
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	// Initialize DB connection
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer db.Close() // Ensure connection is closed on shutdown

	// Initialize Echo web server
	e := echo.New()

	// Setup routes here...
	// e.POST("/api/adverts", advertHandler.CreateAdvert)
	// e.GET("/api/adverts/:id", advertHandler.GetAdvertByID)
	// e.GET("/api/adverts", advertHandler.GetAdverts)

	// Start HTTP server
	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s...", address)
	if err := e.Start(address); err != nil {
		log.Fatalf("error starting server: %v", err)
		os.Exit(1)
	}
}
