package auth

import (
	"time"

	"gorm.io/gorm"
)

// UserRole maps to user_roles
type UserRole struct {
	ID                  uint           `gorm:"primaryKey;autoIncrement"`
	RoleName            string         `gorm:"size:50;unique;not null"` // tenant, devotee, volunteer, super_admin
	Description         string         `gorm:"type:text"`
	CanRegisterPublicly bool           `gorm:"default:true"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// User maps to users table
type User struct {
	ID              uint           `gorm:"primaryKey;autoIncrement"`
	FullName        string         `gorm:"size:255;not null"`
	Email           string         `gorm:"size:255;unique;not null"`
	PasswordHash    string         `gorm:"size:255;not null"`
	Phone           *string        `gorm:"size:20"`
	RoleID          uint           `gorm:"not null"`
	Role            UserRole       `gorm:"foreignKey:RoleID;references:ID"`
	EntityID        *uint          `gorm:"index"` // Optional: linked only after approval
	Status          string         `gorm:"size:20;default:'active'"` // active, pending, rejected, inactive
	EmailVerified   bool           `gorm:"default:false"`
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// ApprovalRequest maps to approval_requests
type ApprovalRequest struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	UserID      uint           `gorm:"not null"`
	User        User           `gorm:"foreignKey:UserID;references:ID"`
	RequestType string         `gorm:"size:50;not null"` // tenant_approval, temple_approval
	EntityID    *uint
	Status      string         `gorm:"size:20;default:'pending'"` // pending, approved, rejected
	AdminNotes  *string        `gorm:"type:text"`
	ApprovedBy  *uint
	ApprovedAt  *time.Time
	RejectedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
