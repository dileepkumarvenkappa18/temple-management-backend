package userprofile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
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

// POST /profiles
func (h *Handler) CreateOrUpdateProfile(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := user.(auth.User)

	if currentUser.EntityID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Entity ID missing"})
		return
	}
	entityID := *currentUser.EntityID

	var input DevoteeProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	profile, err := h.service.CreateOrUpdateProfile(currentUser.ID, entityID, input)
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

	membership, err := h.service.JoinTemple(currentUser.ID, input.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not join temple"})
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
