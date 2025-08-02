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
	AdminToken  string
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
		DBName:      viper.GetString("DB_NAME"),
		DBUsername:  viper.GetString("DB_USERNAME"),
		DBPassword:  viper.GetString("DB_PASSWORD"),
		SSLMode:     viper.GetString("DB_SSLMODE"),
		StorageAddr: defaultStorageAddr,
		AdminToken:  viper.GetString("ADMIN_TOKEN"),
	}
	return cfg, nil
}
