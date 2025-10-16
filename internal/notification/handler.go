package notification

import (
	"net/http"
	"strconv"
	"strings"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/internal/auditlog"
	"github.com/sharath018/temple-management-backend/middleware"
	"github.com/sharath018/temple-management-backend/utils"
)

// Handler wraps the service
type Handler struct {
	Service  Service
	AuditSvc auditlog.Service
}

func NewHandler(s Service, auditSvc auditlog.Service) *Handler {
	return &Handler{
		Service:  s,
		AuditSvc: auditSvc,
	}
}

// POST /api/v1/notifications/templates
func (h *Handler) CreateTemplate(c *gin.Context) {
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	if !ctx.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	ip := middleware.GetIPFromContext(c)

	var input NotificationTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = ctx.UserID
	input.EntityID = *entityID

	if err := h.Service.CreateTemplate(c.Request.Context(), &input, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create template"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GET /api/v1/notifications/templates
func (h *Handler) GetTemplates(c *gin.Context) {
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
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
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
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
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	if !ctx.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	ip := middleware.GetIPFromContext(c)

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

	if err := h.Service.UpdateTemplate(c.Request.Context(), &input, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update template"})
		return
	}

	c.JSON(http.StatusOK, input)
}

// DELETE /api/v1/notifications/templates/:id
func (h *Handler) DeleteTemplate(c *gin.Context) {
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	if !ctx.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	ip := middleware.GetIPFromContext(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	if err := h.Service.DeleteTemplate(c.Request.Context(), uint(id), *entityID, ctx.UserID, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template deleted"})
}

// POST /api/v1/notifications/send
func (h *Handler) SendNotification(c *gin.Context) {
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}

	ctx := accessContext.(middleware.AccessContext)
	
	if !ctx.CanWrite() {
		c.JSON(http.StatusForbidden, gin.H{"error": "write access denied"})
		return
	}

	entityID := ctx.GetAccessibleEntityID()
	if entityID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "no accessible temple"})
		return
	}

	ip := middleware.GetIPFromContext(c)

	var req struct {
		TemplateID *uint    `json:"template_id"`
		Channel    string   `json:"channel" binding:"required"`
		Subject    string   `json:"subject"`
		Body       string   `json:"body" binding:"required"`
		Recipients []string `json:"recipients"`
		Audience   string   `json:"audience"`
	}

	log.Printf("📦 Raw request body for SendNotification")
	
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("📋 Parsed request: Channel=%s, Subject=%s, Body=%s, Recipients=%v, Audience=%s",
		req.Channel, req.Subject, req.Body, req.Recipients, req.Audience)

	// ✅ ENHANCED: Resolve recipients with better error messages
	if len(req.Recipients) == 0 {
		log.Printf("🔍 No recipients provided, resolving via audience: %s", req.Audience)
		
		if req.Audience == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "either recipients or audience must be provided",
			})
			return
		}

		switch req.Audience {
		case "all", "devotees", "volunteers":
			emails, err := h.Service.GetEmailsByAudience(*entityID, req.Audience)
			if err != nil {
				log.Printf("❌ Failed to fetch emails for '%s': %v", req.Audience, err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to fetch users for audience",
					"details": err.Error(),
				})
				return
			}
			
			// ✅ CRITICAL: Check if no recipients found
			if len(emails) == 0 {
				log.Printf("⚠️ No recipients found for audience '%s' in entity %d", req.Audience, *entityID)
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "no_recipients",
					"message": "No recipients found for the selected audience",
					"details": "Please add devotees or volunteers to this temple before sending notifications",
					"audience": req.Audience,
					"entity_id": *entityID,
				})
				return
			}
			
			log.Printf("✅ Resolved %d emails for '%s'", len(emails), req.Audience)
			req.Recipients = emails

		default:
			log.Printf("❌ Invalid audience: %s", req.Audience)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid audience",
				"message": "Audience must be one of: all, devotees, volunteers",
			})
			return
		}
	}

	log.Printf("📤 Calling service.SendNotification with %d recipients", len(req.Recipients))

	// ✅ ENHANCED: Better error handling from service
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
		log.Printf("❌ Service.SendNotification failed: %v", err)
		
		// ✅ Parse error message for better user feedback
		errMsg := err.Error()
		statusCode := http.StatusInternalServerError
		response := gin.H{"error": "failed to send notification"}
		
		// Check for specific error types
		if strings.Contains(errMsg, "no recipients") {
			statusCode = http.StatusBadRequest
			response = gin.H{
				"error": "no_recipients",
				"message": "No recipients found for notification",
				"details": errMsg,
			}
		} else if strings.Contains(errMsg, "not configured") || strings.Contains(errMsg, "SMTP") {
			statusCode = http.StatusServiceUnavailable
			response = gin.H{
				"error": "service_not_configured",
				"message": "Email service is not properly configured",
				"details": errMsg,
			}
		} else if strings.Contains(errMsg, "template") {
			statusCode = http.StatusBadRequest
			response = gin.H{
				"error": "template_error",
				"message": "Template error",
				"details": errMsg,
			}
		} else {
			// Generic error with details
			response["details"] = errMsg
		}
		
		c.JSON(statusCode, response)
		return
	}

	log.Printf("✅ Notification queued successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "notification queued successfully",
		"recipients_count": len(req.Recipients),
		"channel": req.Channel,
		"status": "processing",
	})
}

// GET /api/v1/notifications/logs
func (h *Handler) GetMyNotifications(c *gin.Context) {
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

// GET /api/v1/notifications/inapp
func (h *Handler) GetMyInApp(c *gin.Context) {
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}
	ctx := accessContext.(middleware.AccessContext)
	var entityIDPtr *uint
	if id := ctx.GetAccessibleEntityID(); id != nil {
		entityIDPtr = id
	}
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	items, err := h.Service.ListInAppByUser(c.Request.Context(), ctx.UserID, entityIDPtr, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch in-app notifications"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// PUT /api/v1/notifications/inapp/:id/read
func (h *Handler) MarkInAppRead(c *gin.Context) {
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}
	ctx := accessContext.(middleware.AccessContext)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.Service.MarkInAppAsRead(c.Request.Context(), uint(id), ctx.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

// GET /api/v1/notifications/stream (SSE)
func (h *Handler) StreamInApp(c *gin.Context) {
	accessContext, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
		return
	}
	ctx := accessContext.(middleware.AccessContext)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	channel := "notifications:user:" + strconv.FormatUint(uint64(ctx.UserID), 10)
	sub := utils.RedisClient.Subscribe(c, channel)
	defer sub.Close()

	_, _ = c.Writer.Write([]byte(":ok\n\n"))
	flusher.Flush()

	ch := sub.Channel()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			_, _ = c.Writer.Write([]byte("event: inapp\n"))
			_, _ = c.Writer.Write([]byte("data: " + msg.Payload + "\n\n"))
			flusher.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}

// GET /api/v1/notifications/stream-token?token=JWT (SSE without auth middleware)
func (h *Handler) StreamInAppWithToken(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	cfg := config.Load()
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTAccessSecret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}
	uidFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id missing"})
		return
	}
	userID := uint(uidFloat)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	channel := "notifications:user:" + strconv.FormatUint(uint64(userID), 10)
	sub := utils.RedisClient.Subscribe(c, channel)
	defer sub.Close()

	_, _ = c.Writer.Write([]byte(":ok\n\n"))
	flusher.Flush()

	ch := sub.Channel()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			_, _ = c.Writer.Write([]byte("event: inapp\n"))
			_, _ = c.Writer.Write([]byte("data: " + msg.Payload + "\n\n"))
			flusher.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}