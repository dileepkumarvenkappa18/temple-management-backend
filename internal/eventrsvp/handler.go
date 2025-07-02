package eventrsvp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/internal/event"
)

type Handler struct {
	Service      *Service
	EventService *event.Service
}

func NewHandler(service *Service, eventService *event.Service) *Handler {
	return &Handler{
		Service:      service,
		EventService: eventService,
	}
}

// POST /event-rsvps/:eventID
func (h *Handler) CreateRSVP(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user, ok := userData.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user from context"})
		return
	}

	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil || eventID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// 🔍 Check if event exists
	_, err = h.Service.EventService.GetEventByID(uint(eventID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// 👇 Try update first
	err = h.Service.UpdateRSVPStatus(uint(eventID), user.ID, req.Status, req.Notes)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "RSVP updated successfully"})
		return
	}

	// 👇 Fallback: create new
	rsvp := &RSVP{
		EventID: uint(eventID),
		UserID:  user.ID,
		Status:  req.Status,
		Notes:   req.Notes,
	}

	if err := h.Service.CreateRSVP(rsvp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to RSVP: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "RSVP submitted successfully"})
}

// GET /event-rsvps/my
func (h *Handler) GetMyRSVPs(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	user, ok := userData.(auth.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user from context"})
		return
	}

	rsvps, err := h.Service.GetMyRSVPs(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get RSVPs"})
		return
	}

	c.JSON(http.StatusOK, rsvps)
}

// GET /event-rsvps/:eventID
func (h *Handler) GetRSVPsByEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil || eventID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	rsvps, err := h.Service.GetRSVPsByEvent(uint(eventID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get RSVPs"})
		return
	}

	c.JSON(http.StatusOK, rsvps)
}
