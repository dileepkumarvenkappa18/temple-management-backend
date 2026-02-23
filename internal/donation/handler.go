package donation

import (
	//"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/middleware"
)

// Handler represents the donation HTTP handler
type Handler struct {
	svc Service
}

// NewHandler creates a new donation handler
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// ===========================
// 📌 Extract Access Context
func getAccessContextFromContext(c *gin.Context) (middleware.AccessContext, bool) {
	accessContextRaw, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return middleware.AccessContext{}, false
	}

	accessContext, ok := accessContextRaw.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access context"})
		return middleware.AccessContext{}, false
	}

	return accessContext, true
}


func getEntityIDFromRequest(c *gin.Context, accessContext middleware.AccessContext) (uint, error) {
	// Try URL parameter first
	entityIDParam := c.Param("entity_id")
	if entityIDParam != "" {
		id, err := strconv.ParseUint(entityIDParam, 10, 32)
		if err == nil {
			return uint(id), nil
		}
	}

	// Try query parameter next
	entityIDQuery := c.Query("entity_id")
	if entityIDQuery != "" {
		id, err := strconv.ParseUint(entityIDQuery, 10, 32)
		if err == nil {
			return uint(id), nil
		}
	}

	// Fall back to access context (from JWT token)
	contextEntityID := accessContext.GetAccessibleEntityID()
	if contextEntityID != nil {
		return *contextEntityID, nil
	}

	return 0, nil
}

// ==============================
// 🌟 1. Create Donation
// ==============================
func (h *Handler) CreateDonation(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID, err := getEntityIDFromRequest(c, accessContext)
	if err != nil || entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
		return
	}

	var req CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = accessContext.UserID
	req.EntityID = entityID
	req.IPAddress = middleware.GetIPFromContext(c)

	order, err := h.svc.StartDonation(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    order,
		"success": true,
	})
}

func (h *Handler) VerifyDonation(c *gin.Context) {
	var req VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.IPAddress = middleware.GetIPFromContext(c)

	if err := h.svc.VerifyAndUpdateDonation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Donation verification initiated",
		"success": true,
	})
}

// ==============================
// 🔍 3. Get My Donations
// ==============================
func (h *Handler) GetMyDonations(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID, err := getEntityIDFromRequest(c, accessContext)
	if err != nil || entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
		return
	}

	donations, err := h.svc.GetDonationsByUserAndEntity(accessContext.UserID, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    donations,
		"count":   len(donations),
		"success": true,
	})
}

// ==============================
// 🔍 4. Get All Donations for Temple - UPDATED: Removed header fallback
// ==============================
func (h *Handler) GetDonationsByEntity(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// ✅ EXTRACT entity_id from route param OR query param ONLY
	var entityID uint

	// 1️⃣ Try route param: /entity/:entityId/...
	if entityIDStr := c.Param("entityId"); entityIDStr != "" {
		id, err := strconv.ParseUint(entityIDStr, 10, 64)
		if err != nil || id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id in path"})
			return
		}
		entityID = uint(id)
	}

	// 2️⃣ Fallback to query param ?entity_id=
	if entityID == 0 {
		if id := c.Query("entity_id"); id != "" {
			parsedID, err := strconv.ParseUint(id, 10, 64)
			if err != nil || parsedID == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id in query"})
				return
			}
			entityID = uint(parsedID)
		}
	}

	// 3️⃣ Final fallback: Access context (from JWT)
	if entityID == 0 {
		contextEntityID := accessContext.GetAccessibleEntityID()
		if contextEntityID != nil {
			entityID = *contextEntityID
		}
	}

	// 🔒 REQUIRE entity_id for devotees and templeadmins
	if (accessContext.RoleName == "devotee" || accessContext.RoleName == "templeadmin") && entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "entity_id is required",
		})
		return
	}

	// Initialize filters
	filters := DonationFilters{
		Page:   parseIntQuery(c, "page", 1),
		Limit:  parseIntQuery(c, "limit", 20),
		Status: c.Query("status"),
		Type:   c.Query("type"),
		Method: c.Query("method"),
		Search: c.Query("search"),
	}

	// 🔒 ROLE-BASED FILTERING WITH ENTITY ISOLATION
	switch accessContext.RoleName {
	case "devotee":
		// Devotees see ONLY their own donations for the CURRENT entity
		filters.UserID = accessContext.UserID
		filters.EntityID = entityID

	case "admin", "superadmin":
		// Admins can see all donations
		if entityID != 0 {
			// If entity specified, filter by entity
			filters.EntityID = entityID
			filters.UserID = 0
		} else {
			// If no entity specified, show all (admin privilege)
			filters.EntityID = 0
			filters.UserID = 0
		}

	case "trustee", "staff", "templeadmin":
		// Trustees/Staff can see entity donations
		if entityID != 0 {
			filters.EntityID = entityID
			filters.UserID = 0
		} else {
			// No entity specified, show their own donations only
			filters.UserID = accessContext.UserID
			filters.EntityID = 0
		}

	default:
		// Unknown role: show only their own donations
		filters.UserID = accessContext.UserID
		if entityID != 0 {
			filters.EntityID = entityID
		}
	}

	// Date filters
	if fromStr := c.Query("from"); fromStr != "" {
		if from, err := time.Parse("2006-01-02", fromStr); err == nil {
			filters.From = &from
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if to, err := time.Parse("2006-01-02", toStr); err == nil {
			filters.To = &to
		}
	}

	// Amount filters
	if minStr := c.Query("min"); minStr != "" {
		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			filters.MinAmount = &min
		}
	}
	if maxStr := c.Query("max"); maxStr != "" {
		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			filters.MaxAmount = &max
		}
	}

	// Date presets
	switch c.Query("dateRange") {
	case "today":
		today := time.Now().Truncate(24 * time.Hour)
		tomorrow := today.Add(24 * time.Hour)
		filters.From = &today
		filters.To = &tomorrow
	case "week":
		weekAgo := time.Now().AddDate(0, 0, -7)
		filters.From = &weekAgo
	case "month":
		monthAgo := time.Now().AddDate(0, -1, 0)
		filters.From = &monthAgo
	case "year":
		yearAgo := time.Now().AddDate(-1, 0, 0)
		filters.From = &yearAgo
	}

	// 📝 DEBUG LOG
	log.Printf("🔍 Donation Filters Applied - Role: %s, UserID: %d, EntityID: %d",
		accessContext.RoleName, filters.UserID, filters.EntityID)

	donations, total, err := h.svc.GetDonationsWithFilters(filters, accessContext)
	if err != nil {
		log.Printf("❌ Error fetching donations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        donations,
		"total":       total,
		"page":        filters.Page,
		"limit":       filters.Limit,
		"total_pages": (total + filters.Limit - 1) / filters.Limit,
		"success":     true,
	})
}

// ==============================
// 📊 5. Get Donation Dashboard Stats
// ==============================
func (h *Handler) GetDashboard(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID, err := getEntityIDFromRequest(c, accessContext)
	if err != nil || entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
		return
	}

	stats, err := h.svc.GetDashboardStats(entityID, accessContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    stats,
		"success": true,
	})
}

// ==============================
// 🏆 6. Get Top Donors
// ==============================
func (h *Handler) GetTopDonors(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID, err := getEntityIDFromRequest(c, accessContext)
	if err != nil || entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
		return
	}

	limit := parseIntQuery(c, "limit", 5)
	topDonors, err := h.svc.GetTopDonors(entityID, limit, accessContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    topDonors,
		"success": true,
	})
}

// ==============================
// 📄 7. Generate Receipt
// ==============================
func (h *Handler) GenerateReceipt(c *gin.Context) {
	donationIDStr := c.Param("id")
	donationID, err := strconv.ParseUint(donationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid donation ID"})
		return
	}

	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Extract entity ID
	var entityID uint
	if accessContext.RoleName == "devotee" {
		extractedEntityID, err := getEntityIDFromRequest(c, accessContext)
		if err != nil || extractedEntityID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
			return
		}
		entityID = extractedEntityID
	} else {
		extractedEntityID, err := getEntityIDFromRequest(c, accessContext)
		if err != nil || extractedEntityID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
			return
		}
		entityID = extractedEntityID
	}

	receipt, err := h.svc.GenerateReceipt(uint(donationID), accessContext.UserID, &accessContext, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    receipt,
		"success": true,
	})
}

// ==============================
// 📈 8. Get Donation Analytics
// ==============================
func (h *Handler) GetAnalytics(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID, err := getEntityIDFromRequest(c, accessContext)
	if err != nil || entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
		return
	}

	days := parseIntQuery(c, "days", 30)
	analytics, err := h.svc.GetAnalytics(entityID, days, accessContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    analytics,
		"success": true,
	})
}

// ==============================
// 📊 9. Export Donations
// ==============================
func (h *Handler) ExportDonations(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID, err := getEntityIDFromRequest(c, accessContext)
	if err != nil || entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
		return
	}

	format := c.DefaultQuery("format", "csv")

	// Build filters for export
	filters := DonationFilters{
		EntityID: entityID,
		Status:   c.Query("status"),
		Type:     c.Query("type"),
		Method:   c.Query("method"),
		Search:   c.Query("search"),
		Page:     1,
		Limit:    10000, // Large limit for export
	}

	// Parse date filters
	if fromStr := c.Query("from"); fromStr != "" {
		if from, err := time.Parse("2006-01-02", fromStr); err == nil {
			filters.From = &from
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if to, err := time.Parse("2006-01-02", toStr); err == nil {
			filters.To = &to
		}
	}

	fileContent, filename, err := h.svc.ExportDonations(filters, format, accessContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", fileContent)
}

// ==============================
// 🕐 10. Get Recent Donations
// ==============================
func (h *Handler) GetRecentDonations(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	entityID, err := getEntityIDFromRequest(c, accessContext)
	if err != nil || entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
		return
	}

	// For devotees, get their own donations within the entity
	if accessContext.RoleName == "devotee" {
		recent, err := h.svc.GetRecentDonationsByUserAndEntity(c.Request.Context(), accessContext.UserID, entityID, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch recent donations"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"recent_donations": recent,
			"success":          true,
		})
		return
	}

	// For temple admins, get donations for their accessible entity
	recent, err := h.svc.GetRecentDonationsByEntity(c.Request.Context(), entityID, limit, accessContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch recent donations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recent_donations": recent,
		"success":          true,
	})
}

// Helper function to parse integer query parameters
func parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	if str := c.Query(key); str != "" {
		if val, err := strconv.Atoi(str); err == nil && val > 0 {
			return val
		}
	}
	return defaultValue
}

// ==============================
// 🔔 11. Razorpay Webhook Handler
// ==============================
func (h *Handler) HandleWebhook(c *gin.Context) {
	// Step 1: Verify webhook signature
	webhookSecret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("⚠️ RAZORPAY_WEBHOOK_SECRET not configured")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "webhook secret not configured"})
		return
	}

	// Get signature from header
	signature := c.GetHeader("X-Razorpay-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing signature"})
		return
	}

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Verify signature
	expectedSignature := hmac.New(sha256.New, []byte(webhookSecret))
	expectedSignature.Write(body)
	computedSignature := hex.EncodeToString(expectedSignature.Sum(nil))

	if computedSignature != signature {
		log.Printf("❌ Webhook signature mismatch - Expected: %s, Got: %s", computedSignature, signature)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	// Step 2: Parse webhook payload
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	event, ok := payload["event"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing event type"})
		return
	}

	// Step 3: Handle different webhook events
	switch event {
	case "payment.captured":
		h.handlePaymentCaptured(c, payload)
	case "payment.failed":
		h.handlePaymentFailed(c, payload)
	default:
		log.Printf("⚠️ Unhandled webhook event: %s", event)
		c.JSON(http.StatusOK, gin.H{"message": "event ignored"})
	}
}

// Helper: Handle successful payment webhook
func (h *Handler) handlePaymentCaptured(c *gin.Context, payload map[string]interface{}) {
	paymentData, ok := payload["payload"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload structure"})
		return
	}

	payment, ok := paymentData["payment"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing payment data"})
		return
	}

	entity, ok := payment["entity"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing payment entity"})
		return
	}

	// Extract payment details
	orderID, _ := entity["order_id"].(string)
	paymentID, _ := entity["id"].(string)
	method, _ := entity["method"].(string)
	amountPaise, _ := entity["amount"].(float64)

	if orderID == "" || paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing order_id or payment_id"})
		return
	}

	amount := amountPaise / 100 // Convert paise to rupees

	// Update donation via service
	if err := h.svc.HandleRazorpayWebhook(orderID, paymentID, method, amount); err != nil {
		log.Printf("❌ Webhook processing failed for order %s: %v", orderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update donation"})
		return
	}

	log.Printf("✅ Webhook processed successfully for order %s", orderID)
	c.JSON(http.StatusOK, gin.H{"message": "webhook processed"})
}

// Helper: Handle failed payment webhook
func (h *Handler) handlePaymentFailed(c *gin.Context, payload map[string]interface{}) {
	paymentData, ok := payload["payload"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload structure"})
		return
	}

	payment, ok := paymentData["payment"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing payment data"})
		return
	}

	entity, ok := payment["entity"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing payment entity"})
		return
	}

	orderID, _ := entity["order_id"].(string)
	paymentID, _ := entity["id"].(string)

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing order_id"})
		return
	}

	err := h.svc.HandleFailedPaymentWebhook(orderID, paymentID)
	if err != nil {
		log.Printf("❌ Webhook: Failed to process failed payment for order %s: %v", orderID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update donation"})
		return
	}

	log.Printf("✅ Failed payment webhook processed for order %s", orderID)
	c.JSON(http.StatusOK, gin.H{"message": "webhook processed"})
}