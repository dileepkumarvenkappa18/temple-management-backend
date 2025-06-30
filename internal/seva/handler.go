package seva

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create a new seva (tenant only)
func (h *Handler) CreateSeva(c *gin.Context) {
	var seva Seva
	if err := c.ShouldBindJSON(&seva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(map[string]interface{})
	roleID := uint(user["role_id"].(float64))

	tenantVal, ok := user["tenant_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant_id missing in token"})
		return
	}
	tenantFloat, ok := tenantVal.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid tenant_id format"})
		return
	}
	entityID := uint(tenantFloat)

	var role string
	if roleID == 2 {
		role = "tenant"
	}

	err := h.service.CreateSeva(c, &seva, role, entityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, seva)
}

// Get all sevas for a temple
func (h *Handler) GetSevas(c *gin.Context) {
	entityIDParam := c.Query("entity_id")
	entityID, err := strconv.Atoi(entityIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
		return
	}

	sevas, err := h.service.GetSevasByEntity(c, uint(entityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch sevas"})
		return
	}
	c.JSON(http.StatusOK, sevas)
}

// ✅ Updated Book a seva (devotee only)
func (h *Handler) BookSeva(c *gin.Context) {
	type SevaBookingInput struct {
		SevaID          uint    `json:"seva_id"`
		BookingDate     string  `json:"booking_date"`     // "2025-07-02"
		BookingTime     string  `json:"booking_time"`     // "9:00"
		SpecialRequests string  `json:"special_requests"`
		AmountPaid      float64 `json:"amount_paid"`
		PaymentStatus   string  `json:"payment_status"`
	}

	var input SevaBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(map[string]interface{})
	userID := uint(user["user_id"].(float64))
	roleID := uint(user["role_id"].(float64))

	tenantVal, ok := user["tenant_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant_id missing in token"})
		return
	}
	tenantFloat, ok := tenantVal.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid tenant_id format"})
		return
	}
	entityID := uint(tenantFloat)

	// Parse date and time
	parsedDate, err := time.Parse("2006-01-02", input.BookingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking_date format, expected YYYY-MM-DD"})
		return
	}

	parsedTime, err := time.Parse("15:04", input.BookingTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking_time format, expected HH:MM in 24hr format"})
		return
	}

	role := ""
	if roleID == 3 {
		role = "devotee"
	}

	booking := SevaBooking{
		SevaID:          input.SevaID,
		UserID:          userID,
		EntityID:        entityID,
		BookingDate:     parsedDate,
		BookingTime:     parsedTime,
		SpecialRequests: input.SpecialRequests,
		AmountPaid:      input.AmountPaid,
		PaymentStatus:   input.PaymentStatus,
		Status:          "pending",
	}

	err = h.service.BookSeva(c, &booking, role, userID, entityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, booking)
}

// Get devotee's bookings
func (h *Handler) GetMyBookings(c *gin.Context) {
	user := c.MustGet("user").(map[string]interface{})
	userID := uint(user["user_id"].(float64))

	bookings, err := h.service.GetBookingsForUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// Get temple bookings (tenant only)
func (h *Handler) GetEntityBookings(c *gin.Context) {
	user := c.MustGet("user").(map[string]interface{})
	roleID := uint(user["role_id"].(float64))

	tenantVal, ok := user["tenant_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant_id missing in token"})
		return
	}
	tenantFloat, ok := tenantVal.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid tenant_id format"})
		return
	}
	entityID := uint(tenantFloat)

	if roleID != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "only tenant can view entity bookings"})
		return
	}

	bookings, err := h.service.GetBookingsForEntity(c, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// Cancel a seva booking (devotee)
func (h *Handler) CancelBooking(c *gin.Context) {
	bookingIDParam := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	user := c.MustGet("user").(map[string]interface{})
	userID := uint(user["user_id"].(float64))
	roleID := uint(user["role_id"].(float64))

	if roleID != 3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "only devotee can cancel bookings"})
		return
	}

	err = h.service.CancelBooking(c, uint(bookingID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled successfully"})
}

// Update booking status (temple admin)
func (h *Handler) UpdateBookingStatus(c *gin.Context) {
	bookingIDParam := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	var body struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing status"})
		return
	}

	user := c.MustGet("user").(map[string]interface{})
	roleID := uint(user["role_id"].(float64))

	if roleID != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "only temple admin can update status"})
		return
	}

	err = h.service.UpdateBookingStatus(c, uint(bookingID), body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking status updated"})
}
