package middleware

import (
	"context"
	"strings"

	"arishem/internal/db"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ResolveOrCreateOrg(ctx context.Context, clerkOrgID string) (uuid.UUID, error) {
	pool := db.GetPool()
	if pool == nil {
		if u, err := uuid.Parse(clerkOrgID); err == nil {
			return u, nil
		}
		return uuid.NewMD5(uuid.NameSpaceDNS, []byte(clerkOrgID)), nil
	}
	var dbOrgID uuid.UUID

	err := pool.QueryRow(ctx, "SELECT id FROM organizations WHERE clerk_org_id = $1", clerkOrgID).Scan(&dbOrgID)
	if err == nil {
		return dbOrgID, nil
	}

	dbOrgID = uuid.New()
	name := "Default Workspace"
	if strings.HasPrefix(clerkOrgID, "user_") {
		name = "Personal Workspace"
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO organizations (id, clerk_org_id, name)
		VALUES ($1, $2, $3)
		ON CONFLICT (clerk_org_id) DO NOTHING
	`, dbOrgID, clerkOrgID, name)
	if err != nil {
		return uuid.Nil, err
	}

	err = pool.QueryRow(ctx, "SELECT id FROM organizations WHERE clerk_org_id = $1", clerkOrgID).Scan(&dbOrgID)
	return dbOrgID, err
}

func TenantMiddleware(c *fiber.Ctx) error {
	userID := GetUserID(c)

	if userID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no user in token"})
	}

	dbOrgID, err := ResolveOrCreateOrg(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to resolve organization"})
	}

	c.Locals("org_id", dbOrgID.String())
	c.Locals("role", "org:admin")

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