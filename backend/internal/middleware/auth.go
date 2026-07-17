package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ClerkUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
	}

	tokenString := parts[1]

	secretKey := os.Getenv("CLERK_SECRET_KEY")
	if secretKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "server configuration error"})
	}

	userInfo, err := validateClerkToken(secretKey, tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token: " + err.Error()})
	}

	c.Locals("user_id", userInfo.ID)
	c.Locals("user_email", userInfo.Email)
	c.Locals("user_first_name", userInfo.FirstName)
	c.Locals("user_last_name", userInfo.LastName)

	return c.Next()
}

func validateClerkToken(secretKey, token string) (*ClerkUserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.clerk.com/v1/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed with status: %d", resp.StatusCode)
	}

	var userInfo ClerkUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

func GetClaims(c *fiber.Ctx) map[string]interface{} {
	return map[string]interface{}{
		"user_id":     GetUserID(c),
		"user_email":  GetUserEmail(c),
		"first_name":  GetUserFirstName(c),
		"last_name":   GetUserLastName(c),
	}
}

func GetUserID(c *fiber.Ctx) string {
	userID, _ := c.Locals("user_id").(string)
	return userID
}

func GetUserEmail(c *fiber.Ctx) string {
	email, _ := c.Locals("user_email").(string)
	return email
}

func GetUserFirstName(c *fiber.Ctx) string {
	firstName, _ := c.Locals("user_first_name").(string)
	return firstName
}

func GetUserLastName(c *fiber.Ctx) string {
	lastName, _ := c.Locals("user_last_name").(string)
	return lastName
}