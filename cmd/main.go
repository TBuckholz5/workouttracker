package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/TBuckholz5/workouttracker/internal/config"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	// Read env config.
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	db, err := sql.Open("pgx", fmt.Sprintf("host=%s user=%s port=%d dbname=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPort, config.DBName))
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
	if err := r.Run(fmt.Sprintf(":%d", config.ServerPort)); err != nil {
		log.Fatal(err)
	}
}
