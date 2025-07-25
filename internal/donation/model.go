package donation

import (
	"time"

	"gorm.io/gorm"
)

// ========================
// 📊 Dashboard Response Struct
// ========================
// This struct is used only for returning aggregated donation data to the frontend.
// GORM will ignore the `RecentDonors` field since it's not a DB column.
type DonationDashboardResponse struct {
	TotalDonations  float64 `json:"total_donations"`
	TotalCount      int     `json:"total_count"`
	TotalDonors     int     `json:"total_donors"`
	AverageDonation float64 `json:"average_donation"`
	ThisMonth       float64 `json:"this_month"`
	RecentDonors    []Donor `json:"recent_donors" gorm:"-"` // Ignored by GORM
}


// ========================
// 💡 Constants
// ========================

const (
	// Donation statuses
	StatusPending = "PENDING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"

	// Payment methods
	MethodUPI        = "UPI"
	MethodCard       = "CARD"
	MethodNetbanking = "NETBANKING"
	MethodWallet     = "WALLET"
	MethodCash       = "CASH"
	MethodCheque     = "CHEQUE"

	// Donation types
	TypeGeneral  = "general"
	TypeSeva     = "seva"
	TypeEvent    = "event"
	TypeFestival = "festival"
)

// ========================
// 🧾 Donor Summary Struct (used in response only)
// ========================

type Donor struct {
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
	Method string    `json:"method"`
	Status string    `json:"status"`
}

// ========================
// 💳 Donation Table Model
// ========================

type Donation struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID   uint `gorm:"not null;index" json:"user_id"`   // Donor's user ID
	EntityID uint `gorm:"not null;index" json:"entity_id"` // Temple or entity ID

	Amount       float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	DonationType string  `gorm:"size:50;index" json:"donation_type"`
	ReferenceID  string  `gorm:"index" json:"reference_id,omitempty"`

	Method string `gorm:"size:50;not null;index" json:"method"`
	Status string `gorm:"size:20;default:'PENDING';index" json:"status"`

	OrderID   string  `gorm:"size:100;uniqueIndex" json:"order_id"`
	PaymentID *string `gorm:"size:100;index" json:"payment_id,omitempty"`

	ReceiptURL *string `gorm:"type:text" json:"receipt_url,omitempty"`
	Note       *string `gorm:"type:text" json:"note,omitempty"`

	DonatedAt *time.Time     `json:"donated_at,omitempty"` // Set when status is success
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}






// package donation

// import (
// 	"time"

// 	"gorm.io/gorm"
// )

// // Donation model represents a donation transaction
// type Donation struct {
// 	ID            uint           `gorm:"primaryKey" json:"id"`
// 	UserID        uint           `gorm:"not null;index" json:"user_id"`        // Devotee who donated
// 	EntityID      uint           `gorm:"not null;index" json:"entity_id"`      // Temple ID
// 	Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
// 	DonationType  string         `gorm:"size:50" json:"donation_type"`         // general, seva, event, festival
// 	ReferenceID   *uint          `gorm:"index" json:"reference_id"`            // Optional Seva/Event ID
// 	Method        string         `gorm:"size:50;not null" json:"method"`       // UPI, CARD, NETBANKING, etc.
// 	Status        string         `gorm:"size:20;default:'PENDING'" json:"status"` // PENDING, SUCCESS, FAILED
// 	OrderID       string         `gorm:"size:100;index" json:"order_id"`       // Razorpay Order ID
// 	PaymentID     *string        `gorm:"size:100;index" json:"payment_id"`     // Razorpay Payment ID (nullable until success)
// 	ReceiptURL    *string        `gorm:"type:text" json:"receipt_url"`         // Optional link to PDF/email receipt
// 	Note          *string        `gorm:"type:text" json:"note"`                // Optional donor message / intention
//     DonatedAt time.Time `json:"donated_at"` // Set this manually when payment is verified
// 	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
// 	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
// 	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
// }