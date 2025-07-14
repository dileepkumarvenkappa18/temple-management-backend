package donation

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	razorpay "github.com/razorpay/razorpay-go"
	"github.com/sharath018/temple-management-backend/config"
)

type Service interface {
	CreateDonation(userID uint, entityID uint, amount float64, donationType, referenceID, note string) (*Donation, string, error)
	VerifyDonation(paymentID string, orderID string) (*Donation, error)
	StartDonation(req CreateDonationRequest) (*CreateDonationResponse, error)
	VerifyAndUpdateDonation(req VerifyPaymentRequest) error
	GetDonationsByUser(userID uint) ([]Donation, error)
	GetDonationsByEntity(entityID uint, page int, limit int, status string) ([]Donation, error)
}

type service struct {
	repo   Repository
	client *razorpay.Client
	cfg    *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	client := razorpay.NewClient(cfg.RazorpayKey, cfg.RazorpaySecret)
	return &service{
		repo:   repo,
		client: client,
		cfg:    cfg,
	}
}

// StartDonation initializes the Razorpay order and creates a pending donation entry
func (s *service) StartDonation(req CreateDonationRequest) (*CreateDonationResponse, error) {
	donation, orderID, err := s.CreateDonation(
		req.UserID,
		req.EntityID,
		req.Amount,
		req.DonationType,
		func() string {
			if req.ReferenceID != nil {
				return fmt.Sprintf("%d", *req.ReferenceID)
			}
			return ""
		}(),
		func() string {
			if req.Note != nil {
				return *req.Note
			}
			return ""
		}(),
	)
	if err != nil {
		return nil, err
	}

	return &CreateDonationResponse{
		OrderID:     orderID,
		Amount:      donation.Amount,
		Currency:    "INR",
		RazorpayKey: s.cfg.RazorpayKey,
	}, nil
}

// Securely verifies Razorpay signature and updates payment status
func (s *service) VerifyAndUpdateDonation(req VerifyPaymentRequest) error {
	// Step 1: Verify HMAC Signature
	expected := hmac.New(sha256.New, []byte(s.cfg.RazorpaySecret))
	expected.Write([]byte(req.OrderID + "|" + req.PaymentID))
	computedSignature := hex.EncodeToString(expected.Sum(nil))

	if computedSignature != req.RazorpaySig {
		return fmt.Errorf("invalid payment signature")
	}

	// Step 2: Proceed with Razorpay fetch and DB update
	_, err := s.VerifyDonation(req.PaymentID, req.OrderID)
	return err
}

func (s *service) GetDonationsByUser(userID uint) ([]Donation, error) {
	return s.repo.ListByUserID(context.Background(), userID)
}

func (s *service) GetDonationsByEntity(entityID uint, page int, limit int, status string) ([]Donation, error) {
	return s.repo.ListByEntityID(context.Background(), entityID, page, limit, status)
}

// CreateDonation calls Razorpay order API and creates a donation entry with status = PENDING
func (s *service) CreateDonation(userID uint, entityID uint, amount float64, donationType, referenceID, note string) (*Donation, string, error) {
	amountInPaise := int(amount * 100)

	data := map[string]interface{}{
		"amount":          amountInPaise,
		"currency":        "INR",
		"payment_capture": 1,
		"notes": map[string]interface{}{
			"user_id":       userID,
			"entity_id":     entityID,
			"donation_type": donationType,
			"reference_id":  referenceID,
		},
	}

	order, err := s.client.Order.Create(data, nil)
	if err != nil {
		return nil, "", fmt.Errorf("razorpay order creation failed: %w", err)
	}

	orderID, ok := order["id"].(string)
	if !ok {
		return nil, "", errors.New("unable to extract order_id from Razorpay response")
	}

	donation := &Donation{
		UserID:   userID,
		EntityID: entityID,
		Amount:   amount,
		Method:   "CARD", // Test value for logging
		Status:   "PENDING",
		OrderID:  orderID,
		Note:     &note,
	}

	if err := s.repo.Create(context.Background(), donation); err != nil {
		return nil, "", err
	}

	return donation, orderID, nil
}

// VerifyDonation fetches the Razorpay payment and updates donation based on status
func (s *service) VerifyDonation(paymentID string, orderID string) (*Donation, error) {
	payment, err := s.client.Payment.Fetch(paymentID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("razorpay payment fetch failed: %w", err)
	}

	status, ok := payment["status"].(string)
	if !ok {
		return nil, errors.New("invalid payment status format")
	}

	var amount float64
	switch val := payment["amount"].(type) {
	case float64:
		amount = val / 100
	case json.Number:
		amountPaise, _ := val.Float64()
		amount = amountPaise / 100
	default:
		return nil, fmt.Errorf("unsupported amount type: %T", val)
	}

	donation, err := s.repo.GetByOrderID(context.Background(), orderID)
	if err != nil {
		return nil, errors.New("donation record not found for given order ID")
	}

	if donation.Status == "SUCCESS" {
		return donation, nil
	}

	if donation.OrderID != orderID {
		return nil, errors.New("order ID mismatch")
	}

	newStatus := "FAILED"
	if status == "captured" {
		newStatus = "SUCCESS"
		now := time.Now()
		donation.DonatedAt = &now
	}

	paymentIDStr := payment["id"].(string)
	donation.Amount = amount

	if err := s.repo.UpdatePaymentStatus(context.Background(), orderID, newStatus, &paymentIDStr); err != nil {
		return nil, err
	}

	return donation, nil
}








// package donation

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"time"

// 	razorpay "github.com/razorpay/razorpay-go"
// 	"github.com/sharath018/temple-management-backend/config"
// )

// type Service interface {
// 	CreateDonation(userID uint, entityID uint, amount float64, donationType, referenceID, note string) (*Donation, string, error)
// 	VerifyDonation(paymentID string, orderID string) (*Donation, error)
// 	StartDonation(req CreateDonationRequest) (*CreateDonationResponse, error)
// 	VerifyAndUpdateDonation(req VerifyPaymentRequest) error
// 	GetDonationsByUser(userID uint) ([]Donation, error)
// 	GetDonationsByEntity(entityID uint) ([]Donation, error)
// }

// type service struct {
// 	repo   Repository
// 	client *razorpay.Client
// 	cfg    *config.Config
// }

// func NewService(repo Repository, cfg *config.Config) Service {
// 	client := razorpay.NewClient(cfg.RazorpayKey, cfg.RazorpaySecret)
// 	return &service{
// 		repo:   repo,
// 		client: client,
// 		cfg:    cfg,
// 	}
// }

// // Start donation with Razorpay and log pending entry
// func (s *service) StartDonation(req CreateDonationRequest) (*CreateDonationResponse, error) {
// 	donation, orderID, err := s.CreateDonation(
// 		req.UserID,
// 		req.EntityID,
// 		req.Amount,
// 		req.DonationType,
// 		func() string {
// 			if req.ReferenceID != nil {
// 				return fmt.Sprintf("%d", *req.ReferenceID)
// 			}
// 			return ""
// 		}(),
// 		func() string {
// 			if req.Note != nil {
// 				return *req.Note
// 			}
// 			return ""
// 		}(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &CreateDonationResponse{
// 		OrderID:     orderID,
// 		Amount:      donation.Amount,
// 		Currency:    "INR",
// 		RazorpayKey: s.cfg.RazorpayKey,
// 	}, nil
// }

// func (s *service) VerifyAndUpdateDonation(req VerifyPaymentRequest) error {
// 	_, err := s.VerifyDonation(req.PaymentID, req.OrderID)
// 	return err
// }

// func (s *service) GetDonationsByUser(userID uint) ([]Donation, error) {
// 	return s.repo.ListByUserID(context.Background(), userID)
// }

// func (s *service) GetDonationsByEntity(entityID uint) ([]Donation, error) {
// 	return s.repo.ListByEntityID(context.Background(), entityID)
// }

// // CreateDonation creates a Razorpay order and logs a pending donation
// func (s *service) CreateDonation(userID uint, entityID uint, amount float64, donationType, referenceID, note string) (*Donation, string, error) {
// 	amountInPaise := int(amount * 100)

// 	// 🔧 No need to specify "method" here — Razorpay lets user choose at checkout (test mode too)
// 	data := map[string]interface{}{
// 		"amount":          amountInPaise,
// 		"currency":        "INR",
// 		"payment_capture": 1,
// 		"notes": map[string]interface{}{
// 			"user_id":       userID,
// 			"entity_id":     entityID,
// 			"donation_type": donationType,
// 			"reference_id":  referenceID,
// 		},
// 	}

// 	order, err := s.client.Order.Create(data, nil)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("razorpay order creation failed: %w", err)
// 	}
// 	orderID := order["id"].(string)

// 	// Set method to "CARD" or "NETBANKING" for test logs
// 	method := "CARD" // default for test logs

// 	donation := &Donation{
// 		UserID:   userID,
// 		EntityID: entityID,
// 		Amount:   amount,
// 		Method:   method,
// 		Status:   "PENDING",
// 		OrderID:  orderID,
// 		Note:     &note,
// 	}

// 	if err := s.repo.Create(context.Background(), donation); err != nil {
// 		return nil, "", err
// 	}

// 	return donation, orderID, nil
// }


// // VerifyDonation confirms the payment and updates donation record
// func (s *service) VerifyDonation(paymentID string, orderID string) (*Donation, error) {
// 	payment, err := s.client.Payment.Fetch(paymentID, nil, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("payment fetch failed: %w", err)
// 	}

// 	status := payment["status"].(string)
// 	var amount float64
// switch val := payment["amount"].(type) {
// case float64:
// 	amount = val / 100
// case json.Number:
// 	amountPaise, _ := val.Float64()
// 	amount = amountPaise / 100
// default:
// 	return nil, fmt.Errorf("unsupported amount type: %T", payment["amount"])
// }


// 	donation, err := s.repo.GetByOrderID(context.Background(), orderID)
// 	if err != nil {
// 		return nil, errors.New("donation record not found")
// 	}

// 	newStatus := "FAILED"
// 	if status == "captured" {
// 		newStatus = "SUCCESS"
// 		donation.DonatedAt = time.Now()
// 	}

// 	paymentIDStr := payment["id"].(string)
// 	donation.Amount = amount

// 	err = s.repo.UpdatePaymentStatus(context.Background(), orderID, newStatus, &paymentIDStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return donation, nil
// }