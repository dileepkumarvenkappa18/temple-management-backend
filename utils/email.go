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
	smtpTimeout   = 10 * time.Second
)

// ======================
// Low-level sendEmail
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
	smtpFromEmail = strings.TrimSuffix(smtpFromEmail, "i")

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
	resetURL := fmt.Sprintf("%s/auth-pages/reset-password?token=%s", frontendURL, resetToken)
	subject := "Reset your password"
	body := fmt.Sprintf(
		"Hello,\n\nClick the link below to reset your password:\n%s\n\nIf you did not request this, please ignore this email.\n\nRegards,\nTemple Management Team",
		resetURL,
	)
	return sendEmail(toEmail, subject, body)
}

// ======================
// Tenant (Temple Admin) Emails
// ======================

// SendTenantApprovalEmail is called when superadmin approves a temple admin account
// SendTenantApprovalEmail is called when superadmin approves a temple admin account
func SendTenantApprovalEmail(toEmail, fullName string) {
	subject := "Your Temple Admin Account Has Been Approved"
	body := fmt.Sprintf(
		"Hello %s,\n\n"+
			"Congratulations! Your Temple Admin account has been approved by the Super Admin.\n\n"+
			"You can now log in and start managing your temple.\n\n"+
			// "Your Password : %s\n\n"+ // TODO: uncomment and pass password when needed
			"Login here: %s\n\n"+
			"Regards,\n"+
			"Temple Management Team",
		fullName,
		frontendURL,
	)
	if err := sendEmail(toEmail, subject, body); err != nil {
		fmt.Printf("❌ Failed to send tenant approval email to %s: %v\n", toEmail, err)
	}
}

// SendTenantRejectionEmail is called when superadmin rejects a temple admin account
func SendTenantRejectionEmail(toEmail, fullName, reason string) {
	subject := "Your Temple Admin Account Request Was Rejected"
	body := fmt.Sprintf(
		"Hello %s,\n\n"+
			"We regret to inform you that your Temple Admin account request has been rejected by the Super Admin.\n\n"+
			"Reason: %s\n\n"+
			// "Your Password : %s\n\n"+ // TODO: uncomment and pass password when needed
			"If you believe this is a mistake or would like to re-apply, please visit:\n%s\n\n"+
			"Regards,\n"+
			"Temple Management Team",
		fullName,
		reason,
		frontendURL,
	)
	if err := sendEmail(toEmail, subject, body); err != nil {
		fmt.Printf("❌ Failed to send tenant rejection email to %s: %v\n", toEmail, err)
	}
}

// ======================
// Entity (Temple) Emails
// ======================

// SendEntityApprovalEmail is called when superadmin approves a temple (entity)
// toEmail and templeName come directly from entity.Email and entity.Name
func SendEntityApprovalEmail(toEmail, templeName, _ string) {
	subject := fmt.Sprintf("Your Temple \"%s\" Has Been Approved", templeName)
	body := fmt.Sprintf(
		"Hello,\n\n"+
			"Great news! Your temple \"%s\" has been approved by the Super Admin.\n\n"+
			"You can now manage your temple on the platform.\n\n"+
			"Login here: %s/login\n\n"+
			"Regards,\n"+
			"Temple Management Team",
		templeName,
		frontendURL,
	)
	if err := sendEmail(toEmail, subject, body); err != nil {
		fmt.Printf("❌ Failed to send entity approval email to %s: %v\n", toEmail, err)
	}
}

// SendEntityRejectionEmail is called when superadmin rejects a temple (entity)
// toEmail and templeName come directly from entity.Email and entity.Name
func SendEntityRejectionEmail(toEmail, templeName, _ string, reason string) {
	subject := fmt.Sprintf("Your Temple \"%s\" Was Rejected", templeName)
	body := fmt.Sprintf(
		"Hello,\n\n"+
			"We regret to inform you that your temple \"%s\" has been rejected by the Super Admin.\n\n"+
			"Reason: %s\n\n"+
			"If you believe this is a mistake or would like to re-submit, please log in at:\n%s/login\n\n"+
			"Regards,\n"+
			"Temple Management Team",
		templeName,
		reason,
		frontendURL,
	)
	if err := sendEmail(toEmail, subject, body); err != nil {
		fmt.Printf("❌ Failed to send entity rejection email to %s: %v\n", toEmail, err)
	}
}

// ======================
// Password Reset Notification
// ======================
func SendPasswordResetNotification(toEmail, userName, adminName, newPassword string) error {
	subject := "Your Password Has Been Reset"
	body := fmt.Sprintf(
		"Hello %s,\n\n"+
			"Your password has been reset by %s.\n\n"+
			"New Password: %s\n\n"+
			"Please log in and change your password immediately:\n%s\n\n"+
			"Regards,\n"+
			"Temple Management Team",
		userName,
		adminName,
		newPassword,
		frontendURL,
	)
	return sendEmail(toEmail, subject, body)
}

// ======================
// Welcome Email for Admin-Created Users
// ======================
func SendWelcomeEmail(toEmail, fullName, password, role string) {
	subject := "Welcome - Your Account Has Been Created"
	body := fmt.Sprintf(
		"Hello %s,\n\n"+
			"Your account has been created by the administrator.\n\n"+
			"Login Credentials:\n"+
			"  Email    : %s\n"+
			"  Password : %s\n"+
			"  Role     : %s\n\n"+
			"Login here: %s\n\n"+
			"For your security, please change your password after your first login.\n\n"+
			"Regards,\n"+
			"Temple Management Team",
		fullName,
		toEmail,
		password,
		role,
		frontendURL,
	)
	if err := sendEmail(toEmail, subject, body); err != nil {
		fmt.Printf("❌ Failed to send welcome email to %s: %v\n", toEmail, err)
	}
}

// SendRegistrationEmail is called when any user self-registers
// status will be "active" for normal users, "pending" for templeadmin
func SendRegistrationEmail(toEmail, fullName, role, status string) {
	var subject string
	var body string

	loginURL := fmt.Sprintf("%s/login", frontendURL)

	if status == "pending" {
		// Temple admin — account needs superadmin approval
		subject = "Your Registration Is Under Review"
		body = fmt.Sprintf(
			"Hello %s,\n\n"+
				"Thank you for registering as a Temple Admin on our platform.\n\n"+
				"Your account is currently under review by the Super Admin.\n"+
				"You will receive another email once your account is approved.\n\n"+
				// "Your Password : %s\n\n"+ // TODO: uncomment and pass password when needed
				"Once approved, you can log in here: %s\n\n"+
				"Regards,\n"+
				"Temple Management Team",
			fullName,
			loginURL,
		)
	} else {
		// All other roles — account is immediately active
		subject = "Welcome! Your Account Has Been Created"
		body = fmt.Sprintf(
			"Hello %s,\n\n"+
				"Welcome to the Temple Management Platform!\n\n"+
				"Your account has been successfully created.\n\n"+
				"Role  : %s\n"+
				"Email : %s\n\n"+
				// "Password : %s\n\n"+ // TODO: uncomment and pass password when needed
				"Login here: %s\n\n"+
				"Regards,\n"+
				"Temple Management Team",
			fullName,
			role,
			toEmail,
			loginURL,
		)
	}

	if err := sendEmail(toEmail, subject, body); err != nil {
		fmt.Printf("❌ Failed to send registration email to %s: %v\n", toEmail, err)
	}
}