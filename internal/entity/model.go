<<<<<<< HEAD
// ==================== UPDATED ENTITY MODEL ====================
// entity/models.go
package entity

import (
	"time"
)
=======
package entity

import "time"
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6

type Entity struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Step 1: Temple Basic Information
<<<<<<< HEAD
	Name            string  `gorm:"not null" json:"name"`
	MainDeity       *string `json:"main_deity"`
	TempleType      string  `gorm:"not null" json:"temple_type"`
	EstablishedYear *uint   `json:"established_year"`
	Email           string  `gorm:"unique;not null" json:"email"`
	Phone           string  `gorm:"not null" json:"phone"`
	Description     string  `json:"description"`

	// Step 2: Address Information
	StreetAddress string `gorm:"not null" json:"street_address"`
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
	
	// 🆕 NEW FIELD: Track the role_id of the user who created this temple
	// This is used for auto-approval logic (role_id = 1 for superadmin)
	CreatorRoleID *uint  `json:"creator_role_id" gorm:"index"` // Role ID of creator (1=superadmin for auto-approval)

	// 🆕 NEW FIELD: Active/Inactive status
	IsActive bool `gorm:"default:true" json:"isactive"` // Active/Inactive toggle

	// Meta
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// FileInfo represents metadata about an uploaded file
type FileInfo struct {
	FileName     string    `json:"file_name"`
	FileURL      string    `json:"file_url"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	OriginalName string    `json:"original_name"`
}

// TableName specifies the table name for the Entity model
func (Entity) TableName() string {
	return "entities"
}
=======
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
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
