package entity

import "time"

type Entity struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Name              string     `gorm:"not null" json:"name"`                      // Temple Name
	Email             string     `gorm:"unique;not null" json:"email"`
	Phone             string     `json:"phone"`                                     // Mobile
	LogoURL           string     `json:"logo_url"`
	Description       string     `json:"description"`
	Status            string     `gorm:"default:'pending'" json:"status"`          // pending / approved / rejected
	CreatedBy         uint       `gorm:"not null" json:"created_by"`               // templeadmin who created this
	ProofDocumentsURL string     `json:"proof_documents_url"`                      // optional (JSON array if needed)

	// Address fields (flattened from EntityAddress table)
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	District     string `json:"district"`
	State        string `json:"state"`
	Pincode      string `json:"pincode"`
	Country      string `gorm:"default:'India'" json:"country"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	AddressType  string `json:"address_type"`

	// Document fields (optional, flattened if you don’t want a separate join)
	DocumentType   string `json:"document_type"`   // Aadhaar / PAN
	DocumentTitle  string `json:"document_title"`  // eg. "Main Trust Deed"
	DocumentURL    string `json:"document_url"`    // file path
	DocumentNumber string `json:"document_number"` // PAN/Aadhaar number

	// Meta fields (from old custom fields)
	TempleType      string     `json:"temple_type"`
	EstablishedDate *time.Time `json:"established_date"`
	ContactPerson   string     `json:"contact_person"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
