package seva

import (
	"net/http"
	"strconv"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/middleware"
)

type Handler struct {
	service  Service
	auditSvc auditlog.Service
}

func NewHandler(service Service, auditSvc auditlog.Service) *Handler {
	return &Handler{
		service:  service,
		auditSvc: auditSvc,
	}
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

// ========================= REQUEST STRUCTS =============================

type CreateSevaRequest struct {
	EntityID          uint      `json:"entity_id" binding:"required"`
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

type UpdateSevaRequest struct {
	Name              *string   `json:"name,omitempty"`
	SevaType          *string   `json:"seva_type,omitempty"`
	Description       *string   `json:"description,omitempty"`
	Price             *float64  `json:"price,omitempty"`
	Date              *string   `json:"date,omitempty"`
	StartTime         *string   `json:"start_time,omitempty"`
	EndTime           *string   `json:"end_time,omitempty"`
	Duration          *int      `json:"duration,omitempty"`
	MaxBookingsPerDay *int      `json:"max_bookings_per_day,omitempty"`
	Status            *string   `json:"status,omitempty"`
}

type BookSevaRequest struct {
	SevaID uint `json:"seva_id" binding:"required"`
}

// ========================= SEVA HANDLERS =============================

// 🎯 Create Seva - POST /sevas
func (h *Handler) CreateSeva(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Check write permissions
	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	var input CreateSevaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Get IP address for audit logging
	ip := middleware.GetIPFromContext(c)

	seva := Seva{
		EntityID:          input.EntityID,
		Name:              input.Name,
		SevaType:          input.SevaType,
		Description:       input.Description,
		Price:             input.Price,
		Date:              input.Date,
		StartTime:         input.StartTime,
		EndTime:           input.EndTime,
		Duration:          input.Duration,
		MaxBookingsPerDay: input.MaxBookingsPerDay,
		Status:            "upcoming",
	}

	if err := h.service.CreateSeva(c, &seva, accessContext, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seva: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Seva created successfully", "seva": seva})
}

// 📄 List all sevas - GET /sevas
func (h *Handler) ListEntitySevas(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Get entity ID from query parameter or URL path
	var entityID uint
	entityIDParam := c.Query("entity_id")
	if entityIDParam != "" {
		id, err := strconv.ParseUint(entityIDParam, 10, 32)
		if err == nil {
			entityID = uint(id)
		}
	} else {
		entityIDPath := c.Param("entity_id")
		if entityIDPath != "" {
			id, err := strconv.ParseUint(entityIDPath, 10, 32)
			if err == nil {
				entityID = uint(id)
			}
		}
	}

	// If no entity ID found, fall back to access context
	if entityID == 0 {
		contextEntityID := accessContext.GetAccessibleEntityID()
		if contextEntityID != nil {
			entityID = *contextEntityID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple and no entity_id provided"})
			return
		}
	}

	fmt.Println("entityID for ListEntitySevas:", entityID)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	sevaType := c.Query("seva_type")
	search := c.Query("search")
	status := c.Query("status")

	// Check read permissions
	if !(accessContext.RoleName == "devotee" || accessContext.RoleName == "volunteer") && !accessContext.CanRead() {
		c.JSON(http.StatusForbidden, gin.H{"error": "read access denied"})
		return
	}

	sevas, total, err := h.service.GetSevasWithFilters(
		c,
		entityID,
		sevaType,
		search,
		status,
		limit,
		offset,
	)
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

// 🔍 Get seva by ID - GET /sevas/:id
func (h *Handler) GetSevaByID(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple"})
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

	if seva.EntityID != *entityID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this seva"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"seva": seva})
}

// 🛠 Update seva - PUT /sevas/:id
func (h *Handler) UpdateSeva(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not linked to a temple"})
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

	if existingSeva.EntityID != *entityID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: cannot update this seva"})
		return
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
	if input.MaxBookingsPerDay != nil {
		updatedSeva.MaxBookingsPerDay = *input.MaxBookingsPerDay
	}
	if input.Status != nil {
		validStatuses := map[string]bool{"upcoming": true, "ongoing": true, "completed": true}
		if !validStatuses[*input.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be 'upcoming', 'ongoing', or 'completed'"})
			return
		}
		updatedSeva.Status = *input.Status
	}

	if err := h.service.UpdateSeva(c, &updatedSeva, accessContext, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seva: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seva updated successfully", "seva": updatedSeva})
}

// ❌ Delete seva - DELETE /sevas/:id
func (h *Handler) DeleteSeva(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not linked to a temple"})
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

	if existingSeva.EntityID != *entityID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: cannot delete this seva"})
		return
	}

	if err := h.service.DeleteSeva(c, uint(id), accessContext, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete seva: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seva deleted permanently"})
}

// 📋 Get Sevas for Devotees - GET /sevas/public
func (h *Handler) GetSevas(c *gin.Context) {
	var entityID uint
	entityIDParam := c.Query("entity_id")
	if entityIDParam != "" {
		id, err := strconv.ParseUint(entityIDParam, 10, 32)
		if err == nil {
			entityID = uint(id)
		}
	} else {
		entityIDPath := c.Param("entity_id")
		if entityIDPath != "" {
			id, err := strconv.ParseUint(entityIDPath, 10, 32)
			if err == nil {
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

	sevas, err := h.service.GetPaginatedSevas(
		c,
		entityID,
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

	ip := middleware.GetIPFromContext(c)

	booking := SevaBooking{
		SevaID:      input.SevaID,
		UserID:      user.ID,
		EntityID:    *user.EntityID,
		BookingTime: time.Now(),
		Status:      "pending",
	}

	if err := h.service.BookSeva(c, &booking, "devotee", user.ID, *user.EntityID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Booking failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Seva booked successfully",
		"booking": booking,
	})
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

// 📊 Get Entity Bookings - GET /bookings
func (h *Handler) GetEntityBookings(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	var entityID uint
	entityIDParam := c.Query("entity_id")
	if entityIDParam != "" {
		id, err := strconv.ParseUint(entityIDParam, 10, 32)
		if err == nil {
			entityID = uint(id)
		}
	} else {
		entityIDPath := c.Param("entity_id")
		if entityIDPath != "" {
			id, err := strconv.ParseUint(entityIDPath, 10, 32)
			if err == nil {
				entityID = uint(id)
			}
		}
	}

	if entityID == 0 {
		contextEntityID := accessContext.GetAccessibleEntityID()
		if contextEntityID != nil {
			entityID = *contextEntityID
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

// 📊 Get Booking Counts - GET /bookings/counts
func (h *Handler) GetBookingCounts(c *gin.Context) {
	user := c.MustGet("user").(auth.User)
	var entityID uint

	if user.Role.RoleName == "devotee" {
		entityIDParam := c.Query("entity_id")
		if entityIDParam != "" {
			id, err := strconv.ParseUint(entityIDParam, 10, 32)
			if err == nil {
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
		
		contextEntityID := accessContext.GetAccessibleEntityID()
		if contextEntityID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "No accessible entity"})
			return
		}
		entityID = *contextEntityID
	}

	counts, err := h.service.GetBookingStatusCounts(c, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch counts: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"counts": counts})
}

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