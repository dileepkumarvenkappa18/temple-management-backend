package event

import (
	"time"
)

// ============================
// 🔷 GORM Event Model
type Event struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	EntityID    uint       `gorm:"not null;index" json:"entity_id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	EventType   string     `gorm:"type:varchar(100);not null" json:"event_type"`
	EventDate   time.Time  `gorm:"not null;index" json:"event_date"`
	EventTime   *time.Time `json:"event_time,omitempty"`
	Location    string     `gorm:"type:text" json:"location"`
	CreatedBy   uint       `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`

	RSVPCount int `gorm:"-" json:"rsvp_count"`
}

// ============================
// 🟡 Create Event Request
type CreateEventRequest struct {
	EntityID  uint   `json:"entity_id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	EventType   string `json:"event_type" binding:"required"`
	EventDate   string `json:"event_date" binding:"required"` // 🛠 string format: "2006-01-02"
	EventTime   string `json:"event_time,omitempty"`          // 🛠 string format: "15:04"
	Location    string `json:"location" binding:"required"`
	IsActive *bool `json:"is_active,omitempty"`
}

// ============================
// 🟠 Update Event Request
type UpdateEventRequest struct {
	EntityID  uint   `json:"entity_id"`
	ID          uint   `json:"-"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	EventType   string `json:"event_type" binding:"required"`
	EventDate   string `json:"event_date" binding:"required"` // 🛠 string
	EventTime   string `json:"event_time,omitempty"`          // 🛠 string
	Location    string `json:"location" binding:"required"`
	IsActive *bool `json:"is_active,omitempty"`
}
