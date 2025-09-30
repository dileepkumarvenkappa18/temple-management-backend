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

	Category string `gorm:"size:20;not null" json:"category"` // email, sms, whatsapp

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

// 3. InAppNotification - per-user, in-app bell notifications
type InAppNotification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	EntityID  uint      `gorm:"not null;index" json:"entity_id"`
	Title     string    `gorm:"size:150;not null" json:"title"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Category  string    `gorm:"size:30;not null" json:"category"` // event, seva, donation, system
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}