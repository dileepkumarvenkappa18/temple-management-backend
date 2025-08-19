package event

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/middleware"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

// ===========================
// 📌 Extract Authenticated User
func getUserFromContext(c *gin.Context) (*auth.User, bool) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return nil, false
	}
	user, ok := userRaw.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user"})
		return nil, false
	}
	return &user, true
}

// ===========================
// 🎯 Create Event - POST /events
func (h *Handler) CreateEvent(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok || user.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple"})
		return
	}

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// Basic validation for date being too far in the past
	if req.EventDate != "" {
		if eventDate, err := time.Parse("2006-01-02", req.EventDate); err == nil {
			if eventDate.Before(time.Now().AddDate(-10, 0, 0)) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "event date is too far in the past"})
				return
			}
		}
	}

	// Get IP address for audit logging
	ip := middleware.GetIPFromContext(c)

	// Use the service method that includes validation and audit logging
	if err := h.Service.CreateEvent(&req, user.ID, *user.EntityID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created successfully"})
}

// ===========================
// 🔍 Get Event - GET /events/:id
func (h *Handler) GetEventByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	event, err := h.Service.GetEventByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// ===========================
// 📆 Upcoming Events - GET /events/upcoming
func (h *Handler) GetUpcomingEvents(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok || user.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple"})
		return
	}

	events, err := h.Service.GetUpcomingEvents(*user.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// ===========================
// 📄 List Events - GET /events?limit=&offset=&search=
func (h *Handler) ListEvents(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok || user.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	search := c.Query("search")

	events, err := h.Service.ListEventsByEntity(*user.EntityID, limit, offset, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// ===========================
// 📊 Event Stats - GET /events/stats
func (h *Handler) GetEventStats(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok || user.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple"})
		return
	}

	stats, err := h.Service.GetEventStats(*user.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ===========================
// 🛠 Update Event - PUT /events/:id
func (h *Handler) UpdateEvent(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok || user.EntityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// Get IP address for audit logging
	ip := middleware.GetIPFromContext(c)

	// Use the updated service method that includes audit logging
	if err := h.Service.UpdateEvent(uint(id), &req, *user.EntityID, user.ID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update event: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event updated successfully"})
}

// ===========================
// ❌ Delete Event - DELETE /events/:id
func (h *Handler) DeleteEvent(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok || user.EntityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	// Get IP address for audit logging
	ip := middleware.GetIPFromContext(c)

	// Use the updated service method that includes audit logging
	if err := h.Service.DeleteEvent(uint(id), *user.EntityID, user.ID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete event: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}