package superadmin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =========================== TENANT APPROVAL ===========================

// GET /superadmin/tenants/pending
func (h *Handler) GetPendingTenants(c *gin.Context) {
	tenants, err := h.service.GetPendingTenants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending tenants"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pending_tenants": tenants})
}

// POST /superadmin/tenants/:id/approve
func (h *Handler) ApproveTenant(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	adminID := c.GetUint("userID") // Extracted from JWT
	if err := h.service.ApproveTenant(c.Request.Context(), uint(userID), adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve tenant"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tenant approved successfully"})
}

// POST /superadmin/tenants/:id/reject
func (h *Handler) RejectTenant(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rejection reason required"})
		return
	}

	adminID := c.GetUint("userID")
	if err := h.service.RejectTenant(c.Request.Context(), uint(userID), adminID, body.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject tenant"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tenant rejected successfully"})
}

// =========================== ENTITY APPROVAL ===========================

// GET /superadmin/entities/pending
func (h *Handler) GetPendingEntities(c *gin.Context) {
	entities, err := h.service.GetPendingEntities(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending entities"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pending_entities": entities})
}

// POST /superadmin/entities/:id/approve
func (h *Handler) ApproveEntity(c *gin.Context) {
	idStr := c.Param("id")
	entityID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	adminID := c.GetUint("userID")
	if err := h.service.ApproveEntity(c.Request.Context(), uint(entityID), adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve entity"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Entity approved successfully",
	})
}

// POST /superadmin/entities/:id/reject
func (h *Handler) RejectEntity(c *gin.Context) {
	idStr := c.Param("id")
	entityID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rejection reason required"})
		return
	}

	adminID := c.GetUint("userID")
	if err := h.service.RejectEntity(c.Request.Context(), uint(entityID), adminID, body.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject entity"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Entity rejected successfully",
	})
}
