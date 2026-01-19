package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service Service }

func NewHandler(s Service) *Handler { return &Handler{s} }

// ===============================
// Registration
// ===============================

type RegisterRequest struct {
	FullName string `form:"fullName" json:"fullName" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6"`
	Role     string `form:"role" json:"role" binding:"required"`
	Phone    string `form:"phone" json:"phone" binding:"required"`

	TempleName        string `form:"templeName" json:"templeName"`
	TemplePlace       string `form:"templePlace" json:"templePlace"`
	TempleAddress     string `form:"templeAddress" json:"templeAddress"`
	TemplePhoneNo     string `form:"templePhoneNo" json:"templePhoneNo"`
	TempleDescription string `form:"templeDescription" json:"templeDescription"`

	    // 🆕 Bank Account Details
    AccountHolderName string `form:"accountHolderName" json:"accountHolderName"`
    AccountNumber     string `form:"accountNumber" json:"accountNumber"`
    BankName          string `form:"bankName" json:"bankName"`
    BranchName        string `form:"branchName" json:"branchName"`
    IFSCCode          string `form:"ifscCode" json:"ifscCode"`
    AccountType       string `form:"accountType" json:"accountType"`
    UPIID             string `form:"upiId" json:"upiId"`

	LogoURL       string `json:"logo_url"`
	IntroVideoURL string `json:"intro_video_url"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ❌ Block superadmin
	if strings.ToLower(req.Role) == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Super Admin registration is not allowed"})
		return
	}

	// ✅ Gmail only
	if !isGmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only @gmail.com emails are allowed"})
		return
	}

	// ✅ Temple admin validation
	if strings.ToLower(req.Role) == "templeadmin" {
		if req.TempleName == "" ||
			req.TemplePlace == "" ||
			req.TempleAddress == "" ||
			req.TemplePhoneNo == "" ||
			req.TempleDescription == "" {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "All temple details are required for Temple Admin registration",
			})
			return
		}
		   // 🆕 Validate bank details
    if req.AccountHolderName == "" ||
        req.AccountNumber == "" ||
        req.BankName == "" ||
        req.BranchName == "" ||
        req.IFSCCode == "" ||
        req.AccountType == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "All bank account details are required for Temple Admin registration",
        })
        return
    }

		// Check for logo file
		if _, err := c.FormFile("logo"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Temple logo is required for Temple Admin registration",
			})
			return
		}

		// Check for video file
		if _, err := c.FormFile("video"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Temple intro video is required for Temple Admin registration",
			})
			return
		}
	}

	// 1️⃣ Create user + tenant FIRST
	user, err := h.service.RegisterAndReturnUser(RegisterInput{
		FullName:          req.FullName,
		Email:             req.Email,
		Password:          req.Password,
		Role:              req.Role,
		Phone:             req.Phone,
		TempleName:        req.TempleName,
		TemplePlace:       req.TemplePlace,
		TempleAddress:     req.TempleAddress,
		TemplePhoneNo:     req.TemplePhoneNo,
		TempleDescription: req.TempleDescription,
		 // 🆕 Add bank details
    AccountHolderName: req.AccountHolderName,
    AccountNumber:     req.AccountNumber,
    BankName:          req.BankName,
    BranchName:        req.BranchName,
    IFSCCode:          req.IFSCCode,
    AccountType:       req.AccountType,
    UPIID:             req.UPIID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2️⃣ Only process files for temple admin
	if strings.ToLower(req.Role) == "templeadmin" {
		// Create tenant directory - use /data/uploads to match main.go
		tenantDir := filepath.Join("/data/uploads", "tenants", fmt.Sprint(user.ID))
		if err := os.MkdirAll(tenantDir, os.ModePerm); err != nil {
			log.Printf("⚠️ Failed to create tenant directory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create storage directory",
			})
			return
		}

		var publicLogoURL, publicVideoURL string

		// 3️⃣ Save LOGO
		logoFile, err := c.FormFile("logo")
		if err == nil {
			logoFilename := "logo" + filepath.Ext(logoFile.Filename)
			logoPath := filepath.Join(tenantDir, logoFilename)
			
			if err := c.SaveUploadedFile(logoFile, logoPath); err != nil {
				log.Printf("❌ Failed to save logo: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to save temple logo",
				})
				return
			}
			
			// Make URL: /uploads/tenants/{userID}/logo.ext
			publicLogoURL = fmt.Sprintf("/uploads/tenants/%d/%s", user.ID, logoFilename)
			log.Printf("✅ Logo saved to: %s (Public URL: %s)", logoPath, publicLogoURL)
		} else {
			log.Printf("❌ Logo file error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Temple logo is required",
			})
			return
		}

		// 4️⃣ Save VIDEO
		videoFile, err := c.FormFile("video")
		if err == nil {
			videoFilename := "intro" + filepath.Ext(videoFile.Filename)
			videoPath := filepath.Join(tenantDir, videoFilename)
			
			if err := c.SaveUploadedFile(videoFile, videoPath); err != nil {
				log.Printf("❌ Failed to save video: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to save temple video",
				})
				return
			}
			
			// Make URL: /uploads/tenants/{userID}/intro.ext
			publicVideoURL = fmt.Sprintf("/uploads/tenants/%d/%s", user.ID, videoFilename)
			log.Printf("✅ Video saved to: %s (Public URL: %s)", videoPath, publicVideoURL)
		} else {
			log.Printf("❌ Video file error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Temple intro video is required",
			})
			return
		}

		// 5️⃣ Update tenant_details with public URLs - THIS IS THE CRITICAL FIX
		if publicLogoURL != "" || publicVideoURL != "" {
			log.Printf("🔄 Updating database - UserID: %d, Logo: %s, Video: %s", 
				user.ID, publicLogoURL, publicVideoURL)
			
			if err := h.service.UpdateTenantMedia(user.ID, publicLogoURL, publicVideoURL); err != nil {
				log.Printf("❌ Failed to update tenant media URLs in database: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to update media information in database",
				})
				return
			}
			
			log.Printf("✅ Database updated successfully - Logo: %s, Video: %s", 
				publicLogoURL, publicVideoURL)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Temple Admin registered. Awaiting approval.",
		"tenantId": user.ID,
	})
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
// Forgot Password - FIXED
// ===============================

type forgotPasswordReq struct {
	Email string `json:"email" binding:"required,email" example:"sharath@example.com"`
}

// Custom error types for better error handling
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailNotSent    = errors.New("failed to send email")
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrEmailService    = errors.New("email service unavailable")
	ErrRateLimitExceed = errors.New("too many requests")
)

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req forgotPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": "Please provide a valid email address",
		})
		return
	}

	// Validate email format
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid email",
			"message": "Please provide a valid email address",
		})
		return
	}

	// Call service layer
	err := h.service.RequestPasswordReset(req.Email)
	
	if err != nil {
		// 🔍 Determine the type of error and respond accordingly
		switch {
		case errors.Is(err, ErrUserNotFound):
			// ⚠️ Security: Don't reveal if user exists or not
			// Return success message but log the attempt
			c.JSON(http.StatusOK, gin.H{
				"message": "If an account exists with this email, a password reset link has been sent",
			})
			return

		case errors.Is(err, ErrEmailNotSent), errors.Is(err, ErrEmailService):
			// 🚨 Email service failure - return 500
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to send email",
				"message": "Email service is currently unavailable. Please try again later or contact support.",
			})
			return

		case errors.Is(err, ErrRateLimitExceed):
			// 🚫 Rate limit exceeded
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "too many requests",
				"message": "You have requested too many password resets. Please wait 15 minutes and try again.",
			})
			return

		case strings.Contains(err.Error(), "email"):
			// Generic email-related error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to send email",
				"message": "Unable to send password reset email. Please contact support.",
			})
			return

		default:
			// Unknown error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal server error",
				"message": "An unexpected error occurred. Please try again later.",
			})
			return
		}
	}

	// ✅ Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "If an account exists with this email, a password reset link has been sent",
	})
}

// ===============================
// Reset Password - FIXED
// ===============================

type resetPasswordReq struct {
	Token       string `json:"token" binding:"required" example:"reset_token_abc123"`
	NewPassword string `json:"newPassword" binding:"required,min=6" example:"newsecret123"`
}

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token expired")
	ErrWeakPassword     = errors.New("password too weak")
	ErrTokenNotFound    = errors.New("token not found")
)

func (h *Handler) ResetPassword(c *gin.Context) {
	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": "Please provide both token and new password",
		})
		return
	}

	// Validate password strength
	if len(req.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "weak password",
			"message": "Password must be at least 6 characters long",
		})
		return
	}

	// Call service layer
	err := h.service.ResetPassword(req.Token, req.NewPassword)
	
	if err != nil {
		// 🔍 Determine the type of error and respond accordingly
		switch {
		case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrTokenNotFound):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid token",
				"message": "This password reset link is invalid. Please request a new one.",
			})
			return

		case errors.Is(err, ErrExpiredToken):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "expired token",
				"message": "This password reset link has expired. Please request a new one.",
			})
			return

		case errors.Is(err, ErrWeakPassword):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "weak password",
				"message": "Password does not meet security requirements. Please use a stronger password.",
			})
			return

		case strings.Contains(err.Error(), "token"):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid token",
				"message": "This password reset link is invalid or has expired. Please request a new one.",
			})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal server error",
				"message": "Unable to reset password. Please try again later.",
			})
			return
		}
	}

	// ✅ Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset successfully. You can now login with your new password.",
	})
}

// ===============================
// Logout
// ===============================

func (h *Handler) Logout(c *gin.Context) {
	_ = h.service.Logout() // stateless
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ===============================
// Public Roles
// ===============================

func (h *Handler) GetPublicRoles(c *gin.Context) {
	roles, err := h.service.GetPublicRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}


func (h *Handler) GetAccountDetails(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := userVal.(User)

	if strings.ToLower(user.Role.RoleName) != "templeadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	data, err := h.service.GetAccountDetails(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
func (h *Handler) UpdateAccountDetails(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := userVal.(User)

	if strings.ToLower(user.Role.RoleName) != "templeadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req struct {
		FullName          string `json:"full_name"`
		Phone             string `json:"phone"`
		TempleName        string `json:"temple_name"`
		TemplePlace       string `json:"temple_place"`
		TempleAddress     string `json:"temple_address"`
		TemplePhoneNo     string `json:"temple_phone_no"`
		TempleDescription string `json:"temple_description"`
		LogoURL           string `json:"logo_url"`
		IntroVideoURL     string `json:"intro_video_url"`
		AccountHolderName string `json:"account_holder_name"`
		AccountNumber     string `json:"account_number"`
		BankName          string `json:"bank_name"`
		BranchName        string `json:"branch_name"`
		IFSCCode          string `json:"ifsc_code"`
		AccountType       string `json:"account_type"`
		UPIID             string `json:"upi_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := UpdateAccountDetailsInput{
		FullName:          req.FullName,
		Phone:             req.Phone,
		TempleName:        req.TempleName,
		TemplePlace:       req.TemplePlace,
		TempleAddress:     req.TempleAddress,
		TemplePhoneNo:     req.TemplePhoneNo,
		TempleDescription: req.TempleDescription,
		LogoURL:           req.LogoURL,
		IntroVideoURL:     req.IntroVideoURL,
		AccountHolderName: req.AccountHolderName,
		AccountNumber:     req.AccountNumber,
		BankName:          req.BankName,
		BranchName:        req.BranchName,
		IFSCCode:          req.IFSCCode,
		AccountType:       req.AccountType,
		UPIID:             req.UPIID,
	}

	data, err := h.service.UpdateAccountDetails(user.ID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Account details updated successfully"})
}