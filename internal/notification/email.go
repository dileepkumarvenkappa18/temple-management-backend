package notification

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
	"strings"

	"github.com/sharath018/temple-management-backend/config"
)

// EmailSender implements Channel interface using SMTP
type EmailSender struct {
	Host     string
	Port     string
	Username string
	Password string
	FromName string
	FromAddr string
}

// ✅ Accept config instead of using os.Getenv
func NewEmailSender(cfg *config.Config) *EmailSender {
	return &EmailSender{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		FromName: cfg.SMTPFromName,
		FromAddr: cfg.SMTPFromEmail,
	}
}

// Send renders the HTML template and sends the email
func (e *EmailSender) Send(to []string, subject string, body string) error {
	// Step 1: Load and parse the template
	tmplPath := filepath.Join("templates", "example.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		fmt.Println("❌ Failed to parse email template:", err)
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Step 2: Inject subject + body
	var htmlBody bytes.Buffer
	err = tmpl.Execute(&htmlBody, map[string]string{
		"Subject": subject,
		"Body":    body,
	})
	if err != nil {
		fmt.Println("❌ Failed to render email template:", err)
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Step 3: Build headers
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	from := fmt.Sprintf("%s <%s>", e.FromName, e.FromAddr)
	toHeader := strings.Join(to, ", ")

	headers := map[string]string{
		"From":         from,
		"To":           toHeader,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}

	var msgBuilder strings.Builder
	for k, v := range headers {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n" + htmlBody.String())

	// Step 4: Send email
	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)
	fmt.Println("📤 Sending email to:", to, "via", addr)

	err = smtp.SendMail(addr, auth, e.FromAddr, to, []byte(msgBuilder.String()))
	if err != nil {
		fmt.Println("❌ smtp.SendMail failed:", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("✅ Email sent successfully to:", to)
	return nil
}
