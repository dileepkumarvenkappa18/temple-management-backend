package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// Handler wraps the service
type Handler struct {
	Service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

// POST /api/v1/notifications/templates
func (h *Handler) CreateTemplate(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	var input NotificationTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = user.ID
	input.EntityID = *user.EntityID

	if err := h.Service.CreateTemplate(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create template"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GET /api/v1/notifications/templates
func (h *Handler) GetTemplates(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	templates, err := h.Service.GetTemplates(c.Request.Context(), *user.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GET /api/v1/notifications/templates/:id
func (h *Handler) GetTemplateByID(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	template, err := h.Service.GetTemplateByID(c.Request.Context(), uint(id), *user.EntityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// PUT /api/v1/notifications/templates/:id
func (h *Handler) UpdateTemplate(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	var input NotificationTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = uint(id)
	input.UserID = user.ID
	input.EntityID = *user.EntityID

	if err := h.Service.UpdateTemplate(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update template"})
		return
	}

	c.JSON(http.StatusOK, input)
}

// DELETE /api/v1/notifications/templates/:id
func (h *Handler) DeleteTemplate(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	if err := h.Service.DeleteTemplate(c.Request.Context(), uint(id), *user.EntityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template deleted"})
}

// POST /api/v1/notifications/send
func (h *Handler) SendNotification(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	var req struct {
		TemplateID *uint    `json:"template_id"`                 // Optional
		Channel    string   `json:"channel" binding:"required"` // email, sms, whatsapp
		Subject    string   `json:"subject"`                     // only for email
		Body       string   `json:"body" binding:"required"`
		Recipients []string `json:"recipients"`                  // Optional now
		Audience   string   `json:"audience"`                    // all, devotees, volunteers
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entityID := *user.EntityID

	// If recipients not provided, resolve using audience
	if len(req.Recipients) == 0 {
		switch req.Audience {
		case "all":
			emails, err := h.Service.GetEmailsByAudience(entityID, "all")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users for audience"})
				return
			}
			req.Recipients = emails

		case "devotees", "volunteers":
			emails, err := h.Service.GetEmailsByAudience(entityID, req.Audience)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users for audience"})
				return
			}
			req.Recipients = emails

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing audience if no recipients are provided"})
			return
		}
	}

	if err := h.Service.SendNotification(
		c.Request.Context(),
		user.ID,
		entityID,
		req.TemplateID,
		req.Channel,
		req.Subject,
		req.Body,
		req.Recipients,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification sent"})
}


// GET /api/v1/notifications/logs
func (h *Handler) GetMyNotifications(c *gin.Context) {
	user := c.MustGet("user").(auth.User)

	logs, err := h.Service.GetNotificationsByUser(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
