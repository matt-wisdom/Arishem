package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"arishem/internal/api"
	"arishem/internal/db"
	"arishem/internal/middleware"
	"arishem/internal/reports"
)

func main() {
	if err := db.Init(); err != nil {
		log.Printf("Warning: failed to initialize database: %v", err)
	}

	if err := reports.InitS3(); err != nil {
		log.Printf("Warning: failed to initialize S3: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "Arishem",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
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