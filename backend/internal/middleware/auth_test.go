package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWKSCache(t *testing.T) {
	cache := NewJWKSCache("https://example.com/.well-known/jwks.json")
	if cache == nil {
		t.Error("expected non-nil cache")
	}
	if cache.jwksURL != "https://example.com/.well-known/jwks.json" {
		t.Errorf("expected jwksURL, got %s", cache.jwksURL)
	}
}

func TestGetJWKS(t *testing.T) {
	defer ResetJWKSCache()
	cache := getJWKS()
	if cache == nil {
		t.Error("expected non-nil cache")
	}
}

func TestAuthMiddleware(t *testing.T) {
	// Generate mock RSA key pair for signing and verification
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	// Prepare public key parts in base64url encoding
	nBytes := privateKey.PublicKey.N.Bytes()
	nStr := base64.RawURLEncoding.EncodeToString(nBytes)

	eBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(eBytes, uint32(privateKey.PublicKey.E))
	start := 0
	for start < len(eBytes) && eBytes[start] == 0 {
		start++
	}
	eStr := base64.RawURLEncoding.EncodeToString(eBytes[start:])

	// Setup mock JWKS server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"keys": [
				{
					"kty": "RSA",
					"use": "sig",
					"alg": "RS256",
					"kid": "test-key-id",
					"n": "%s",
					"e": "%s"
				}
			]
		}`, nStr, eStr)
	}))
	defer ts.Close()

	// Configure environment variables
	os.Setenv("CLERK_JWKS_URL", ts.URL)
	os.Setenv("CLERK_SECRET_KEY", "test-secret-key")
	defer func() {
		os.Unsetenv("CLERK_JWKS_URL")
		os.Unsetenv("CLERK_SECRET_KEY")
		ResetJWKSCache()
	}()

	// Reset JWKS Cache before using it to ensure it picks up CLERK_JWKS_URL
	ResetJWKSCache()

	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		claims := GetClaims(c)
		return c.JSON(fiber.Map{
			"user_id": userID,
			"claims":  claims,
		})
	})

	// 1. Valid Token Test
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "https://clerk.com",
		"aud": "arishem",
		"sub": "user_12345",
		"org_id": "org_123",
		"role": "org:admin",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	token.Header["kid"] = "test-key-id"
	validTokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validTokenString)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed app test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Response body: %s", string(body))
	} else {
		var res struct {
			UserID string                 `json:"user_id"`
			Claims map[string]interface{} `json:"claims"`
		}
		json.NewDecoder(resp.Body).Decode(&res)
		if res.UserID != "user_12345" {
			t.Errorf("expected user_id user_12345, got %s", res.UserID)
		}
	}

	// 2. Missing Header Test
	reqMissing := httptest.NewRequest("GET", "/test", nil)
	respMissing, err := app.Test(reqMissing)
	if err != nil {
		t.Fatalf("failed app test: %v", err)
	}
	if respMissing.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for missing auth header, got %d", respMissing.StatusCode)
	}

	// 3. Invalid Authorization Format (no Bearer prefix)
	reqFormat := httptest.NewRequest("GET", "/test", nil)
	reqFormat.Header.Set("Authorization", "Invalid "+validTokenString)
	respFormat, err := app.Test(reqFormat)
	if err != nil {
		t.Fatalf("failed app test: %v", err)
	}
	if respFormat.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for invalid format, got %d", respFormat.StatusCode)
	}

	// 4. Invalid Signature Test (different key)
	otherKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	invalidTokenString, _ := token.SignedString(otherKey)

	reqInvalidSig := httptest.NewRequest("GET", "/test", nil)
	reqInvalidSig.Header.Set("Authorization", "Bearer "+invalidTokenString)
	respInvalidSig, err := app.Test(reqInvalidSig)
	if err != nil {
		t.Fatalf("failed app test: %v", err)
	}
	if respInvalidSig.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for invalid signature, got %d", respInvalidSig.StatusCode)
	}
}

func TestAuthMiddleware_ServerConfigError(t *testing.T) {
	os.Unsetenv("CLERK_SECRET_KEY")
	defer ResetJWKSCache()

	app := fiber.New()
	app.Use(AuthMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer dummy.token.here")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed app test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 when CLERK_SECRET_KEY is missing, got %d", resp.StatusCode)
	}
}