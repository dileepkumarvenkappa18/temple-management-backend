// ==================== UPDATED ENTITY MODEL ====================
// entity/models.go
package entity

import (
	"time"
)

type Entity struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Step 1: Temple Basic Information
	Name            string  `gorm:"not null" json:"name"`
	MainDeity       *string `json:"main_deity"`
	TempleType      string  `gorm:"not null" json:"temple_type"`
	EstablishedYear *uint   `json:"established_year"`
	Email           string  `gorm:"unique;not null" json:"email"`
	Phone           string  `gorm:"not null" json:"phone"`
	Description     string  `json:"description"`

	// Step 2: Address Information
	StreetAddress string `gorm:"not null" json:"street_address"` // Make it required
	Landmark      string `json:"landmark"`
	City          string `gorm:"not null" json:"city"`
	District      string `gorm:"not null" json:"district"`
	State         string `gorm:"not null" json:"state"`
	Pincode       string `gorm:"not null" json:"pincode"`
	MapLink       string `json:"map_link"`

	// Step 3: Document Uploads (URLs to stored files)
	RegistrationCertURL string `json:"registration_cert_url"`
	TrustDeedURL        string `json:"trust_deed_url"`
	PropertyDocsURL     string `json:"property_docs_url"`
	AdditionalDocsURLs  string `json:"additional_docs_urls"` // JSON string of array

	// File metadata (new fields to store file information)
	RegistrationCertInfo string `json:"registration_cert_info"` // JSON metadata
	TrustDeedInfo        string `json:"trust_deed_info"`        // JSON metadata
	PropertyDocsInfo     string `json:"property_docs_info"`     // JSON metadata
	AdditionalDocsInfo   string `json:"additional_docs_info"`   // JSON metadata

	// Terms and Verification
	AcceptedTerms bool   `gorm:"default:false" json:"accepted_terms"`
	Status        string `gorm:"default:'pending'" json:"status"`
	CreatedBy     uint   `gorm:"not null" json:"created_by"`

	// Meta
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// FileInfo represents metadata about an uploaded file
type FileInfo struct {
	FileName     string `json:"file_name"`
	FileURL      string `json:"file_url"`
	FileSize     int64  `json:"file_size"`
	FileType     string `json:"file_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	OriginalName string `json:"original_name"`
}