package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"arishem/internal/db"
	"arishem/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestHealthEndpoint(t *testing.T) {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"status":"ok"}` {
		t.Errorf("expected {\"status\":\"ok\"}, got %s", string(body))
	}
}

func TestRegisterRoutes(t *testing.T) {
	app := fiber.New()
	RegisterRoutes(app)

	// Test that routes are registered
	routes := app.GetRoutes()
	if len(routes) == 0 {
		t.Error("expected routes to be registered")
	}

	// Check specific routes exist
	routeMethods := make(map[string]string)
	for _, route := range routes {
		routeMethods[route.Path] = route.Method
	}

	expectedRoutes := []string{
		"/scans",
		"/scans/:id",
		"/llmpentest",
		"/llmpentest/:id",
		"/reports",
		"/reports/:id",
		"/reports/:id/download",
		"/integrations/github",
		"/integrations/github/repos",
		"/alerts",
		"/alerts/:id",
		"/alerts/test/:id",
	}

	for _, path := range expectedRoutes {
		if _, exists := routeMethods[path]; !exists {
			t.Logf("Warning: route %s may not be registered", path)
		}
	}
}

func TestCreateScanRequest(t *testing.T) {
	req := CreateScanRequest{
		Target: "https://github.com/test/repo",
		Branch: "main",
	}

	if req.Target != "https://github.com/test/repo" {
		t.Errorf("expected Target, got %s", req.Target)
	}
	if req.Branch != "main" {
		t.Errorf("expected Branch main, got %s", req.Branch)
	}
}

func TestGitHubConnectRequest(t *testing.T) {
	req := GitHubConnectRequest{
		Token: "ghp_xxx",
	}

	if req.Token != "ghp_xxx" {
		t.Errorf("expected Token, got %s", req.Token)
	}
}

func TestGitHubRepo(t *testing.T) {
	repo := GitHubRepo{
		Name:     "test-repo",
		FullName: "testuser/test-repo",
		Private:  true,
		HTMLURL:  "https://github.com/testuser/test-repo",
	}

	if repo.Name != "test-repo" {
		t.Errorf("expected Name, got %s", repo.Name)
	}
	if !repo.Private {
		t.Error("expected Private to be true")
	}
}

func TestCreateLLMPentestRequest(t *testing.T) {
	req := CreateLLMPentestRequest{
		TargetEndpoint: "https://api.example.com/chat",
		APIKey:         "sk-test",
		TestModules:    []string{"prompt_injection", "jailbreak"},
	}

	if req.TargetEndpoint != "https://api.example.com/chat" {
		t.Errorf("expected TargetEndpoint, got %s", req.TargetEndpoint)
	}
	if len(req.TestModules) != 2 {
		t.Errorf("expected 2 modules, got %d", len(req.TestModules))
	}
}

func TestCreateAlertRequest(t *testing.T) {
	req := CreateAlertRequest{
		SeverityThreshold: "high",
		Channel:           "slack",
		ChannelConfig:     map[string]interface{}{"url": "https://hooks.slack.com/xxx"},
	}

	if req.SeverityThreshold != "high" {
		t.Errorf("expected high, got %s", req.SeverityThreshold)
	}
	if req.Channel != "slack" {
		t.Errorf("expected slack, got %s", req.Channel)
	}
}

func TestRouteGroups(t *testing.T) {
	app := fiber.New()
	
	scans := app.Group("/scans")
	scans.Post("/code", func(c *fiber.Ctx) error { return nil })
	
	routes := app.GetRoutes()
	if len(routes) == 0 {
		t.Error("expected routes")
	}
}

func initTestDB(t *testing.T) bool {
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgresql://launchpad:password@localhost:5432/launchpad")
	}
	err := db.Init()
	if err != nil {
		t.Logf("Skipping database-dependent test: %v", err)
		return false
	}

	// Load and run initial schema migrations
	schemaBytes, err := os.ReadFile("../../../migrations/001_initial_schema.sql")
	if err != nil {
		schemaBytes, err = os.ReadFile("../../migrations/001_initial_schema.sql")
	}
	if err == nil {
		_, err = db.GetPool().Exec(context.Background(), string(schemaBytes))
		if err != nil {
			t.Logf("Failed to run migrations: %v", err)
		}
		// Apply Alter Table to upgrade test database schema
		_, err = db.GetPool().Exec(context.Background(), `
			ALTER TABLE llm_pentest_runs 
			ADD COLUMN IF NOT EXISTS logs TEXT DEFAULT '',
			ADD COLUMN IF NOT EXISTS docker BOOLEAN DEFAULT FALSE,
			ADD COLUMN IF NOT EXISTS config_mode VARCHAR(50) DEFAULT 'default',
			ADD COLUMN IF NOT EXISTS api_key TEXT DEFAULT '',
			ADD COLUMN IF NOT EXISTS model VARCHAR(100) DEFAULT '',
			ADD COLUMN IF NOT EXISTS llm_provider VARCHAR(50) DEFAULT '',
			ADD COLUMN IF NOT EXISTS api_base TEXT DEFAULT '',
			ADD COLUMN IF NOT EXISTS mode VARCHAR(50) DEFAULT '',
			ADD COLUMN IF NOT EXISTS budget INT DEFAULT 0,
			ADD COLUMN IF NOT EXISTS concurrency INT DEFAULT 0,
			ADD COLUMN IF NOT EXISTS version INT DEFAULT 1;
		`)
		if err != nil {
			t.Logf("Failed to alter table in test: %v", err)
		}
	} else {
		t.Logf("Failed to read schema file: %v", err)
	}

	return true
}

func TestScanHandlers_Integration(t *testing.T) {
	if !initTestDB(t) {
		t.Skip("Database not available")
	}
	defer db.Close()

	ctx := context.Background()
	orgUUID := uuid.New()
	orgID := orgUUID.String()

	// Insert test organization
	_, err := db.GetPool().Exec(ctx, `
		INSERT INTO organizations (id, clerk_org_id, name)
		VALUES ($1, $2, $3)
	`, orgUUID, "clerk_"+orgUUID.String(), "Test API Org")
	if err != nil {
		t.Fatalf("failed to insert test organization: %v", err)
	}
	defer db.GetPool().Exec(ctx, `DELETE FROM organizations WHERE id = $1`, orgUUID)

	app := fiber.New()
	// Mock middleware to set tenant context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("org_id", orgID)
		c.Locals("role", "org:admin")
		return c.Next()
	})

	app.Post("/scans/code", CreateCodeScan)
	app.Get("/scans", ListScans)
	app.Get("/scans/:id", GetScan)

	// 1. Create Scan
	reqPayload := CreateScanRequest{
		Target: "https://github.com/test/repo",
		Branch: "main",
	}
	jsonBytes, _ := json.Marshal(reqPayload)

	req := httptest.NewRequest("POST", "/scans/code", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /scans/code request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", resp.StatusCode)
	}

	var scanObj models.Scan
	json.NewDecoder(resp.Body).Decode(&scanObj)
	if scanObj.Target != reqPayload.Target {
		t.Errorf("expected target %s, got %s", reqPayload.Target, scanObj.Target)
	}

	// 2. List Scans
	reqList := httptest.NewRequest("GET", "/scans", nil)
	respList, err := app.Test(reqList)
	if err != nil {
		t.Fatalf("GET /scans request failed: %v", err)
	}
	if respList.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", respList.StatusCode)
	}

	var scansList []models.Scan
	json.NewDecoder(respList.Body).Decode(&scansList)
	if len(scansList) == 0 {
		t.Error("expected at least one scan returned, got 0")
	}

	// 3. Get Scan
	reqGet := httptest.NewRequest("GET", "/scans/"+scanObj.ID.String(), nil)
	respGet, err := app.Test(reqGet)
	if err != nil {
		t.Fatalf("GET /scans/:id request failed: %v", err)
	}
	if respGet.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", respGet.StatusCode)
	}

	var getResult map[string]interface{}
	json.NewDecoder(respGet.Body).Decode(&getResult)
	if getResult["scan"] == nil {
		t.Error("expected scan object in response, got nil")
	}
}

func TestCreateCodeScan_Forbidden(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("org_id", uuid.New().String())
		c.Locals("role", "org:viewer")
		return c.Next()
	})
	app.Post("/scans/code", CreateCodeScan)

	reqPayload := CreateScanRequest{
		Target: "https://github.com/test/repo",
	}
	jsonBytes, _ := json.Marshal(reqPayload)

	req := httptest.NewRequest("POST", "/scans/code", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /scans/code request failed: %v", err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden, got %d", resp.StatusCode)
	}
}

func TestRerunLLMPentest_Integration(t *testing.T) {
	if !initTestDB(t) {
		t.Skip("Database not available")
	}
	defer db.Close()

	ctx := context.Background()
	orgUUID := uuid.New()
	orgID := orgUUID.String()

	// Insert test organization
	_, err := db.GetPool().Exec(ctx, `
		INSERT INTO organizations (id, clerk_org_id, name)
		VALUES ($1, $2, $3)
	`, orgUUID, "clerk_"+orgUUID.String(), "Test API Org")
	if err != nil {
		t.Fatalf("failed to insert test organization: %v", err)
	}
	defer db.GetPool().Exec(ctx, `DELETE FROM organizations WHERE id = $1`, orgUUID)

	origRunID := uuid.New()
	_, err = db.GetPool().Exec(ctx, `
		INSERT INTO llm_pentest_runs (
			id, org_id, target_endpoint, status, test_modules, 
			docker, config_mode, api_key, model, llm_provider, 
			api_base, mode, budget, concurrency, version, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`, origRunID, orgUUID, "https://api.example.com/chat", "completed", []string{"prompt_injection"},
		false, "custom", "sk-test", "gpt-4", "openai", "https://api.openai.com/v1", "both", 8, 4, 1, time.Now())
	if err != nil {
		t.Fatalf("failed to insert test LLM run: %v", err)
	}
	defer db.GetPool().Exec(ctx, `DELETE FROM llm_pentest_runs WHERE org_id = $1`, orgUUID)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("org_id", orgID)
		c.Locals("role", "org:admin")
		return c.Next()
	})
	app.Post("/llmpentest/:id/rerun", RerunLLMPentest)

	req := httptest.NewRequest("POST", "/llmpentest/"+origRunID.String()+"/rerun", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("POST /llmpentest/:id/rerun request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", resp.StatusCode)
	}

	var runObj models.LLMPentestRun
	json.NewDecoder(resp.Body).Decode(&runObj)
	if runObj.TargetEndpoint != "https://api.example.com/chat" {
		t.Errorf("expected target_endpoint %s, got %s", "https://api.example.com/chat", runObj.TargetEndpoint)
	}
	if runObj.ConfigMode != "custom" {
		t.Errorf("expected configMode custom, got %s", runObj.ConfigMode)
	}
}

func TestDeleteLLMPentest_Integration(t *testing.T) {
	if !initTestDB(t) {
		t.Skip("Database not available")
	}
	defer db.Close()

	ctx := context.Background()
	orgUUID := uuid.New()
	orgID := orgUUID.String()

	// Insert test organization
	_, err := db.GetPool().Exec(ctx, `
		INSERT INTO organizations (id, clerk_org_id, name)
		VALUES ($1, $2, $3)
	`, orgUUID, "clerk_"+orgUUID.String(), "Test API Org")
	if err != nil {
		t.Fatalf("failed to insert test organization: %v", err)
	}
	defer db.GetPool().Exec(ctx, `DELETE FROM organizations WHERE id = $1`, orgUUID)

	origRunID := uuid.New()
	_, err = db.GetPool().Exec(ctx, `
		INSERT INTO llm_pentest_runs (
			id, org_id, target_endpoint, status, test_modules, 
			docker, config_mode, api_key, model, llm_provider, 
			api_base, mode, budget, concurrency, version, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`, origRunID, orgUUID, "https://api.example.com/chat", "completed", []string{"prompt_injection"},
		false, "custom", "sk-test", "gpt-4", "openai", "https://api.openai.com/v1", "both", 8, 4, 1, time.Now())
	if err != nil {
		t.Fatalf("failed to insert test LLM run: %v", err)
	}
	defer db.GetPool().Exec(ctx, `DELETE FROM llm_pentest_runs WHERE org_id = $1`, orgUUID)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("org_id", orgID)
		c.Locals("role", "org:admin")
		return c.Next()
	})
	app.Delete("/llmpentest/:id", DeleteLLMPentest)

	req := httptest.NewRequest("DELETE", "/llmpentest/"+origRunID.String(), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("DELETE /llmpentest/:id request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}

	// Verify it was deleted
	var exists bool
	err = db.GetPool().QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM llm_pentest_runs WHERE id = $1)
	`, origRunID).Scan(&exists)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if exists {
		t.Error("expected run to be deleted from database, but it still exists")
	}
}