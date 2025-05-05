package main

import (
	"Advertising/configs"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	// Load environment variables first
	configs.LoadEnvConfig()

	// Load the main application configuration from config.yaml
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatal("error loading config:", err)
	}

	// Initialize Echo
	e := echo.New()

	// Set up routes
	// e.POST("/api/adverts", advertHandler.CreateAdvert)
	// e.GET("api/advets/:id", advertHandler.GetAdvertByID)
	// e.GET("/api/adverts", advertHandler.GetAdverts)

	// Start the server
	log.Printf("Starting server on port %d...", cfg.Server.Port)
	if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal("Error starting server: %v", err)
		os.Exit(1)
	}
}
