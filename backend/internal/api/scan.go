package api

import (
	"time"

	"arishem/internal/db"
	"arishem/internal/jobs"
	"arishem/internal/middleware"
	"arishem/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateScanRequest struct {
	Target string `json:"target"`
	Branch string `json:"branch"`
}

func CreateCodeScan(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	role := middleware.GetRole(c)
	if role != "org:admin" && role != "org:engineer" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
	}

	var req CreateScanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Target == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "target is required"})
	}

	orgUUID, _ := uuid.Parse(orgID)
	scan := models.Scan{
		ID:        uuid.New(),
		OrgID:     orgUUID,
		Type:      models.ScanTypeCode,
		Status:    models.ScanStatusQueued,
		Target:    req.Target,
		Branch:    req.Branch,
		CreatedAt: time.Now(),
	}

	if scan.Branch == "" {
		scan.Branch = "main"
	}

	_, err := db.GetPool().Exec(c.Context(), `
		INSERT INTO scans (id, org_id, type, status, target, branch, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, scan.ID, scan.OrgID, scan.Type, scan.Status, scan.Target, scan.Branch, scan.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create scan"})
	}

	task := jobs.NewScanTask(orgID, scan.ID.String(), string(scan.Type), scan.Target, scan.Branch)
	if err := jobs.EnqueueScanTask(c.Context(), task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to enqueue scan job"})
	}

	return c.Status(fiber.StatusCreated).JSON(scan)
}

func CreateWebappScan(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	role := middleware.GetRole(c)
	if role != "org:admin" && role != "org:engineer" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
	}

	var req CreateScanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Target == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "target is required"})
	}

	orgUUID, _ := uuid.Parse(orgID)
	scan := models.Scan{
		ID:        uuid.New(),
		OrgID:     orgUUID,
		Type:      models.ScanTypeWebapp,
		Status:    models.ScanStatusQueued,
		Target:    req.Target,
		CreatedAt: time.Now(),
	}

	_, err := db.GetPool().Exec(c.Context(), `
		INSERT INTO scans (id, org_id, type, status, target, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, scan.ID, scan.OrgID, scan.Type, scan.Status, scan.Target, scan.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create scan"})
	}

	task := jobs.NewScanTask(orgID, scan.ID.String(), string(scan.Type), scan.Target, "")
	if err := jobs.EnqueueScanTask(c.Context(), task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to enqueue scan job"})
	}

	return c.Status(fiber.StatusCreated).JSON(scan)
}

func ListScans(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)

	rows, err := db.GetPool().Query(c.Context(), `
		SELECT id, org_id, type, status, target, branch, created_at, completed_at
		FROM scans WHERE org_id = $1 ORDER BY created_at DESC LIMIT 50
	`, orgUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch scans"})
	}
	defer rows.Close()

	var scans []models.Scan
	for rows.Next() {
		var scan models.Scan
		if err := rows.Scan(&scan.ID, &scan.OrgID, &scan.Type, &scan.Status, &scan.Target, &scan.Branch, &scan.CreatedAt, &scan.CompletedAt); err != nil {
			continue
		}
		scans = append(scans, scan)
	}

	if scans == nil {
		scans = []models.Scan{}
	}

	return c.JSON(scans)
}

func GetScan(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)
	scanID := c.Params("id")
	scanUUID, err := uuid.Parse(scanID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid scan id"})
	}

	var scan models.Scan
	err = db.GetPool().QueryRow(c.Context(), `
		SELECT id, org_id, type, status, target, branch, created_at, completed_at
		FROM scans WHERE id = $1 AND org_id = $2
	`, scanUUID, orgUUID).Scan(&scan.ID, &scan.OrgID, &scan.Type, &scan.Status, &scan.Target, &scan.Branch, &scan.CreatedAt, &scan.CompletedAt)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "scan not found"})
	}

	rows, err := db.GetPool().Query(c.Context(), `
		SELECT id, org_id, scan_id, run_id, title, severity, description, remediation, created_at
		FROM findings WHERE scan_id = $1 AND org_id = $2
	`, scanUUID, orgUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch findings"})
	}
	defer rows.Close()

	var findings []models.Finding
	for rows.Next() {
		var f models.Finding
		if err := rows.Scan(&f.ID, &f.OrgID, &f.ScanID, &f.RunID, &f.Title, &f.Severity, &f.Description, &f.Remediation, &f.CreatedAt); err != nil {
			continue
		}
		findings = append(findings, f)
	}

	return c.JSON(fiber.Map{"scan": scan, "findings": findings})
}

func CancelScan(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)
	scanID := c.Params("id")
	scanUUID, err := uuid.Parse(scanID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid scan id"})
	}

	// Try canceling the active task in the worker pool
	jobs.CancelActiveTask(scanID)

	result, err := db.GetPool().Exec(c.Context(), `
		UPDATE scans SET status = $1 WHERE id = $2 AND org_id = $3 AND status IN ($4, $5)
	`, models.ScanStatusCancelled, scanUUID, orgUUID, models.ScanStatusQueued, models.ScanStatusRunning)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to cancel scan"})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "scan not found or cannot be cancelled"})
	}

	return c.JSON(fiber.Map{"message": "scan cancelled"})
}