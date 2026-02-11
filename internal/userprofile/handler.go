package userprofile

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
	"github.com/sharath018/temple-management-backend/middleware"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// ===========================
// 🔹 PROFILE ENDPOINTS
// ===========================

// GET /profiles/me
func (h *Handler) GetMyProfile(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(auth.User)

	profile, err := h.service.Get(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GET /entities/:id/devotees/:userId/profile
func (h *Handler) GetDevoteeProfileByEntity(c *gin.Context) {

	// -------- Parse route params --------
	entityIDUint, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	userIDUint, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	entityID := uint(entityIDUint)
	userID := uint(userIDUint)

	// -------- Get authenticated user --------
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser, ok := userVal.(auth.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user object"})
		return
	}

	// -------- Get access context --------
	accessVal, exists := c.Get("access_context")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access context"})
		return
	}
	accessCtx, ok := accessVal.(middleware.AccessContext)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access context"})
		return
	}

	// -------- Role-based access check --------
	hasAccess := false

	switch currentUser.Role.RoleName {

	case "superadmin":
		hasAccess = true

	case "templeadmin":
		hasAccess = accessCtx.DirectEntityID != nil &&
			*accessCtx.DirectEntityID == entityID

	case "standarduser", "monitoringuser":
		hasAccess = accessCtx.AssignedEntityID != nil &&
			*accessCtx.AssignedEntityID == entityID

	default:
		hasAccess = false
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied to view devotee profile for this entity",
		})
		return
	}
	profile, err := h.service.GetByUserIDAndEntity(userID, entityID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Devotee profile not found for this entity",
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// POST /profiles
func (h *Handler) CreateOrUpdateProfile(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(auth.User)

	var input DevoteeProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var entityID uint

	// ✅ Try to extract EntityID from user context first
	if currentUser.EntityID != nil {
		entityID = *currentUser.EntityID
	} else {
		// 🔍 fallback: look up from memberships
		memberships, err := h.service.ListMemberships(currentUser.ID)
		if err != nil || len(memberships) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No associated temple found. Please join a temple first."})
			return
		}
		entityID = memberships[0].EntityID // default to first membership
	}

	// ✅ EXTRACT IP ADDRESS FROM CONTEXT
	ip := middleware.GetIPFromContext(c)

	// ✅ PASS CONTEXT AND IP TO SERVICE
	profile, err := h.service.CreateOrUpdateProfile(c.Request.Context(), currentUser.ID, entityID, input, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ===========================
// 🔹 MEMBERSHIP ENDPOINTS
// ===========================

// POST /memberships
func (h *Handler) JoinTemple(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(auth.User)

	var input struct {
		EntityID uint `json:"entity_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// ✅ EXTRACT IP ADDRESS FROM CONTEXT
	ip := middleware.GetIPFromContext(c)

	// ✅ GET USER ROLE FOR AUDIT LOGGING
	userRole := "unknown"
	if currentUser.Role.RoleName != "" {
		userRole = currentUser.Role.RoleName
	}

	// ✅ PASS CONTEXT, USER ROLE AND IP TO SERVICE
	membership, err := h.service.JoinTemple(c.Request.Context(), currentUser.ID, input.EntityID, userRole, ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, membership)
}

// GET /memberships
func (h *Handler) ListMemberships(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(auth.User)

	memberships, err := h.service.ListMemberships(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch memberships"})
		return
	}

	c.JSON(http.StatusOK, memberships)
}

// ===========================
// 🔹 SEARCH TEMPLES ENDPOINT
// ===========================

// GET /temples/search?query=&state=&temple_type=
func (h *Handler) SearchTemples(c *gin.Context) {
	query := c.Query("query")            // name/city/state search text
	state := c.Query("state")            // optional filter
	templeType := c.Query("temple_type") // optional filter

	results, err := h.service.SearchTemples(query, state, templeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search temples"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GET /profiles/recent-temples
func (h *Handler) GetRecentTemples(c *gin.Context) {
	temples, err := h.service.GetRecentTemples()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch recent temples"})
		return
	}
	c.JSON(http.StatusOK, temples)
}
