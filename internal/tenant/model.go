package tenant

import (
    "time"

    "gorm.io/gorm"
)

// TenantUser represents a user in a multi-tenant system.
type TenantUser struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    TenantID  uint           `json:"-"` // hidden from request/response; set from tenant context
    Name      string         `gorm:"size:100;not null" json:"name"`
    Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
    Phone     string         `gorm:"size:20" json:"phone"`
    Role      string         `gorm:"size:50;not null" json:"role"` // "Standard User" or "Monitoring User"
    Password  string         `gorm:"size:255;not null" json:"-"`   // stored hashed
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // optional soft delete
}
