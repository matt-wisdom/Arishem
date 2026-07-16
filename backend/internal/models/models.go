package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `json:"id"`
	ClerkOrgID  string    `json:"clerk_org_id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
}

type User struct {
	ID          uuid.UUID `json:"id"`
	ClerkUserID string    `json:"clerk_user_id"`
	OrgID       uuid.UUID `json:"org_id"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type ScanType string
type ScanStatus string

const (
	ScanTypeCode   ScanType = "code"
	ScanTypeWebapp ScanType = "webapp"

	ScanStatusQueued    ScanStatus = "queued"
	ScanStatusRunning   ScanStatus = "running"
	ScanStatusCompleted ScanStatus = "completed"
	ScanStatusFailed    ScanStatus = "failed"
	ScanStatusCancelled ScanStatus = "cancelled"
)

type Scan struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	Title       string     `json:"title"`
	Type        ScanType   `json:"type"`
	Status      ScanStatus `json:"status"`
	Target      string     `json:"target"`
	Branch      string     `json:"branch,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type LLMPentestRun struct {
	ID             uuid.UUID  `json:"id"`
	OrgID          uuid.UUID  `json:"org_id"`
	Title          string     `json:"title"`
	TargetEndpoint string     `json:"target_endpoint"`
	Status         ScanStatus `json:"status"`
	TestModules    []string   `json:"test_modules"`
	Logs           string     `json:"logs"`
	Docker         bool       `json:"docker"`
	ConfigMode     string     `json:"config_mode"`
	APIKey         string     `json:"api_key"`
	Model          string     `json:"model"`
	LLMProvider    string     `json:"llm_provider"`
	APIBase        string     `json:"api_base"`
	Mode           string     `json:"mode"`
	Budget         int        `json:"budget"`
	Concurrency    int        `json:"concurrency"`
	Version        int        `json:"version"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

type Finding struct {
	ID           uuid.UUID `json:"id"`
	OrgID        uuid.UUID `json:"org_id"`
	ScanID       *uuid.UUID `json:"scan_id,omitempty"`
	RunID        *uuid.UUID `json:"run_id,omitempty"`
	Title        string    `json:"title"`
	Severity     Severity  `json:"severity"`
	Description  string    `json:"description"`
	Remediation  string    `json:"remediation"`
	CreatedAt    time.Time `json:"created_at"`
}

type ReportFormat string

const (
	ReportFormatHTML  ReportFormat = "html"
	ReportFormatMD    ReportFormat = "md"
	ReportFormatSARIF ReportFormat = "sarif"
)

type Report struct {
	ID         uuid.UUID     `json:"id"`
	OrgID      uuid.UUID     `json:"org_id"`
	ScanID     *uuid.UUID    `json:"scan_id,omitempty"`
	RunID      *uuid.UUID    `json:"run_id,omitempty"`
	Format     ReportFormat  `json:"format"`
	StorageKey string        `json:"storage_key"`
	CreatedAt  time.Time     `json:"created_at"`
}

type IntegrationProvider string

const (
	IntegrationProviderGitHub IntegrationProvider = "github"
)

type Integration struct {
	ID                 uuid.UUID          `json:"id"`
	OrgID              uuid.UUID          `json:"org_id"`
	Provider           IntegrationProvider `json:"provider"`
	CredentialsEncrypted string           `json:"credentials_encrypted"`
	CreatedAt          time.Time          `json:"created_at"`
}

type AlertChannel string

const (
	AlertChannelEmail   AlertChannel = "email"
	AlertChannelSlack   AlertChannel = "slack"
	AlertChannelWebhook AlertChannel = "webhook"
)

type AlertRule struct {
	ID               uuid.UUID          `json:"id"`
	OrgID            uuid.UUID          `json:"org_id"`
	SeverityThreshold Severity          `json:"severity_threshold"`
	Channel          AlertChannel       `json:"channel"`
	ChannelConfig    map[string]interface{} `json:"channel_config"`
	Active           bool               `json:"active"`
	CreatedAt        time.Time          `json:"created_at"`
}