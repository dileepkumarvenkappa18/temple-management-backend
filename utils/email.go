package utils

import (
	"fmt"
)

// ====================== 🔐 Password Reset ======================

// SendResetLink simulates sending a password reset link to the user's email.
// In production, you'd integrate SMTP, SendGrid, Mailgun, etc.
func SendResetLink(toEmail string, resetToken string) error {
	resetURL := fmt.Sprintf("https://yourfrontend.com/reset-password?token=%s", resetToken)

	fmt.Println("📬 Password Reset Email")
	fmt.Printf("To      : %s\n", toEmail)
	fmt.Printf("Subject : Reset your password\n")
	fmt.Printf("Body    : Click here to reset your password: %s\n", resetURL)

	return nil
}

// ====================== ✅ Tenant Approval ======================

func SendTenantApprovalEmail(toEmail, fullName string) {
	fmt.Println("✅ Tenant Approved Email")
	fmt.Printf("To      : %s\n", toEmail)
	fmt.Printf("Subject : Your account has been approved\n")
	fmt.Printf("Body    : Hello %s, your account has been approved by the Super Admin. You can now log in and manage your temple.\n", fullName)
}

// ====================== ❌ Tenant Rejection ======================

func SendTenantRejectionEmail(toEmail, fullName, reason string) {
	fmt.Println("❌ Tenant Rejected Email")
	fmt.Printf("To      : %s\n", toEmail)
	fmt.Printf("Subject : Your account request was rejected\n")
	fmt.Printf("Body    : Hello %s, your account request was rejected by the Super Admin.\nReason: %s\n", fullName, reason)
}

// ====================== 🏛️ Entity Approval ======================

func SendEntityApprovalEmail(toEmail, fullName, templeName string) {
	fmt.Println("🏛️ Entity Approved Email")
	fmt.Printf("To      : %s\n", toEmail)
	fmt.Printf("Subject : Your Temple \"%s\" Has Been Approved\n", templeName)
	fmt.Printf("Body    : Hello %s, your temple \"%s\" has been successfully approved. You can now manage it on the platform.\n", fullName, templeName)
}

// ====================== 🏛️ Entity Rejection ======================

func SendEntityRejectionEmail(toEmail, fullName, templeName, reason string) {
	fmt.Println("🚫 Entity Rejected Email")
	fmt.Printf("To      : %s\n", toEmail)
	fmt.Printf("Subject : Your Temple \"%s\" Was Rejected\n", templeName)
	fmt.Printf("Body    : Hello %s, your temple \"%s\" was rejected.\nReason: %s\n", fullName, templeName, reason)
}
