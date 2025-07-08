package superadmin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =========================== Helper ===========================

func parseUintParam(c *gin.Context, param string) (uint, bool) {
	idStr := c.Param(param)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid ID in path"})
		return 0, false
	}
	return uint(id), true
}

func parsePaginationParams(c *gin.Context) (int, int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	return page, limit
}

// =========================== TENANT APPROVAL ===========================

func (h *Handler) GetPendingTenants(c *gin.Context) {
	page, limit := parsePaginationParams(c)
	search := c.DefaultQuery("search", "")
	status := c.DefaultQuery("status", "pending")
	
	// Sanitize search
	search = strings.TrimSpace(search)

	tenants, total, err := h.service.GetPendingTenants(c.Request.Context(), page, limit, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not fetch tenants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tenants,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"search": search,
			"status": status,
		},
	})
}

func (h *Handler) ApproveTenant(c *gin.Context) {
	userID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	adminID := c.GetUint("userID")
	if err := h.service.ApproveTenant(c.Request.Context(), userID, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tenant approved"})
}

func (h *Handler) RejectTenant(c *gin.Context) {
	userID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reason is required for rejection"})
		return
	}

	adminID := c.GetUint("userID")
	err := h.service.RejectTenant(c.Request.Context(), userID, adminID, body.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tenant rejected"})
}

// =========================== ENTITY APPROVAL ===========================

func (h *Handler) GetPendingEntities(c *gin.Context) {
	page, limit := parsePaginationParams(c)
	search := c.DefaultQuery("search", "")
	status := c.DefaultQuery("status", "pending")
	search = strings.TrimSpace(search)

	entities, total, err := h.service.GetPendingEntities(c.Request.Context(), page, limit, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not fetch entities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    entities,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"search": search,
			"status": status,
		},
	})
}

func (h *Handler) ApproveEntity(c *gin.Context) {
	entityID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	adminID := c.GetUint("userID")
	if err := h.service.ApproveEntity(c.Request.Context(), entityID, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Entity approved"})
}

func (h *Handler) RejectEntity(c *gin.Context) {
	entityID, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Reason is required for rejection"})
		return
	}

	adminID := c.GetUint("userID")
	if err := h.service.RejectEntity(c.Request.Context(), entityID, adminID, body.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Entity rejected"})
}
