package donation

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// Handler represents the donation HTTP handler
type Handler struct {
	svc Service
}

// NewHandler creates a new donation handler
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// ==============================
// 🌟 1. Create Razorpay Order & Log Donation Intent
// ==============================
// @Summary Create a donation order
// @Tags Donations
// @Accept json
// @Produce json
// @Param body body CreateDonationRequest true "Donation request"
// @Success 200 {object} CreateDonationResponse
// @Failure 400 {object} gin.H
// @Router /v1/donations [post]
func (h *Handler) CreateDonation(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	if currentUser.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

	var req CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = currentUser.ID
	req.EntityID = *currentUser.EntityID

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

// ==============================
// ✅ 2. Verify Razorpay Signature (Frontend calls this after payment success)
// ==============================
// @Summary Verify Razorpay donation payment
// @Tags Donations
// @Accept json
// @Produce json
// @Param body body VerifyPaymentRequest true "Payment verification"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /v1/donations/verify [post]
func (h *Handler) VerifyDonation(c *gin.Context) {
	var req VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.VerifyAndUpdateDonation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Donation verified successfully",
		"success": true,
	})
}

// ==============================
// 🔍 3. Get My Donations (Devotee View)
// ==============================
// @Summary Get logged-in user's donations
// @Tags Donations
// @Produce json
// @Success 200 {object} gin.H
// @Router /v1/donations/my [get]
func (h *Handler) GetMyDonations(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	donations, err := h.svc.GetDonationsByUser(currentUser.ID)
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
// 🔍 4. Get All Donations for Temple (Temple Admin View) - Enhanced with Filters
// ==============================
// @Summary Get donations by temple entity ID with filters and pagination
// @Tags Donations
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(20)
// @Param status query string false "Payment status filter"
// @Param from query string false "Start date filter (YYYY-MM-DD)"
// @Param to query string false "End date filter (YYYY-MM-DD)"
// @Param type query string false "Donation type filter"
// @Param method query string false "Payment method filter"
// @Param min query float64 false "Minimum amount filter"
// @Param max query float64 false "Maximum amount filter"
// @Param search query string false "Search by donor name, email, or transaction ID"
// @Success 200 {object} DonationListResponse
// @Failure 400 {object} gin.H
// @Router /v1/donations [get]
func (h *Handler) GetDonationsByEntity(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}
	if currentUser.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

	// Parse query parameters
	filters := DonationFilters{
		EntityID: *currentUser.EntityID,
		Page:     parseIntQuery(c, "page", 1),
		Limit:    parseIntQuery(c, "limit", 20),
		Status:   c.Query("status"),
		Type:     c.Query("type"),
		Method:   c.Query("method"),
		Search:   c.Query("search"),
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

	// Parse amount filters
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

	// Handle date range presets
	switch c.Query("dateRange") {
	case "today":
		today := time.Now().Truncate(24 * time.Hour)
		filters.From = &today
		tomorrow := today.Add(24 * time.Hour)
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

	donations, total, err := h.svc.GetDonationsWithFilters(filters)
	if err != nil {
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
// @Summary Get donation dashboard statistics
// @Tags Donations
// @Produce json
// @Success 200 {object} DashboardStats
// @Router /v1/donations/dashboard [get]
func (h *Handler) GetDashboard(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}
	if currentUser.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

	stats, err := h.svc.GetDashboardStats(*currentUser.EntityID)
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
// @Summary Get top donors for entity
// @Tags Donations
// @Produce json
// @Param limit query int false "Number of top donors" default(5)
// @Success 200 {object} gin.H
// @Router /v1/donations/top-donors [get]
func (h *Handler) GetTopDonors(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}
	if currentUser.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

	limit := parseIntQuery(c, "limit", 5)
	topDonors, err := h.svc.GetTopDonors(*currentUser.EntityID, limit)
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
// @Summary Generate donation receipt
// @Tags Donations
// @Produce json
// @Param id path int true "Donation ID"
// @Success 200 {object} gin.H
// @Router /v1/donations/{id}/receipt [get]
func (h *Handler) GenerateReceipt(c *gin.Context) {
	donationIDStr := c.Param("id")
	donationID, err := strconv.ParseUint(donationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid donation ID"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	receipt, err := h.svc.GenerateReceipt(uint(donationID), currentUser.ID)
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
// @Summary Get donation analytics and trends
// @Tags Donations
// @Produce json
// @Param days query int false "Number of days for trends" default(30)
// @Success 200 {object} gin.H
// @Router /v1/donations/analytics [get]
func (h *Handler) GetAnalytics(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}
	if currentUser.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

	days := parseIntQuery(c, "days", 30)
	analytics, err := h.svc.GetAnalytics(*currentUser.EntityID, days)
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
// @Summary Export donations as CSV
// @Tags Donations
// @Produce application/octet-stream
// @Param format query string false "Export format" default(csv)
// @Success 200 {file} file
// @Router /v1/donations/export [get]
func (h *Handler) ExportDonations(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	currentUser, ok := user.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}
	if currentUser.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

	format := c.DefaultQuery("format", "csv")

	// Build filters for export
	filters := DonationFilters{
		EntityID: *currentUser.EntityID,
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

	fileContent, filename, err := h.svc.ExportDonations(filters, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", fileContent)
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



// package donation

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/sharath018/temple-management-backend/internal/auth"
// )

// // Handler represents the donation HTTP handler
// type Handler struct {
// 	svc Service
// }

// // NewHandler creates a new donation handler
// func NewHandler(svc Service) *Handler {
// 	return &Handler{svc: svc}
// }

// // ==============================
// // 🌟 1. Create Razorpay Order & Log Donation
// // ==============================
// func (h *Handler) CreateDonation(c *gin.Context) {
// 	user, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	currentUser := user.(auth.User)

// 	var req CreateDonationRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	req.UserID = currentUser.ID

// 	// ✅ REMOVE this in test mode. Let Razorpay choose method (card, netbanking, etc.)
// 	// req.Method = "UPI"

// 	if currentUser.EntityID == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
// 		return
// 	}
// 	req.EntityID = *currentUser.EntityID

// 	order, err := h.svc.StartDonation(req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, order)
// }


// // ==============================
// // ✅ 2. Verify Payment Signature (Client-side call after success)
// // ==============================
// func (h *Handler) VerifyDonation(c *gin.Context) {
// 	var req VerifyPaymentRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := h.svc.VerifyAndUpdateDonation(req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Donation verified successfully"})
// }

// // ==============================
// // 🔍 3. Get My Donations
// // ==============================
// func (h *Handler) GetMyDonations(c *gin.Context) {
// 	user, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	currentUser := user.(auth.User)

// 	donations, err := h.svc.GetDonationsByUser(currentUser.ID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, donations)
// }

// // ==============================
// // 🔍 4. Get Donations for Temple Admin by entity_id
// // ==============================
// func (h *Handler) GetDonationsByEntity(c *gin.Context) {
// 	user, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}
// 	currentUser := user.(auth.User)

// 	if currentUser.EntityID == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
// 		return
// 	}

// 	donations, err := h.svc.GetDonationsByEntity(*currentUser.EntityID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, donations)
// }