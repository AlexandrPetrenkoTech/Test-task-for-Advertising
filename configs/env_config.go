package configs

import (
	"github.com/spf13/viper"
	"log"
)

// LoadEnvConfig loads environment variables from .env file and automatically binds them.
func LoadEnvConfig() {
	// Set the config file to .env
	viper.SetConfigFile(".env")

	// Read the .env file if exists
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading .env file: %v", err)
	}

	// Automatically bind environment variables
	viper.AutomaticEnv()
}
