package seva

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	//"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/middleware"
)

type Handler struct {
	service  Service
	auditSvc auditlog.Service
	repo     Repository
}

func NewHandler(service Service, auditSvc auditlog.Service, repo Repository) *Handler {
	return &Handler{
		service:  service,
		auditSvc: auditSvc,
		repo:     repo,
	}
}

// ===========================
// 📌 Extract Access Context
// Returns a pointer so pointer-receiver methods (CanWrite, CanRead, etc.) work correctly.
func getAccessContextFromContext(c *gin.Context) (*middleware.AccessContext, bool) {
	accessContextRaw, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return nil, false
	}

	accessContext, ok := accessContextRaw.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access context"})
		return nil, false
	}

	return &accessContext, true
}

func (h *Handler) canAccessSeva(accessContext *middleware.AccessContext, sevaEntityID uint) bool {
	switch accessContext.RoleName {
	case "superadmin":
		return true

	case "templeadmin":
		// Fast path: direct entity match
		if accessContext.DirectEntityID != nil && *accessContext.DirectEntityID == sevaEntityID {
			return true
		}
		// Ownership check via TenantID
		if accessContext.TenantID > 0 && h.repo != nil {
			tenantID, err := h.repo.GetTenantIDByEntityID(sevaEntityID)
			if err == nil && tenantID == accessContext.TenantID {
				return true
			}
		}
		return false

	case "standarduser", "monitoringuser":
		// Direct entity match
		if accessContext.AssignedEntityID != nil && *accessContext.AssignedEntityID == sevaEntityID {
			return true
		}
		// Ownership check via TenantID
		if accessContext.TenantID > 0 && h.repo != nil {
			tenantID, err := h.repo.GetTenantIDByEntityID(sevaEntityID)
			if err == nil && tenantID == accessContext.TenantID {
				return true
			}
		}
		return false
	}
	return false
}

// resolveEntityID resolves entity ID from: query param → path param → access context
func resolveEntityID(c *gin.Context, accessContext *middleware.AccessContext) (uint, bool) {
	if entityIDParam := c.Query("entity_id"); entityIDParam != "" {
		if id, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
			return uint(id), true
		}
	}
	if entityIDPath := c.Param("entity_id"); entityIDPath != "" {
		if id, err := strconv.ParseUint(entityIDPath, 10, 32); err == nil {
			return uint(id), true
		}
	}
	if ctxID := accessContext.GetAccessibleEntityID(); ctxID != nil {
		return *ctxID, true
	}
	return 0, false
}

// ========================= REQUEST STRUCTS =============================

type CreateSevaRequest struct {
	EntityID       uint    `json:"entity_id"`
	Name           string  `json:"name" binding:"required"`
	SevaType       string  `json:"seva_type" binding:"required"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Date           string  `json:"date"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	Duration       int     `json:"duration"`
	AvailableSlots int     `json:"available_slots"`
}

type UpdateSevaRequest struct {
	EntityID       uint     `json:"entity_id"`
	Name           *string  `json:"name,omitempty"`
	SevaType       *string  `json:"seva_type,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Price          *float64 `json:"price,omitempty"`
	Date           *string  `json:"date,omitempty"`
	StartTime      *string  `json:"start_time,omitempty"`
	EndTime        *string  `json:"end_time,omitempty"`
	Duration       *int     `json:"duration,omitempty"`
	AvailableSlots *int     `json:"available_slots,omitempty"`
	Status         *string  `json:"status,omitempty"`
}

type BookSevaRequest struct {
	SevaID uint `json:"seva_id" binding:"required"`
}

type BookSevaWithPaymentRequest struct {
	SevaID   uint    `json:"seva_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	EntityID uint    `json:"entity_id"`
	SevaName string  `json:"seva_name"`
	SevaType string  `json:"seva_type"`
}

type VerifySevaPaymentRequest struct {
	RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
	RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
	RazorpaySignature string `json:"razorpay_signature" binding:"required"`
	SevaID            uint   `json:"seva_id" binding:"required"`
}

// ========================= SEVA HANDLERS =============================

func (h *Handler) CreateSeva(c *gin.Context) {
    accessContext, ok := getAccessContextFromContext(c)
    if !ok {
        return
    }

    var input CreateSevaRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    // Resolve entity ID:
    // Priority 1: URL path param
    // Priority 2: Request body entity_id (sent from frontend route.params.id) ← TRUST THIS
    // Priority 3: Access context (JWT) — LAST RESORT ONLY
    var entityID uint
    if entityIDParam := c.Param("entity_id"); entityIDParam != "" {
        if id, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
            entityID = uint(id)
        }
    }
    if entityID == 0 && input.EntityID != 0 {
        entityID = input.EntityID // ← use exactly what frontend sent (route.params.id = 1)
    }
    if entityID == 0 {
        if ctxID := accessContext.GetAccessibleEntityID(); ctxID != nil {
            entityID = *ctxID
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
            return
        }
    }

    // FIX: Check access AFTER entityID is resolved from body, NOT before.
    // Previously canAccessSeva was called with JWT entity (3), causing it to
    // override the correctly-passed body entity_id (1).
    if !isSuperAdmin(accessContext.RoleName) {
        if !h.canAccessSeva(accessContext, entityID) {
            c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this entity"})
            return
        }
    }

    if !accessContext.CanWrite() {
        c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
        return
    }

    ip := middleware.GetIPFromContext(c)

    seva := Seva{
        EntityID:       entityID,       // ← now correctly 1, not 3
        Name:           input.Name,
        SevaType:       input.SevaType,
        Description:    input.Description,
        Price:          input.Price,
        Date:           input.Date,
        StartTime:      input.StartTime,
        EndTime:        input.EndTime,
        Duration:       input.Duration,
        AvailableSlots: input.AvailableSlots,
        BookedSlots:    0,
        RemainingSlots: input.AvailableSlots,
        Status:         "upcoming",
    }

    if err := h.service.CreateSeva(c, &seva, *accessContext, ip); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seva: " + err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Seva created successfully", "seva": seva})
}

// 📄 List all sevas for temple admin with filters and pagination
func (h *Handler) ListEntitySevas(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	var entityID uint
	entityID, _ = resolveEntityID(c, accessContext)

	if entityID == 0 && (accessContext.RoleName == "standarduser" || accessContext.RoleName == "monitoringuser") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_id is required"})
		return
	}
	if entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple and no entity_id provided"})
		return
	}

	fmt.Println("entityID for ListEntitySevas:", entityID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1000"))
	offset := (page - 1) * limit

	sevaType := c.Query("seva_type")
	search := c.Query("search")
	status := c.Query("status")

	if !(accessContext.RoleName == "devotee" || accessContext.RoleName == "volunteer") && !accessContext.CanRead() {
		c.JSON(http.StatusForbidden, gin.H{"error": "read access denied"})
		return
	}

	sevas, total, err := h.service.GetSevasWithFilters(c, entityID, sevaType, search, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sevas: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sevas": sevas,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// 📊 Get Approved Booking Counts Per Seva
func (h *Handler) GetApprovedBookingCounts(c *gin.Context) {
	var entityID uint

	if entityIDParam := c.Query("entity_id"); entityIDParam != "" {
		if id, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
			entityID = uint(id)
		}
	}
	if entityID == 0 {
		if entityIDPath := c.Param("entity_id"); entityIDPath != "" {
			if id, err := strconv.ParseUint(entityIDPath, 10, 32); err == nil {
				entityID = uint(id)
			}
		}
	}
	if entityID == 0 {
		if user, exists := c.Get("user"); exists {
			if authUser, ok := user.(auth.User); ok && authUser.EntityID != nil {
				entityID = *authUser.EntityID
			}
		}
	}
	if entityID == 0 {
		if accessContext, ok := getAccessContextFromContext(c); ok {
			if ctxID := accessContext.GetAccessibleEntityID(); ctxID != nil {
				entityID = *ctxID
			}
		}
	}
	if entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_id is required"})
		return
	}

	counts, err := h.service.GetApprovedBookingCountsPerSeva(c, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch booking counts: " + err.Error()})
		return
	}

	type SevaCountResponse struct {
		SevaID        uint  `json:"seva_id"`
		ApprovedCount int64 `json:"approved_count"`
	}

	var response []SevaCountResponse
	for sevaID, count := range counts {
		response = append(response, SevaCountResponse{
			SevaID:        sevaID,
			ApprovedCount: count,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// 🔍 Get seva by ID
func (h *Handler) GetSevaByID(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seva ID"})
		return
	}

	seva, err := h.service.GetSevaByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seva not found"})
		return
	}

	if !isSuperAdmin(accessContext.RoleName) {
		if !h.canAccessSeva(accessContext, seva.EntityID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied to this seva"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"seva": seva})
}

// 🛠 Update seva
func (h *Handler) UpdateSeva(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seva ID"})
		return
	}

	var input UpdateSevaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	ip := middleware.GetIPFromContext(c)

	existingSeva, err := h.service.GetSevaByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seva not found"})
		return
	}

	if !isSuperAdmin(accessContext.RoleName) {
		if !h.canAccessSeva(accessContext, existingSeva.EntityID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: cannot update this seva"})
			return
		}
	}

	updatedSeva := *existingSeva
	if input.Name != nil {
		updatedSeva.Name = *input.Name
	}
	if input.SevaType != nil {
		updatedSeva.SevaType = *input.SevaType
	}
	if input.Description != nil {
		updatedSeva.Description = *input.Description
	}
	if input.Price != nil {
		updatedSeva.Price = *input.Price
	}
	if input.Date != nil {
		updatedSeva.Date = *input.Date
	}
	if input.StartTime != nil {
		updatedSeva.StartTime = *input.StartTime
	}
	if input.EndTime != nil {
		updatedSeva.EndTime = *input.EndTime
	}
	if input.Duration != nil {
		updatedSeva.Duration = *input.Duration
	}
	if input.AvailableSlots != nil {
		updatedSeva.AvailableSlots = *input.AvailableSlots
		updatedSeva.RemainingSlots = updatedSeva.AvailableSlots - updatedSeva.BookedSlots
		if updatedSeva.RemainingSlots < 0 {
			updatedSeva.RemainingSlots = 0
		}
	}
	if input.Status != nil {
		validStatuses := map[string]bool{"upcoming": true, "ongoing": true, "completed": true}
		if !validStatuses[*input.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be 'upcoming', 'ongoing', or 'completed'"})
			return
		}
		updatedSeva.Status = *input.Status
	}

	if err := h.service.UpdateSeva(c, &updatedSeva, *accessContext, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seva: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seva updated successfully", "seva": updatedSeva})
}

// ❌ Delete seva
func (h *Handler) DeleteSeva(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seva ID"})
		return
	}

	ip := middleware.GetIPFromContext(c)

	existingSeva, err := h.service.GetSevaByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seva not found"})
		return
	}

	if !isSuperAdmin(accessContext.RoleName) {
		if !h.canAccessSeva(accessContext, existingSeva.EntityID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: cannot delete this seva"})
			return
		}
	}

	if err := h.service.DeleteSeva(c, uint(id), *accessContext, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete seva: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seva deleted permanently"})
}

// 📋 Get Sevas for Devotees
func (h *Handler) GetSevas(c *gin.Context) {
	var entityID uint
	if entityIDParam := c.Query("entity_id"); entityIDParam != "" {
		if id, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
			entityID = uint(id)
		}
	} else {
		if entityIDPath := c.Param("entity_id"); entityIDPath != "" {
			if id, err := strconv.ParseUint(entityIDPath, 10, 32); err == nil {
				entityID = uint(id)
			}
		}
	}

	user := c.MustGet("user").(auth.User)
	if user.Role.RoleName != "devotee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	if entityID == 0 {
		if user.EntityID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "user not linked to a temple and no entity_id provided"})
			return
		}
		entityID = *user.EntityID
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	sevaType := c.Query("seva_type")
	search := c.Query("search")

	sevas, err := h.service.GetPaginatedSevas(c, entityID, sevaType, search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sevas: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sevas": sevas})
}

// ========================= BOOKING HANDLERS =============================

// 🎫 Book Seva
func (h *Handler) BookSeva(c *gin.Context) {
	var input BookSevaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user := c.MustGet("user").(auth.User)
	if user.Role.RoleName != "devotee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	seva, err := h.service.GetSevaByID(c, input.SevaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seva not found"})
		return
	}

	ip := middleware.GetIPFromContext(c)

	booking := SevaBooking{
		SevaID:      input.SevaID,
		UserID:      user.ID,
		EntityID:    seva.EntityID,
		BookingTime: time.Now(),
		Status:      "pending",
	}

	if err := h.service.BookSeva(c, &booking, "devotee", user.ID, seva.EntityID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Seva booked successfully",
		"booking": booking,
	})
}

// 💳 Book Seva with Razorpay Payment
func (h *Handler) BookSevaWithPayment(c *gin.Context) {
	var input BookSevaWithPaymentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user := c.MustGet("user").(auth.User)
	if user.Role.RoleName != "devotee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	seva, err := h.service.GetSevaByID(c, input.SevaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seva not found"})
		return
	}

	if seva.RemainingSlots <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No slots available for this seva"})
		return
	}

	// Fetch Razorpay keys from DB (per-temple, stored in tenant_bank_account_details)
	razorpayKey, razorpaySecret, err := h.repo.GetRazorpayKeysByEntityID(input.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment gateway not configured for this temple: " + err.Error()})
		return
	}

	client := razorpay.NewClient(razorpayKey, razorpaySecret)
	data := map[string]interface{}{
		"amount":   int(input.Amount * 100),
		"currency": "INR",
		"receipt":  fmt.Sprintf("seva_%d_%d", input.SevaID, time.Now().Unix()),
		"notes": map[string]interface{}{
			"seva_id":   input.SevaID,
			"user_id":   user.ID,
			"entity_id": input.EntityID,
			"seva_name": input.SevaName,
			"seva_type": input.SevaType,
		},
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment order: " + err.Error()})
		return
	}

	orderID := body["id"].(string)

	ip := middleware.GetIPFromContext(c)
	booking := SevaBooking{
		SevaID:          input.SevaID,
		UserID:          user.ID,
		EntityID:        input.EntityID,
		Amount:          input.Amount,
		BookingTime:     time.Now(),
		Status:          "pending",
		RazorpayOrderID: orderID,
	}

	if err := h.service.CreateSevaBookingWithPayment(c, &booking, user.ID, input.EntityID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":      true,
		"order_id":     orderID,
		"razorpay_key": razorpayKey,
		"amount":       input.Amount,
		"booking_id":   booking.ID,
		"message":      "Razorpay order created successfully",
	})
}

// ✅ Verify Seva Payment
func (h *Handler) VerifySevaPayment(c *gin.Context) {
	var input VerifySevaPaymentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user := c.MustGet("user").(auth.User)

	// Fetch booking first to get EntityID for per-temple Razorpay secret lookup
	booking, err := h.repo.GetBookingByOrderID(c, input.RazorpayOrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Booking not found for this order",
		})
		return
	}

	// Fetch Razorpay secret from DB using the booking's entity
	_, razorpaySecret, err := h.repo.GetRazorpayKeysByEntityID(booking.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Payment gateway not configured for this temple: " + err.Error(),
		})
		return
	}

	message := input.RazorpayOrderID + "|" + input.RazorpayPaymentID
	mac := hmac.New(sha256.New, []byte(razorpaySecret))
	mac.Write([]byte(message))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	if expectedSignature != input.RazorpaySignature {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid payment signature",
		})
		return
	}

	ip := middleware.GetIPFromContext(c)

	if err := h.service.VerifySevaPayment(
		c,
		input.RazorpayOrderID,
		input.RazorpayPaymentID,
		input.RazorpaySignature,
		input.SevaID,
		user.ID,
		ip,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Payment verification failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment verified successfully",
		"status":  "approved",
	})
}

func (h *Handler) Write(b []byte) {
	panic("unimplemented")
}

func (h *Handler) GetMyBookings(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	bookings, err := h.service.GetBookingsForUser(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

// 📊 Get Entity Bookings
func (h *Handler) GetEntityBookings(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	var entityID uint
	if entityIDParam := c.Query("entity_id"); entityIDParam != "" {
		if id, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
			entityID = uint(id)
		}
	} else {
		if entityIDPath := c.Param("entity_id"); entityIDPath != "" {
			if id, err := strconv.ParseUint(entityIDPath, 10, 32); err == nil {
				entityID = uint(id)
			}
		}
	}
	if entityID == 0 {
		if ctxID := accessContext.GetAccessibleEntityID(); ctxID != nil {
			entityID = *ctxID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple and no entity_id provided"})
			return
		}
	}

	status := c.Query("status")
	sevaType := c.Query("seva_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	bookings, err := h.service.GetDetailedBookingsWithFilters(
		c, entityID, status, sevaType, startDate, endDate, search, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch detailed bookings: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

func (h *Handler) GetBookingByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := h.service.GetBookingByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"booking": booking})
}

// 📊 Get Booking Counts
func (h *Handler) GetBookingCounts(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	var entityID uint

	if user.Role.RoleName == "devotee" {
		if entityIDParam := c.Query("entity_id"); entityIDParam != "" {
			if id, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
				entityID = uint(id)
			}
		} else if user.EntityID != nil {
			entityID = *user.EntityID
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}
	} else {
		accessContext, ok := getAccessContextFromContext(c)
		if !ok {
			return
		}
		if ctxID := accessContext.GetAccessibleEntityID(); ctxID != nil {
			entityID = *ctxID
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "No accessible entity"})
			return
		}
	}

	counts, err := h.service.GetBookingStatusCounts(c, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch counts: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"counts": counts})
}

// ✏️ UpdateBookingStatus
func (h *Handler) UpdateBookingStatus(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var input struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status field"})
		return
	}

	ip := middleware.GetIPFromContext(c)

	if err := h.service.UpdateBookingStatus(c, uint(id), input.Status, user.ID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Status update failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated successfully"})
}