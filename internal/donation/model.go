package donation

import (
	"time"

	"gorm.io/gorm"
)

// DonationStatus represents valid payment states
const (
	StatusPending = "PENDING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)

const (
	MethodUPI        = "UPI"
	MethodCard       = "CARD"
	MethodNetbanking = "NETBANKING"
	MethodWallet     = "WALLET"
)

const (
	TypeGeneral      = "general"
	TypeSeva         = "seva"
	TypeEvent        = "event"
	TypeFestival     = "festival"
	TypeConstruction = "construction"
	TypeAnnadanam    = "annadanam"
	TypeEducation    = "education"
	TypeMaintenance  = "maintenance"
)

type Donation struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID   uint `gorm:"not null;index" json:"user_id"`   // Devotee who donated
	EntityID uint `gorm:"not null;index" json:"entity_id"` // Temple ID

	Amount       float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	DonationType string  `gorm:"size:50;index" json:"donation_type"`
	ReferenceID  *uint   `gorm:"index" json:"reference_id,omitempty"`

	Method string `gorm:"size:50;not null;index" json:"method"`
	Status string `gorm:"size:20;default:'PENDING';index" json:"status"`

	// FIX: json tags match what DonationWithUser returns so both are consistent
	OrderID   string  `gorm:"size:100;uniqueIndex" json:"transactionId"`
	PaymentID *string `gorm:"size:100;index" json:"paymentId,omitempty"`

	Note *string `gorm:"type:text" json:"note,omitempty"`

	AccountHolderName string `gorm:"size:255" json:"account_holder_name"`
	AccountNumber     string `gorm:"size:30" json:"account_number"`
	AccountType       string `gorm:"size:20" json:"account_type"`
	IFSCCode          string `gorm:"size:11" json:"ifsc_code"`
	UPIID             string `gorm:"size:100" json:"upi_id"` // FIX: was *string, now string

	DonatedAt *time.Time     `json:"donated_at,omitempty"` // Set only on successful payment
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Donation model
func (Donation) TableName() string {
	return "donations"
}