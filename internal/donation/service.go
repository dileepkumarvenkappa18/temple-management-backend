package donation

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	razorpay "github.com/razorpay/razorpay-go"
	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/middleware"
)

type Service interface {
	// Core donation operations (DEVOTEE - UNCHANGED)
	StartDonation(req CreateDonationRequest) (*CreateDonationResponse, error)
	VerifyAndUpdateDonation(req VerifyPaymentRequest) error
	HandleRazorpayWebhook(orderID, paymentID, method string, amount float64) error // ✅ ADD THIS LINE
	HandleFailedPaymentWebhook(orderID, paymentID string) error // ✅ ADD THIS
	// Data retrieval - UPDATED for entity-based approach
	GetDonationsByUser(userID uint) ([]DonationWithUser, error)
	GetDonationsByUserAndEntity(userID uint, entityID uint) ([]DonationWithUser, error) // NEW
	GetDonationsWithFilters(filters DonationFilters, accessContext middleware.AccessContext) ([]DonationWithUser, int, error)

	// Analytics and reporting - TEMPLE ADMIN (UPDATED)
	GetDashboardStats(entityID uint, accessContext middleware.AccessContext) (*DashboardStats, error)
	GetTopDonors(entityID uint, limit int, accessContext middleware.AccessContext) ([]TopDonor, error)
	GetAnalytics(entityID uint, days int, accessContext middleware.AccessContext) (*AnalyticsData, error)

	// Receipt and export - BOTH (UPDATED)
	GenerateReceipt(donationID uint, userID uint, accessContext *middleware.AccessContext, entityID uint) (*Receipt, error) // NEW: Added entityID
	ExportDonations(filters DonationFilters, format string, accessContext middleware.AccessContext) ([]byte, string, error)

	// Recent donations - BOTH (UPDATED)
	GetRecentDonationsByUser(ctx context.Context, userID uint, limit int) ([]RecentDonation, error)
	GetRecentDonationsByUserAndEntity(ctx context.Context, userID uint, entityID uint, limit int) ([]RecentDonation, error) // NEW
	GetRecentDonationsByEntity(ctx context.Context, entityID uint, limit int, accessContext middleware.AccessContext) ([]RecentDonation, error)
}

type service struct {
	repo     Repository
	client   *razorpay.Client
	cfg      *config.Config
	auditSvc auditlog.Service
}

func NewService(repo Repository, cfg *config.Config, auditSvc auditlog.Service) Service {
	client := razorpay.NewClient(cfg.RazorpayKey, cfg.RazorpaySecret)
	return &service{
		repo:     repo,
		client:   client,
		cfg:      cfg,
		auditSvc: auditSvc,
	}
}

// ==============================
// Core Donation Operations (DEVOTEE - UNCHANGED)
// ==============================

// StartDonation initializes the Razorpay order and creates a pending donation entry
func (s *service) StartDonation(req CreateDonationRequest) (*CreateDonationResponse, error) {
	ctx := context.Background()

	// Create Razorpay order
	amountInPaise := int(req.Amount * 100)

	data := map[string]interface{}{
		"amount":          amountInPaise,
		"currency":        "INR",
		"payment_capture": 1,
		"notes": map[string]interface{}{
			"user_id":       req.UserID,
			"entity_id":     req.EntityID,
			"donation_type": req.DonationType,
		},
	}

	if req.ReferenceID != nil {
		data["notes"].(map[string]interface{})["reference_id"] = *req.ReferenceID
	}

	order, err := s.client.Order.Create(data, nil)
	if err != nil {
		s.auditSvc.LogAction(ctx, &req.UserID, &req.EntityID, "DONATION_INITIATED", map[string]interface{}{
			"amount":        req.Amount,
			"donation_type": req.DonationType,
			"error":         err.Error(),
		}, req.IPAddress, "failure")

		return nil, fmt.Errorf("razorpay order creation failed: %w", err)
	}

	orderID, ok := order["id"].(string)
	if !ok {
		s.auditSvc.LogAction(ctx, &req.UserID, &req.EntityID, "DONATION_INITIATED", map[string]interface{}{
			"amount":        req.Amount,
			"donation_type": req.DonationType,
			"error":         "unable to extract order_id from Razorpay response",
		}, req.IPAddress, "failure")

		return nil, errors.New("unable to extract order_id from Razorpay response")
	}

	// Create pending donation record
	donation := &Donation{
		UserID:       req.UserID,
		EntityID:     req.EntityID,
		Amount:       req.Amount,
		DonationType: req.DonationType,
		ReferenceID:  req.ReferenceID,
		Method:       " ", // Will be updated after payment
		Status:       StatusPending,
		OrderID:      orderID,
		Note:         req.Note,
	}

	if err := s.repo.Create(context.Background(), donation); err != nil {
		s.auditSvc.LogAction(ctx, &req.UserID, &req.EntityID, "DONATION_INITIATED", map[string]interface{}{
			"amount":        req.Amount,
			"donation_type": req.DonationType,
			"order_id":      orderID,
			"error":         err.Error(),
		}, req.IPAddress, "failure")

		return nil, fmt.Errorf("failed to create donation record: %w", err)
	}

	s.auditSvc.LogAction(ctx, &req.UserID, &req.EntityID, "DONATION_INITIATED", map[string]interface{}{
		"amount":        req.Amount,
		"donation_type": req.DonationType,
		"order_id":      orderID,
		"reference_id":  req.ReferenceID,
	}, req.IPAddress, "success")

	return &CreateDonationResponse{
		OrderID:     orderID,
		Amount:      req.Amount,
		Currency:    "INR",
		RazorpayKey: s.cfg.RazorpayKey,
	}, nil
}


// VerifyAndUpdateDonation securely verifies Razorpay signature and updates payment status
// ✅ Works for UPI / CARD / WALLET
// ❌ Skips NETBANKING (handled via webhook)
func (s *service) VerifyAndUpdateDonation(req VerifyPaymentRequest) error {
	ctx := context.Background()

	// Step 1: Get donation record FIRST
	donation, err := s.repo.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		s.auditSvc.LogAction(ctx, nil, nil, "DONATION_VERIFICATION_FAILED", map[string]interface{}{
			"order_id": req.OrderID,
			"reason":   "donation record not found",
		}, req.IPAddress, "failure")

		return errors.New("donation record not found for given order ID")
	}

	// ✅ CRITICAL: Idempotency check - prevent duplicate processing
	if donation.Status == StatusSuccess {
		log.Printf("⚠️ Donation already successful for order %s (payment_id: %v), skipping verification",
			req.OrderID, donation.PaymentID)
		return nil // Return success, payment already processed
	}

	// 🚨 CRITICAL FIX: DO NOT verify NETBANKING here
	if donation.Method == MethodNetbanking {
		log.Printf("⚠️ Netbanking payment for order %s will be confirmed via webhook", req.OrderID)
		// Netbanking confirmation comes via Razorpay webhook
		return nil
	}

	// Step 2: Verify HMAC Signature (ONLY for instant methods)
	expected := hmac.New(sha256.New, []byte(s.cfg.RazorpaySecret))
	expected.Write([]byte(req.OrderID + "|" + req.PaymentID))
	computedSignature := hex.EncodeToString(expected.Sum(nil))

	if computedSignature != req.RazorpaySig {
		s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "DONATION_VERIFICATION_FAILED", map[string]interface{}{
			"order_id":   req.OrderID,
			"payment_id": req.PaymentID,
			"reason":     "invalid payment signature",
		}, req.IPAddress, "failure")

		return fmt.Errorf("invalid payment signature")
	}

	// Step 3: Fetch payment details from Razorpay
	payment, err := s.client.Payment.Fetch(req.PaymentID, nil, nil)
	if err != nil {
		s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "DONATION_VERIFICATION_FAILED", map[string]interface{}{
			"order_id":   req.OrderID,
			"payment_id": req.PaymentID,
			"reason":     "razorpay payment fetch failed",
			"error":      err.Error(),
		}, req.IPAddress, "failure")

		return fmt.Errorf("razorpay payment fetch failed: %w", err)
	}

	status, ok := payment["status"].(string)
	if !ok {
		return errors.New("invalid payment status format")
	}

	// Step 4: Double-check idempotency (in case of race condition)
	// Re-fetch donation to ensure status hasn't changed
	donation, err = s.repo.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		return fmt.Errorf("failed to re-fetch donation: %w", err)
	}

	if donation.Status == StatusSuccess {
		log.Printf("⚠️ Race condition detected: Donation already successful for order %s", req.OrderID)
		return nil
	}

	// Step 5: Extract amount
	var amount float64
	switch val := payment["amount"].(type) {
	case float64:
		amount = val / 100
	case json.Number:
		amountPaise, _ := val.Float64()
		amount = amountPaise / 100
	default:
		return fmt.Errorf("unsupported amount type: %T", val)
	}

	// Step 6: Determine final status
	newStatus := StatusFailed
	var donatedAt *time.Time
	auditAction := "DONATION_FAILED"
	auditStatus := "failure"

	if status == "captured" {
		newStatus = StatusSuccess
		now := time.Now()
		donatedAt = &now
		auditAction = "DONATION_SUCCESS"
		auditStatus = "success"
	}

	// Step 7: Extract payment method
	method := "UNKNOWN"
	if paymentMethod, ok := payment["method"].(string); ok {
		method = paymentMethod
	}

	// ✅ NEW STEP 8: Extract actual recipient/payment details from Razorpay
	var accountHolderName, accountNumber, accountType, ifscCode, upiID string

	switch method {
	case "upi":
		// Extract UPI details
		if vpa, ok := payment["vpa"].(map[string]interface{}); ok {
			// UPI account holder name
			if name, ok := vpa["name"].(string); ok {
				accountHolderName = name
			}
			// UPI ID (e.g., xyz@paytm)
			if handle, ok := vpa["handle"].(string); ok {
				upiID = handle
				accountNumber = handle // Store UPI ID as account number for UPI payments
			}
		}
		accountType = "UPI"
		
		log.Printf("📱 UPI Payment Details - Holder: %s, UPI ID: %s", accountHolderName, upiID)

	case "netbanking":
		// Extract netbanking details
		if acquirer, ok := payment["acquirer_data"].(map[string]interface{}); ok {
			if bankName, ok := acquirer["bank_name"].(string); ok {
				accountHolderName = bankName
			}
			if bankAccount, ok := acquirer["bank_account_number"].(string); ok {
				accountNumber = bankAccount
			}
		}
		// Try to get bank from payment object
		if bank, ok := payment["bank"].(string); ok && accountHolderName == "" {
			accountHolderName = bank
		}
		accountType = "BANK_TRANSFER"
		
		log.Printf("🏦 Netbanking Payment Details - Bank: %s, Account: %s", accountHolderName, accountNumber)

	case "card":
		// Extract card details (limited info for security)
		if card, ok := payment["card"].(map[string]interface{}); ok {
			// Card holder name
			if name, ok := card["name"].(string); ok {
				accountHolderName = name
			}
			// Last 4 digits of card
			if last4, ok := card["last4"].(string); ok {
				accountNumber = "XXXX-XXXX-XXXX-" + last4
			}
			// Card network
			if network, ok := card["network"].(string); ok {
				accountType = "CARD_" + network
			}
		}
		
		log.Printf("💳 Card Payment Details - Holder: %s, Card: %s, Type: %s", 
			accountHolderName, accountNumber, accountType)

	case "wallet":
		// Extract wallet details
		if wallet, ok := payment["wallet"].(string); ok {
			accountHolderName = wallet // Wallet provider name (e.g., "paytm", "phonepe")
			accountType = "WALLET"
		}
		// Try to get wallet from acquirer_data
		if acquirer, ok := payment["acquirer_data"].(map[string]interface{}); ok {
			if walletName, ok := acquirer["wallet_name"].(string); ok {
				accountHolderName = walletName
			}
		}
		
		log.Printf("👛 Wallet Payment Details - Provider: %s", accountHolderName)

	default:
		log.Printf("⚠️ Unknown payment method: %s - limited details captured", method)
		// Try to extract any available info
		if acquirer, ok := payment["acquirer_data"].(map[string]interface{}); ok {
			if name, ok := acquirer["name"].(string); ok {
				accountHolderName = name
			}
		}
	}

	// Fallback: If we still don't have account holder name, use email from payment
	if accountHolderName == "" {
		if email, ok := payment["email"].(string); ok {
			accountHolderName = email
		}
	}

	// Step 9: Update donation record with payment details
	updateParams := UpdatePaymentDetailsParams{
		Status:    newStatus,
		PaymentID: &req.PaymentID,
		Method:    method,
		Amount:    amount,
		DonatedAt: donatedAt,
		
		// ✅ CRITICAL: Store actual recipient details from Razorpay
		AccountHolderName: accountHolderName,
		AccountNumber:     accountNumber,
		AccountType:       accountType,
		IFSCCode:          ifscCode,
		UPIID:             upiID,
	}

	err = s.repo.UpdatePaymentDetails(ctx, req.OrderID, updateParams)

	if err != nil {
		s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "DONATION_UPDATE_FAILED", map[string]interface{}{
			"order_id":   req.OrderID,
			"payment_id": req.PaymentID,
			"error":      err.Error(),
		}, req.IPAddress, "failure")

		return err
	}

	// Step 10: Audit success with recipient details
	s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, auditAction, map[string]interface{}{
		"order_id":           req.OrderID,
		"payment_id":         req.PaymentID,
		"amount":             amount,
		"method":             method,
		"razorpay_status":    status,
		"account_holder":     accountHolderName,
		"payment_account":    accountNumber,
		"upi_id":             upiID,
	}, req.IPAddress, auditStatus)

	log.Printf("✅ Payment verification completed for order %s: status=%s, method=%s, amount=%.2f, recipient=%s",
		req.OrderID, newStatus, method, amount, accountHolderName)

	return nil
}

// ==============================
// Data Retrieval - UPDATED for entity-based approach
// ==============================

func (s *service) GetDonationsByUser(userID uint) ([]DonationWithUser, error) {
	donations, err := s.repo.ListByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	for i := range donations {
		donations[i].Date = donations[i].CreatedAt
		donations[i].Type = donations[i].DonationType
		donations[i].DonorName = donations[i].UserName
		donations[i].DonorEmail = donations[i].UserEmail
		donations[i].PaymentMethod = donations[i].Method

		if donations[i].DonatedAt == nil {
			donations[i].DonatedAt = &donations[i].CreatedAt
		}
	}

	return donations, nil
}

// NEW: Get donations by user filtered by entity
func (s *service) GetDonationsByUserAndEntity(userID uint, entityID uint) ([]DonationWithUser, error) {
	donations, err := s.repo.ListByUserIDAndEntity(context.Background(), userID, entityID)
	if err != nil {
		return nil, err
	}

	for i := range donations {
		donations[i].Date = donations[i].CreatedAt
		donations[i].Type = donations[i].DonationType
		donations[i].DonorName = donations[i].UserName
		donations[i].DonorEmail = donations[i].UserEmail
		donations[i].PaymentMethod = donations[i].Method

		if donations[i].DonatedAt == nil {
			donations[i].DonatedAt = &donations[i].CreatedAt
		}
	}

	return donations, nil
}

func (s *service) GetDonationsWithFilters(filters DonationFilters, accessContext middleware.AccessContext) ([]DonationWithUser, int, error) {
	// Check permissions
	if !accessContext.CanRead() {
		return nil, 0, errors.New("read access denied")
	}

	// 🔒 UPDATED LOGIC: Verify access based on filter combination
	if filters.EntityID != 0 && filters.UserID != 0 {
		// Both filters: verify entity access AND user is requesting their own data
		entityID := accessContext.GetAccessibleEntityID()
		if entityID == nil || *entityID != filters.EntityID {
			return nil, 0, errors.New("access denied to requested entity")
		}
		if filters.UserID != accessContext.UserID {
			return nil, 0, errors.New("access denied: cannot view other users' donations")
		}
	} else if filters.EntityID != 0 {
		// Entity-based filtering: verify user has access to this entity
		entityID := accessContext.GetAccessibleEntityID()
		if entityID == nil || *entityID != filters.EntityID {
			return nil, 0, errors.New("access denied to requested entity")
		}
	} else if filters.UserID != 0 {
		// User-based filtering: verify user can only see their own donations
		if filters.UserID != accessContext.UserID {
			return nil, 0, errors.New("access denied: cannot view other users' donations")
		}
	} else {
		return nil, 0, errors.New("either entity_id or user_id must be specified")
	}

	donations, total, err := s.repo.ListWithFilters(context.Background(), filters)
	if err != nil {
		return nil, 0, err
	}

	for i := range donations {
		donations[i].Date = donations[i].CreatedAt
		donations[i].Type = donations[i].DonationType
		donations[i].DonorName = donations[i].UserName
		donations[i].DonorEmail = donations[i].UserEmail
		donations[i].PaymentMethod = donations[i].Method

		if donations[i].DonatedAt == nil {
			donations[i].DonatedAt = &donations[i].CreatedAt
		}
	}

	return donations, total, nil
}

// ==============================
// Analytics and Reporting - TEMPLE ADMIN (UPDATED)
// ==============================

func (s *service) GetDashboardStats(entityID uint, accessContext middleware.AccessContext) (*DashboardStats, error) {
	// Check permissions
	if !accessContext.CanRead() {
		return nil, errors.New("read access denied")
	}

	// Verify entity access
	accessibleEntityID := accessContext.GetAccessibleEntityID()
	if accessibleEntityID == nil || *accessibleEntityID != entityID {
		return nil, errors.New("access denied to requested entity")
	}

	ctx := context.Background()

	// Get overall stats
	totalStats, err := s.repo.GetTotalStats(ctx, entityID)
	if err != nil {
		return nil, err
	}

	// Get this month stats
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthStats, err := s.repo.GetStatsInDateRange(ctx, entityID, monthStart, now)
	if err != nil {
		return nil, err
	}

	// Get today stats
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStats, err := s.repo.GetStatsInDateRange(ctx, entityID, todayStart, now)
	if err != nil {
		return nil, err
	}

	// Get unique donor count
	donorCount, err := s.repo.GetUniqueDonorCount(ctx, entityID)
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		TotalAmount:    totalStats.Amount,
		TotalCount:     totalStats.Count,
		CompletedCount: totalStats.CompletedCount,
		PendingCount:   totalStats.PendingCount,
		FailedCount:    totalStats.FailedCount,
		ThisMonth:      monthStats.Amount,
		Today:          todayStats.Amount,
		TotalDonors:    donorCount,
		AverageAmount: func() float64 {
			if totalStats.CompletedCount > 0 {
				return totalStats.Amount / float64(totalStats.CompletedCount)
			}
			return 0
		}(),
	}, nil
}

func (s *service) GetTopDonors(entityID uint, limit int, accessContext middleware.AccessContext) ([]TopDonor, error) {
	// Check permissions
	if !accessContext.CanRead() {
		return nil, errors.New("read access denied")
	}

	// Verify entity access
	accessibleEntityID := accessContext.GetAccessibleEntityID()
	if accessibleEntityID == nil || *accessibleEntityID != entityID {
		return nil, errors.New("access denied to requested entity")
	}

	return s.repo.GetTopDonors(context.Background(), entityID, limit)
}

func (s *service) GetAnalytics(entityID uint, days int, accessContext middleware.AccessContext) (*AnalyticsData, error) {
	// Check permissions
	if !accessContext.CanRead() {
		return nil, errors.New("read access denied")
	}

	// Verify entity access
	accessibleEntityID := accessContext.GetAccessibleEntityID()
	if accessibleEntityID == nil || *accessibleEntityID != entityID {
		return nil, errors.New("access denied to requested entity")
	}

	ctx := context.Background()

	// Get donation trends
	trends, err := s.repo.GetDonationTrends(ctx, entityID, days)
	if err != nil {
		return nil, err
	}

	// Get donation by type
	byType, err := s.repo.GetDonationsByType(ctx, entityID)
	if err != nil {
		return nil, err
	}

	// Get donation by method
	byMethod, err := s.repo.GetDonationsByMethod(ctx, entityID)
	if err != nil {
		return nil, err
	}

	return &AnalyticsData{
		Trends:   trends,
		ByType:   byType,
		ByMethod: byMethod,
	}, nil
}

// ==============================
// Receipt and Export - BOTH (UPDATED)
// ==============================

// NEW: Updated to include entity validation
func (s *service) GenerateReceipt(donationID uint, userID uint, accessContext *middleware.AccessContext, entityID uint) (*Receipt, error) {
	ctx := context.Background()

	donation, err := s.repo.GetByIDWithUser(ctx, donationID)
	if err != nil {
		return nil, err
	}

	// Check access permissions
	hasAccess := false

	// For devotees - can only access their own donations within their entity
	if donation.UserID == userID && donation.EntityID == entityID {
		hasAccess = true
	}

	// For temple admins - can access donations for their entity
	if accessContext != nil {
		accessibleEntityID := accessContext.GetAccessibleEntityID()
		if accessibleEntityID != nil && *accessibleEntityID == donation.EntityID && accessContext.CanRead() {
			hasAccess = true
		}
	}

	if !hasAccess {
		return nil, errors.New("unauthorized to access this donation")
	}

	if donation.Status != StatusSuccess {
		return nil, errors.New("receipt can only be generated for successful donations")
	}

	transactionID := donation.OrderID
	if donation.PaymentID != nil {
		transactionID = *donation.PaymentID
	}

	donatedAt := donation.CreatedAt
	if donation.DonatedAt != nil {
		donatedAt = *donation.DonatedAt
	}

	return &Receipt{
		ID:             donation.ID,
		DonationAmount: donation.Amount,
		DonationType:   donation.DonationType,
		DonorName:      donation.UserName,
		DonorEmail:     donation.UserEmail,
		TransactionID:  transactionID,
		DonatedAt:      donatedAt,
		Method:         donation.Method,
		EntityName:     donation.EntityName,
		ReceiptNumber:  fmt.Sprintf("RCP-%d-%d", donation.EntityID, donation.ID),
		GeneratedAt:    time.Now(),
	}, nil
}

func (s *service) ExportDonations(filters DonationFilters, format string, accessContext middleware.AccessContext) ([]byte, string, error) {
	// Check permissions
	if !accessContext.CanRead() {
		return nil, "", errors.New("read access denied")
	}

	// Verify entity access
	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil || *entityID != filters.EntityID {
		return nil, "", errors.New("access denied to requested entity")
	}

	ctx := context.Background()

	// Get all donations matching filters
	donations, _, err := s.repo.ListWithFilters(ctx, filters)
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "csv":
		return s.exportAsCSV(donations)
	default:
		return nil, "", errors.New("unsupported export format")
	}
}

func (s *service) exportAsCSV(donations []DonationWithUser) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"ID", "Date", "Donor Name", "Donor Email", "Amount", "Type",
		"Method", "Status", "Transaction ID", "Reference ID", "Note",
	}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	// Write data
	for _, donation := range donations {
		donatedAt := donation.CreatedAt
		if donation.DonatedAt != nil {
			donatedAt = *donation.DonatedAt
		}

		record := []string{
			strconv.FormatUint(uint64(donation.ID), 10),
			donatedAt.Format("2006-01-02 15:04:05"),
			donation.UserName,
			donation.UserEmail,
			fmt.Sprintf("%.2f", donation.Amount),
			donation.DonationType,
			donation.Method,
			donation.Status,
			func() string {
				if donation.PaymentID != nil {
					return *donation.PaymentID
				}
				return donation.OrderID
			}(),
			func() string {
				if donation.ReferenceID != nil {
					return strconv.FormatUint(uint64(*donation.ReferenceID), 10)
				}
				return ""
			}(),
			func() string {
				if donation.Note != nil {
					return *donation.Note
				}
				return ""
			}(),
		}
		if err := writer.Write(record); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("donations_%d.csv", time.Now().Unix())
	return buf.Bytes(), filename, nil
}

// ==============================
// Recent Donations - BOTH (UPDATED)
// ==============================
func (s *service) GetRecentDonationsByUser(ctx context.Context, userID uint, limit int) ([]RecentDonation, error) {
	return s.repo.GetRecentDonationsByUser(ctx, userID, limit)
}

// NEW: Get recent donations by user filtered by entity
func (s *service) GetRecentDonationsByUserAndEntity(ctx context.Context, userID uint, entityID uint, limit int) ([]RecentDonation, error) {
	return s.repo.GetRecentDonationsByUserAndEntity(ctx, userID, entityID, limit)
}

func (s *service) GetRecentDonationsByEntity(ctx context.Context, entityID uint, limit int, accessContext middleware.AccessContext) ([]RecentDonation, error) {
	// Check permissions
	if !accessContext.CanRead() {
		return nil, errors.New("read access denied")
	}

	// Verify entity access
	accessibleEntityID := accessContext.GetAccessibleEntityID()
	if accessibleEntityID == nil || *accessibleEntityID != entityID {
		return nil, errors.New("access denied to requested entity")
	}

	return s.repo.GetRecentDonationsByEntity(ctx, entityID, limit)
}

// ✅ NEW: HandleRazorpayWebhook processes webhook callbacks from Razorpay
func (s *service) HandleRazorpayWebhook(orderID, paymentID, method string, amount float64) error {
	ctx := context.Background()
	
	// Get donation record to verify it exists and get user/entity info
	donation, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("❌ Webhook: Donation not found for order %s: %v", orderID, err)
		return fmt.Errorf("donation not found for order_id: %s", orderID)
	}

	// ✅ Idempotency check: Skip if already successful
	if donation.Status == StatusSuccess {
		log.Printf("⚠️ Webhook: Donation already successful for order %s, skipping", orderID)
		return nil
	}

	// Update donation to success
	now := time.Now()
	err = s.repo.UpdatePaymentDetails(ctx, orderID, UpdatePaymentDetailsParams{
		Status:    StatusSuccess,
		PaymentID: &paymentID,
		Method:    method,
		Amount:    amount,
		DonatedAt: &now,
	})

	if err != nil {
		log.Printf("❌ Webhook: Failed to update donation for order %s: %v", orderID, err)
		
		// Audit the failure
		s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "WEBHOOK_UPDATE_FAILED", map[string]interface{}{
			"order_id":   orderID,
			"payment_id": paymentID,
			"error":      err.Error(),
		}, "razorpay_webhook", "failure")
		
		return err
	}

	// Audit success
	s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "DONATION_SUCCESS_WEBHOOK", map[string]interface{}{
		"order_id":   orderID,
		"payment_id": paymentID,
		"amount":     amount,
		"method":     method,
	}, "razorpay_webhook", "success")

	log.Printf("✅ Webhook: Successfully updated donation for order %s (method: %s, amount: %.2f)", 
		orderID, method, amount)

	return nil
}
// ✅ NEW: HandleFailedPaymentWebhook processes failed payment webhooks
func (s *service) HandleFailedPaymentWebhook(orderID, paymentID string) error {
	ctx := context.Background()
	
	// Get donation record
	donation, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("❌ Webhook (Failed): Donation not found for order %s: %v", orderID, err)
		return fmt.Errorf("donation not found for order_id: %s", orderID)
	}

	// Skip if already marked as failed
	if donation.Status == StatusFailed {
		log.Printf("⚠️ Webhook (Failed): Already marked failed for order %s", orderID)
		return nil
	}

	// Update to failed status
	err = s.repo.UpdatePaymentDetails(ctx, orderID, UpdatePaymentDetailsParams{
		Status:    StatusFailed,
		PaymentID: &paymentID,
		Method:    "UNKNOWN",
		Amount:    donation.Amount, // Keep original amount
	})

	if err != nil {
		log.Printf("❌ Webhook (Failed): Failed to update for order %s: %v", orderID, err)
		
		s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "WEBHOOK_FAILED_UPDATE_ERROR", map[string]interface{}{
			"order_id":   orderID,
			"payment_id": paymentID,
			"error":      err.Error(),
		}, "razorpay_webhook", "failure")
		
		return err
	}

	// Audit
	s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "DONATION_FAILED_WEBHOOK", map[string]interface{}{
		"order_id":   orderID,
		"payment_id": paymentID,
	}, "razorpay_webhook", "success")

	log.Printf("✅ Webhook (Failed): Updated order %s to FAILED status", orderID)
	return nil
}