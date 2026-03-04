package donation

import "time"

// ==============================
// DTOs and Request/Response Models
// ==============================

// CreateDonationRequest is sent by frontend to initiate a donation
type CreateDonationRequest struct {
	UserID       uint    `json:"-"`                                                                                                  // Filled from JWT claims
	EntityID     uint    `json:"-"`                                                                                                  // Set from user context
	Amount       float64 `json:"amount" binding:"required,gt=0"`                                                                     // Donation amount in INR
	DonationType string  `json:"donationType" binding:"required,oneof=general seva event festival construction annadanam education maintenance"` // Donation type
	ReferenceID  *uint   `json:"referenceID,omitempty"`                                                                              // Optional: SevaID or EventID
	Note         *string `json:"note,omitempty"`                                                                                     // Optional donor message
	IPAddress    string  `json:"-"`                                                                                                  // For audit logging (filled from middleware)
}

// CreateDonationResponse is returned to frontend after creating Razorpay order
// CreateDonationResponse is returned to frontend after creating Razorpay order
type CreateDonationResponse struct {
	OrderID     string             `json:"order_id"`
	Amount      float64            `json:"amount"`
	Currency    string             `json:"currency"`
	RazorpayKey string             `json:"razorpay_key"`
	Tenant      TenantPaymentInfo  `json:"tenant"`
}

// TenantPaymentInfo holds the tenant's registered bank details
// shown to the devotee on the payment page
type TenantPaymentInfo struct {
	AccountHolderName string `json:"account_holder_name"`
	AccountNumber     string `json:"account_number"`
	BankName          string `json:"bank_name"`
	BranchName        string `json:"branch_name"`
	IFSCCode          string `json:"ifsc_code"`
	AccountType       string `json:"account_type"`
	UPIID             string `json:"upi_id"`
}

// VerifyPaymentRequest is used by frontend to confirm payment success.
// Frontend sends: { paymentID, orderID, razorpaySig } — these match the json tags below ✅
type VerifyPaymentRequest struct {
	OrderID     string `json:"orderID" binding:"required"`     // Razorpay order ID from frontend
	PaymentID   string `json:"paymentID" binding:"required"`   // Razorpay payment ID from frontend
	RazorpaySig string `json:"razorpaySig" binding:"required"` // Signature to verify payment
	IPAddress   string `json:"-"`                              // For audit logging (filled from middleware)
}

// DonationWithUser includes user and entity information returned to the frontend.
//
// JSON tag mapping (critical — frontend relies on these exact keys):
//
//	OrderID    -> "transactionId"  (frontend getOrderId() checks donation.transactionId ✅)
//	PaymentID  -> "paymentId"      (frontend getPaymentId() checks donation.paymentId ✅)
//	Method     -> "paymentMethod"  AND "method" (both exported for frontend compatibility)
//	UPIID      -> "upi_id"        (changed to string to avoid null JSON, empty string is safe ✅)
type DonationWithUser struct {
	ID           uint      `json:"id" db:"id"`
	UserID       uint      `json:"user_id" db:"user_id"`
	EntityID     uint      `json:"entity_id" db:"entity_id"`
	Amount       float64   `json:"amount" db:"amount"`
	DonationType string    `json:"donationType" db:"donation_type"`
	ReferenceID  *uint     `json:"referenceID,omitempty" db:"reference_id"`
	Method       string    `json:"paymentMethod" db:"method"` // primary JSON key: paymentMethod
	Status       string    `json:"status" db:"status"`
	OrderID      string    `json:"transactionId" db:"order_id"` // frontend getOrderId() looks for transactionId ✅
	PaymentID    *string   `json:"paymentId,omitempty" db:"payment_id"` // frontend getPaymentId() looks for paymentId ✅
	Note         *string   `json:"note,omitempty" db:"note"`
	DonatedAt    *time.Time `json:"donatedAt,omitempty" db:"donated_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	// User information
	UserName  string `json:"userName" db:"user_name"`
	UserEmail string `json:"userEmail" db:"user_email"`

	// Entity information
	EntityName string `json:"entityName" db:"entity_name"`

	// Computed aliases for frontend compatibility
	Date          time.Time `json:"date" db:"created_at"`
	Type          string    `json:"type" db:"donation_type"`
	DonorName     string    `json:"donorName" db:"user_name"`
	DonorEmail    string    `json:"donorEmail" db:"user_email"`
	PaymentMethod string    `json:"method" db:"method"` // secondary alias: method

	// FIX: Payment/bank details stored about the PAYEE (temple/merchant), NOT the payer.
	// These are populated from the temple's registered bank details via entityRepo.
	// account_holder_name = temple's registered account holder name (NOT the donor's name/UPI)
	AccountHolderName string `json:"account_holder_name" db:"account_holder_name"`
	AccountNumber     string `json:"account_number" db:"account_number"`
	AccountType       string `json:"account_type" db:"account_type"`
	IFSCCode          string `json:"ifsc_code" db:"ifsc_code"`
	// FIX: Changed from *string to string so JSON never emits null (frontend getField checks for null/empty) ✅
	UPIID string `json:"upi_id" db:"upi_id"`
}

// DonationFilters for filtering and pagination
type DonationFilters struct {
	EntityID  uint       `json:"entity_id"`
	UserID    uint       `json:"user_id,omitempty"`
	Status    string     `json:"status,omitempty"`
	Type      string     `json:"type,omitempty"`
	Method    string     `json:"method,omitempty"`
	From      *time.Time `json:"from,omitempty"`
	To        *time.Time `json:"to,omitempty"`
	MinAmount *float64   `json:"min_amount,omitempty"`
	MaxAmount *float64   `json:"max_amount,omitempty"`
	Search    string     `json:"search,omitempty"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}

// UpdatePaymentDetailsParams for updating payment information after Razorpay callback
type UpdatePaymentDetailsParams struct {
	Status    string
	PaymentID *string
	Method    string
	Amount    float64
	DonatedAt *time.Time
	// Payee details — must be the TEMPLE's info, not the payer/donor's info
	AccountHolderName string // Temple's registered account holder name
	AccountNumber     string // Temple's account number or UPI ID
	AccountType       string // UPI / CARD / BANK_TRANSFER / WALLET etc.
	IFSCCode          string // Temple's IFSC code
	UPIID             string // Temple's UPI ID
}

// ==============================
// Analytics and Reporting Models
// ==============================

// DashboardStats represents overall donation statistics
type DashboardStats struct {
	TotalAmount    float64 `json:"totalAmount"`
	TotalCount     int     `json:"total_count"`
	CompletedCount int     `json:"completed"`
	PendingCount   int     `json:"pending"`
	FailedCount    int     `json:"failed"`
	ThisMonth      float64 `json:"thisMonth"`
	Today          float64 `json:"today"`
	TotalDonors    int     `json:"totalDonors"`
	AverageAmount  float64 `json:"averageAmount"`
}

// StatsResult for database aggregation queries
type StatsResult struct {
	Amount         float64 `json:"amount"`
	Count          int     `json:"count"`
	CompletedCount int     `json:"completed_count"`
	PendingCount   int     `json:"pending_count"`
	FailedCount    int     `json:"failed_count"`
}

// TopDonor represents a top donor
type TopDonor struct {
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	TotalAmount   float64 `json:"total_amount"`
	DonationCount int     `json:"donation_count"`
}

// TrendData for donation trends over time
type TrendData struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
	Count  int       `json:"count"`
}

// TypeData for donations by type
type TypeData struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}

// MethodData for donations by payment method
type MethodData struct {
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}

// AnalyticsData combines all analytics information
type AnalyticsData struct {
	Trends   []TrendData  `json:"trends"`
	ByType   []TypeData   `json:"byType"`
	ByMethod []MethodData `json:"byMethod"`
}

// Receipt represents a donation receipt
type Receipt struct {
	ID             uint      `json:"id"`
	DonationAmount float64   `json:"donationAmount"`
	DonationType   string    `json:"donationType"`
	DonorName      string    `json:"donorName"`
	DonorEmail     string    `json:"donorEmail"`
	TransactionID  string    `json:"transactionId"`
	DonatedAt      time.Time `json:"donatedAt"`
	Method         string    `json:"method"`
	EntityName     string    `json:"entityName"`
	ReceiptNumber  string    `json:"receiptNumber"`
	GeneratedAt    time.Time `json:"generatedAt"`
}

// DonationListResponse represents paginated donation list response
type DonationListResponse struct {
	Data       []DonationWithUser `json:"data"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
	Success    bool               `json:"success"`
}

// RecentDonation represents recent donation info
type RecentDonation struct {
	Amount       float64   `json:"amount" db:"amount"`
	DonationType string    `json:"donation_type" db:"donation_type"`
	Method       string    `json:"method" db:"method"`
	Status       string    `json:"status" db:"status"`
	DonatedAt    time.Time `json:"donated_at" db:"donated_at"`
	UserName     string    `json:"user_name" db:"user_name"`
	EntityName   string    `json:"entity_name" db:"entity_name"`
}