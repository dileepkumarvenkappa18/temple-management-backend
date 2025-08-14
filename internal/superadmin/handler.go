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

// =========================== TENANT APPROVAL ===========================

// GET /superadmin/tenants?status=pending&limit=10&page=1
func (h *Handler) GetTenantsWithFilters(c *gin.Context) {
	status := strings.ToLower(c.DefaultQuery("status", "pending"))
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	tenants, total, err := h.service.GetTenantsWithFilters(c.Request.Context(), status, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tenants,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// PATCH /superadmin/tenants/:id
func (h *Handler) UpdateTenantApprovalStatus(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	adminID := c.GetUint("userID")
	action := strings.ToLower(body.Status)

	switch action {
	case "approved":
		err = h.service.ApproveTenant(c.Request.Context(), uint(userID), adminID)
	case "rejected":
		if body.Reason == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rejection reason required"})
			return
		}
		err = h.service.RejectTenant(c.Request.Context(), uint(userID), adminID, body.Reason)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Use APPROVED or REJECTED"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant status updated successfully"})
}

// =========================== ENTITY APPROVAL ===========================

// GET /superadmin/entities?status=pending&limit=10&page=1
func (h *Handler) GetEntitiesWithFilters(c *gin.Context) {
	status := strings.ToUpper(c.DefaultQuery("status", "PENDING"))
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	entities, total, err := h.service.GetEntitiesWithFilters(c.Request.Context(), status, limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  entities,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// PATCH /superadmin/entities/:id
func (h *Handler) UpdateEntityApprovalStatus(c *gin.Context) {
	idStr := c.Param("id")
	entityID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	adminID := c.GetUint("userID")
	action := strings.ToLower(body.Status)

	switch action {
	case "approved":
		err = h.service.ApproveEntity(c.Request.Context(), uint(entityID), adminID)
	case "rejected":
		if body.Reason == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rejection reason required"})
			return
		}
		err = h.service.RejectEntity(c.Request.Context(), uint(entityID), adminID, body.Reason)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Use APPROVED or REJECTED"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entity status updated successfully"})
}

// GET /superadmin/tenant-approval-counts
func (h *Handler) GetTenantApprovalCounts(c *gin.Context) {
	ctx := c.Request.Context()

	counts, err := h.service.GetTenantApprovalCounts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tenant approval counts"})
		return
	}

	c.JSON(http.StatusOK, counts)
}

// GET /superadmin/temple-approval-counts
func (h *Handler) GetTempleApprovalCounts(c *gin.Context) {
	ctx := c.Request.Context()

	counts, err := h.service.GetTempleApprovalCounts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch temple approval counts"})
		return
	}

	c.JSON(http.StatusOK, counts)
}

// =========================== USER MANAGEMENT ===========================

// POST /superadmin/users - Create new user (admin-created users)
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate templeadmin details if role is templeadmin
	if strings.ToLower(req.Role) == "templeadmin" {
		if req.TempleName == "" || req.TemplePlace == "" || req.TempleAddress == "" ||
			req.TemplePhoneNo == "" || req.TempleDescription == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All temple details are required for Temple Admin role"})
			return
		}
	}

	adminID := c.GetUint("userID")
	if err := h.service.CreateUser(c.Request.Context(), req, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// GET /superadmin/users - Get all users with pagination (excluding devotee and volunteer)
func (h *Handler) GetUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	search := c.DefaultQuery("search", "")
	roleFilter := c.DefaultQuery("role", "")
	statusFilter := c.DefaultQuery("status", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	users, total, err := h.service.GetUsers(c.Request.Context(), limit, page, search, roleFilter, statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GET /superadmin/users/:id - Get user by ID
func (h *Handler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// PUT /superadmin/users/:id - Update user
func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("userID")
	if err := h.service.UpdateUser(c.Request.Context(), uint(userID), req, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DELETE /superadmin/users/:id - Delete user
func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	adminID := c.GetUint("userID")
	if err := h.service.DeleteUser(c.Request.Context(), uint(userID), adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// PATCH /superadmin/users/:id/status - Activate/Deactivate user
func (h *Handler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	validStatuses := []string{"active", "inactive"}
	status := strings.ToLower(body.Status)
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Use 'active' or 'inactive'"})
		return
	}

	adminID := c.GetUint("userID")
	if err := h.service.UpdateUserStatus(c.Request.Context(), uint(userID), status, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

// GET /superadmin/user-roles - Get all available user roles
func (h *Handler) GetUserRoles(c *gin.Context) {
	roles, err := h.service.GetUserRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}