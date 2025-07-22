package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service Service }

func NewHandler(s Service) *Handler { return &Handler{s} }

// ===============================
// Registration
// ===============================

type RegisterRequest struct {
	FullName string `json:"fullName" binding:"required" example:"Sharath Kumar"`
	Email    string `json:"email" binding:"required,email" example:"example@gmail.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
	Role     string `json:"role" binding:"required" example:"templeadmin"`
	Phone    string `json:"phone" binding:"required" example:"+919876543210"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ❌ Block superadmin registration
	if req.Role == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Super Admin registration is not allowed"})
		return
	}

	// ✅ Validate Gmail only
	if !isGmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only @gmail.com emails are allowed for registration"})
		return
	}

	if err := h.service.Register(RegisterInput(req)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "templeadmin" {
		c.JSON(http.StatusCreated, gin.H{"message": "Temple Admin registered. Awaiting approval."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

// 🔍 Email helper
func isGmail(email string) bool {
	return strings.HasSuffix(strings.ToLower(email), "@gmail.com")
}

// ===============================
// Login
// ===============================

type loginReq struct {
	Email    string `json:"email" binding:"required,email" example:"sharath@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

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

	userPayload := gin.H{
		"id":       user.ID,
		"fullName": user.FullName,
		"email":    user.Email,
		"roleId":   user.RoleID,
	}
	if user.EntityID != nil {
		userPayload["entityId"] = user.EntityID
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"user":         userPayload,
	})
}

// ===============================
// Refresh Token
// ===============================

type refreshReq struct {
	RefreshToken string `json:"refreshToken" binding:"required" example:"your_refresh_token_here"`
}

func (h *Handler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": token})
}

// ===============================
// Forgot Password
// ===============================

type forgotPasswordReq struct {
	Email string `json:"email" binding:"required,email" example:"sharath@example.com"`
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req forgotPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.RequestPasswordReset(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reset link sent to email (if account exists)"})
}

// ===============================
// Reset Password
// ===============================

type resetPasswordReq struct {
	Token       string `json:"token" binding:"required" example:"reset_token_abc123"`
	NewPassword string `json:"newPassword" binding:"required,min=6" example:"newsecret123"`
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}

// ===============================
// Logout
// ===============================

func (h *Handler) Logout(c *gin.Context) {
	_ = h.service.Logout() // stateless
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}