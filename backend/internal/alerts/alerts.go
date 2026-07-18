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
	"github.com/resend/resend-go/v3"
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

	subject := fmt.Sprintf("[Arishem] Scan Completed - %d findings", payload.FindingCount)
	htmlBody := fmt.Sprintf(`
		<h2>Arishem Scan Complete</h2>
		<p>A scan has completed with the following results:</p>
		<ul>
			<li><strong>Scan ID:</strong> %s</li>
			<li><strong>Findings:</strong> %d</li>
			<li><strong>Severity:</strong> %s</li>
		</ul>
		<p><a href="%s" style="background: #00e599; color: #000; padding: 10px 20px; text-decoration: none; border-radius: 4px;">View Report</a></p>
	`, payload.ScanID, payload.FindingCount, payload.Severity, payload.ReportURL)

	textBody := fmt.Sprintf(`Arishem Scan Complete

Scan ID: %s
Findings: %d
Severity: %s
Report: %s

Log in to Arishem to view the full report.
`, payload.ScanID, payload.FindingCount, payload.Severity, payload.ReportURL)

	// Try Resend first if API key is set
	resendKey := os.Getenv("RESEND_API_KEY")
	if resendKey != "" {
		return sendEmailResend(resendKey, addr, subject, htmlBody, textBody)
	}

	// Fall back to SMTP
	return sendEmailSMTP(addr, subject, textBody)
}

func sendEmailResend(apiKey, to, subject, htmlBody, textBody string) error {
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Arishem Alerts <alerts@arishem.site>",
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
		Text:    textBody,
	}

	_, err := client.Emails.Send(params)
	return err
}

func sendEmailSMTP(addr, subject, body string) error {
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

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, addr, subject, body)

	auth := smtp.PlainAuth("", user, pass, host)
	return smtp.SendMail(host+":"+port, auth, from, []string{addr}, []byte(msg))
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