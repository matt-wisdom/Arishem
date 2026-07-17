package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ClerkVerifyResponse struct {
	Sub string `json:"sub"`
}

var clerkHTTPClient = &http.Client{Timeout: 15 * time.Second}

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

	userID, err := validateClerkToken(secretKey, tokenString)
	if err != nil {
		fmt.Printf("Clerk auth error: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token: " + err.Error()})
	}

	c.Locals("user_id", userID)

	return c.Next()
}

func validateClerkToken(secretKey, token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.clerk.com/v1/tokens/verify", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := clerkHTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token validation failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	var verifyResp ClerkVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return "", fmt.Errorf("failed to decode verify response: %w", err)
	}

	if verifyResp.Sub == "" {
		return "", fmt.Errorf("no subject in token response")
	}

	return verifyResp.Sub, nil
}

func GetClaims(c *fiber.Ctx) map[string]interface{} {
	return map[string]interface{}{
		"user_id": GetUserID(c),
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