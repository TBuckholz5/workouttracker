package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

func main() {
	// Read env config.
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

	db, err := sql.Open("pgx", fmt.Sprintf("host=%s user=%s port=%d dbname=%s sslmode=disable", dbHost, dbUser, dbPort, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}

	// Start server.
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})
	if err := r.Run(fmt.Sprintf(":%d", serverPort)); err != nil {
		log.Fatal(err)
	}
}
