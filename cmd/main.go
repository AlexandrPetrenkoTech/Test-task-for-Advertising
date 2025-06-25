package main

import (
	"fmt"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/repository/postgres"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/service"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	// auto-generated package with documentation
	_ "github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/docs"

	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/configs"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/handler"
	"github.com/AlexandrPetrenkoTech/Test-task-for-Advertising/pkg/repository"
)

// @title Advertising API
// @version 1.0
// @description A service for submitting and storing advertisements
// @host localhost:8080
// @BasePath /api

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

	// Initialize web server
	e := echo.New()

	// Swagger UI available at: /swagger/index.html
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Register routes
	// let's assume you're creating the service and passing it directly to the handler:
	advertSvc := service.NewAdvertService(postgres.NewPostgresAdvertRepo(db), postgres.NewPostgresPhotoRepo(db))
	handler.NewAdvertHandler(e, advertSvc)

	// Start HTTP server
	address := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s...", address)
	if err := e.Start(address); err != nil {
		log.Fatalf("error starting server: %v", err)
		os.Exit(1)
	}
}
