package configs

import (
	"log"

	"github.com/spf13/viper"
)

// LoadEnvConfig loads .env variables and binds them to environment
func LoadEnvConfig() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found: %v", err)
	}
	viper.AutomaticEnv()
}
