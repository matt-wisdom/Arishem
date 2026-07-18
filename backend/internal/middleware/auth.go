package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWKSCache struct {
	mu       sync.RWMutex
	keys     map[string]interface{}
	jwksURL  string
	client   *http.Client
}

func NewJWKSCache(jwksURL string) *JWKSCache {
	return &JWKSCache{
		keys:    make(map[string]interface{}),
		jwksURL: jwksURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *JWKSCache) GetKey(kid string) (interface{}, error) {
	c.mu.RLock()
	if key, ok := c.keys[kid]; ok {
		c.mu.RUnlock()
		return key, nil
	}
	c.mu.RUnlock()

	resp, err := c.client.Get(c.jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			Use string `json:"use"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" {
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				continue
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				continue
			}
			var eVal int
			if len(eBytes) < 4 {
				padded := make([]byte, 4)
				copy(padded[4-len(eBytes):], eBytes)
				eVal = int(binary.BigEndian.Uint32(padded))
			} else {
				eVal = int(binary.BigEndian.Uint32(eBytes))
			}
			pubKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: eVal,
			}
			c.keys[key.Kid] = pubKey
		} else {
			c.keys[key.Kid] = key
		}
	}

	if key, ok := c.keys[kid]; ok {
		return key, nil
	}
	return nil, fmt.Errorf("key not found: %s", kid)
}

var (
	jwksCache     *JWKSCache
	jwksCacheOnce sync.Once
)

func getJWKS() *JWKSCache {
	jwksCacheOnce.Do(func() {
		jwksURL := os.Getenv("CLERK_JWKS_URL")
		if jwksURL == "" {
			jwksURL = "https://clerk.com/.well-known/jwks.json"
		}
		jwksCache = NewJWKSCache(jwksURL)
	})
	return jwksCache
}

func getJWKSForIssuer(issuer string) *JWKSCache {
	jwksURL := issuer + "/.well-known/jwks.json"
	return NewJWKSCache(jwksURL)
}

func ResetJWKSCache() {
	jwksCache = nil
	jwksCacheOnce = sync.Once{}
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

	claims, err := validateClerkTokenJWT(secretKey, tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token: " + err.Error()})
	}

	c.Locals("user_id", claims["sub"])
	
	// Try to get email from claims
	email := ""
	if e, ok := claims["email"].(string); ok {
		email = e
	} else if e, ok := claims["email_address"].(string); ok {
		email = e
	}
	
	// If no email in token, fetch from Clerk API
	if email == "" {
		userID, _ := claims["sub"].(string)
		if userID != "" {
			email = fetchUserEmailFromClerk(userID)
		}
	}
	
	if email != "" {
		c.Locals("user_email", email)
	}
	if fn, ok := claims["first_name"].(string); ok {
		c.Locals("user_first_name", fn)
	}
	if ln, ok := claims["last_name"].(string); ok {
		c.Locals("user_last_name", ln)
	}

	return c.Next()
}

func validateClerkTokenJWT(secretKey, tokenString string) (map[string]interface{}, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	iss, ok := claims["iss"].(string)
	if !ok || iss == "" {
		return nil, fmt.Errorf("missing issuer in token")
	}

	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, fmt.Errorf("missing kid in token header")
	}

	jwks := getJWKSForIssuer(iss)
	pubKey, err := jwks.GetKey(kid)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("JWT signature validation failed: %v", err)
	}

	finalClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims after validation")
	}

	return finalClaims, nil
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

func fetchUserEmailFromClerk(userID string) string {
	clerkSecret := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecret == "" {
		return ""
	}
	
	url := "https://api.clerk.com/v1/users/" + userID
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+clerkSecret)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to fetch Clerk user: %v", err)
		return ""
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		log.Printf("Clerk API returned status %d", resp.StatusCode)
		return ""
	}
	
	var userData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return ""
	}
	
	// Try email_addresses array
	if emails, ok := userData["email_addresses"].([]interface{}); ok && len(emails) > 0 {
		if emailObj, ok := emails[0].(map[string]interface{}); ok {
			if email, ok := emailObj["email_address"].(string); ok {
				return email
			}
		}
	}
	
	// Try primary_email_address_id
	if primaryID, ok := userData["primary_email_address_id"].(string); ok && primaryID != "" {
		for _, e := range userData["email_addresses"].([]interface{}) {
			if emailObj, ok := e.(map[string]interface{}); ok {
				if id, ok := emailObj["id"].(string); ok && id == primaryID {
					if email, ok := emailObj["email_address"].(string); ok {
						return email
					}
				}
			}
		}
	}
	
	return ""
}