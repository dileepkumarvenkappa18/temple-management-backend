package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// SMTP Configuration
var (
	smtpHost      = os.Getenv("SMTP_HOST")
	smtpPort      = os.Getenv("SMTP_PORT")
	smtpUsername  = os.Getenv("SMTP_USERNAME")
	smtpPassword  = os.Getenv("SMTP_PASSWORD")
	smtpFromName  = os.Getenv("SMTP_FROM_NAME")
	smtpFromEmail = os.Getenv("SMTP_FROM_EMAIL")
	frontendURL   = os.Getenv("FRONTEND_URL") // New environment variable for frontend URL
)

// sendEmail handles the actual SMTP connection and sending
func sendEmail(to, subject, body string) error {
	// First, always log the email for debugging
	fmt.Println("📧 Email Details:")
	fmt.Printf("To      : %s\n", to)
	fmt.Printf("Subject : %s\n", subject)
	fmt.Printf("Body    : %s\n", body)
	
	// Check if SMTP is configured
	if smtpHost == "" || smtpUsername == "" || smtpPassword == "" {
		fmt.Println("⚠️ SMTP not configured. Email not sent.")
		return nil
	}
	
	// Fix any configuration issues
	if smtpFromEmail == "" {
		smtpFromEmail = smtpUsername
	}
	
	// Remove any typos in email addresses
	smtpFromEmail = strings.TrimSuffix(smtpFromEmail, "i") // Fix common typo
	
	// Setup auth
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	
	// Compose email
	from := smtpFromName
	if from == "" {
		from = smtpFromEmail
	} else {
		from = fmt.Sprintf("%s <%s>", smtpFromName, smtpFromEmail)
	}
	
	// Format email message with headers
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", from, to, subject, body))
	
	// Connect and send
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	
	fmt.Println("📤 Attempting to send email via SMTP...")
	err := smtp.SendMail(addr, auth, smtpFromEmail, []string{to}, msg)
	if err != nil {
		fmt.Printf("❌ SMTP Error: %v\n", err)
		return err
	}
	
	fmt.Println("✅ Email sent successfully via SMTP!")
	return nil
}

// ====================== 🔐 Password Reset ======================

// SendResetLink sends a password reset link to the user's email
func SendResetLink(toEmail string, resetToken string) error {
	// Use environment variable with fallback to default development URL
	baseURL := frontendURL
	if baseURL == "" {
		// Default to localhost for development environment
		baseURL = "http://localhost:8080"
		fmt.Println("⚠️ FRONTEND_URL not set, using default:", baseURL)
	}
	
	// Changed from /reset-password to /#/reset-password to work with Vue Router
	resetURL := fmt.Sprintf("%s/auth-pages/reset-password?token=%s", baseURL, resetToken)
	subject := "Reset your password"
	body := fmt.Sprintf("Click here to reset your password: %s\n\nIf you did not request this password reset, please ignore this email.", resetURL)
	
	return sendEmail(toEmail, subject, body)
}

// ====================== ✅ Tenant Approval ======================

func SendTenantApprovalEmail(toEmail, fullName string) {
	subject := "Your account has been approved"
	body := fmt.Sprintf("Hello %s, your account has been approved by the Super Admin. You can now log in and manage your temple.", fullName)
	
	_ = sendEmail(toEmail, subject, body)
}

// ====================== ❌ Tenant Rejection ======================

func SendTenantRejectionEmail(toEmail, fullName, reason string) {
	subject := "Your account request was rejected"
	body := fmt.Sprintf("Hello %s, your account request was rejected by the Super Admin.\nReason: %s", fullName, reason)
	
	_ = sendEmail(toEmail, subject, body)
}

// SendPasswordResetNotification sends an email to notify user about password reset
func SendPasswordResetNotification(toEmail, userName, adminName, newPassword string) error {
    subject := "Your password has been reset"
    body := fmt.Sprintf("Hello %s, your password has been reset by %s. If you did not request this change, please contact support immediately.\n\nYour new password is: %s\n\nFor security reasons, please change your password after logging in.", userName, adminName, newPassword)
    
    return sendEmail(toEmail, subject, body)
}

// ====================== 🏛️ Entity Approval ======================

func SendEntityApprovalEmail(toEmail, fullName, templeName string) {
	subject := fmt.Sprintf("Your Temple \"%s\" Has Been Approved", templeName)
	body := fmt.Sprintf("Hello %s, your temple \"%s\" has been successfully approved. You can now manage it on the platform.", fullName, templeName)
	
	_ = sendEmail(toEmail, subject, body)
}

// ====================== 🏛️ Entity Rejection ======================

func SendEntityRejectionEmail(toEmail, fullName, templeName, reason string) {
	subject := fmt.Sprintf("Your Temple \"%s\" Was Rejected", templeName)
	body := fmt.Sprintf("Hello %s, your temple \"%s\" was rejected.\nReason: %s", fullName, templeName, reason)
	
	_ = sendEmail(toEmail, subject, body)
}
