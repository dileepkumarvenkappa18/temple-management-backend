package notification

import (
	"time"

	"gorm.io/datatypes"
)

// 1. NotificationTemplate - reusable messages
type NotificationTemplate struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `gorm:"not null;index" json:"user_id"`
	Name     string `gorm:"size:100;not null;index:idx_name_entity,unique" json:"name"`
	EntityID uint   `gorm:"not null;index:idx_name_entity,unique" json:"entity_id"`

	Channel   string    `gorm:"size:20;not null" json:"channel"`   // email, sms, whatsapp
	Subject   string    `gorm:"size:255" json:"subject,omitempty"` // optional for email
	Body      string    `gorm:"type:text;not null" json:"body"`    // Go template format
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 2. NotificationLog - each actual message sent
type NotificationLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`      // sender
	EntityID   uint           `gorm:"not null;index" json:"entity_id"`    // temple context
	TemplateID *uint          `gorm:"index" json:"template_id,omitempty"` // optional
	Channel    string         `gorm:"size:20;not null" json:"channel"`    // email, sms, whatsapp
	Subject    string         `gorm:"size:255" json:"subject,omitempty"`
	Body       string         `gorm:"type:text;not null" json:"body"`
	Recipients datatypes.JSON `gorm:"type:jsonb;not null" json:"recipients"` // email/phone array
	Status     string         `gorm:"size:20;default:'pending'" json:"status"`
	Error      *string        `json:"error,omitempty"`
	IsRead     bool           `gorm:"default:false" json:"is_read"` // ✅ optional if you want read-tracking
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
