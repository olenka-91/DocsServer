package config

import "github.com/spf13/viper"

type Config struct {
	HTTPPort    string
	DBHost      string
	DBPort      string
	DBUsername  string
	DBPassword  string
	DBName      string
	SSLMode     string
	StorageAddr string
}

const (
	defaultStorageAddr = "./storage"
)

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := &Config{
		HTTPPort:    viper.GetString("APP_PORT"),
		DBHost:      viper.GetString("DB_HOST"),
		DBPort:      viper.GetString("DB_PORT"),
		DBUsername:  viper.GetString("DB_USERNAME"),
		DBPassword:  viper.GetString("DB_PASSWORD"),
		DBName:      viper.GetString("DB_NAME"),
		SSLMode:     viper.GetString("DB_SSLMODE"),
		StorageAddr: defaultStorageAddr,
	}
	return cfg, nil
}
