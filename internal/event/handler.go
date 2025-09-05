package event

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/middleware"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
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

// ===========================
// 🎯 Create Event - POST /events
func (h *Handler) CreateEvent(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Check if user has access to an entity
	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple"})
		return
	}

	// Check write permissions (handled by middleware, but double-check)
	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
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

	// Use the service method with access context
	if err := h.Service.CreateEvent(&req, accessContext, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created successfully"})
}

// ===========================
// 🔍 Get Event - GET /events/:id
func (h *Handler) GetEventByID(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Check if user has access to an entity
	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	event, err := h.Service.GetEventByID(uint(id), accessContext)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// ===========================
// 📆 Upcoming Events - GET /events/upcoming
// GetUpcomingEvents - GET /events/upcoming
// GetUpcomingEvents - GET /events/upcoming
func (h *Handler) GetUpcomingEvents(c *gin.Context) {
    accessContext, ok := getAccessContextFromContext(c)
    if !ok {
        return
    }

    // Get tenant ID from request header for cross-tenant access
    tenantIDHeader := c.GetHeader("X-Tenant-ID")
    var tenantID *uint
    if tenantIDHeader != "" {
        id, err := strconv.ParseUint(tenantIDHeader, 10, 32)
        if err == nil {
            uintID := uint(id)
            tenantID = &uintID
        }
    }

    // Use tenant ID from header if present, otherwise use entity ID from access context
    var targetEntityID *uint
    if tenantID != nil {
        targetEntityID = tenantID
    } else {
        targetEntityID = accessContext.GetAccessibleEntityID()
    }

    if targetEntityID == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple"})
        return
    }

    // Allow access for devotees and volunteers regardless of entity
    if accessContext.RoleName == "devotee" || accessContext.RoleName == "volunteer" || accessContext.CanRead() {
        events, err := h.Service.GetUpcomingEvents(accessContext)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
            return
        }

        c.JSON(http.StatusOK, events)
        return
    }

    // If not a devotee/volunteer and doesn't have read access
    c.JSON(http.StatusForbidden, gin.H{"error": "read access denied"})
}


// ===========================
// 📄 List Events - GET /events?limit=&offset=&search=
func (h *Handler) ListEvents(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Check if user has access to an entity
	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	search := c.Query("search")

	// Allow access for devotees and volunteers regardless of entity, similar to GetUpcomingEvents
	if accessContext.RoleName == "devotee" || accessContext.RoleName == "volunteer" || accessContext.CanRead() {
		events, err := h.Service.ListEventsByEntity(accessContext, limit, offset, search)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
			return
		}

		c.JSON(http.StatusOK, events)
		return
	}

	// If not a devotee/volunteer and doesn't have read access
	c.JSON(http.StatusForbidden, gin.H{"error": "read access denied"})
}

// ===========================
// 📊 Event Stats - GET /events/stats
func (h *Handler) GetEventStats(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Check if user has access to an entity
	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple"})
		return
	}

	stats, err := h.Service.GetEventStats(accessContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ===========================
// 🛠 Update Event - PUT /events/:id
func (h *Handler) UpdateEvent(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Check if user has access to an entity
	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not linked to a temple"})
		return
	}

	// Check write permissions (handled by middleware, but double-check)
	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
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

	// Use the updated service method with access context
	if err := h.Service.UpdateEvent(uint(id), &req, accessContext, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update event: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event updated successfully"})
}

// ===========================
// ❌ Delete Event - DELETE /events/:id
func (h *Handler) DeleteEvent(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Check if user has access to an entity
	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not linked to a temple"})
		return
	}

	// Check write permissions (handled by middleware, but double-check)
	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	// Get IP address for audit logging
	ip := middleware.GetIPFromContext(c)

	// Use the updated service method with access context
	if err := h.Service.DeleteEvent(uint(id), accessContext, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete event: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}