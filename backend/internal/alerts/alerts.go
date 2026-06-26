package alerts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"arishem/internal/db"
	"arishem/internal/models"

	"github.com/google/uuid"
)

type AlertPayload struct {
	OrgID       string
	ScanID      string
	Severity    string
	FindingCount int
	ReportURL   string
}

func DispatchAlerts(ctx context.Context, orgID string, payload AlertPayload) {
	orgUUID, _ := uuid.Parse(orgID)

	rows, err := db.GetPool().Query(ctx, `
		SELECT id, severity_threshold, channel, channel_config, active
		FROM alert_rules WHERE org_id = $1 AND active = true
	`, orgUUID)
	if err != nil {
		log.Printf("Failed to fetch alert rules: %v", err)
		return
	}
	defer rows.Close()

	severityOrder := map[models.Severity]int{
		models.SeverityCritical: 4,
		models.SeverityHigh:     3,
		models.SeverityMedium:   2,
		models.SeverityLow:      1,
		models.SeverityInfo:     0,
	}

	for rows.Next() {
		var rule models.AlertRule
		if err := rows.Scan(&rule.ID, &rule.SeverityThreshold, &rule.Channel, &rule.ChannelConfig, &rule.Active); err != nil {
			continue
		}

		payloadSeverity := payload.Severity
		var payloadSeverityLevel int
		for s, level := range severityOrder {
			if string(s) == payloadSeverity {
				payloadSeverityLevel = level
				break
			}
		}

		thresholdLevel := severityOrder[rule.SeverityThreshold]
		if payloadSeverityLevel >= thresholdLevel {
			go sendAlert(rule, payload)
		}
	}
}

func sendAlert(rule models.AlertRule, payload AlertPayload) {
	var err error
	switch rule.Channel {
	case models.AlertChannelEmail:
		err = sendEmail(rule.ChannelConfig, payload)
	case models.AlertChannelSlack:
		err = sendSlack(rule.ChannelConfig, payload)
	case models.AlertChannelWebhook:
		err = sendWebhook(rule.ChannelConfig, payload)
	}

	if err != nil {
		log.Printf("Failed to send alert via %s: %v", rule.Channel, err)
	}
}

func sendEmail(config map[string]interface{}, payload AlertPayload) error {
	addr, _ := config["address"].(string)
	if addr == "" {
		return fmt.Errorf("email address not configured")
	}

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "587"
	}

	from := "alerts@arishem.com"
	if user != "" {
		from = user
	}

	subject := fmt.Sprintf("[Arishem] Scan Completed - %d findings", payload.FindingCount)
	body := fmt.Sprintf(`A scan has completed with the following results:

Org ID: %s
Scan ID: %s
Finding Count: %d
Severity: %s
Report URL: %s

Log in to Arishem to view the full report.
`, payload.OrgID, payload.ScanID, payload.FindingCount, payload.Severity, payload.ReportURL)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, addr, subject, body)

	auth := smtp.PlainAuth("", user, pass, host)
	err := smtp.SendMail(host+":"+port, auth, from, []string{addr}, []byte(msg))

	return err
}

func sendSlack(config map[string]interface{}, payload AlertPayload) error {
	webhookURL, _ := config["url"].(string)
	if webhookURL == "" {
		return fmt.Errorf("Slack webhook URL not configured")
	}

	slackMsg := map[string]interface{}{
		"text": fmt.Sprintf("*Arishem Alert*\nScan completed with %d findings\nReport: %s", payload.FindingCount, payload.ReportURL),
	}

	jsonData, _ := json.Marshal(slackMsg)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Slack webhook returned status %d", resp.StatusCode)
	}

	return nil
}

func sendWebhook(config map[string]interface{}, payload AlertPayload) error {
	webhookURL, _ := config["url"].(string)
	if webhookURL == "" {
		return fmt.Errorf("webhook URL not configured")
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

func SendTestAlert(ctx context.Context, ruleID string, payload AlertPayload) error {
	ruleUUID, err := uuid.Parse(ruleID)
	if err != nil {
		return fmt.Errorf("invalid rule ID")
	}

	var rule models.AlertRule
	err = db.GetPool().QueryRow(ctx, `
		SELECT id, severity_threshold, channel, channel_config, active
		FROM alert_rules WHERE id = $1
	`, ruleUUID).Scan(&rule.ID, &rule.SeverityThreshold, &rule.Channel, &rule.ChannelConfig, &rule.Active)

	if err != nil {
		return fmt.Errorf("rule not found")
	}

	sendAlert(rule, payload)
	return nil
}