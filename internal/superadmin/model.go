package superadmin

import (
	"time"
)

type TenantApprovalCount struct {
	Approved int64 `json:"approved"`
	Pending  int64 `json:"pending"`
	Rejected int64 `json:"rejected"`
}

// ================ TEMPLE APPROVAL COUNTS ================

type TempleApprovalCount struct {
	PendingApproval int64 `json:"pending_approval"`
	ActiveTemples   int64 `json:"active_temples"`
	Rejected        int64 `json:"rejected"`
	TotalDevotees   int64 `json:"total_devotees"`
}

// ================ TENANT WITH TEMPLE DETAILS ================

type TenantWithDetails struct {
	// User details
	ID           uint      `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	RoleID       uint      `json:"role_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Temple details
	TempleDetails *TenantTempleDetails `json:"temple_details,omitempty"`
}

type TenantTempleDetails struct {
	ID                uint      `json:"id"`
	TempleName        string    `json:"temple_name"`
	TemplePlace       string    `json:"temple_place"`
	TempleAddress     string    `json:"temple_address"`
	TemplePhoneNo     string    `json:"temple_phone_no"`
	TempleDescription string    `json:"temple_description"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}