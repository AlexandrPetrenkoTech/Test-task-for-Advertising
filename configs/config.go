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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
