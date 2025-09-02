package tenant

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "log"
)

// Handler handles HTTP requests
type Handler struct {
    service *Service
}

// NewHandler creates a new handler instance
func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// GetUsers handles the GET request to fetch tenant users
func (h *Handler) GetUsers(c *gin.Context) {
    // CRITICAL DEBUGGING
    log.Printf("🔴 GET USERS - Request path: %s", c.Request.URL.Path)
    log.Printf("🔴 GET USERS - All params: %v", c.Params)
    
    // Get tenant ID from route parameter
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
    
    // Always return an array, even if empty
    if users == nil {
        users = []UserResponse{}
    }
    
    log.Printf("Returning %d users for tenant %d", len(users), tenantID)
    c.JSON(http.StatusOK, users)
}

// CreateOrUpdateUser handles the POST request to create or update a tenant user
// CreateOrUpdateUser handles the POST request to create or update a tenant user
func (h *Handler) CreateOrUpdateUser(c *gin.Context) {
    // CRITICAL DEBUGGING
    log.Printf("🔴 CREATE USER - Request path: %s", c.Request.URL.Path)
    log.Printf("🔴 CREATE USER - All params: %v", c.Params)
    
    // Get tenant ID preferring the X-Tenant-ID header over route parameter
    var tenantID uint64
    var err error
    
    // First try to get tenant ID from header
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
    
    // If header parsing failed, fall back to route parameter
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
    
    log.Printf("Creating/updating user %s (%s) for tenant %d", input.Name, input.Email, tenantID)
    
    user, err := h.service.CreateOrUpdateTenantUser(uint(tenantID), input)
    if err != nil {
        log.Printf("Failed to create/update user: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user: " + err.Error()})
        return
    }
    
    log.Printf("User created/updated successfully: %s (ID: %d) for tenant ID: %d", user.Email, user.ID, tenantID)
    c.JSON(http.StatusOK, gin.H{
        "message": "User created and assigned successfully",
        "user": user,
    })
}