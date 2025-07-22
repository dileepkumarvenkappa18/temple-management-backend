package entity

import "time"

type Entity struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Step 1: Temple Basic Information
	Name            string  `gorm:"not null" json:"name"`        // Temple Name
	MainDeity       *string `json:"main_deity"`                  // Nullable
	TempleType      string  `gorm:"not null" json:"temple_type"` // Temple Type
	EstablishedYear *uint   `json:"established_year"`
	// Optional
	Email       string `gorm:"unique;not null" json:"email"` // Temple Email
	Phone       string `gorm:"not null" json:"phone"`        // Contact Phone
	Description string `json:"description"`                  // Temple Description
  
	// Step 2: Address Information
	StreetAddress string `json:"street_address"` // No 'not null' tag
	// Street Address
	Landmark string `json:"landmark"` // Optional landmark
	City     string `gorm:"not null" json:"city"`
	District string `gorm:"not null" json:"district"`
	State    string `gorm:"not null" json:"state"`
	Pincode  string `gorm:"not null" json:"pincode"`
	MapLink  string `json:"map_link"` // Optional Google Map link

	// Step 3: Document Uploads
	RegistrationCertURL string `json:"registration_cert_url"` // Required – S3/Cloudinary URL
	TrustDeedURL        string `json:"trust_deed_url"`        // Required – S3/Cloudinary URL
	PropertyDocsURL     string `json:"property_docs_url"`     // Optional
	AdditionalDocsURLs  string `json:"additional_docs_urls"`  // Optional – JSON string

	// Terms and Verification
	AcceptedTerms bool   `gorm:"default:false" json:"accepted_terms"`
	Status        string `gorm:"default:'pending'" json:"status"` // pending / approved / rejected
	CreatedBy     uint   `gorm:"not null" json:"created_by"`

	// Meta (optional)
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}