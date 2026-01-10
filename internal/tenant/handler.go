package tenant

import (
	//"bytes"
	"fmt"
	"io"
    "os"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/google/uuid"

	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/s3"
)


// Handler handles HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new handler instance
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// helper: extract userID from context
func getUserIDFromContext(c *gin.Context) (uint, bool) {
	if v, ok := c.Get("userID"); ok {
		switch id := v.(type) {
		case uint:
			return id, true
		case int:
			return uint(id), true
		case float64:
			return uint(id), true
		}
	}

	if v, ok := c.Get("user_id"); ok {
		switch id := v.(type) {
		case uint:
			return id, true
		case int:
			return uint(id), true
		case float64:
			return uint(id), true
		}
	}

	return 0, false
}

// =========================
// GET TENANT PROFILE
// =========================
func (h *Handler) GetTenantProfile(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	log.Printf("📋 Getting tenant profile for user ID: %d", userID)

	profile, err := h.service.GetTenantProfile(userID)
	if err != nil {
		log.Printf("❌ Failed to fetch tenant profile: %v", err)

		// Business errors → NOT 500
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Tenant profile not available",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ Tenant profile fetched successfully for user %d", userID)
	c.JSON(http.StatusOK, profile)
}

// =========================
// UPDATE TENANT PROFILE
// =========================
func (h *Handler) UpdateTenantProfile(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var input UpdateTenantProfileRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("❌ Invalid update payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	log.Printf("📝 Updating tenant profile for user ID: %d", userID)

	profile, err := h.service.UpdateTenantProfile(userID, input)
	if err != nil {
		log.Printf("❌ Failed to update tenant profile: %v", err)
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Failed to update tenant profile",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ Tenant profile updated successfully for user %d", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}

// UpdateUser handles the PUT request to update a user
func (h *Handler) UpdateUser(c *gin.Context) {
    // Get tenant ID and user ID from route parameters
    tenantIDStr := c.Param("id")
    userIDStr := c.Param("userId")
    
    log.Printf("🔵 Updating user %s for tenant %s", userIDStr, tenantIDStr)
    
    tenantID, err := strconv.ParseUint(tenantIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
        return
    }
    
    userID, err := strconv.ParseUint(userIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    // First log raw data for debugging
    var rawData map[string]interface{}
    if err := c.ShouldBindJSON(&rawData); err != nil {
        log.Printf("🔴 Error binding raw JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return
    }
    
    log.Printf("🔵 Received raw update data: %+v", rawData)
    
    // Check if this is a status-only update
    if status, exists := rawData["Status"].(string); exists && len(rawData) <= 2 {
        log.Printf("🔵 Processing as status update: %s for user %d", status, userID)
        
        // Check if user belongs to this tenant
        exists, err := h.service.repo.CheckUserBelongsToTenant(uint(userID), uint(tenantID))
        if err != nil {
            log.Printf("🔴 Error checking user-tenant relationship: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user-tenant relationship"})
            return
        }
        
        if !exists {
            c.JSON(http.StatusBadRequest, gin.H{"error": "User does not belong to this tenant"})
            return
        }
        
        // Update status in both tables
        err = h.service.repo.UpdateUserStatus(uint(userID), uint(tenantID), status)
        if err != nil {
            log.Printf("🔴 Error updating status: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status: " + err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message": "Status updated successfully",
        })
        return
    }
    
    // Handle as a full user update
    input := UserInput{}
    
    if name, ok := rawData["Name"].(string); ok {
        input.Name = name
    }
    if email, ok := rawData["Email"].(string); ok {
        input.Email = email
    }
    if phone, ok := rawData["Phone"].(string); ok {
        input.Phone = phone
    }
    if role, ok := rawData["Role"].(string); ok {
        input.Role = role
    }
    if password, ok := rawData["Password"].(string); ok {
        input.Password = password
    }
    if status, ok := rawData["Status"].(string); ok {
        input.Status = status
    }
    
    log.Printf("🔵 Updating user %d for tenant %d: %s (%s), Role: %s, Status: %s", 
        userID, tenantID, input.Name, input.Email, input.Role, input.Status)
    
    if input.Name == "" || input.Email == "" {
        log.Printf("🔴 Missing required fields: Name or Email")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Email are required fields"})
        return
    }
    
    user, err := h.service.UpdateUser(uint(tenantID), uint(userID), input, input.Status)
    if err != nil {
        log.Printf("🔴 Error updating user: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
        return
    }
    
    log.Printf("✅ User updated successfully: %+v", user)
    c.JSON(http.StatusOK, gin.H{
        "message": "User updated successfully",
        "user": user,
    })
}

// UpdateUserStatus updates only a user's status
func (h *Handler) UpdateUserStatus(c *gin.Context) {
    tenantIDStr := c.Param("id")
    userIDStr := c.Param("userId")

    log.Printf("🔵 Updating status for user %s in tenant %s", userIDStr, tenantIDStr)

    tenantID, err := strconv.ParseUint(tenantIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
        return
    }

    userID, err := strconv.ParseUint(userIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var statusData struct {
        Status string `json:"Status"`
    }

    if err := c.ShouldBindJSON(&statusData); err == nil && statusData.Status != "" {
        log.Printf("🔵 Received status update: %s for user %d", statusData.Status, userID)

        user, err := h.service.UpdateUser(uint(tenantID), uint(userID), UserInput{}, statusData.Status)
        if err != nil {
            log.Printf("🔴 Error updating status: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status: " + err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "Status updated successfully",
            "user": user,
        })
        return
    }
}

// GetUsers handles the GET request to fetch tenant users
func (h *Handler) GetUsers(c *gin.Context) {
    log.Printf("🔴 GET USERS - Request path: %s", c.Request.URL.Path)
    log.Printf("🔴 GET USERS - All params: %v", c.Params)
    
    tenantIDStr := c.Param("id")
    log.Printf("🔴 GET USERS - Raw tenant ID from route param: %s", tenantIDStr)
    
    tenantID, err := strconv.ParseUint(tenantIDStr, 10, 64)
    if err != nil {
        log.Printf("🔴 ERROR - Invalid tenant ID: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
        return
    }
    
    log.Printf("🔴 GET USERS - Using tenant ID: %d", tenantID)
    
    role := c.Query("role")
    
    users, err := h.service.GetTenantUsers(uint(tenantID), role)
    if err != nil {
        log.Printf("Failed to fetch users: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users: " + err.Error()})
        return
    }
    
    if users == nil {
        users = []UserResponse{}
    }
    
    log.Printf("Returning %d users for tenant %d", len(users), tenantID)
    c.JSON(http.StatusOK, users)
}

// CreateOrUpdateUser handles the POST request to create or update a tenant user
func (h *Handler) CreateOrUpdateUser(c *gin.Context) {
    log.Printf("🔴 CREATE USER - Request path: %s", c.Request.URL.Path)
    log.Printf("🔴 CREATE USER - All params: %v", c.Params)
    
    var tenantID uint64
    var err error
    
    tenantIDHeader := c.GetHeader("X-Tenant-ID")
    log.Printf("🔴 CREATE USER - X-Tenant-ID header: %s", tenantIDHeader)
    
    if tenantIDHeader != "" {
        tenantID, err = strconv.ParseUint(tenantIDHeader, 10, 64)
        if err == nil {
            log.Printf("🔴 CREATE USER - Using tenant ID from header: %d", tenantID)
        } else {
            log.Printf("🔴 ERROR - Invalid tenant ID in header: %v", err)
        }
    }
    
    if err != nil || tenantIDHeader == "" {
        tenantIDStr := c.Param("id")
        log.Printf("🔴 CREATE USER - Raw tenant ID from route param: %s", tenantIDStr)
        
        tenantID, err = strconv.ParseUint(tenantIDStr, 10, 64)
        if err != nil {
            log.Printf("🔴 ERROR - Invalid tenant ID in route: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
            return
        }
        log.Printf("🔴 CREATE USER - Using tenant ID from route param: %d", tenantID)
    }
    
    var input UserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("Invalid input: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }
    
    creatorID, exists := c.Get("user_id")
    if !exists {
        creatorIDStr := c.GetHeader("X-User-ID")
        if creatorIDStr != "" {
            creatorIDUint, err := strconv.ParseUint(creatorIDStr, 10, 64)
            if err == nil {
                creatorID = uint(creatorIDUint)
            }
        }
    }
    
    creatorIDUint := uint(1)
    if id, ok := creatorID.(uint); ok {
        creatorIDUint = id
    } else if id, ok := creatorID.(float64); ok {
        creatorIDUint = uint(id)
    } else if id, ok := creatorID.(int); ok {
        creatorIDUint = uint(id)
    } else if id, ok := creatorID.(uint64); ok {
        creatorIDUint = uint(id)
    }
    
    log.Printf("Creating/updating user %s (%s) for tenant %d by creator %d", 
               input.Name, input.Email, tenantID, creatorIDUint)
    
    user, err := h.service.CreateOrUpdateTenantUser(uint(tenantID), input, creatorIDUint)
    if err != nil {
        log.Printf("Failed to create/update user: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user: " + err.Error()})
        return
    }
    
    log.Printf("User created/updated successfully: %s (ID: %d) for tenant ID: %d", 
               user.Email, user.ID, tenantID)
    c.JSON(http.StatusOK, gin.H{
        "message": "User created and assigned successfully",
        "user": user,
    })
}
func (h *Handler) UploadFile(c *gin.Context) {
	// Get logged-in user ID
	userIDVal, exists := c.Get("userID")
	if !exists {
		userIDVal, exists = c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
	}
	userID := userIDVal.(uint)

	// Get tenant profile to resolve tenant_id
	profile, err := h.service.GetTenantProfile(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant profile not found"})
		return
	}
	tenantID := profile.TenantID

	// Get file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	fileType := c.PostForm("type")
	if fileType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type is required"})
		return
	}

	// Validation rules
	var allowedExtensions []string
	var maxSize int64
	var folder string

	switch fileType {
	case "logo":
		allowedExtensions = []string{".jpg", ".jpeg", ".png", ".webp"}
		maxSize = 5 * 1024 * 1024
		folder = "logo"
	case "video":
		allowedExtensions = []string{".mp4", ".webm", ".mov"}
		maxSize = 50 * 1024 * 1024
		folder = "video"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Max file size is %d MB", maxSize/(1024*1024)),
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	valid := false
	for _, e := range allowedExtensions {
		if ext == e {
			valid = true
			break
		}
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	// Create directory
	baseDir := fmt.Sprintf("uploads/tenants/%d/%s", tenantID, folder)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Save file
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), ext)
	fullPath := filepath.Join(baseDir, filename)

	out, err := os.Create(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	fileURL := "/" + fullPath

	c.JSON(http.StatusOK, gin.H{
		"url": fileURL,
	})
}
