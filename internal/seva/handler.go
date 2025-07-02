package seva

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// Create a new seva (templeadmin only)
func (h *Handler) CreateSeva(c *gin.Context) {
	var seva Seva
	if err := c.ShouldBindJSON(&seva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(auth.User)

	if user.Role.RoleName != "templeadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only temple admin can create sevas"})
		return
	}

	if user.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

	err := h.service.CreateSeva(c, &seva, "templeadmin", *user.EntityID)
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

// Book a seva (devotee only)
func (h *Handler) BookSeva(c *gin.Context) {
	type SevaBookingInput struct {
		SevaID          uint    `json:"seva_id"`
		BookingDate     string  `json:"booking_date"` // YYYY-MM-DD
		BookingTime     string  `json:"booking_time"` // HH:MM
		SpecialRequests string  `json:"special_requests"`
		AmountPaid      float64 `json:"amount_paid"`
		PaymentStatus   string  `json:"payment_status"`
	}

	var input SevaBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(auth.User)

	if user.Role.RoleName != "devotee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only devotee can book seva"})
		return
	}

	if user.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to any entity"})
		return
	}

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

	booking := SevaBooking{
		SevaID:          input.SevaID,
		UserID:          user.ID,
		EntityID:        *user.EntityID,
		BookingDate:     parsedDate,
		BookingTime:     parsedTime,
		SpecialRequests: input.SpecialRequests,
		AmountPaid:      input.AmountPaid,
		PaymentStatus:   input.PaymentStatus,
		Status:          "pending",
	}

	err = h.service.BookSeva(c, &booking, "devotee", user.ID, *user.EntityID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, booking)
}

// Get devotee's bookings
func (h *Handler) GetMyBookings(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	bookings, err := h.service.GetBookingsForUser(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// Get all bookings for a temple (templeadmin only)
func (h *Handler) GetEntityBookings(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	if user.Role.RoleName != "templeadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only templeadmin can view entity bookings"})
		return
	}

	if user.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to entity"})
		return
	}

	bookings, err := h.service.GetBookingsForEntity(c, *user.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// Cancel a seva booking (devotee only)
func (h *Handler) CancelBooking(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	if user.Role.RoleName != "devotee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only devotee can cancel bookings"})
		return
	}

	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	err = h.service.CancelBooking(c, uint(bookingID), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled successfully"})
}

// Update booking status (templeadmin only)
func (h *Handler) UpdateBookingStatus(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	if user.Role.RoleName != "templeadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only templeadmin can update booking status"})
		return
	}

	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	var input struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing status"})
		return
	}

	err = h.service.UpdateBookingStatus(c, uint(bookingID), input.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking status updated"})
}
