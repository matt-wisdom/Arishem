package jobs

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"arishem/internal/alerts"
	"arishem/internal/db"
	"arishem/internal/models"

	"github.com/google/uuid"
)

const TypeScan = "scan"

type ScanTask struct {
	OrgID  string `json:"org_id"`
	ScanID string `json:"scan_id"`
	Type   string `json:"type"`
	Target string `json:"target"`
	Branch string `json:"branch"`
}

func NewScanTask(orgID, scanID, scanType, target, branch string) *ScanTask {
	return &ScanTask{
		OrgID:  orgID,
		ScanID: scanID,
		Type:   scanType,
		Target: target,
		Branch: branch,
	}
}

func EnqueueScanTask(ctx context.Context, task *ScanTask) error {
	select {
	case scanQueue <- task:
	default:
		go func() { scanQueue <- task }()
	}
	return nil
}

func processScanTask(ctx context.Context, task *ScanTask) {
	scanUUID, _ := uuid.Parse(task.ScanID)

	log.Printf("Processing scan job: %s (type=%s target=%s)", task.ScanID, task.Type, task.Target)

	pool := db.GetPool()
	if pool == nil {
		log.Printf("Database not initialized, cannot update scan status")
		return
	}

	pool.Exec(ctx, `UPDATE scans SET status = $1 WHERE id = $2`, models.ScanStatusRunning, scanUUID)

	var scanErr error
	if task.Type == "code" {
		scanErr = runCodeScan(ctx, task)
	} else if task.Type == "webapp" {
		scanErr = runWebappScan(ctx, task)
	}

	now := time.Now()
	if scanErr != nil {
		if ctx.Err() == context.Canceled {
			log.Printf("Scan cancelled by user: %s", task.ScanID)
			pool.Exec(ctx, `UPDATE scans SET status = $1, completed_at = $2 WHERE id = $3`, models.ScanStatusCancelled, now, scanUUID)
			return
		}
		pool.Exec(ctx, `UPDATE scans SET status = $1, completed_at = $2 WHERE id = $3`, models.ScanStatusFailed, now, scanUUID)
		log.Printf("Scan failed: %s error=%v", task.ScanID, scanErr)
	} else {
		pool.Exec(ctx, `UPDATE scans SET status = $1, completed_at = $2 WHERE id = $3`, models.ScanStatusCompleted, now, scanUUID)

		alertPayload := alerts.AlertPayload{
			OrgID:        task.OrgID,
			ScanID:       task.ScanID,
			Severity:     "completed",
			FindingCount: 0,
			ReportURL:    fmt.Sprintf("/api/reports/%s/download?format=html", task.ScanID),
		}
		alerts.DispatchAlerts(ctx, task.OrgID, alertPayload)
	}
}

func runCodeScan(ctx context.Context, task *ScanTask) error {
	scannerPath := "../../scanner/code_scanner/runner.py"
	if _, err := os.Stat(scannerPath); os.IsNotExist(err) {
		scannerPath = "../scanner/code_scanner/runner.py"
	}

	cmd := exec.CommandContext(ctx, "python3", scannerPath,
		"--target", task.Target,
		"--output-dir", "/tmp/arishem/reports",
		"--formats", "html,md,sarif",
	)
	cmd.Dir = "/mnt/C6EE65A1EE658B0F/WORKEST/Arishem"

	output, err := cmd.CombinedOutput()
	log.Printf("Code scan output: %s", string(output))

	if err != nil {
		return fmt.Errorf("code scan failed: %w", err)
	}

	reportUUID := uuid.New()
	storageKey := fmt.Sprintf("scans/%s/report.html", task.ScanID)

	_, err = db.GetPool().Exec(ctx, `
		INSERT INTO reports (id, org_id, scan_id, format, storage_key, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, reportUUID, task.OrgID, task.ScanID, models.ReportFormatHTML, storageKey, time.Now())

	return err
}

func runWebappScan(ctx context.Context, task *ScanTask) error {
	log.Printf("Webapp scan not yet implemented for target: %s", task.Target)
	time.Sleep(2 * time.Second)

	scanUUID, _ := uuid.Parse(task.ScanID)
	now := time.Now()
	db.GetPool().Exec(ctx, `
		INSERT INTO findings (id, org_id, scan_id, title, severity, description, remediation, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, uuid.New(), task.OrgID, scanUUID, "Sample Finding", models.SeverityLow, "This is a placeholder finding for webapp scans.", "Implement actual DAST scanning.", now)

	return nil
}