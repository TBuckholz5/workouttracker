package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int
	ServerHost string
	DBUser     string
	DBPort     int
	DBName     string
	DBHost     string
	DBPassword string
	SslMode    string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	serverPort := viper.GetInt("SERVER_PORT")
	serverHost := viper.GetString("SERVER_HOST")

	databasePort := viper.GetInt("DATABASE_PORT")
	databaseUser := viper.GetString("DATABASE_USER")
	databaseName := viper.GetString("DATABASE_NAME")
	databaseHost := viper.GetString("DATABASE_HOST")
	databasePassword := viper.GetString("DATABASE_PASSWORD")
	databaseSslMode := viper.GetString("DATABASE_SSLMODE")

	return &Config{
		ServerPort: serverPort,
		ServerHost: serverHost,
		DBUser:     databaseUser,
		DBPort:     databasePort,
		DBName:     databaseName,
		DBHost:     databaseHost,
		DBPassword: databasePassword,
		SslMode:    databaseSslMode,
	}, nil
}
