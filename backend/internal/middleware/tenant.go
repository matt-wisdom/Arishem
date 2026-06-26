package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func TenantMiddleware(c *fiber.Ctx) error {
	claims := GetClaims(c)

	orgID, ok := claims["org_id"].(string)
	if !ok || orgID == "" {
		orgs, ok := claims["orgs"].(map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no organization found"})
		}
		for id := range orgs {
			orgID = id
			break
		}
	}

	if orgID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no organization in token"})
	}

	c.Locals("org_id", orgID)

	role := "org:viewer"
	if roleClaim, ok := claims["role"].(string); ok {
		role = roleClaim
	} else if orgs, ok := claims["orgs"].(map[string]interface{}); ok {
		if orgData, ok := orgs[orgID].(map[string]interface{}); ok {
			if r, ok := orgData["role"].(string); ok {
				role = r
			}
		}
	}

	c.Locals("role", role)

	return c.Next()
}

func GetOrgID(c *fiber.Ctx) string {
	orgID, _ := c.Locals("org_id").(string)
	return orgID
}

func GetRole(c *fiber.Ctx) string {
	role, _ := c.Locals("role").(string)
	return role
}

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := GetRole(c)
		for _, r := range allowedRoles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
	}
}