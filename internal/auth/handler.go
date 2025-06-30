package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Service handler wrapper
type Handler struct{ service Service }

func NewHandler(s Service) *Handler { return &Handler{s} }

// ✅ Swagger-compatible request struct for registration
type RegisterRequest struct {
	FullName string `json:"fullName" binding:"required" example:"Sharath Kumar"`
	Email    string `json:"email" binding:"required,email" example:"sharath@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
	Role     string `json:"role" binding:"required" example:"templeadmin"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email" example:"sharath@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken" binding:"required" example:"your_refresh_token_here"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "User details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot register as a Super Admin"})
		return
	}

	if err := h.service.Register(RegisterInput(req)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ Templeadmins go through approval flow
	if req.Role == "templeadmin" {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Temple Admin registered successfully. Awaiting Super Admin approval.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}


// Login godoc
// @Summary Login and get tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginReq body loginReq true "Login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, user, err := h.service.Login(LoginInput(req))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"user": gin.H{
			"id":        user.ID,
			"fullName":  user.FullName,
			"email":     user.Email,
			"roleId":    user.RoleID,
			// "tenantId": optional → we'll add this if token logic supports it later
		},
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refreshReq body refreshReq true "Refresh payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newAccessToken, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": newAccessToken})
}
