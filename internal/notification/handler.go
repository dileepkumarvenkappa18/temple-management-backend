package notification

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/middleware"
)

// Handler wraps the service
type Handler struct {
	Service  Service
	AuditSvc auditlog.Service // ✅ Audit service for IP extraction
}

// ✅ Updated constructor to accept audit service
func NewHandler(s Service, auditSvc auditlog.Service) *Handler {
	return &Handler{
		Service:  s,
		AuditSvc: auditSvc,
	}
}

// POST /api/v1/notifications/templates
func (h *Handler) CreateTemplate(c *gin.Context) {
	// Get access context instead of direct user access
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	// Check if user can write (templeadmin and standarduser only)
	if !ctx.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	// Get accessible entity ID
	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	ip := middleware.GetIPFromContext(c) // ✅ Extract IP for audit

	var input NotificationTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = ctx.UserID
	input.EntityID = *entityID

	// ✅ Pass IP to service for audit logging
	if err := h.Service.CreateTemplate(c.Request.Context(), &input, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create template"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GET /api/v1/notifications/templates
func (h *Handler) GetTemplates(c *gin.Context) {
	// Get access context
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	// Get accessible entity ID
	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	templates, err := h.Service.GetTemplates(c.Request.Context(), *entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// GET /api/v1/notifications/templates/:id
func (h *Handler) GetTemplateByID(c *gin.Context) {
	// Get access context
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	// Get accessible entity ID
	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	template, err := h.Service.GetTemplateByID(c.Request.Context(), uint(id), *entityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// PUT /api/v1/notifications/templates/:id
func (h *Handler) UpdateTemplate(c *gin.Context) {
	// Get access context
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	// Check if user can write (templeadmin and standarduser only)
	if !ctx.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	// Get accessible entity ID
	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	ip := middleware.GetIPFromContext(c) // ✅ Extract IP for audit

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
	input.UserID = ctx.UserID
	input.EntityID = *entityID

	// ✅ Pass IP to service for audit logging
	if err := h.Service.UpdateTemplate(c.Request.Context(), &input, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update template"})
		return
	}

	c.JSON(http.StatusOK, input)
}

// DELETE /api/v1/notifications/templates/:id
func (h *Handler) DeleteTemplate(c *gin.Context) {
	// Get access context
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	// Check if user can write (templeadmin and standarduser only)
	if !ctx.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	// Get accessible entity ID
	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	ip := middleware.GetIPFromContext(c) // ✅ Extract IP for audit

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	// ✅ Pass user ID and IP to service for audit logging
	if err := h.Service.DeleteTemplate(c.Request.Context(), uint(id), *entityID, ctx.UserID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template deleted"})
}

// POST /api/v1/notifications/send
func (h *Handler) SendNotification(c *gin.Context) {
	// Get access context
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	// Check if user can write (templeadmin and standarduser only)
	if !ctx.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	// Get accessible entity ID
	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	ip := middleware.GetIPFromContext(c) // ✅ Extract IP for audit

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

	// If recipients not provided, resolve using audience
	if len(req.Recipients) == 0 {
		switch req.Audience {
		case "all":
			emails, err := h.Service.GetEmailsByAudience(*entityID, "all")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users for audience"})
				return
			}
			req.Recipients = emails

		case "devotees", "volunteers":
			emails, err := h.Service.GetEmailsByAudience(*entityID, req.Audience)
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

	// ✅ Pass IP to service for audit logging
	if err := h.Service.SendNotification(
		c.Request.Context(),
		ctx.UserID,
		*entityID,
		req.TemplateID,
		req.Channel,
		req.Subject,
		req.Body,
		req.Recipients,
		ip,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification sent"})
}

// GET /api/v1/notifications/logs
func (h *Handler) GetMyNotifications(c *gin.Context) {
	// Get access context
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)

	logs, err := h.Service.GetNotificationsByUser(c.Request.Context(), ctx.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}