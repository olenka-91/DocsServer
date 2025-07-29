package config

import "github.com/spf13/viper"

type Config struct {
	HTTPPort       string
	DBHost, DBName string
	DBPort, DBUser string
	DBPassword     string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := &Config{
		HTTPPort:   viper.GetString("APP_PORT"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBName:     viper.GetString("DB_NAME"),
		DBUser:     viper.GetString("DB_USERNAME"),
		DBPassword: viper.GetString("DB_PASSWORD"),
	}
	return cfg, nil
}
