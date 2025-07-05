package utils

import (
	"fmt"
)

// SendResetLink simulates sending a password reset link to the user's email.
// In production, you'd integrate SMTP, SendGrid, Mailgun, etc.
func SendResetLink(toEmail string, resetToken string) error {
	resetURL := fmt.Sprintf("https://yourfrontend.com/reset-password?token=%s", resetToken)

	// Simulated email log — replace with real SMTP/email logic later
	fmt.Println("📬 Sending password reset email:")
	fmt.Printf("To      : %s\n", toEmail)
	fmt.Printf("Subject : Reset your password\n")
	fmt.Printf("Body    : Click the link to reset your password: %s\n", resetURL)

	// TODO: Replace this print with actual email-sending logic
	return nil
}
