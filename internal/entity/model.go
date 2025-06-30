package entity

import "time"

// ========== ENTITY (Temple Info Only) ==========
type Entity struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	TempleName      string     `gorm:"not null" json:"temple_name"`        // Temple full name
	EntityCode      string     `gorm:"unique;not null" json:"entity_code"` // Internal code
	Email           string     `gorm:"unique;not null" json:"email"`
	Mobile          string     `json:"mobile"`
	TempleType      string     `json:"temple_type"` // Hindu, Jain, etc.
	EstablishedDate *time.Time `json:"established_date"`
	ContactPerson   string     `json:"contact_person"` // Admin or manager

	LogoURL   string    `json:"logo_url"` // Path to logo/image
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Address   EntityAddress    `gorm:"foreignKey:EntityID" json:"address,omitempty"`
	Documents []EntityDocument `gorm:"foreignKey:EntityID" json:"documents,omitempty"`
}

// ========== ENTITY ADDRESS ==========
type EntityAddress struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	EntityID     uint   `gorm:"not null;uniqueIndex" json:"entity_id"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	District     string `json:"district"`
	State        string `json:"state"`
	Pincode      string `json:"pincode"`
	Country      string `gorm:"default:'India'" json:"country"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	AddressType  string `json:"address_type"` // Primary / Permanent etc.
}

// ========== ENTITY DOCUMENT ==========
type EntityDocument struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EntityID       uint      `gorm:"not null;index" json:"entity_id"`
	DocumentType   string    `gorm:"not null" json:"document_type"` // Aadhaar, PAN
	DocumentTitle  string    `json:"document_title"`                // "Main Trust Deed"
	DocumentURL    string    `gorm:"not null" json:"document_url"`  // File path
	DocumentNumber string    `json:"document_number"`               // PAN/Aadhaar number etc.
	UploadedAt     time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
}
