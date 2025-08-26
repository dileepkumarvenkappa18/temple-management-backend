package middleware

import (
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// Role constants to avoid string typos
const (
	RoleSuperAdmin     = "superadmin"
	RoleTempleAdmin    = "templeadmin"
	RoleStandardUser   = "standarduser"
	RoleMonitoringUser = "monitoringuser"
	RoleDevotee        = "devotee"
	RoleVolunteer      = "volunteer"
)

// AccessContext stores user access information
type AccessContext struct {
	UserID           uint
	RoleName         string
	DirectEntityID   *uint  // User's own entity (for templeadmin)
	AssignedEntityID *uint  // Assigned tenant entity (for standarduser/monitoringuser)
	PermissionType   string // "full" or "readonly"
}

// GetAccessibleEntityID returns the entity ID the user can access
func (ac *AccessContext) GetAccessibleEntityID() *uint {
	if ac.AssignedEntityID != nil {
		return ac.AssignedEntityID
	}
	return ac.DirectEntityID
}

// CanWrite returns true if the user has write permissions
func (ac *AccessContext) CanWrite() bool {
	return ac.PermissionType == "full"
}

// CanRead returns true if the user has read permissions
func (ac *AccessContext) CanRead() bool {
	return ac.PermissionType == "full" || ac.PermissionType == "readonly"
}

// ResolveAccessContext helper to create access context from user and assignment
func ResolveAccessContext(user auth.User, assignedTenantID *uint) AccessContext {
	accessContext := AccessContext{
		UserID:         user.ID,
		RoleName:       user.Role.RoleName,
		DirectEntityID: user.EntityID,
		PermissionType: "full", // default
	}
	
	// If user is standarduser or monitoringuser with assigned tenant
	if assignedTenantID != nil {
		accessContext.AssignedEntityID = assignedTenantID
		
		// Set permission type based on role
		if user.Role.RoleName == RoleMonitoringUser {
			accessContext.PermissionType = "readonly"
		}
	}
	
	return accessContext
}

// ExtractTenantIDFromContext extracts tenant ID from request headers
func ExtractTenantIDFromContext(c *gin.Context) *uint {
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		return nil
	}
	
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 64)
	if err != nil {
		return nil
	}
	
	id := uint(tenantID)
	return &id
}

