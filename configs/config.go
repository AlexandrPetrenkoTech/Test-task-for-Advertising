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

// LoadConfig reads config.yaml and overrides with ENV
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Override with environment variables
	if viper.IsSet("PORT") {
		cfg.Server.Port = viper.GetInt("PORT")
	}
	if viper.IsSet("DB_HOST") {
		cfg.DB.Host = viper.GetString("DB_HOST")
	}
	if viper.IsSet("DB_PORT") {
		cfg.DB.Port = viper.GetInt("DB_PORT")
	}
	if viper.IsSet("DB_USER") {
		cfg.DB.User = viper.GetString("DB_USER")
	}
	if viper.IsSet("DB_PASSWORD") {
		cfg.DB.Password = viper.GetString("DB_PASSWORD")
	}
	if viper.IsSet("DB_NAME") {
		cfg.DB.Name = viper.GetString("DB_NAME")
	}

	return &cfg, nil
}
