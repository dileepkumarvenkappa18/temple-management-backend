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
// ✅ 2. Verify Razorpay Signature
// ==============================
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
// 🔍 3. Get My Donations (Devotee View) - FIXED
// ==============================
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
// 🔍 4. Get All Donations for Temple (Temple Admin View) - FIXED
// ==============================
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

// ==============================
// 🕐 10. Get Recent Donations (FIXED - Only for logged in user)
// ==============================
func (h *Handler) GetRecentDonations(c *gin.Context) {
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

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	recent, err := h.svc.GetRecentDonationsByUser(c.Request.Context(), currentUser.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch recent donations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recent_donations": recent,
		"success": true,
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