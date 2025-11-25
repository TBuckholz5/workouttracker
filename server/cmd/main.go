package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	exerciseApi "github.com/TBuckholz5/workouttracker/internal/api/v1/exercise"
	userApi "github.com/TBuckholz5/workouttracker/internal/api/v1/user"
	"github.com/TBuckholz5/workouttracker/internal/config"
	"github.com/TBuckholz5/workouttracker/internal/middleware/auth"
	exerciseRepo "github.com/TBuckholz5/workouttracker/internal/repository/exercise"
	userRepo "github.com/TBuckholz5/workouttracker/internal/repository/user"
	exerciseServ "github.com/TBuckholz5/workouttracker/internal/service/exercise"
	userServ "github.com/TBuckholz5/workouttracker/internal/service/user"
	"github.com/TBuckholz5/workouttracker/internal/util/hash"
	"github.com/TBuckholz5/workouttracker/internal/util/jwt"
	workoutSessionApi "github.com/TBuckholz5/workouttracker/internal/workoutsession/api/v1"
	workoutSessionRepo "github.com/TBuckholz5/workouttracker/internal/workoutsession/repository"
	workoutSessionServ "github.com/TBuckholz5/workouttracker/internal/workoutsession/service"
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

	// Generate JWT secret.
	// TODO: Store in a better way than in memory - means server restarts will invalidate all tokens.
	jwtSecret := make([]byte, 32)
	_, err = rand.Read(jwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	// Start server.
	jwtService := jwt.NewJwtService(jwtSecret)
	authMiddleware := auth.AuthMiddleware(jwtService)
	r := gin.Default()
	userRepository := userRepo.NewRepository(pool)
	userService := userServ.NewService(userRepository, hash.NewBcryptHasher(), jwtService)
	userHandler := userApi.NewHandler(userService)
	apiV1 := r.Group("/api/v1")
	userApi.RegisterUserRoutes(apiV1, userHandler)
	exerciseRepository := exerciseRepo.NewRepository(pool)
	exerciseService := exerciseServ.NewService(exerciseRepository)
	exerciseHandler := exerciseApi.NewHandler(exerciseService)
	exerciseApi.RegisterExerciseRoutes(apiV1, exerciseHandler, authMiddleware)
	workoutSessionRepository := workoutSessionRepo.NewRepository(pool)
	workoutSessionService := workoutSessionServ.NewService(workoutSessionRepository)
	workoutSessionHandler := workoutSessionApi.NewHandler(workoutSessionService)
	workoutSessionApi.RegisterWorkoutSessionRoutes(apiV1, workoutSessionHandler, authMiddleware)
	if err := r.Run(fmt.Sprintf(":%d", config.ServerPort)); err != nil {
		log.Fatal(err)
	}
}
