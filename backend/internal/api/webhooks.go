package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"arishem/internal/db"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClerkEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ClerkOrganization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClerkUser struct {
	ID             string `json:"id"`
	EmailAddresses []struct {
		EmailAddress string `json:"email_address"`
	} `json:"email_addresses"`
}

type ClerkMembership struct {
	PublicUserData struct {
		UserID string `json:"user_id"`
	} `json:"public_user_data"`
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
}

func verifyClerkSignature(payload []byte, signature string) bool {
	secret := os.Getenv("CLERK_WEBHOOK_SECRET")
	if secret == "" {
		log.Println("Warning: CLERK_WEBHOOK_SECRET not set")
		return true
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSig := "v1=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedSig), []byte(signature))
}

func HandleClerkWebhook(c *fiber.Ctx) error {
	signature := c.Get("svix-signature")
	timestamp := c.Get("svix-timestamp")

	if signature == "" || timestamp == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing signature"})
	}

	body := c.Body()
	if !verifyClerkSignature(body, signature) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid signature"})
	}

	var event ClerkEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	switch event.Type {
	case "organization.created":
		return handleOrganizationCreated(c, event.Data)
	case "organizationMembership.created":
		return handleMembershipCreated(c, event.Data)
	case "organizationMembership.deleted":
		return handleMembershipDeleted(c, event.Data)
	case "user.deleted":
		return handleUserDeleted(c, event.Data)
	default:
		log.Printf("Unhandled Clerk event type: %s", event.Type)
	}

	return c.SendStatus(http.StatusOK)
}

func handleOrganizationCreated(c *fiber.Ctx, data json.RawMessage) error {
	var org ClerkOrganization
	if err := json.Unmarshal(data, &org); err != nil {
		return err
	}

	orgUUID := uuid.New()
	_, err := db.GetPool().Exec(c.Context(), `
		INSERT INTO organizations (id, clerk_org_id, name, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (clerk_org_id) DO NOTHING
	`, orgUUID, org.ID, org.Name, time.Now())

	if err != nil {
		log.Printf("Failed to create organization: %v", err)
		return err
	}

	log.Printf("Created organization: %s (%s)", org.Name, org.ID)
	return nil
}

func handleMembershipCreated(c *fiber.Ctx, data json.RawMessage) error {
	var membership ClerkMembership
	if err := json.Unmarshal(data, &membership); err != nil {
		return err
	}

	orgID, _ := uuid.Parse(membership.OrganizationID)
	userUUID := uuid.New()

	role := "org:viewer"
	if membership.Role == "org_admin" {
		role = "org:admin"
	} else if membership.Role == "org_member" {
		role = "org:engineer"
	}

	_, err := db.GetPool().Exec(c.Context(), `
		INSERT INTO users (id, clerk_user_id, org_id, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (clerk_user_id, org_id) DO NOTHING
	`, userUUID, membership.PublicUserData.UserID, orgID, role, time.Now())

	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return err
	}

	log.Printf("Created membership: user=%s org=%s role=%s", membership.PublicUserData.UserID, membership.OrganizationID, role)
	return nil
}

func handleMembershipDeleted(c *fiber.Ctx, data json.RawMessage) error {
	var membership ClerkMembership
	if err := json.Unmarshal(data, &membership); err != nil {
		return err
	}

	orgID, _ := uuid.Parse(membership.OrganizationID)

	_, err := db.GetPool().Exec(c.Context(), `
		DELETE FROM users WHERE clerk_user_id = $1 AND org_id = $2
	`, membership.PublicUserData.UserID, orgID)

	if err != nil {
		log.Printf("Failed to delete user membership: %v", err)
		return err
	}

	log.Printf("Deleted membership: user=%s org=%s", membership.PublicUserData.UserID, membership.OrganizationID)
	return nil
}

func handleUserDeleted(c *fiber.Ctx, data json.RawMessage) error {
	var user ClerkUser
	if err := json.Unmarshal(data, &user); err != nil {
		return err
	}

	_, err := db.GetPool().Exec(c.Context(), `
		DELETE FROM users WHERE clerk_user_id = $1
	`, user.ID)

	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return err
	}

	log.Printf("Deleted user: %s", user.ID)
	return nil
}