package donation

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
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

// EntityBankDetails - temple's registered bank/UPI info including Razorpay credentials.
type EntityBankDetails struct {
	AccountHolderName string
	AccountNumber     string
	IFSCCode          string
	UPIID             string
	BankName          string
	BranchName        string
	AccountType       string
	RazorpayKeyID     string // tenant's own Razorpay key — stored in tenant_bank_account_details
	RazorpaySecret    string // tenant's own Razorpay secret — stored in tenant_bank_account_details
}

// EntityRepository - minimal interface to fetch temple bank details.
type EntityRepository interface {
	GetBankDetailsByEntityID(ctx context.Context, entityID uint) (*EntityBankDetails, error)
}

// Service Interface
type Service interface {
	StartDonation(req CreateDonationRequest) (*CreateDonationResponse, error)
	VerifyAndUpdateDonation(req VerifyPaymentRequest) error
	HandleRazorpayWebhook(orderID, paymentID, method string, amount float64) error
	HandleFailedPaymentWebhook(orderID, paymentID string) error

	GetDonationsByUser(userID uint) ([]DonationWithUser, error)
	GetDonationsByUserAndEntity(userID uint, entityID uint) ([]DonationWithUser, error)
	GetDonationsWithFilters(filters DonationFilters, accessContext middleware.AccessContext) ([]DonationWithUser, int, error)

	GetDashboardStats(entityID uint, accessContext middleware.AccessContext) (*DashboardStats, error)
	GetTopDonors(entityID uint, limit int, accessContext middleware.AccessContext) ([]TopDonor, error)
	GetAnalytics(entityID uint, days int, accessContext middleware.AccessContext) (*AnalyticsData, error)

	GenerateReceipt(donationID uint, userID uint, accessContext *middleware.AccessContext, entityID uint) (*Receipt, error)
	ExportDonations(filters DonationFilters, format string, accessContext middleware.AccessContext) ([]byte, string, error)

	GetRecentDonationsByUser(ctx context.Context, userID uint, limit int) ([]RecentDonation, error)
	GetRecentDonationsByUserAndEntity(ctx context.Context, userID uint, entityID uint, limit int) ([]RecentDonation, error)
	GetRecentDonationsByEntity(ctx context.Context, entityID uint, limit int, accessContext middleware.AccessContext) ([]RecentDonation, error)
}

type service struct {
	repo       Repository
	entityRepo EntityRepository
	cfg        *config.Config // kept for non-Razorpay config
	auditSvc   auditlog.Service
}

// NewService creates a donation service without entity repo
func NewService(repo Repository, cfg *config.Config, auditSvc auditlog.Service) Service {
	return &service{repo: repo, cfg: cfg, auditSvc: auditSvc}
}

// NewServiceWithEntityRepo creates a donation service with entity repo (required for Razorpay)
func NewServiceWithEntityRepo(repo Repository, entityRepo EntityRepository, cfg *config.Config, auditSvc auditlog.Service) Service {
	return &service{repo: repo, entityRepo: entityRepo, cfg: cfg, auditSvc: auditSvc}
}

// ==============================
// getTenantBank — fetches tenant bank details once and reuses
// ==============================
func (s *service) getTenantBank(ctx context.Context, entityID uint) *EntityBankDetails {
	if s.entityRepo == nil {
		return nil
	}
	eb, err := s.entityRepo.GetBankDetailsByEntityID(ctx, entityID)
	if err != nil || eb == nil {
		log.Printf("⚠️ Could not fetch tenant bank for entity=%d: %v", entityID, err)
		return nil
	}
	return eb
}

// ==============================
// StartDonation
// ==============================
func (s *service) StartDonation(req CreateDonationRequest) (*CreateDonationResponse, error) {
	ctx := context.Background()

	// ── Fetch tenant bank details (includes Razorpay keys) ───────────────
	eb := s.getTenantBank(ctx, req.EntityID)

	// ── Require tenant's own Razorpay key — NO platform fallback ─────────
	// Keys are stored per-tenant in tenant_bank_account_details table.
	// If tenant hasn't configured Razorpay, they should use UPI Direct instead.
	if eb == nil || eb.RazorpayKeyID == "" || eb.RazorpaySecret == "" {
		log.Printf("❌ [StartDonation] Tenant entity=%d has no Razorpay keys in DB", req.EntityID)
		return nil, fmt.Errorf("temple has not configured Razorpay payment gateway yet. Please use UPI Direct or contact the temple admin")
	}

	// ── Create per-tenant Razorpay client using tenant's own credentials ──
	razorpayClient := razorpay.NewClient(eb.RazorpayKeyID, eb.RazorpaySecret)
	log.Printf("🔑 [StartDonation] Using tenant Razorpay key for entity=%d key=%s",
		req.EntityID, maskKey(eb.RazorpayKeyID))

	// ── Create Razorpay order ─────────────────────────────────────────────
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

	order, err := razorpayClient.Order.Create(data, nil)
	if err != nil {
		s.auditSvc.LogAction(ctx, &req.UserID, &req.EntityID, "DONATION_INITIATED",
			map[string]interface{}{"amount": req.Amount, "error": err.Error()},
			req.IPAddress, "failure")
		return nil, fmt.Errorf("razorpay order creation failed: %w", err)
	}

	orderID, ok := order["id"].(string)
	if !ok {
		return nil, errors.New("unable to extract order_id from Razorpay response")
	}

	// ── Save donation record ─────────────────────────────────────────────
	donation := &Donation{
		UserID:       req.UserID,
		EntityID:     req.EntityID,
		Amount:       req.Amount,
		DonationType: req.DonationType,
		ReferenceID:  req.ReferenceID,
		Method:       " ",
		Status:       StatusPending,
		OrderID:      orderID,
		Note:         req.Note,
	}
	if err := s.repo.Create(context.Background(), donation); err != nil {
		return nil, fmt.Errorf("failed to create donation record: %w", err)
	}

	s.auditSvc.LogAction(ctx, &req.UserID, &req.EntityID, "DONATION_INITIATED",
		map[string]interface{}{"amount": req.Amount, "order_id": orderID},
		req.IPAddress, "success")

	// ── Build tenant info for frontend display ───────────────────────────
	tenantInfo := TenantPaymentInfo{
		AccountHolderName: eb.AccountHolderName,
		AccountNumber:     eb.AccountNumber,
		BankName:          eb.BankName,
		BranchName:        eb.BranchName,
		IFSCCode:          eb.IFSCCode,
		AccountType:       eb.AccountType,
		UPIID:             eb.UPIID,
	}
	log.Printf("🏦 Tenant bank for entity=%d: holder=%s upi=%s bank=%s",
		req.EntityID, eb.AccountHolderName, eb.UPIID, eb.BankName)

	return &CreateDonationResponse{
		OrderID:     orderID,
		Amount:      req.Amount,
		Currency:    "INR",
		RazorpayKey: eb.RazorpayKeyID, // ✅ tenant's own key sent to frontend
		Tenant:      tenantInfo,
	}, nil
}

// ==============================
// VerifyAndUpdateDonation
// ==============================
func (s *service) VerifyAndUpdateDonation(req VerifyPaymentRequest) error {
	ctx := context.Background()

	// ── Step 1: Find donation to get entityID ────────────────────────────
	donation, err := s.repo.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		log.Printf("❌ Donation not found for order=%s: %v", req.OrderID, err)
		return fmt.Errorf("donation not found: %w", err)
	}

	// ── Step 2: Fetch tenant bank (must have Razorpay secret) ────────────
	eb := s.getTenantBank(ctx, donation.EntityID)

	// ── Step 3: Require tenant's Razorpay secret for HMAC ────────────────
	// Must use the SAME secret that was used to CREATE the order.
	// No platform fallback — if tenant has no secret, verification cannot proceed.
	if eb == nil || eb.RazorpaySecret == "" {
		log.Printf("❌ [Verify] Tenant entity=%d has no Razorpay secret in DB", donation.EntityID)
		return fmt.Errorf("temple Razorpay secret not configured — cannot verify payment")
	}
	log.Printf("🔑 [Verify] Using tenant secret for entity=%d", donation.EntityID)

	// ── Step 4: HMAC Signature Verification ──────────────────────────────
	// Razorpay format: "<order_id>|<payment_id>" signed with secret
	signatureData := req.OrderID + "|" + req.PaymentID
	mac := hmac.New(sha256.New, []byte(eb.RazorpaySecret))
	mac.Write([]byte(signatureData))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	log.Printf("🔐 Verify: order=%s payment=%s", req.OrderID, req.PaymentID)
	log.Printf("🔐 Expected sig: %s", expectedSig)
	log.Printf("🔐 Received sig: %s", req.RazorpaySig)

	if expectedSig != req.RazorpaySig {
		log.Printf("❌ Signature mismatch for order=%s", req.OrderID)
		return fmt.Errorf("invalid payment signature")
	}
	log.Printf("✅ Signature verified for order=%s", req.OrderID)

	// ── Step 5: Extract tenant bank fields for donation record ────────────
	accountHolder := eb.AccountHolderName
	accountNumber := eb.AccountNumber
	ifscCode      := eb.IFSCCode
	upiID         := eb.UPIID
	accountType   := eb.AccountType
	log.Printf("🏦 Tenant bank: entity=%d holder=%s upi=%s", donation.EntityID, accountHolder, upiID)

	// ── Step 6: Update donation → SUCCESS ────────────────────────────────
	now       := time.Now()
	paymentID := req.PaymentID

	if err := s.repo.UpdatePaymentDetails(ctx, req.OrderID, UpdatePaymentDetailsParams{
		Status:            StatusSuccess,
		PaymentID:         &paymentID,
		Method:            "upi",
		Amount:            donation.Amount,
		DonatedAt:         &now,
		AccountHolderName: accountHolder,
		AccountNumber:     accountNumber,
		IFSCCode:          ifscCode,
		UPIID:             upiID,
		AccountType:       accountType,
	}); err != nil {
		log.Printf("❌ Failed to update donation order=%s: %v", req.OrderID, err)
		return fmt.Errorf("failed to update donation: %w", err)
	}

	log.Printf("✅ Donation SUCCESS order=%s payment=%s tenant=%s", req.OrderID, paymentID, accountHolder)
	return nil
}

// ==============================
// Data Retrieval
// ==============================

func (s *service) GetDonationsByUser(userID uint) ([]DonationWithUser, error) {
	donations, err := s.repo.ListByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return s.enrichDonations(donations), nil
}

func (s *service) GetDonationsByUserAndEntity(userID uint, entityID uint) ([]DonationWithUser, error) {
	donations, err := s.repo.ListByUserIDAndEntity(context.Background(), userID, entityID)
	if err != nil {
		return nil, err
	}
	return s.enrichDonations(donations), nil
}

func (s *service) GetDonationsWithFilters(filters DonationFilters, accessContext middleware.AccessContext) ([]DonationWithUser, int, error) {
	if !accessContext.CanRead() {
		return nil, 0, errors.New("read access denied")
	}

	if filters.EntityID != 0 && filters.UserID != 0 {
		entityID := accessContext.GetAccessibleEntityID()
		if entityID == nil || *entityID != filters.EntityID {
			return nil, 0, errors.New("access denied to requested entity")
		}
		if filters.UserID != accessContext.UserID {
			return nil, 0, errors.New("access denied: cannot view other users' donations")
		}
	} else if filters.EntityID != 0 {
		entityID := accessContext.GetAccessibleEntityID()
		if entityID == nil || *entityID != filters.EntityID {
			return nil, 0, errors.New("access denied to requested entity")
		}
	} else if filters.UserID != 0 {
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
	return s.enrichDonations(donations), total, nil
}

func (s *service) enrichDonations(donations []DonationWithUser) []DonationWithUser {
	for i := range donations {
		donations[i].Date          = donations[i].CreatedAt
		donations[i].Type          = donations[i].DonationType
		donations[i].DonorName     = donations[i].UserName
		donations[i].DonorEmail    = donations[i].UserEmail
		donations[i].PaymentMethod = donations[i].Method
		if donations[i].DonatedAt == nil {
			donations[i].DonatedAt = &donations[i].CreatedAt
		}
	}
	return donations
}

// ==============================
// Analytics and Reporting
// ==============================

func (s *service) GetDashboardStats(entityID uint, accessContext middleware.AccessContext) (*DashboardStats, error) {
	if !accessContext.CanRead() {
		return nil, errors.New("read access denied")
	}
	accessibleEntityID := accessContext.GetAccessibleEntityID()
	if accessibleEntityID == nil || *accessibleEntityID != entityID {
		return nil, errors.New("access denied to requested entity")
	}

	ctx        := context.Background()
	totalStats, err := s.repo.GetTotalStats(ctx, entityID)
	if err != nil {
		return nil, err
	}

	now        := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthStats, err := s.repo.GetStatsInDateRange(ctx, entityID, monthStart, now)
	if err != nil {
		return nil, err
	}

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayStats, err := s.repo.GetStatsInDateRange(ctx, entityID, todayStart, now)
	if err != nil {
		return nil, err
	}

	donorCount, err := s.repo.GetUniqueDonorCount(ctx, entityID)
	if err != nil {
		return nil, err
	}

	avgAmount := 0.0
	if totalStats.CompletedCount > 0 {
		avgAmount = totalStats.Amount / float64(totalStats.CompletedCount)
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
		AverageAmount:  avgAmount,
	}, nil
}

func (s *service) GetTopDonors(entityID uint, limit int, accessContext middleware.AccessContext) ([]TopDonor, error) {
	if !accessContext.CanRead() {
		return nil, errors.New("read access denied")
	}
	accessibleEntityID := accessContext.GetAccessibleEntityID()
	if accessibleEntityID == nil || *accessibleEntityID != entityID {
		return nil, errors.New("access denied to requested entity")
	}
	return s.repo.GetTopDonors(context.Background(), entityID, limit)
}

func (s *service) GetAnalytics(entityID uint, days int, accessContext middleware.AccessContext) (*AnalyticsData, error) {
	if !accessContext.CanRead() {
		return nil, errors.New("read access denied")
	}
	accessibleEntityID := accessContext.GetAccessibleEntityID()
	if accessibleEntityID == nil || *accessibleEntityID != entityID {
		return nil, errors.New("access denied to requested entity")
	}

	ctx := context.Background()
	trends, err := s.repo.GetDonationTrends(ctx, entityID, days)
	if err != nil {
		return nil, err
	}
	byType, err := s.repo.GetDonationsByType(ctx, entityID)
	if err != nil {
		return nil, err
	}
	byMethod, err := s.repo.GetDonationsByMethod(ctx, entityID)
	if err != nil {
		return nil, err
	}
	return &AnalyticsData{Trends: trends, ByType: byType, ByMethod: byMethod}, nil
}

// ==============================
// Receipt and Export
// ==============================

func (s *service) GenerateReceipt(donationID uint, userID uint, accessContext *middleware.AccessContext, entityID uint) (*Receipt, error) {
	ctx := context.Background()
	donation, err := s.repo.GetByIDWithUser(ctx, donationID)
	if err != nil {
		return nil, err
	}

	hasAccess := (donation.UserID == userID && donation.EntityID == entityID)
	if accessContext != nil {
		if accessibleEntityID := accessContext.GetAccessibleEntityID(); accessibleEntityID != nil &&
			*accessibleEntityID == donation.EntityID && accessContext.CanRead() {
			hasAccess = true
		}
	}
	if !hasAccess {
		return nil, errors.New("unauthorized to access this donation")
	}
	if donation.Status != StatusSuccess {
		return nil, errors.New("receipt can only be generated for successful donations")
	}

	txnID := donation.OrderID
	if donation.PaymentID != nil {
		txnID = *donation.PaymentID
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
		TransactionID:  txnID,
		DonatedAt:      donatedAt,
		Method:         donation.Method,
		EntityName:     donation.EntityName,
		ReceiptNumber:  fmt.Sprintf("RCP-%d-%d", donation.EntityID, donation.ID),
		GeneratedAt:    time.Now(),
	}, nil
}

func (s *service) ExportDonations(filters DonationFilters, format string, accessContext middleware.AccessContext) ([]byte, string, error) {
	if !accessContext.CanRead() {
		return nil, "", errors.New("read access denied")
	}
	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil || *entityID != filters.EntityID {
		return nil, "", errors.New("access denied to requested entity")
	}

	donations, _, err := s.repo.ListWithFilters(context.Background(), filters)
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

	_ = writer.Write([]string{"ID", "Date", "Donor Name", "Donor Email", "Amount", "Type", "Method", "Status", "Transaction ID", "Reference ID", "Note"})

	for _, d := range donations {
		donatedAt := d.CreatedAt
		if d.DonatedAt != nil {
			donatedAt = *d.DonatedAt
		}
		txnID := d.OrderID
		if d.PaymentID != nil {
			txnID = *d.PaymentID
		}
		refID, note := "", ""
		if d.ReferenceID != nil {
			refID = strconv.FormatUint(uint64(*d.ReferenceID), 10)
		}
		if d.Note != nil {
			note = *d.Note
		}
		_ = writer.Write([]string{
			strconv.FormatUint(uint64(d.ID), 10),
			donatedAt.Format("2006-01-02 15:04:05"),
			d.UserName, d.UserEmail,
			fmt.Sprintf("%.2f", d.Amount),
			d.DonationType, d.Method, d.Status,
			txnID, refID, note,
		})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), fmt.Sprintf("donations_%d.csv", time.Now().Unix()), nil
}

// ==============================
// Recent Donations
// ==============================

func (s *service) GetRecentDonationsByUser(ctx context.Context, userID uint, limit int) ([]RecentDonation, error) {
	return s.repo.GetRecentDonationsByUser(ctx, userID, limit)
}

func (s *service) GetRecentDonationsByUserAndEntity(ctx context.Context, userID uint, entityID uint, limit int) ([]RecentDonation, error) {
	return s.repo.GetRecentDonationsByUserAndEntity(ctx, userID, entityID, limit)
}

func (s *service) GetRecentDonationsByEntity(ctx context.Context, entityID uint, limit int, accessContext middleware.AccessContext) ([]RecentDonation, error) {
	if !accessContext.CanRead() {
		return nil, errors.New("read access denied")
	}
	accessibleEntityID := accessContext.GetAccessibleEntityID()
	if accessibleEntityID == nil || *accessibleEntityID != entityID {
		return nil, errors.New("access denied to requested entity")
	}
	return s.repo.GetRecentDonationsByEntity(ctx, entityID, limit)
}

// ==============================
// Webhook Handlers
// ==============================

func (s *service) HandleRazorpayWebhook(orderID, paymentID, method string, amount float64) error {
	ctx := context.Background()

	donation, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("❌ Webhook: Donation not found for order %s: %v", orderID, err)
		return fmt.Errorf("donation not found for order_id: %s", orderID)
	}
	if donation.Status == StatusSuccess {
		log.Printf("⚠️ Webhook: Already successful for order %s", orderID)
		return nil
	}

	// Fetch tenant bank for account holder name and UPI
	eb := s.getTenantBank(ctx, donation.EntityID)

	var accountHolderName, upiID string
	if eb != nil {
		accountHolderName = eb.AccountHolderName
		upiID             = eb.UPIID
	}

	now := time.Now()
	if err = s.repo.UpdatePaymentDetails(ctx, orderID, UpdatePaymentDetailsParams{
		Status:            StatusSuccess,
		PaymentID:         &paymentID,
		Method:            method,
		Amount:            amount,
		DonatedAt:         &now,
		AccountHolderName: accountHolderName,
		UPIID:             upiID,
	}); err != nil {
		log.Printf("❌ Webhook: Failed to update order %s: %v", orderID, err)
		s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "WEBHOOK_UPDATE_FAILED",
			map[string]interface{}{"order_id": orderID, "payment_id": paymentID, "error": err.Error()},
			"razorpay_webhook", "failure")
		return err
	}

	s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "DONATION_SUCCESS_WEBHOOK",
		map[string]interface{}{"order_id": orderID, "payment_id": paymentID, "amount": amount, "method": method},
		"razorpay_webhook", "success")
	log.Printf("✅ Webhook: Updated order %s (method=%s amount=%.2f payee=%s)", orderID, method, amount, accountHolderName)
	return nil
}

func (s *service) HandleFailedPaymentWebhook(orderID, paymentID string) error {
	ctx := context.Background()

	donation, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("❌ Webhook(Failed): Donation not found for order %s: %v", orderID, err)
		return fmt.Errorf("donation not found for order_id: %s", orderID)
	}
	if donation.Status == StatusFailed {
		log.Printf("⚠️ Webhook(Failed): Already failed for order %s", orderID)
		return nil
	}

	if err = s.repo.UpdatePaymentDetails(ctx, orderID, UpdatePaymentDetailsParams{
		Status:    StatusFailed,
		PaymentID: &paymentID,
		Method:    "UNKNOWN",
		Amount:    donation.Amount,
	}); err != nil {
		log.Printf("❌ Webhook(Failed): Failed to update order %s: %v", orderID, err)
		s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "WEBHOOK_FAILED_UPDATE_ERROR",
			map[string]interface{}{"order_id": orderID, "error": err.Error()},
			"razorpay_webhook", "failure")
		return err
	}

	s.auditSvc.LogAction(ctx, &donation.UserID, &donation.EntityID, "DONATION_FAILED_WEBHOOK",
		map[string]interface{}{"order_id": orderID, "payment_id": paymentID},
		"razorpay_webhook", "success")
	log.Printf("✅ Webhook(Failed): Updated order %s to FAILED", orderID)
	return nil
}

// maskKey returns first 8 chars of key for safe logging (never log full secret)
func maskKey(key string) string {
	if len(key) <= 8 {
		return key
	}
	return key[:8] + "..."
}