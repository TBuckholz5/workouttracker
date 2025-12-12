package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"github.com/TBuckholz5/workouttracker/internal/config"
	exerciseApi "github.com/TBuckholz5/workouttracker/internal/domains/exercise/api/v1"
	exerciseRepo "github.com/TBuckholz5/workouttracker/internal/domains/exercise/repository"
	exerciseServ "github.com/TBuckholz5/workouttracker/internal/domains/exercise/service"
	userApi "github.com/TBuckholz5/workouttracker/internal/domains/user/api/v1"
	userRepo "github.com/TBuckholz5/workouttracker/internal/domains/user/repository"
	userServ "github.com/TBuckholz5/workouttracker/internal/domains/user/service"
	workoutSessionApi "github.com/TBuckholz5/workouttracker/internal/domains/workoutsession/api/v1"
	workoutSessionRepo "github.com/TBuckholz5/workouttracker/internal/domains/workoutsession/repository"
	workoutSessionServ "github.com/TBuckholz5/workouttracker/internal/domains/workoutsession/service"
	"github.com/TBuckholz5/workouttracker/internal/routing"
	"github.com/TBuckholz5/workouttracker/internal/routing/middleware"
	"github.com/TBuckholz5/workouttracker/internal/routing/middleware/auth"
	"github.com/TBuckholz5/workouttracker/internal/routing/middleware/logging"
	"github.com/TBuckholz5/workouttracker/internal/util/hash"
	"github.com/TBuckholz5/workouttracker/internal/util/jwt"
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
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName, config.SslMode))
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

	// Define dependencies.
	jwtService := jwt.NewJwtService(jwtSecret)
	authMiddleware := auth.NewAuthMiddleware(jwtService)
	loggingMiddleware := logging.NewLoggingMiddleware()

	userRepository := userRepo.NewRepository(pool)
	userService := userServ.NewService(userRepository, hash.NewBcryptHasher(), jwtService)
	userHandler := userApi.NewHandler(userService)

	// Register routes.
	mux := http.NewServeMux()

	apiMux := routing.RegisterRouterGroup(routing.Config{
		Mux:        mux,
		GroupRoute: "/api/v1/",
	})

	userMux := routing.RegisterRouterGroup(routing.Config{
		Mux:         apiMux,
		Middlewares: []middleware.Middleware{loggingMiddleware},
		GroupRoute:  "/user/",
	})
	routing.RegisterRoute(routing.Config{
		Mux:     userMux,
		Handler: http.HandlerFunc(userHandler.Register),
		Route:   "/register",
		Method:  "POST",
	})
	routing.RegisterRoute(routing.Config{
		Mux:     userMux,
		Handler: http.HandlerFunc(userHandler.Login),
		Route:   "/login",
		Method:  "POST",
	})

	exerciseRepository := exerciseRepo.NewRepository(pool)
	exerciseService := exerciseServ.NewService(exerciseRepository)
	exerciseHandler := exerciseApi.NewHandler(exerciseService)
	exerciseMux := routing.RegisterRouterGroup(routing.Config{
		Mux:         apiMux,
		Middlewares: []middleware.Middleware{loggingMiddleware, authMiddleware},
		GroupRoute:  "/exercise/",
	})
	routing.RegisterRoute(routing.Config{
		Mux:     exerciseMux,
		Handler: http.HandlerFunc(exerciseHandler.CreateExercise),
		Route:   "/create",
		Method:  "POST",
	})
	routing.RegisterRoute(routing.Config{
		Mux:     exerciseMux,
		Handler: http.HandlerFunc(exerciseHandler.GetExerciseForUser),
		Route:   "/getForUser",
		Method:  "GET",
	})

	workoutSessionRepository := workoutSessionRepo.NewRepository(pool)
	workoutSessionService := workoutSessionServ.NewService(workoutSessionRepository)
	workoutSessionHandler := workoutSessionApi.NewHandler(workoutSessionService)
	workoutSessionMux := routing.RegisterRouterGroup(routing.Config{
		Mux:         apiMux,
		Middlewares: []middleware.Middleware{loggingMiddleware, authMiddleware},
		GroupRoute:  "/workoutsession/",
	})
	routing.RegisterRoute(routing.Config{
		Mux:     workoutSessionMux,
		Handler: http.HandlerFunc(workoutSessionHandler.Create),
		Route:   "/create",
		Method:  "POST",
	})

	// Start server.
	fmt.Println("Starting server on port", config.ServerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), mux); err != nil {
		log.Fatal(err)
	}
}
