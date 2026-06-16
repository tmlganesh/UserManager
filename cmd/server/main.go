package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/ganesh/ainyx/configs"
	"github.com/ganesh/ainyx/db/sqlc"
	"github.com/ganesh/ainyx/internal/handler"
	"github.com/ganesh/ainyx/internal/middleware"
	"github.com/ganesh/ainyx/internal/repository"
	"github.com/ganesh/ainyx/internal/service"
	"github.com/ganesh/ainyx/internal/validator"
)

func main() {
	// ── Configuration ────────────────────────────────────────────────
	cfg := configs.Load()

	// ── Logger ───────────────────────────────────────────────────────
	logger, err := initLogger(cfg.AppEnv)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("starting ainyx server",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.Port),
	)

	// ── Database ─────────────────────────────────────────────────────
	pool, err := connectDB(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// ── Run migrations ───────────────────────────────────────────────
	if err := runMigrations(pool, logger); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	// ── Dependency Injection ─────────────────────────────────────────
	queries := sqlc.New(pool)
	userRepo := repository.NewUserRepository(queries)
	userSvc := service.NewUserService(userRepo, logger)
	userValidator := validator.New()
	userHandler := handler.NewUserHandler(userSvc, userValidator, logger)

	// ── Fiber App ────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      "ainyx",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	// Global middleware.
	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(logger))

	// Health check.
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Register routes.
	userHandler.RegisterRoutes(app)

	// ── Graceful Shutdown ────────────────────────────────────────────
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			logger.Fatal("server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}
	logger.Info("server stopped")
}

// initLogger creates a Zap logger configured for the given environment.
// Production uses JSON output; development uses colored console output.
func initLogger(env string) (*zap.Logger, error) {
	if env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}

// connectDB establishes a connection pool to PostgreSQL.
func connectDB(databaseURL string, logger *zap.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	logger.Info("database connection established")
	return pool, nil
}

// runMigrations applies the schema migration directly.
// In production, you would use a migration tool like golang-migrate.
func runMigrations(pool *pgxpool.Pool, logger *zap.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	migration := `
	CREATE TABLE IF NOT EXISTS users (
		id   SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		dob  DATE NOT NULL
	);`

	_, err := pool.Exec(ctx, migration)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	logger.Info("database migrations applied")
	return nil
}
