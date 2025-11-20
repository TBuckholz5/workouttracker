package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int
	DBUser     string
	DBPort     int
	DBName     string
	DBHost     string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("env")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	serverPort := viper.GetInt("server.port")
	dbUser := viper.GetString("database.user")
	dbPort := viper.GetInt("database.port")
	dbName := viper.GetString("database.name")
	dbHost := viper.GetString("database.host")

	return &Config{
		ServerPort: serverPort,
		DBUser:     dbUser,
		DBPort:     dbPort,
		DBName:     dbName,
		DBHost:     dbHost,
	}, nil
}
