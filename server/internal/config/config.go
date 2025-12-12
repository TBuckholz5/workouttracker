package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int
	ServerHost string
	JWTSecret  string
	DBUser     string
	DBPort     int
	DBName     string
	DBHost     string
	DBPassword string
	SslMode    string
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("DATABASE_PORT", 5432)
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_SSLMODE", "disable")

	viper.AutomaticEnv()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No .env file found, using environment variables")
		} else {
			return nil, err
		}
	}
	serverPort := viper.GetInt("SERVER_PORT")
	serverHost := viper.GetString("SERVER_HOST")
	jwtSecret := viper.GetString("JWT_SECRET")

	databasePort := viper.GetInt("DATABASE_PORT")
	databaseUser := viper.GetString("DATABASE_USER")
	databaseName := viper.GetString("DATABASE_NAME")
	databaseHost := viper.GetString("DATABASE_HOST")
	databasePassword := viper.GetString("DATABASE_PASSWORD")
	databaseSslMode := viper.GetString("DATABASE_SSLMODE")

	return &Config{
		ServerPort: serverPort,
		ServerHost: serverHost,
		JWTSecret:  jwtSecret,
		DBUser:     databaseUser,
		DBPort:     databasePort,
		DBName:     databaseName,
		DBHost:     databaseHost,
		DBPassword: databasePassword,
		SslMode:    databaseSslMode,
	}, nil
}
