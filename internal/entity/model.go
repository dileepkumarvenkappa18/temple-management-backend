package entity

import "time"

// Entity represents a temple's core data structure
type Entity struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Step 1: Temple Basic Information
	Name            string  `gorm:"not null" json:"name"`                // Temple Name
	MainDeity       *string `gorm:"type:varchar(100)" json:"main_deity"` // Nullable
	TempleType      string  `gorm:"not null" json:"temple_type"`         // Type (e.g., Shiva, Vishnu)
	EstablishedYear *uint   `json:"established_year"`                    // Optional year

	Email       string `gorm:"unique;not null" json:"email"` // Temple Email
	Phone       string `gorm:"not null" json:"phone"`        // Contact Number
	Description string `gorm:"type:text" json:"description"` // Temple Description (optional)

	// Step 2: Address Information
	StreetAddress string `gorm:"type:varchar(255)" json:"street_address"` // Optional
	Landmark      string `gorm:"type:varchar(255)" json:"landmark"`       // Optional landmark
	City          string `gorm:"not null" json:"city"`                    // Required
	District      string `gorm:"not null" json:"district"`                // Required
	State         string `gorm:"not null" json:"state"`                   // Required
	Pincode       string `gorm:"not null" json:"pincode"`                 // Required
	MapLink       string `gorm:"type:text" json:"map_link"`               // Optional Google Maps URL

	// Step 3: Document Uploads
	RegistrationCertURL string `gorm:"not null" json:"registration_cert_url"` // Required – S3/Cloudinary URL
	TrustDeedURL        string `gorm:"not null" json:"trust_deed_url"`        // Required
	PropertyDocsURL     string `json:"property_docs_url"`                     // Optional
	AdditionalDocsURLs  string `json:"additional_docs_urls"`                  // Optional (comma-separated or JSON string)

	// Terms and Verification
	AcceptedTerms bool   `gorm:"default:false" json:"accepted_terms"` // Must accept terms
	Status        string `gorm:"default:'pending'" json:"status"`     // pending / approved / rejected
	CreatedBy     uint   `gorm:"not null" json:"created_by"`          // Admin/User ID who created

	// Metadata
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserEntity struct {
	UserID   uint `gorm:"column:user_id"`
	EntityID uint `gorm:"column:entity_id"`
}