package event

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

// CreateEvent handles POST /events
func (h *Handler) CreateEvent(c *gin.Context) {
	userRaw, userOk := c.Get("user")
	if !userOk {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - user not found in context"})
		return
	}

	userData, ok := userRaw.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user token structure"})
		return
	}

	userIDRaw, ok1 := userData["user_id"]
	entityIDRaw, ok2 := userData["tenant_id"]

	if !ok1 || !ok2 || userIDRaw == nil || entityIDRaw == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or tenant_id missing in token"})
		return
	}

	userID := uint(userIDRaw.(float64))
	entityID := uint(entityIDRaw.(float64))

	var req Event
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// PATCH: Validate EventDate is not empty/default
	if req.EventDate.IsZero() || req.EventDate.Before(time.Now().AddDate(-10, 0, 0)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event date is required and must be valid"})
		return
	}

	req.CreatedBy = userID
	req.EntityID = entityID

	if err := h.Service.CreateEvent(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully","event_id": req.ID})
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

	entityIDRaw, exists := userMap["tenant_id"]
	if !exists || entityIDRaw == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id missing in token"})
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

	entityIDRaw, exists := userMap["tenant_id"]
	if !exists || entityIDRaw == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id missing in token"})
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
