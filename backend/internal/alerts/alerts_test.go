package alerts

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"arishem/internal/db"
	"arishem/internal/models"

	"github.com/google/uuid"
)

func TestAlertPayload(t *testing.T) {
	payload := AlertPayload{
		OrgID:        "org_123",
		ScanID:       "scan_456",
		Severity:     "high",
		FindingCount: 5,
		ReportURL:    "https://arishem.com/reports/123",
	}

	if payload.OrgID != "org_123" {
		t.Errorf("expected org_123, got %s", payload.OrgID)
	}
	if payload.FindingCount != 5 {
		t.Errorf("expected 5 findings, got %d", payload.FindingCount)
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
	} else {
		t.Logf("Failed to read schema file: %v", err)
	}

	return true
}

func TestDispatchAlerts(t *testing.T) {
	if !initTestDB(t) {
		t.Skip("Database not available")
	}
	defer db.Close()

	// Spin up test server for webhook alerts
	var receivedPayloads []AlertPayload
	var mu sync.Mutex
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		body, _ := io.ReadAll(r.Body)
		var p AlertPayload
		if err := json.Unmarshal(body, &p); err == nil {
			mu.Lock()
			receivedPayloads = append(receivedPayloads, p)
			mu.Unlock()
		}
	}))
	defer ts.Close()

	// Create test org and alert rules
	ctx := context.Background()
	orgUUID := uuid.New()
	orgID := orgUUID.String()

	// Insert test organization
	_, err := db.GetPool().Exec(ctx, `
		INSERT INTO organizations (id, clerk_org_id, name)
		VALUES ($1, $2, $3)
	`, orgUUID, "clerk_"+orgUUID.String(), "Test Alert Org")
	if err != nil {
		t.Fatalf("failed to insert test organization: %v", err)
	}
	defer db.GetPool().Exec(ctx, `DELETE FROM organizations WHERE id = $1`, orgUUID)

	// Insert test alert rule
	ruleUUID := uuid.New()
	config := map[string]interface{}{"url": ts.URL}
	configJSON, _ := json.Marshal(config)

	_, err = db.GetPool().Exec(ctx, `
		INSERT INTO alert_rules (id, org_id, severity_threshold, channel, channel_config, active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, ruleUUID, orgUUID, models.SeverityHigh, models.AlertChannelWebhook, configJSON, true)
	if err != nil {
		t.Fatalf("failed to insert test alert rule: %v", err)
	}
	defer db.GetPool().Exec(ctx, `DELETE FROM alert_rules WHERE id = $1`, ruleUUID)

	// 1. Dispatch alert below threshold (medium) -> should not trigger webhook
	payloadLow := AlertPayload{
		OrgID:        orgID,
		ScanID:       uuid.New().String(),
		Severity:     "medium",
		FindingCount: 2,
		ReportURL:    "http://test.url",
	}
	DispatchAlerts(ctx, orgID, payloadLow)

	// Give async worker time to execute
	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	if len(receivedPayloads) != 0 {
		t.Errorf("expected 0 payloads received for below-threshold alert, got %d", len(receivedPayloads))
	}
	mu.Unlock()

	// 2. Dispatch alert at or above threshold (critical) -> should trigger webhook
	payloadHigh := AlertPayload{
		OrgID:        orgID,
		ScanID:       uuid.New().String(),
		Severity:     "critical",
		FindingCount: 5,
		ReportURL:    "http://test.url",
	}
	DispatchAlerts(ctx, orgID, payloadHigh)

	// Give async worker time to execute
	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	if len(receivedPayloads) != 1 {
		t.Errorf("expected 1 payload received, got %d", len(receivedPayloads))
	} else {
		received := receivedPayloads[0]
		if received.ScanID != payloadHigh.ScanID {
			t.Errorf("expected ScanID %s, got %s", payloadHigh.ScanID, received.ScanID)
		}
	}
	mu.Unlock()
}

func TestSendTestAlert(t *testing.T) {
	if !initTestDB(t) {
		t.Skip("Database not available")
	}
	defer db.Close()

	// Spin up test server for webhook alerts
	var receivedPayloads []AlertPayload
	var mu sync.Mutex
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		body, _ := io.ReadAll(r.Body)
		var p AlertPayload
		if err := json.Unmarshal(body, &p); err == nil {
			mu.Lock()
			receivedPayloads = append(receivedPayloads, p)
			mu.Unlock()
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	orgUUID := uuid.New()

	// Insert test organization
	_, err := db.GetPool().Exec(ctx, `
		INSERT INTO organizations (id, clerk_org_id, name)
		VALUES ($1, $2, $3)
	`, orgUUID, "clerk_"+orgUUID.String(), "Test Alert Org")
	if err != nil {
		t.Fatalf("failed to insert test organization: %v", err)
	}
	defer db.GetPool().Exec(ctx, `DELETE FROM organizations WHERE id = $1`, orgUUID)

	// Insert test alert rule
	ruleUUID := uuid.New()
	config := map[string]interface{}{"url": ts.URL}
	configJSON, _ := json.Marshal(config)

	_, err = db.GetPool().Exec(ctx, `
		INSERT INTO alert_rules (id, org_id, severity_threshold, channel, channel_config, active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, ruleUUID, orgUUID, models.SeverityHigh, models.AlertChannelWebhook, configJSON, true)
	if err != nil {
		t.Fatalf("failed to insert test alert rule: %v", err)
	}
	defer db.GetPool().Exec(ctx, `DELETE FROM alert_rules WHERE id = $1`, ruleUUID)

	testPayload := AlertPayload{
		OrgID:        orgUUID.String(),
		ScanID:       uuid.New().String(),
		Severity:     "info",
		FindingCount: 1,
		ReportURL:    "http://test.url/test",
	}

	err = SendTestAlert(ctx, ruleUUID.String(), testPayload)
	if err != nil {
		t.Errorf("unexpected error in SendTestAlert: %v", err)
	}

	// Give async worker time to execute
	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	if len(receivedPayloads) != 1 {
		t.Errorf("expected test alert to trigger webhook, got %d hits", len(receivedPayloads))
	}
	mu.Unlock()
}

func TestSendTestAlertNotFound(t *testing.T) {
	if !initTestDB(t) {
		t.Skip("Database not available")
	}
	defer db.Close()

	ctx := context.Background()
	randomUUID := uuid.New().String()
	testPayload := AlertPayload{}

	err := SendTestAlert(ctx, randomUUID, testPayload)
	if err == nil {
		t.Error("expected error for non-existent rule ID, got nil")
	} else if err.Error() != "rule not found" {
		t.Errorf("expected 'rule not found' error, got: %v", err)
	}
}

func TestAlertSeverityComparison(t *testing.T) {
	payload := AlertPayload{
		Severity: "critical",
	}

	if payload.Severity != "critical" {
		t.Errorf("expected critical, got %s", payload.Severity)
	}
}