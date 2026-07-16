package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGetOrgID(t *testing.T) {
	app := fiber.New()
	
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("org_id", "org_123")
		orgID := GetOrgID(c)
		if orgID != "org_123" {
			t.Errorf("expected org_123, got %s", orgID)
		}
		return c.SendStatus(200)
	})
	
	req := httptest.NewRequest("GET", "/", nil)
	_, err := app.Test(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetRole(t *testing.T) {
	app := fiber.New()
	
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "org:engineer")
		role := GetRole(c)
		if role != "org:engineer" {
			t.Errorf("expected org:engineer, got %s", role)
		}
		return c.SendStatus(200)
	})
	
	req := httptest.NewRequest("GET", "/", nil)
	_, err := app.Test(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRequireRoleAllowsAdmin(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "org:admin")
		return c.Next()
	})
	app.Use(RequireRole("org:admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRequireRoleDeniesViewer(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "org:viewer")
		return c.Next()
	})
	app.Use(RequireRole("org:admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != 403 {
		t.Errorf("expected 403, got %d", resp.StatusCode)
	}
}



// Let's rewrite the above test cleanly with individual apps to avoid route interference.

func TestTenantMiddleware_DirectClaims(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("claims", jwt.MapClaims{
			"org_id": "org_abc",
			"role":   "org:admin",
		})
		return c.Next()
	})
	app.Use(TenantMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"org_id": GetOrgID(c),
			"role":   GetRole(c),
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var res map[string]string
	json.NewDecoder(resp.Body).Decode(&res)
	if _, err := uuid.Parse(res["org_id"]); err != nil {
		t.Errorf("expected org_id to be a valid UUID, got %s (error: %v)", res["org_id"], err)
	}
	if res["role"] != "org:admin" {
		t.Errorf("expected role org:admin, got %s", res["role"])
	}
}

func TestTenantMiddleware_OrgsMapClaims(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("claims", jwt.MapClaims{
			"orgs": map[string]interface{}{
				"org_xyz": map[string]interface{}{
					"role": "org:engineer",
				},
			},
		})
		return c.Next()
	})
	app.Use(TenantMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"org_id": GetOrgID(c),
			"role":   GetRole(c),
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var res map[string]string
	json.NewDecoder(resp.Body).Decode(&res)
	if _, err := uuid.Parse(res["org_id"]); err != nil {
		t.Errorf("expected org_id to be a valid UUID, got %s (error: %v)", res["org_id"], err)
	}
	if res["role"] != "org:engineer" {
		t.Errorf("expected role org:engineer, got %s", res["role"])
	}
}

func TestTenantMiddlewareNoOrgReturns403(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("claims", jwt.MapClaims{})
		return c.Next()
	})
	app.Use(TenantMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 403 {
		t.Errorf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestGetClaims(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("claims", jwt.MapClaims{"foo": "bar"})
		claims := GetClaims(c)
		if claims["foo"] != "bar" {
			t.Errorf("expected claim foo to be bar, got %v", claims["foo"])
		}
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	_, err := app.Test(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetUserID(t *testing.T) {
	app := fiber.New()
	
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", "user_123")
		userID := GetUserID(c)
		if userID != "user_123" {
			t.Errorf("expected user_123, got %s", userID)
		}
		return c.SendStatus(200)
	})
	
	req := httptest.NewRequest("GET", "/", nil)
	_, err := app.Test(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}