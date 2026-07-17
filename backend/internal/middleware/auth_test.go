package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	os.Setenv("CLERK_SECRET_KEY", "test-secret-key")
	defer os.Unsetenv("CLERK_SECRET_KEY")

	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		return c.JSON(fiber.Map{"user_id": userID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed app test: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for missing auth header, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	os.Setenv("CLERK_SECRET_KEY", "test-secret-key")
	defer os.Unsetenv("CLERK_SECRET_KEY")

	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Invalid token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed app test: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for invalid format, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_ServerConfigError(t *testing.T) {
	os.Unsetenv("CLERK_SECRET_KEY")

	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed app test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 when CLERK_SECRET_KEY is missing, got %d", resp.StatusCode)
	}
}