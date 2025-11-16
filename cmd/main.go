package main

import (
	"context"
	"fmt"
	"log"

	userapi "github.com/TBuckholz5/workouttracker/internal/api/v1/user"
	"github.com/TBuckholz5/workouttracker/internal/config"
	userrepo "github.com/TBuckholz5/workouttracker/internal/repository/user"
	userservice "github.com/TBuckholz5/workouttracker/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	// Read env config.
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Connect to database.
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable",
		config.DBUser, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Run migration.
	db := stdlib.OpenDBFromPool(pool)
	defer func() { _ = db.Close() }()
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}

	// Start server.
	userRepository := userrepo.NewRepository(pool)
	userService := userservice.NewService(userRepository)
	userHandler := userapi.NewHandler(userService)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})
	apiV1 := r.Group("/api/v1")
	userapi.RegisterUserRoutes(apiV1, userHandler)
	if err := r.Run(fmt.Sprintf(":%d", config.ServerPort)); err != nil {
		log.Fatal(err)
	}
}
