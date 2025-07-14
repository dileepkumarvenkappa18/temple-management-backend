package userprofile

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	service Service
}

// NewHandler creates a new Handler
func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

// ===========================
// 🔹 PROFILE ENDPOINTS
// ===========================

// GET /profiles/:userID
func (h *Handler) GetProfile(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	profile, err := h.service.Get(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// PUT /profiles/:userID
func (h *Handler) CreateOrUpdateProfile(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var dto CreateOrUpdateProfileDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	profile, err := h.service.CreateOrUpdateProfile(uint(userID), dto)
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
// POST /memberships
func (h *Handler) JoinTemple(c *gin.Context) {
	var input struct {
		EntityID uint `json:"entity_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// ✅ Extract user ID from token (middleware must set it via `c.Set("user_id", userID)`)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	membership, err := h.service.JoinTemple(userID, input.EntityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not join temple"})
		return
	}

	c.JSON(http.StatusOK, membership)
}


// GET /memberships?user_id=1
func (h *Handler) ListMemberships(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing user_id"})
		return
	}

	memberships, err := h.service.ListMemberships(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch memberships"})
		return
	}

	c.JSON(http.StatusOK, memberships)
}