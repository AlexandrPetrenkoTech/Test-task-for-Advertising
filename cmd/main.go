package main

import (
	"Advertising/configs"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	// Загрузка конфигурации
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatal("error loading config:", err)
	}
	// Инициализация базы данных и других сервисов
	//dbRepo := repository.NewDBRepository(cfg)
	//advertService := service.NewAdvertService(dbRepo)
	//advertHandler := handler.NewAdvertHandler(advertService)

	//Инициализация Echo
	e := echo.New()

	//Настройка маршрутов
	//e.POST("/api/adverts", advertHandler.CreateAdvert)
	//e.GET("api/advets/:id", advertHandler.GetAdvertByID)
	//e.GET("/api/adverts", advertHandler.GetAdverts)

	//Запуск сервера
	log.Printf("Starting server on port %s...", cfg.Server.Port)
	if err := e.Start(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatal("Error starting server: %v", err)
		os.Exit(1)
	}
}
