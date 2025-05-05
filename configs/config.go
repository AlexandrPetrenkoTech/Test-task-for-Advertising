package configs

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Host string
		Port int
	}
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
}

func LoadConfig(path string) (*Config, error) {
	// Set up Viper to read config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	// Load configuration from YAML file
	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal the YAML configuration into the Config struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Override values with environment variables from .env if present
	if viper.IsSet("PORT") {
		cfg.Server.Port = viper.GetInt("PORT")
	}
	if viper.IsSet("DATABASE_URL") {
		cfg.DB.Host = viper.GetString("DATABASE_URL")
	}

	return &cfg, nil
}
