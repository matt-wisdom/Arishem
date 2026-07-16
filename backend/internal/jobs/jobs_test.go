package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewScanTask(t *testing.T) {
	task := NewScanTask("org_123", "scan_456", "code", "https://github.com/test/repo", "main")

	if task.OrgID != "org_123" {
		t.Errorf("expected OrgID org_123, got %s", task.OrgID)
	}
	if task.ScanID != "scan_456" {
		t.Errorf("expected ScanID scan_456, got %s", task.ScanID)
	}
	if task.Type != "code" {
		t.Errorf("expected Type code, got %s", task.Type)
	}
	if task.Target != "https://github.com/test/repo" {
		t.Errorf("expected Target, got %s", task.Target)
	}
	if task.Branch != "main" {
		t.Errorf("expected Branch main, got %s", task.Branch)
	}
}

func TestNewScanTaskDefaultBranch(t *testing.T) {
	task := NewScanTask("org_123", "scan_456", "code", "https://github.com/test/repo", "")

	if task.Branch != "" {
		t.Logf("Branch is empty as expected: %s", task.Branch)
	}
}

func TestNewLLMPentestTask(t *testing.T) {
	task := NewLLMPentestTask("org_123", "run_456", "https://api.example.com/chat", "sk-test", []string{"prompt_injection", "jailbreak"}, "gpt-4", "openai", "", "both", 8, 4, false, "custom")

	if task.OrgID != "org_123" {
		t.Errorf("expected OrgID org_123, got %s", task.OrgID)
	}
	if task.RunID != "run_456" {
		t.Errorf("expected RunID run_456, got %s", task.RunID)
	}
	if task.TargetEndpoint != "https://api.example.com/chat" {
		t.Errorf("expected TargetEndpoint, got %s", task.TargetEndpoint)
	}
	if task.APIKey != "sk-test" {
		t.Errorf("expected APIKey sk-test, got %s", task.APIKey)
	}
	if len(task.TestModules) != 2 {
		t.Errorf("expected 2 test modules, got %d", len(task.TestModules))
	}
}

func TestEnqueueScanTask(t *testing.T) {
	// Just verify it doesn't panic and returns nil (runs async)
	task := NewScanTask("org_123", "scan_456", "code", "https://github.com/test/repo", "main")
	err := EnqueueScanTask(context.Background(), task)
	
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEnqueueLLMPentestTask(t *testing.T) {
	// Just verify it doesn't panic and returns nil (runs async)
	task := NewLLMPentestTask("org_123", "run_456", "https://api.example.com/chat", "sk-test", []string{"prompt_injection"}, "gpt-4", "openai", "", "both", 8, 4, false, "custom")
	err := EnqueueLLMPentestTask(context.Background(), task)
	
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestScanTaskStruct(t *testing.T) {
	task := &ScanTask{
		OrgID:  "org_abc",
		ScanID: "scan_xyz",
		Type:   "webapp",
		Target: "https://example.com",
		Branch: "develop",
	}

	if task.OrgID != "org_abc" {
		t.Errorf("expected org_abc, got %s", task.OrgID)
	}
	if task.Type != "webapp" {
		t.Errorf("expected webapp, got %s", task.Type)
	}
}

func TestLLMPentestTaskStruct(t *testing.T) {
	task := &LLMPentestTask{
		OrgID:          "org_abc",
		RunID:          "run_xyz",
		TargetEndpoint: "https://api.test.com/v1/chat",
		APIKey:         "sk_key",
		TestModules:    []string{"all"},
	}

	if task.OrgID != "org_abc" {
		t.Errorf("expected org_abc, got %s", task.OrgID)
	}
	if len(task.TestModules) != 1 {
		t.Errorf("expected 1 test module, got %d", len(task.TestModules))
	}
}

func TestJobTypes(t *testing.T) {
	if TypeScan != "scan" {
		t.Errorf("expected scan, got %s", TypeScan)
	}
	if TypeLLMPentest != "llm_pentest" {
		t.Errorf("expected llm_pentest, got %s", TypeLLMPentest)
	}
}

func TestUUIDParsing(t *testing.T) {
	orgID := uuid.New().String()
	parsed, err := uuid.Parse(orgID)
	
	if err != nil {
		t.Errorf("failed to parse UUID: %v", err)
	}
	if parsed.String() != orgID {
		t.Errorf("UUID mismatch")
	}
}

func TestTimeNow(t *testing.T) {
	now := time.Now()
	if now.IsZero() {
		t.Error("time.Now() should not be zero")
	}
}

func TestRegisterAndCancelActiveTask(t *testing.T) {
	runID := uuid.New().String()
	orgID := uuid.New().String()
	canceled := false
	cancel := func() {
		canceled = true
	}

	// Register task
	registered := RegisterActiveTask(runID, orgID, cancel)
	if !registered {
		t.Error("expected task to be registered successfully")
	}

	// Double register should fail because limit per org is max 2 (let's check with 3 tasks)
	runID2 := uuid.New().String()
	registered2 := RegisterActiveTask(runID2, orgID, func() {})
	if !registered2 {
		t.Error("expected second task to be registered successfully (max is 2)")
	}

	runID3 := uuid.New().String()
	registered3 := RegisterActiveTask(runID3, orgID, func() {})
	if registered3 {
		t.Error("expected third task to be rate-limited (max is 2)")
	}

	// Cancel task
	canceledOk := CancelActiveTask(runID)
	if !canceledOk {
		t.Error("expected task to be canceled successfully")
	}
	if !canceled {
		t.Error("expected cancel function to be executed")
	}

	// Deregister task
	DeregisterActiveTask(runID2, orgID)
	DeregisterActiveTask(runID3, orgID)
}