package middleware

import (
	"strconv"
<<<<<<< HEAD
	"fmt"
	
	"github.com/gin-gonic/gin"
=======
	
	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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

<<<<<<< HEAD
// GetEntityIDForOperation returns the entity ID to use for operations like create/update
func (ac *AccessContext) GetEntityIDForOperation() *uint {
	// For operations, prefer the most specific entity ID available
	entityID := ac.GetAccessibleEntityID()
	if entityID != nil {
		fmt.Printf("Using entity ID %d for operation (role: %s)\n", *entityID, ac.RoleName)
	} else {
		fmt.Printf("WARNING: No entity ID available for operation (role: %s)\n", ac.RoleName)
	}
	return entityID
}

=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
// CanWrite returns true if the user has write permissions
func (ac *AccessContext) CanWrite() bool {
	return ac.PermissionType == "full"
}

// CanRead returns true if the user has read permissions
func (ac *AccessContext) CanRead() bool {
	return ac.PermissionType == "full" || ac.PermissionType == "readonly"
}

<<<<<<< HEAD
// CanAccessEntity checks if the user can access a specific entity
func (ac *AccessContext) CanAccessEntity(entityID uint) bool {
	// SuperAdmin can access any entity
	if ac.RoleName == RoleSuperAdmin {
		return true
	}
	
	// Check if this entity matches user's accessible entity
	accessibleEntityID := ac.GetAccessibleEntityID()
	if accessibleEntityID != nil && *accessibleEntityID == entityID {
		return true
	}
	
	return false
}

// ResolveAccessContext helper to create access context from user and assignment
func ResolveAccessContext(user interface{}, assignedTenantID *uint) AccessContext {
	// Type assertion to get the auth.User fields we need
	var userID uint
	var roleName string
	var entityID *uint
	
	// This approach allows the function to work without directly importing auth
	// which helps prevent import cycles
	switch u := user.(type) {
	case struct {
		ID       uint
		RoleName string
		EntityID *uint
	}:
		userID = u.ID
		roleName = u.RoleName
		entityID = u.EntityID
	default:
		// Try to extract using reflection or other methods
		// For now, return a minimal context
		return AccessContext{
			PermissionType: "readonly",
		}
	}

	accessContext := AccessContext{
		UserID:         userID,
		RoleName:       roleName,
		DirectEntityID: entityID,
		PermissionType: "full", // default
	}

	// If user is standarduser or monitoringuser with assigned tenant
	if assignedTenantID != nil {
		accessContext.AssignedEntityID = assignedTenantID

		// Set permission type based on role
		if roleName == RoleMonitoringUser {
			accessContext.PermissionType = "readonly"
		}
	}

=======
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
	
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	return accessContext
}

// ExtractTenantIDFromContext extracts tenant ID from request headers
func ExtractTenantIDFromContext(c *gin.Context) *uint {
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		return nil
	}
<<<<<<< HEAD

=======
	
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 64)
	if err != nil {
		return nil
	}
<<<<<<< HEAD

=======
	
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	id := uint(tenantID)
	return &id
}

<<<<<<< HEAD
// GetEntityIDFromContext is a utility function to get the current entity ID from context
func GetEntityIDFromContext(c *gin.Context) *uint {
	// First try to get from context (set by middleware)
	if entityID, exists := c.Get("entity_id"); exists {
		if id, ok := entityID.(uint); ok {
			return &id
		}
	}
	
	// Try to get from access context
	if accessCtx, exists := c.Get("access_context"); exists {
		if ctx, ok := accessCtx.(AccessContext); ok {
			return ctx.GetEntityIDForOperation()
		}
	}
	
	return nil
}

// ValidateEntityAccess checks if the current user can access the specified entity
func ValidateEntityAccess(c *gin.Context, requestedEntityID uint) bool {
	if accessCtx, exists := c.Get("access_context"); exists {
		if ctx, ok := accessCtx.(AccessContext); ok {
			return ctx.CanAccessEntity(requestedEntityID)
		}
	}
	return false
}

// LogEntityResolution logs entity resolution for debugging
func LogEntityResolution(c *gin.Context, operation string) {
	entityID := GetEntityIDFromContext(c)
	if accessCtx, exists := c.Get("access_context"); exists {
		if ctx, ok := accessCtx.(AccessContext); ok {
			fmt.Printf("[%s] Entity Resolution - UserID: %d, Role: %s, EntityID: %v, DirectEntityID: %v, AssignedEntityID: %v\n",
				operation, ctx.UserID, ctx.RoleName, entityID, ctx.DirectEntityID, ctx.AssignedEntityID)
		}
	}
}
=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
