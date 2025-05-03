package main

import (
	"Advertising/configs"
	"fmt"
	"log"
)

func main() {
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatal("error loading config:", err)
	}
	fmt.Printf("Server is running on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
}
