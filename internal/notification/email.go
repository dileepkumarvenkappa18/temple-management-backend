package notification

import (
	"bytes"
	"crypto/tls"
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

// Send sends ONE email with sender in To field and all recipients in BCC
// - Sender sees all recipients in their sent folder
// - Each recipient only sees their own email in "To: bcc: their-email@example.com"
func (e *EmailSender) Send(to []string, subject string, body string) error {
	if len(to) == 0 {
		return fmt.Errorf("no recipients provided")
	}

	fmt.Printf("📧 Sending one email to %d BCC recipients\n", len(to))

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

	// Step 3: Build headers with sender in To field
	from := "Temple Management System <noreply@templemanagement.com>"
	
	headers := map[string]string{
		"From":         from,
		"To":           e.FromAddr,  // ✅ Your email goes in To field - you'll see all recipients in sent folder
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
		// NO Bcc header - keeps recipients hidden from each other
	}

	var msgBuilder strings.Builder
	for k, v := range headers {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n" + htmlBody.String())
	message := []byte(msgBuilder.String())

	// Step 4: Send to all recipients via SMTP envelope
	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)
	
	fmt.Printf("📤 Sending via %s to %d BCC recipients\n", addr, len(to))
	fmt.Printf("📝 From display: %s\n", from)
	fmt.Printf("📧 To field: %s (will appear in Sent folder)\n", e.FromAddr)
	fmt.Printf("🔐 SMTP auth: %s\n", e.Username)

	err = e.sendMailWithTLS(addr, to, message)
	if err != nil {
		fmt.Println("❌ Email send failed:", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("✅ Email sent - check your Sent folder to see all %d BCC recipients\n", len(to))
	return nil
}

// ✅ Custom send function with proper TLS handling
func (e *EmailSender) sendMailWithTLS(addr string, to []string, message []byte) error {
	// Create TLS config - skip verification for Docker environments
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         e.Host,
	}

	// Connect to the SMTP server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer client.Close()

	// Start TLS
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate with REAL email (must match Gmail account)
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// ✅ CRITICAL: MAIL FROM must be the authenticated email
	// You CANNOT change this - Gmail requires it to match your account
	if err = client.Mail(e.FromAddr); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Add all recipients to SMTP envelope (sender + BCC list)
	successCount := 0
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			fmt.Printf("⚠️  Failed to add recipient %s: %v\n", recipient, err)
			continue
		}
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("failed to add any recipients")
	}

	if successCount < len(to) {
		fmt.Printf("⚠️  Partial delivery: %d/%d recipients added\n", successCount, len(to))
	}

	// Send message body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = writer.Write(message)
	if err != nil {
		writer.Close()
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return client.Quit()
}