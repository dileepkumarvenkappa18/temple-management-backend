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
	
	// Stores JSON: {"logo": "url/to/logo.jpg", "video": "url/to/video.mp4"}
	Media string `json:"media" gorm:"type:text"` // JSON string containing logo and video URLs

	// Terms and Verification
	AcceptedTerms bool   `gorm:"default:false" json:"accepted_terms"`
	Status        string `gorm:"default:'pending'" json:"status"`
	CreatedBy     uint   `gorm:"not null" json:"created_by"`
	
	// Track the role_id of the user who created this temple
	CreatorRoleID *uint  `json:"creator_role_id" gorm:"index"` // Role ID of creator (1=superadmin for auto-approval)

	// Active/Inactive status
	IsActive bool `gorm:"default:true" json:"isactive"` // Active/Inactive toggle
	ApprovedAt      *time.Time `json:"approved_at" gorm:"column:approved_at"`
	RejectedAt      *time.Time `json:"rejected_at" gorm:"column:rejected_at"`
	RejectionReason string     `json:"rejection_reason" gorm:"column:rejection_reason;type:text"`

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

// 🆕 MediaInfo represents temple logo and video
type MediaInfo struct {
	Logo  string `json:"logo,omitempty"`  // URL to logo image
	Video string `json:"video,omitempty"` // URL to video file
}

// TableName specifies the table name for the Entity model
func (Entity) TableName() string {
	return "entities"
}
// Add these structs to entity/model.go

// CreatorDetails represents the temple creator's public information


type CreatorTempleInfo struct {
	TempleName        string `json:"temple_name"`
	TemplePlace       string `json:"temple_place"`
	TempleAddress     string `json:"temple_address"`
	TemplePhoneNo     string `json:"temple_phone_no"`
	TempleDescription string `json:"temple_description"`
	LogoURL           string `json:"logo_url"`
	IntroVideoURL     string `json:"intro_video_url"`
}

type CreatorBankInfo struct {
	AccountHolderName string  `json:"account_holder_name"`
	AccountNumber     string  `json:"account_number"`
	BankName          string  `json:"bank_name"`
	BranchName        string  `json:"branch_name"`
	IFSCCode          string  `json:"ifsc_code"`
	AccountType       string  `json:"account_type"`
	UPIID             *string `json:"upi_id,omitempty"`
	// Note: Account number is intentionally excluded for security
}
type CreatorDetails struct {
	Name              string  `json:"name"`                  // maps to: u.full_name AS name
	AccountHolderName string  `json:"account_holder_name"`   // b.account_holder_name
	AccountNumber     string  `json:"account_number"`        // b.account_number
	IFSCCode          string  `json:"ifsc_code"`             // b.ifsc_code
	AccountType       string  `json:"account_type"`          // b.account_type
	UPIID             *string `json:"upi_id,omitempty"`      // b.upi_id
}