package api

import (
	"arishem/internal/db"
	"arishem/internal/middleware"
	"arishem/internal/models"
	"arishem/internal/reports"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ListReports(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)

	rows, err := db.GetPool().Query(c.Context(), `
		SELECT id, org_id, scan_id, run_id, format, storage_key, created_at
		FROM reports WHERE org_id = $1 ORDER BY created_at DESC LIMIT 50
	`, orgUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch reports"})
	}
	defer rows.Close()

	var reportList []models.Report
	for rows.Next() {
		var r models.Report
		if err := rows.Scan(&r.ID, &r.OrgID, &r.ScanID, &r.RunID, &r.Format, &r.StorageKey, &r.CreatedAt); err != nil {
			continue
		}
		reportList = append(reportList, r)
	}

	if reportList == nil {
		reportList = []models.Report{}
	}

	return c.JSON(reportList)
}

func GetReport(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)
	reportID := c.Params("id")
	reportUUID, err := uuid.Parse(reportID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid report id"})
	}

	var report models.Report
	err = db.GetPool().QueryRow(c.Context(), `
		SELECT id, org_id, scan_id, run_id, format, storage_key, created_at
		FROM reports WHERE id = $1 AND org_id = $2
	`, reportUUID, orgUUID).Scan(&report.ID, &report.OrgID, &report.ScanID, &report.RunID, &report.Format, &report.StorageKey, &report.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "report not found"})
	}

	return c.JSON(report)
}

func GetReportDownload(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	reportID := c.Params("id")
	format := c.Query("format", "html")

	url, err := reports.GetSignedURL(c.Context(), orgID, reportID, format)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate download URL"})
	}

	return c.JSON(fiber.Map{"url": url})
}