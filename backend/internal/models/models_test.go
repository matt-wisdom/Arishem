package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestOrganization(t *testing.T) {
	org := Organization{
		ID:         uuid.New(),
		ClerkOrgID: "org_123",
		Name:       "Test Org",
		CreatedAt:  time.Now(),
	}

	if org.ClerkOrgID != "org_123" {
		t.Errorf("expected ClerkOrgID org_123, got %s", org.ClerkOrgID)
	}
	if org.Name != "Test Org" {
		t.Errorf("expected Name Test Org, got %s", org.Name)
	}
}

func TestUser(t *testing.T) {
	orgID := uuid.New()
	user := User{
		ID:          uuid.New(),
		ClerkUserID: "user_123",
		OrgID:       orgID,
		Role:        "org:admin",
		CreatedAt:   time.Now(),
	}

	if user.ClerkUserID != "user_123" {
		t.Errorf("expected ClerkUserID user_123, got %s", user.ClerkUserID)
	}
	if user.Role != "org:admin" {
		t.Errorf("expected Role org:admin, got %s", user.Role)
	}
}

func TestScan(t *testing.T) {
	orgID := uuid.New()
	scan := Scan{
		ID:        uuid.New(),
		OrgID:     orgID,
		Type:      ScanTypeCode,
		Status:    ScanStatusQueued,
		Target:    "https://github.com/test/repo",
		Branch:    "main",
		CreatedAt: time.Now(),
	}

	if scan.Type != ScanTypeCode {
		t.Errorf("expected Type code, got %s", scan.Type)
	}
	if scan.Status != ScanStatusQueued {
		t.Errorf("expected Status queued, got %s", scan.Status)
	}
	if scan.Target != "https://github.com/test/repo" {
		t.Errorf("expected Target, got %s", scan.Target)
	}
}

func TestScanType(t *testing.T) {
	tests := []struct {
		input    ScanType
		expected string
	}{
		{ScanTypeCode, "code"},
		{ScanTypeWebapp, "webapp"},
	}

	for _, tt := range tests {
		if string(tt.input) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.input)
		}
	}
}

func TestScanStatus(t *testing.T) {
	statuses := []ScanStatus{
		ScanStatusQueued,
		ScanStatusRunning,
		ScanStatusCompleted,
		ScanStatusFailed,
		ScanStatusCancelled,
	}

	expected := []string{"queued", "running", "completed", "failed", "cancelled"}

	for i, status := range statuses {
		if string(status) != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], status)
		}
	}
}

func TestLLMPentestRun(t *testing.T) {
	orgID := uuid.New()
	run := LLMPentestRun{
		ID:             uuid.New(),
		OrgID:          orgID,
		TargetEndpoint: "https://api.example.com/chat",
		Status:         ScanStatusQueued,
		TestModules:    []string{"prompt_injection", "jailbreak"},
		CreatedAt:      time.Now(),
	}

	if run.TargetEndpoint != "https://api.example.com/chat" {
		t.Errorf("expected TargetEndpoint, got %s", run.TargetEndpoint)
	}
	if len(run.TestModules) != 2 {
		t.Errorf("expected 2 test modules, got %d", len(run.TestModules))
	}
}

func TestSeverity(t *testing.T) {
	severities := []Severity{
		SeverityCritical,
		SeverityHigh,
		SeverityMedium,
		SeverityLow,
		SeverityInfo,
	}

	expected := []string{"critical", "high", "medium", "low", "info"}

	for i, sev := range severities {
		if string(sev) != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], sev)
		}
	}
}

func TestFinding(t *testing.T) {
	orgID := uuid.New()
	scanID := uuid.New()
	finding := Finding{
		ID:          uuid.New(),
		OrgID:       orgID,
		ScanID:      &scanID,
		Title:       "SQL Injection",
		Severity:    SeverityHigh,
		Description: "User input not properly sanitized",
		Remediation: "Use parameterized queries",
		CreatedAt:   time.Now(),
	}

	if finding.Title != "SQL Injection" {
		t.Errorf("expected Title SQL Injection, got %s", finding.Title)
	}
	if finding.Severity != SeverityHigh {
		t.Errorf("expected Severity high, got %s", finding.Severity)
	}
	if finding.ScanID == nil {
		t.Error("expected ScanID to be set")
	}
}

func TestReport(t *testing.T) {
	orgID := uuid.New()
	runID := uuid.New()
	report := Report{
		ID:         uuid.New(),
		OrgID:      orgID,
		RunID:      &runID,
		Format:     ReportFormatHTML,
		StorageKey: "reports/report.html",
		CreatedAt:  time.Now(),
	}

	if report.Format != ReportFormatHTML {
		t.Errorf("expected Format html, got %s", report.Format)
	}
	if report.RunID == nil {
		t.Error("expected RunID to be set")
	}
}

func TestIntegration(t *testing.T) {
	orgID := uuid.New()
	integration := Integration{
		ID:                   uuid.New(),
		OrgID:                orgID,
		Provider:             IntegrationProviderGitHub,
		CredentialsEncrypted: "encrypted_token",
		CreatedAt:            time.Now(),
	}

	if integration.Provider != IntegrationProviderGitHub {
		t.Errorf("expected Provider github, got %s", integration.Provider)
	}
}

func TestAlertRule(t *testing.T) {
	orgID := uuid.New()
	rule := AlertRule{
		ID:               uuid.New(),
		OrgID:            orgID,
		SeverityThreshold: SeverityHigh,
		Channel:          AlertChannelSlack,
		ChannelConfig:    map[string]interface{}{"url": "https://hooks.slack.com/xxx"},
		Active:           true,
		CreatedAt:        time.Now(),
	}

	if rule.SeverityThreshold != SeverityHigh {
		t.Errorf("expected SeverityThreshold high, got %s", rule.SeverityThreshold)
	}
	if rule.Channel != AlertChannelSlack {
		t.Errorf("expected Channel slack, got %s", rule.Channel)
	}
	if !rule.Active {
		t.Error("expected Active to be true")
	}
}