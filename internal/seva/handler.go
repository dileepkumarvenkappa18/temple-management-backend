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

// ========================= REQUEST STRUCTS =============================

type CreateSevaRequest struct {
	Name              string    `json:"name" binding:"required"`
	SevaType          string    `json:"seva_type" binding:"required"`
	Description       string    `json:"description"`
	Price             float64   `json:"price"`
	Date              string    `json:"date"`
	StartTime         string    `json:"start_time"`
	EndTime           string    `json:"end_time"`
	Duration          int       `json:"duration"`
	MaxBookingsPerDay int       `json:"max_bookings_per_day"`
}

type BookSevaRequest struct {
	SevaID uint `json:"seva_id" binding:"required"`
}

// type BookSevaRequest struct {
// 	SevaID          uint    `json:"seva_id" binding:"required"`
// 	BookingDate     string  `json:"booking_date" binding:"required"` // format: YYYY-MM-DD
// 	BookingTime     string  `json:"booking_time" binding:"required"` // format: HH:MM
// 	SpecialRequests string  `json:"special_requests"`
// 	AmountPaid      float64 `json:"amount_paid"`
// 	PaymentStatus   string  `json:"payment_status"`
// }

// ========================= SEVA HANDLERS =============================

func (h *Handler) CreateSeva(c *gin.Context) {
	var input CreateSevaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user := c.MustGet("user").(auth.User)
	if user.Role.RoleName != "templeadmin" || user.EntityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized or invalid entity"})
		return
	}

	seva := Seva{
		EntityID:          *user.EntityID,
		Name:              input.Name,
		SevaType:          input.SevaType,
		Description:       input.Description,
		Price:             input.Price,
		Date:              input.Date,
		StartTime:         input.StartTime,
		EndTime:           input.EndTime,
		Duration:          input.Duration,
		MaxBookingsPerDay: input.MaxBookingsPerDay,
		Status:            "upcoming",        // ✅ defaulted
		IsActive:          true,              // ✅ defaulted
	}

	if err := h.service.CreateSeva(c, &seva, user.Role.RoleName, *user.EntityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seva: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Seva created successfully", "seva": seva})
}


func (h *Handler) GetSevas(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	if user.Role.RoleName != "devotee" || user.EntityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	sevaType := c.Query("seva_type")
	search := c.Query("search")

	sevas, err := h.service.GetPaginatedSevas(
		c,
		*user.EntityID,
		sevaType,
		search,
		limit,
		offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sevas: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sevas": sevas})
}


// ========================= BOOKING HANDLERS =============================



func (h *Handler) BookSeva(c *gin.Context) {
	var input BookSevaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	user := c.MustGet("user").(auth.User)
	if user.Role.RoleName != "devotee" || user.EntityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized or invalid entity"})
		return
	}

	booking := SevaBooking{
		SevaID:      input.SevaID,
		UserID:      user.ID,
		EntityID:    *user.EntityID,
		BookingTime: time.Now(),       // ⏱️ Auto-generated now
		Status:      "pending",        // default state
	}

	if err := h.service.BookSeva(c, &booking, "devotee", user.ID, *user.EntityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Seva booked successfully",
		"booking": booking,
	})
}


// func (h *Handler) BookSeva(c *gin.Context) {
// 	var input BookSevaRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
// 		return
// 	}

// 	user := c.MustGet("user").(auth.User)
// 	if user.Role.RoleName != "devotee" || user.EntityID == nil {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized or invalid entity"})
// 		return
// 	}

// 	date, err := time.Parse("2006-01-02", input.BookingDate)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking_date format"})
// 		return
// 	}
// 	timeSlot, err := time.Parse("15:04", input.BookingTime)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking_time format"})
// 		return
// 	}

// 	booking := SevaBooking{
// 		SevaID:          input.SevaID,
// 		UserID:          user.ID,
// 		EntityID:        *user.EntityID,
// 		BookingDate:     date,
// 		BookingTime:     timeSlot,
// 		SpecialRequests: input.SpecialRequests,
// 		AmountPaid:      input.AmountPaid,
// 		PaymentStatus:   input.PaymentStatus,
// 		Status:          "pending",
// 	}

// 	if err := h.service.BookSeva(c, &booking, "devotee", user.ID, *user.EntityID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking failed: " + err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Seva booked successfully", "booking": booking})
// }

func (h *Handler) GetMyBookings(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	bookings, err := h.service.GetBookingsForUser(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch bookings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

func (h *Handler) GetEntityBookings(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	if user.Role.RoleName != "templeadmin" || user.EntityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	status := c.Query("status")
	sevaType := c.Query("seva_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	bookings, err := h.service.GetDetailedBookingsWithFilters(
		c, *user.EntityID, status, sevaType, startDate, endDate, search, limit, offset,
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

func (h *Handler) GetBookingCounts(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	if user.Role.RoleName != "templeadmin" || user.EntityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	counts, err := h.service.GetBookingStatusCounts(c, *user.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch counts: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"counts": counts})
}

func (h *Handler) UpdateBookingStatus(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	if user.Role.RoleName != "templeadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

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

	if err := h.service.UpdateBookingStatus(c, uint(id), input.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Status update failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated successfully"})
}
