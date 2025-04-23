package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	handler "7-solutions-test-backend/internal/adapter/http"
	mongoadapter "7-solutions-test-backend/internal/adapter/mongo"
	"7-solutions-test-backend/internal/auth"
	"7-solutions-test-backend/internal/core/user"
	"7-solutions-test-backend/internal/task"

	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Initialize HTTP server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(handler.LoggingMiddleware)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("Authorization Header:", c.Request().Header.Get("Authorization"))
			return next(c)
		}
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	dbName := os.Getenv("MONGO_DB_NAME")
	db := client.Database(dbName)
	fmt.Println("Connected to MongoDB: " + dbName)

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	repo := mongoadapter.NewUserMongoRepo(db)
	service := user.NewService(repo)
	jwtService := auth.NewJWTService(jwtSecret)

	// Setup HTTP handlers
	h := handler.NewHandler(service, jwtService)
	h.RegisterRoutes(e)

	// Background task
	go task.StartUserCountLogger(ctx, repo)

	// Start server in goroutine
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for shutdown signal (Ctrl+C)
	<-ctx.Done()
	log.Println("🛑 Shutdown signal received...")

	// Shutdown logic
	stop() // Stop receiving more signals
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Stopping Echo server")
	if err := e.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Disconnecting from MongoDB")
	if err := client.Disconnect(ctxShutDown); err != nil {
		log.Fatalf("Mongo disconnect failed: %v", err)
	}

	log.Println("✅ Graceful shutdown complete")

}
