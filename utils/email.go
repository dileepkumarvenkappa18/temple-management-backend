package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"sync"
	"time"
)

// ======================
// SMTP Configuration
// ======================
var (
	smtpHost      = os.Getenv("SMTP_HOST")
	smtpPort      = os.Getenv("SMTP_PORT")
	smtpUsername  = os.Getenv("SMTP_USERNAME")
	smtpPassword  = os.Getenv("SMTP_PASSWORD")
	smtpFromName  = os.Getenv("SMTP_FROM_NAME")
	smtpFromEmail = os.Getenv("SMTP_FROM_EMAIL")
	frontendURL   = os.Getenv("FRONTEND_URL")
	smtpTimeout   = 10 * time.Second // Timeout for SMTP connection
)

// ======================
// ✅ FIXED: Low-level sendEmail with proper TLS handling
// ======================
func sendEmail(to, subject, body string) error {
	fmt.Println("📧 Sending Email:")
	fmt.Printf("To      : %s\nSubject : %s\nBody    : %s\n", to, subject, body)

	if smtpHost == "" || smtpUsername == "" || smtpPassword == "" {
		fmt.Println("⚠️ SMTP not configured. Email not sent.")
		return nil
	}

	if smtpFromEmail == "" {
		smtpFromEmail = smtpUsername
	}
	smtpFromEmail = strings.TrimSuffix(smtpFromEmail, "i") // Remove accidental typo

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	client, err := smtp.Dial(addr)
	if err != nil {
		fmt.Printf("❌ Failed to dial SMTP server: %v\n", err)
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer client.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		fmt.Printf("❌ TLS connection error: %v\n", err)
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	if err := client.Auth(auth); err != nil {
		fmt.Printf("❌ SMTP auth error: %v\n", err)
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := client.Mail(smtpFromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	from := smtpFromName
	if from == "" {
		from = smtpFromEmail
	} else {
		from = fmt.Sprintf("%s <%s>", smtpFromName, smtpFromEmail)
	}

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n%s", from, to, subject, body))

	_, err = w.Write(msg)
	if err != nil {
		w.Close()
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	if err := client.Quit(); err != nil {
		fmt.Printf("⚠️ QUIT command error (non-critical): %v\n", err)
	}

	fmt.Println("✅ Email sent successfully!")
	return nil
}

// ======================
// Async bulk email sender
// ======================
func SendBulkEmailsAsync(recipients []string, subject, body string) {
	go func() {
		var wg sync.WaitGroup
		for _, email := range recipients {
			wg.Add(1)
			go func(to string) {
				defer wg.Done()
				if err := sendEmail(to, subject, body); err != nil {
					fmt.Printf("❌ Failed to send email to %s: %v\n", to, err)
				} else {
					fmt.Printf("✅ Email sent to %s\n", to)
				}
			}(email)
		}
		wg.Wait()
	}()
}

// ======================
// Password Reset
// ======================
func SendResetLink(toEmail string, resetToken string) error {
	baseURL := frontendURL

	resetURL := fmt.Sprintf("%s/auth-pages/reset-password?token=%s", baseURL, resetToken)
	subject := "Reset your password"
	body := fmt.Sprintf("Click here to reset your password: %s\n\nIf you did not request this password reset, please ignore this email.", resetURL)

	return sendEmail(toEmail, subject, body)
}

// ======================
// Tenant Emails
// ======================
func SendTenantApprovalEmail(toEmail, fullName string) {
	subject := "Your account has been approved"
	body := fmt.Sprintf("Hello %s, your account has been approved by the Super Admin. You can now log in and manage your temple.", fullName)
	_ = sendEmail(toEmail, subject, body)
}

func SendTenantRejectionEmail(toEmail, fullName, reason string) {
	subject := "Your account request was rejected"
	body := fmt.Sprintf("Hello %s, your account request was rejected by the Super Admin.\nReason: %s", fullName, reason)
	_ = sendEmail(toEmail, subject, body)
}

// Password reset notification
func SendPasswordResetNotification(toEmail, userName, adminName, newPassword string) error {
	subject := "Your password has been reset"
	body := fmt.Sprintf("Hello %s, your password has been reset by %s.\n\nNew password: %s\n\nPlease change it after logging in.", userName, adminName, newPassword)
	return sendEmail(toEmail, subject, body)
}

// ======================
// Entity Emails
// ======================
func SendEntityApprovalEmail(toEmail, fullName, templeName string) {
	subject := fmt.Sprintf("Your Temple \"%s\" Has Been Approved", templeName)
	body := fmt.Sprintf("Hello %s, your temple \"%s\" has been successfully approved. You can now manage it on the platform.", fullName, templeName)
	_ = sendEmail(toEmail, subject, body)
}

func SendEntityRejectionEmail(toEmail, fullName, templeName, reason string) {
	subject := fmt.Sprintf("Your Temple \"%s\" Was Rejected", templeName)
	body := fmt.Sprintf("Hello %s, your temple \"%s\" was rejected.\nReason: %s", fullName, templeName, reason)
	_ = sendEmail(toEmail, subject, body)
}

// ======================
// ✅ NEW: Welcome Email for Admin-Created Users
// ======================
func SendWelcomeEmail(toEmail, fullName, password, role string) {
	subject := "Your account has been created"
	body := fmt.Sprintf(
		"Hello %s,\n\nYour account has been created by the administrator.\n\nHere are your login credentials:\n\nEmail    : %s\nPassword : %s\nRole     : %s\n\nPlease log in at: %s\n\nFor your security, please change your password after your first login.\n\nBest regards,\nTemple Management Team",
		fullName, toEmail, password, role, frontendURL,
	)
	_ = sendEmail(toEmail, subject, body)
}