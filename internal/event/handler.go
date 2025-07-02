package event

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

// CreateEvent handles POST /events
func (h *Handler) CreateEvent(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	user, ok := userRaw.(auth.User) // ✅ use actual struct
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
		return
	}

	if user.EntityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_id missing for user"})
		return
	}

	var req Event
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if req.EventDate.IsZero() || req.EventDate.Before(time.Now().AddDate(-10, 0, 0)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event date is required and must be valid"})
		return
	}

	req.CreatedBy = user.ID
	req.EntityID = *user.EntityID

	if err := h.Service.CreateEvent(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   req,
	})
}


// GetEventByID handles GET /events/:id
func (h *Handler) GetEventByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := h.Service.GetEventByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetUpcomingEvents handles GET /events/upcoming
func (h *Handler) GetUpcomingEvents(c *gin.Context) {
	entityRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	userMap, ok := entityRaw.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		return
	}

	entityIDRaw, exists := userMap["entity_id"]
	if !exists || entityIDRaw == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_id missing in token"})
		return
	}

	entityID := uint(entityIDRaw.(float64))

	events, err := h.Service.GetUpcomingEvents(entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// ListEvents handles GET /events
func (h *Handler) ListEvents(c *gin.Context) {
	entityRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	userMap, ok := entityRaw.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		return
	}

	entityIDRaw, exists := userMap["entity_id"]
	if !exists || entityIDRaw == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_id missing in token"})
		return
	}

	entityID := uint(entityIDRaw.(float64))

	events, err := h.Service.ListEventsByEntity(entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// UpdateEvent handles PUT /events/:id
func (h *Handler) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var req Event
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	req.ID = uint(id)

	if err := h.Service.UpdateEvent(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// DeleteEvent handles DELETE /events/:id
func (h *Handler) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := h.Service.DeleteEvent(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
