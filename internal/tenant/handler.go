package tenant

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

// Handler wraps the service layer.
type Handler struct {
    service *Service
}

// NewHandler returns a new Handler instance.
func NewHandler(service *Service) *Handler {
    return &Handler{service}
}

// Mock function to get tenant ID from logged-in context/session
// Replace with real authentication logic
func GetLoggedInTenantID(c *gin.Context) uint {
    // Example: read from JWT claims or session
    return 1 // assuming tenant ID 1 for demo
}

// GetUsers fetches users for the logged-in tenant, with optional filtering.
func (h *Handler) GetUsers(c *gin.Context) {
    role := c.Query("role")
    name := c.Query("name")

    tenantID := GetLoggedInTenantID(c)

    users, err := h.service.GetUsers(tenantID, role, name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

// CreateOrUpdateUser creates a new user for the logged-in tenant.
func (h *Handler) CreateOrUpdateUser(c *gin.Context) {
    var input TenantUser
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    // Set tenant ID from logged-in tenant
    input.TenantID = GetLoggedInTenantID(c)

    user, err := h.service.CreateOrUpdateUser(input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user: " + err.Error()})
        return
    }

    // Hide password before sending response
    user.Password = ""
    c.JSON(http.StatusOK, user)
}

// UpdateUser updates user details based on user ID in URL.
func (h *Handler) UpdateUser(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID: " + err.Error()})
        return
    }

    var input TenantUser
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    input.ID = uint(id)
    input.TenantID = GetLoggedInTenantID(c) // ensure tenant ownership

    user, err := h.service.CreateOrUpdateUser(input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
        return
    }

    user.Password = ""
    c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user by ID
func (h *Handler) DeleteUser(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID: " + err.Error()})
        return
    }

    tenantID := GetLoggedInTenantID(c)
    err = h.service.DeleteUser(uint(id), tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
