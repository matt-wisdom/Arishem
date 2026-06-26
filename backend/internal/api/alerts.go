package api

import (
	"time"

	"arishem/internal/alerts"
	"arishem/internal/db"
	"arishem/internal/middleware"
	"arishem/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateAlertRequest struct {
	SeverityThreshold models.Severity               `json:"severity_threshold"`
	Channel           models.AlertChannel           `json:"channel"`
	ChannelConfig     map[string]interface{}        `json:"channel_config"`
}

func ListAlerts(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)

	rows, err := db.GetPool().Query(c.Context(), `
		SELECT id, org_id, severity_threshold, channel, channel_config, active, created_at
		FROM alert_rules WHERE org_id = $1 ORDER BY created_at DESC
	`, orgUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch alerts"})
	}
	defer rows.Close()

	var rules []models.AlertRule
	for rows.Next() {
		var r models.AlertRule
		if err := rows.Scan(&r.ID, &r.OrgID, &r.SeverityThreshold, &r.Channel, &r.ChannelConfig, &r.Active, &r.CreatedAt); err != nil {
			continue
		}
		rules = append(rules, r)
	}

	if rules == nil {
		rules = []models.AlertRule{}
	}

	return c.JSON(rules)
}

func CreateAlert(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	role := middleware.GetRole(c)
	if role != "org:admin" && role != "org:engineer" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
	}

	var req CreateAlertRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	orgUUID, _ := uuid.Parse(orgID)
	rule := models.AlertRule{
		ID:               uuid.New(),
		OrgID:            orgUUID,
		SeverityThreshold: req.SeverityThreshold,
		Channel:          req.Channel,
		ChannelConfig:    req.ChannelConfig,
		Active:           true,
		CreatedAt:        time.Now(),
	}

	_, err := db.GetPool().Exec(c.Context(), `
		INSERT INTO alert_rules (id, org_id, severity_threshold, channel, channel_config, active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, rule.ID, rule.OrgID, rule.SeverityThreshold, rule.Channel, rule.ChannelConfig, rule.Active, rule.CreatedAt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create alert rule"})
	}

	return c.Status(fiber.StatusCreated).JSON(rule)
}

func UpdateAlert(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)
	ruleID := c.Params("id")
	ruleUUID, err := uuid.Parse(ruleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid alert id"})
	}

	var req CreateAlertRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	result, err := db.GetPool().Exec(c.Context(), `
		UPDATE alert_rules SET severity_threshold = $1, channel = $2, channel_config = $3 WHERE id = $4 AND org_id = $5
	`, req.SeverityThreshold, req.Channel, req.ChannelConfig, ruleUUID, orgUUID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update alert rule"})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "alert rule not found"})
	}

	return c.JSON(fiber.Map{"message": "alert rule updated"})
}

func DeleteAlert(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	orgUUID, _ := uuid.Parse(orgID)
	ruleID := c.Params("id")
	ruleUUID, err := uuid.Parse(ruleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid alert id"})
	}

	result, err := db.GetPool().Exec(c.Context(), `DELETE FROM alert_rules WHERE id = $1 AND org_id = $2`, ruleUUID, orgUUID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete alert rule"})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "alert rule not found"})
	}

	return c.JSON(fiber.Map{"message": "alert rule deleted"})
}

func TestAlert(c *fiber.Ctx) error {
	orgID := middleware.GetOrgID(c)
	ruleID := c.Params("id")

	alert := alerts.AlertPayload{
		OrgID:       orgID,
		ScanID:      "test-scan-id",
		Severity:    "high",
		FindingCount: 1,
		ReportURL:   "https://example.com/test-report",
	}

	if err := alerts.SendTestAlert(c.Context(), ruleID, alert); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to send test alert"})
	}

	return c.JSON(fiber.Map{"message": "test alert sent"})
}