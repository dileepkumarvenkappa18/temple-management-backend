package event

import "time"

// ============================
// 🔷 GORM Event Model
type Event struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	EntityID             uint       `gorm:"not null;index" json:"entity_id"` // Temple organizing the event
	Title                string     `gorm:"type:varchar(255);not null" json:"title"`
	Description          string     `gorm:"type:text" json:"description"`
	EventDate            time.Time  `gorm:"not null;index" json:"event_date"` // Required
	EventTime            *time.Time `json:"event_time"`                       // Nullable
	EndDate              *time.Time `json:"end_date"`                         // Nullable
	EndTime              *time.Time `json:"end_time"`   
	Category string `gorm:"type:varchar(50)" json:"category"`                      // Nullable
	Location             string     `gorm:"type:text" json:"location"`
	ImageURL             string     `gorm:"type:text" json:"image_url"`
	MaxAttendees         *int       `json:"max_attendees"` // Nullable (nil = unlimited)
	RegistrationRequired bool       `gorm:"default:false" json:"registration_required"`
	IsActive             bool       `gorm:"default:true;index" json:"is_active"`
	CreatedBy            uint       `gorm:"not null" json:"created_by"` // Temple admin ID
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// ✅ Dynamically calculated, not stored in DB
	RSVPCount int `gorm:"-" json:"rsvp_count"`
}

// ============================
// 🟦 Request Structs for Safety

type CreateEventRequest struct {
	Title                string     `json:"title" binding:"required"`
	Description          string     `json:"description" binding:"required"`
	EventDate            time.Time  `json:"event_date" binding:"required"`
	EventTime            *time.Time `json:"event_time,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	EndTime              *time.Time `json:"end_time,omitempty"`
	Category string `gorm:"type:varchar(50)" json:"category"`
	Location             string     `json:"location" binding:"required"`
	ImageURL             string     `json:"image_url,omitempty"`
	MaxAttendees         *int       `json:"max_attendees,omitempty"`
	RegistrationRequired bool       `json:"registration_required"`
}

type UpdateEventRequest struct {
	ID                   uint       `json:"-"`
	Title                string     `json:"title" binding:"required"`
	Description          string     `json:"description" binding:"required"`
	EventDate            time.Time  `json:"event_date" binding:"required"`
	EventTime            *time.Time `json:"event_time,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	EndTime              *time.Time `json:"end_time,omitempty"`
	Category string `gorm:"type:varchar(50)" json:"category"`
	Location             string     `json:"location" binding:"required"`
	ImageURL             string     `json:"image_url,omitempty"`
	MaxAttendees         *int       `json:"max_attendees,omitempty"`
	RegistrationRequired bool       `json:"registration_required"`
	IsActive             *bool      `json:"is_active,omitempty"` // Optional update field
}


// EventRSVP represents a user's RSVP to an event
// type EventRSVP struct {
// 	ID        uint      `gorm:"primaryKey" json:"id"`
// 	EventID   uint      `gorm:"not null;index" json:"event_id"`
// 	UserID    uint      `gorm:"not null;index" json:"user_id"`
// 	Status    string    `gorm:"type:varchar(20);default:'attending'" json:"status"` // 'attending', 'maybe', 'not_attending'
// 	Notes     string    `gorm:"type:text" json:"notes"`
// 	RSVPDate  time.Time `gorm:"autoCreateTime" json:"rsvp_date"`
// }
