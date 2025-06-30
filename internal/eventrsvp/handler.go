package eventrsvp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	userData, _ := c.Get("user")
	userMap := userData.(map[string]interface{})
	userID := uint(userMap["user_id"].(float64))

	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil || eventID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// 🔍 Check if event exists (via eventService)
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
	err = h.Service.UpdateRSVPStatus(uint(eventID), userID, req.Status, req.Notes)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "RSVP updated successfully"})
		return
	}

	// 👇 Fallback: create new
	rsvp := &RSVP{
		EventID: uint(eventID),
		UserID:  userID,
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
	userData, _ := c.Get("user")
	userMap := userData.(map[string]interface{})
	userID := userMap["user_id"].(float64)

	rsvps, err := h.Service.GetMyRSVPs(uint(userID))
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
