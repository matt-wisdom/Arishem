package api

import (
	"encoding/json"
	"net/http"
	"time"

	"arishem/internal/db"
	"arishem/internal/middleware"
	"arishem/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type GitHubConnectRequest struct {
	Token string `json:"token"`
}

type GitHubRepo struct {
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	Private    bool   `json:"private"`
	HTMLURL    string `json:"html_url"`
}

func ConnectGitHub(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	role := middleware.GetRole(c)
	if role != "org:admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only admins can manage integrations"})
	}

	var req GitHubConnectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	orgUUID, _ := uuid.Parse(orgID)
	integration := models.Integration{
		ID:                   uuid.New(),
		OrgID:                orgUUID,
		Provider:             models.IntegrationProviderGitHub,
		CredentialsEncrypted: req.Token,
		CreatedAt:            time.Now(),
	}

	_, err := db.GetPool().Exec(c.Context(), `
		INSERT INTO integrations (id, org_id, provider, credentials_encrypted, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (org_id, provider) DO UPDATE SET credentials_encrypted = $4
	`, integration.ID, integration.OrgID, integration.Provider, integration.CredentialsEncrypted)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save integration"})
	}

	return c.JSON(fiber.Map{"message": "GitHub connected successfully"})
}

func GetGitHubIntegration(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)

	var integration models.Integration
	err := db.GetPool().QueryRow(c.Context(), `
		SELECT id, org_id, provider, created_at FROM integrations WHERE org_id = $1 AND provider = $2
	`, orgUUID, models.IntegrationProviderGitHub).Scan(&integration.ID, &integration.OrgID, &integration.Provider, &integration.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "GitHub integration not found"})
	}

	return c.JSON(fiber.Map{"connected": true, "provider": integration.Provider})
}

func ListGitHubRepos(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)

	var token string
	err := db.GetPool().QueryRow(c.Context(), `
		SELECT credentials_encrypted FROM integrations WHERE org_id = $1 AND provider = $2
	`, orgUUID, models.IntegrationProviderGitHub).Scan(&token)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "GitHub integration not found"})
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(c.Context(), ts)

	req, _ := http.NewRequest("GET", "https://api.github.com/user/repos?sort=updated&per_page=50", nil)
	req = req.WithContext(c.Context())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := tc.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch repos"})
	}
	defer resp.Body.Close()

	var repos []struct {
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Private     bool   `json:"private"`
		HTMLURL     string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to parse repos"})
	}

	var repoList []GitHubRepo
	for _, repo := range repos {
		repoList = append(repoList, GitHubRepo{
			Name:     repo.Name,
			FullName: repo.FullName,
			Private:  repo.Private,
			HTMLURL:  repo.HTMLURL,
		})
	}

	if repoList == nil {
		repoList = []GitHubRepo{}
	}

	return c.JSON(repoList)
}

func DisconnectGitHub(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	role := middleware.GetRole(c)
	if role != "org:admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only admins can manage integrations"})
	}

	orgUUID, _ := uuid.Parse(orgID)
	_, err := db.GetPool().Exec(c.Context(), `
		DELETE FROM integrations WHERE org_id = $1 AND provider = $2
	`, orgUUID, models.IntegrationProviderGitHub)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to disconnect integration"})
	}

	return c.JSON(fiber.Map{"message": "GitHub disconnected successfully"})
}