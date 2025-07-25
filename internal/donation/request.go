package donation

// ==============================
// 📝 Donation Request DTOs
// ==============================

// CreateDonationRequest is sent by frontend to initiate a donation
type CreateDonationRequest struct {
	UserID       uint    `json:"-"`                                // Filled from JWT claims (backend-side)
	EntityID     uint    `json:"-"`                                // Set from user context
	Amount       float64 `json:"amount" binding:"required,gt=0"`   // ₹ > 0
	DonationType string  `json:"donation_type" binding:"required"` // general | seva | event | festival
	ReferenceID  string  `json:"reference_id,omitempty"`           // Optional: SevaID / EventID
	Note         *string `json:"note,omitempty"`                   // Optional message or intention
	Method       string  `json:"-"`                                // Let Razorpay UI decide
}

// CreateDonationResponse is returned to frontend after creating Razorpay order
type CreateDonationResponse struct {
	OrderID     string  `json:"order_id"`     // Razorpay order ID (for checkout.js)
	Amount      float64 `json:"amount"`       // ₹ amount
	Currency    string  `json:"currency"`     // INR always
	RazorpayKey string  `json:"razorpay_key"` // For Razorpay Checkout on client
}

// VerifyPaymentRequest is used by frontend to confirm the success (after modal closes)
type VerifyPaymentRequest struct {
	OrderID     string `json:"order_id" binding:"required"`           // Razorpay order ID
	PaymentID   string `json:"payment_id" binding:"required"`         // Razorpay payment ID
	RazorpaySig string `json:"razorpay_signature" binding:"required"` // Signature to verify authenticity
}

// ==============================
// 📊 Donation Dashboard Response
// ==============================

type RecentDonor struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"` // formatted date string
}


// package donation

// // CreateDonationRequest represents the payload to start a donation
// type CreateDonationRequest struct {
// 	UserID       uint    `json:"-"`                     // Set from middleware
//     EntityID uint `json:"-"`
// 	Amount       float64 `json:"amount" binding:"required,gt=0"`
// 	DonationType string  `json:"donation_type"`         // general | seva | event | festival
// 	ReferenceID  *uint   `json:"reference_id"`          // SevaID or EventID
// 	Method       string  `json:"-"`                     // Default: UPI
// 	Note         *string `json:"note"`
// }

// // CreateDonationResponse holds Razorpay Order and donation info
// type CreateDonationResponse struct {
// 	OrderID     string  `json:"order_id"`
// 	Amount      float64 `json:"amount"`
// 	Currency    string  `json:"currency"`
// 	RazorpayKey string  `json:"razorpay_key"`
// }

// // VerifyPaymentRequest is called after successful Razorpay payment
// type VerifyPaymentRequest struct {
// 	OrderID        string `json:"order_id" binding:"required"`
// 	PaymentID      string `json:"payment_id" binding:"required"`
// 	RazorpaySig    string `json:"razorpay_signature" binding:"required"`
// }