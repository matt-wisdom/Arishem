package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"arishem/internal/api"
	"arishem/internal/db"
	"arishem/internal/jobs"
	"arishem/internal/middleware"
	"arishem/internal/reports"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log/slog"
	"strings"
)

func main() {
	loadEnv()

	// Initialize structured JSON logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	slog.Info("Starting Arishem API server...")

	if err := db.Init(); err != nil {
		slog.Error("Failed to initialize database", slog.Any("error", err))
	} else {
		// Auto-migrate schema: add logs and rerun parameters to llm_pentest_runs if not exists
		if pool := db.GetPool(); pool != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			_, err := pool.Exec(ctx, `
				ALTER TABLE llm_pentest_runs 
				ADD COLUMN IF NOT EXISTS logs TEXT DEFAULT '',
				ADD COLUMN IF NOT EXISTS docker BOOLEAN DEFAULT FALSE,
				ADD COLUMN IF NOT EXISTS config_mode VARCHAR(50) DEFAULT 'default',
				ADD COLUMN IF NOT EXISTS api_key TEXT DEFAULT '',
				ADD COLUMN IF NOT EXISTS model VARCHAR(100) DEFAULT '',
				ADD COLUMN IF NOT EXISTS llm_provider VARCHAR(50) DEFAULT '',
				ADD COLUMN IF NOT EXISTS api_base TEXT DEFAULT '',
				ADD COLUMN IF NOT EXISTS mode VARCHAR(50) DEFAULT 'both',
				ADD COLUMN IF NOT EXISTS budget INTEGER DEFAULT 8,
				ADD COLUMN IF NOT EXISTS concurrency INTEGER DEFAULT 4,
				ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1
			`)
			if err != nil {
				slog.Error("Failed to auto-upgrade DB schema for llm_pentest_runs", slog.Any("error", err))
			} else {
				slog.Info("Database schema verified and upgraded successfully")
			}
		}
	}

	if err := reports.InitS3(); err != nil {
		slog.Error("Failed to initialize S3", slog.Any("error", err))
	} else {
		slog.Info("S3 storage initialized successfully")
	}

	// Start background worker pool for jobs
	jobs.StartWorkerPool(context.Background())

	app := fiber.New(fiber.Config{
		AppName:        "Arishem",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		ReadBufferSize: 32768,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	webhookGroup := app.Group("/webhooks")
	webhookGroup.Post("/clerk", api.HandleClerkWebhook)

	apiGroup := app.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware)
	apiGroup.Use(middleware.TenantMiddleware)

	api.RegisterRoutes(apiGroup)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func loadEnv() {
	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
	}
}