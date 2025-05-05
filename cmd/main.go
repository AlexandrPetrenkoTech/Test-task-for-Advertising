package main

import (
	"Advertising/configs"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	// Loading configuration
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatal("error loading config:", err)
	}
	// Initializing the database and other services
	//dbRepo := repository.NewDBRepository(cfg)
	//advertService := service.NewAdvertService(dbRepo)
	//advertHandler := handler.NewAdvertHandler(advertService)

	//Initializing Echo
	e := echo.New()

	//Setting up routes
	//e.POST("/api/adverts", advertHandler.CreateAdvert)
	//e.GET("api/advets/:id", advertHandler.GetAdvertByID)
	//e.GET("/api/adverts", advertHandler.GetAdverts)

	//Starting the server
	log.Printf("Starting server on port %s...", cfg.Server.Port)
	if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal("Error starting server: %v", err)
		os.Exit(1)
	}
}
