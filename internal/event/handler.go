package event

import (
	"fmt"
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

	// Priority 1: URL path param
	var entityID uint
	if entityIDParam := c.Param("entity_id"); entityIDParam != "" {
		if id, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
			entityID = uint(id)
		}
	}

	// Parse body first so entity_id in body is available as fallback
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// Priority 2: Request body entity_id
	if entityID == 0 && req.EntityID != 0 {
		entityID = req.EntityID
	}

	// Priority 3: Access context fallback
	if entityID == 0 {
		if contextEntityID := accessContext.GetAccessibleEntityID(); contextEntityID != nil {
			entityID = *contextEntityID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is not linked to a temple and no entity_id provided"})
			return
		}
	}

	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	if req.EventDate != "" {
		if eventDate, err := time.Parse("2006-01-02", req.EventDate); err == nil {
			if eventDate.Before(time.Now().AddDate(-10, 0, 0)) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "event date is too far in the past"})
				return
			}
		}
	}

	ip := middleware.GetIPFromContext(c)

	if err := h.Service.CreateEvent(&req, accessContext, entityID, ip); err != nil {
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
func (h *Handler) GetUpcomingEvents(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Priority 1: Query param
	var entityID uint
	if entityIDParam := c.Query("entity_id"); entityIDParam != "" {
		if id, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
			entityID = uint(id)
		}
	}

	// Priority 2: URL path param
	if entityID == 0 {
		if entityIDPath := c.Param("entity_id"); entityIDPath != "" {
			if id, err := strconv.ParseUint(entityIDPath, 10, 32); err == nil {
				entityID = uint(id)
			}
		}
	}

	// Priority 3: Access context fallback
	if entityID == 0 {
		if contextEntityID := accessContext.GetAccessibleEntityID(); contextEntityID != nil {
			entityID = *contextEntityID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple and no entity_id provided"})
			return
		}
	}

	if accessContext.RoleName == "devotee" || accessContext.RoleName == "volunteer" || accessContext.CanRead() {
		events, err := h.Service.GetUpcomingEventsByEntityID(accessContext, entityID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			return
		}
		c.JSON(http.StatusOK, events)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "read access denied"})
}

// ===========================
// 📄 List Events - GET /events
func (h *Handler) ListEvents(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	// Priority 1: Query param
	var entityID uint
	if entityIDParam := c.Query("entity_id"); entityIDParam != "" {
		if id, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
			entityID = uint(id)
		}
	}

	// Priority 2: URL path param
	if entityID == 0 {
		if entityIDPath := c.Param("entity_id"); entityIDPath != "" {
			if id, err := strconv.ParseUint(entityIDPath, 10, 32); err == nil {
				entityID = uint(id)
			}
		}
	}

	// Priority 3: Access context fallback
	if entityID == 0 {
		if contextEntityID := accessContext.GetAccessibleEntityID(); contextEntityID != nil {
			entityID = *contextEntityID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple and no entity_id provided"})
			return
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	search := c.Query("search")

	if accessContext.RoleName == "devotee" || accessContext.RoleName == "volunteer" || accessContext.CanRead() {
		events, err := h.Service.ListEventsByEntityID(accessContext, entityID, limit, offset, search)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
			return
		}
		c.JSON(http.StatusOK, events)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "read access denied"})
}

// ===========================
// 📊 Event Stats - GET /events/stats
func (h *Handler) GetEventStats(c *gin.Context) {
	accessContext, ok := getAccessContextFromContext(c)
	if !ok {
		return
	}

	entityID := accessContext.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not linked to a temple"})
		return
	}
	fmt.Println("entityID :=", *entityID)

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

	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	// Parse body first so entity_id in body is available as fallback
	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// Priority 1: URL path param
	var entityID uint
	if entityIDParam := c.Param("entity_id"); entityIDParam != "" {
		if eid, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
			entityID = uint(eid)
		}
	}

	// Priority 2: Request body entity_id
	if entityID == 0 && req.EntityID != 0 {
		entityID = req.EntityID
	}

	// Priority 3: Access context fallback
	if entityID == 0 {
		if contextEntityID := accessContext.GetAccessibleEntityID(); contextEntityID != nil {
			entityID = *contextEntityID
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not linked to a temple and no entity_id provided"})
			return
		}
	}

	ip := middleware.GetIPFromContext(c)

	// Pass resolved entityID explicitly to service
	if err := h.Service.UpdateEvent(uint(id), &req, accessContext, entityID, ip); err != nil {
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

	if !accessContext.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	// Priority 1: URL path param
	var entityID uint
	if entityIDParam := c.Param("entity_id"); entityIDParam != "" {
		if eid, err := strconv.ParseUint(entityIDParam, 10, 32); err == nil {
			entityID = uint(eid)
		}
	}

	// Priority 2: Query param (DELETE has no body)
	if entityID == 0 {
		if entityIDQuery := c.Query("entity_id"); entityIDQuery != "" {
			if eid, err := strconv.ParseUint(entityIDQuery, 10, 32); err == nil {
				entityID = uint(eid)
			}
		}
	}

	// Priority 3: Access context fallback
	if entityID == 0 {
		if contextEntityID := accessContext.GetAccessibleEntityID(); contextEntityID != nil {
			entityID = *contextEntityID
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not linked to a temple and no entity_id provided"})
			return
		}
	}

	ip := middleware.GetIPFromContext(c)

	// Pass resolved entityID explicitly to service
	if err := h.Service.DeleteEvent(uint(id), accessContext, entityID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete event: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}