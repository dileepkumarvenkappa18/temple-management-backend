package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
	AssignedEntityID *uint  // Assigned entity from URL or user assignment
	PermissionType   string // "full" or "readonly"
	TenantID         uint   // Tenant this user belongs to (for standarduser/monitoringuser)
}

// GetAccessibleEntityID returns the entity ID the user can access
func (ac *AccessContext) GetAccessibleEntityID() *uint {
	// First try DirectEntityID (for templeadmin)
	if ac.DirectEntityID != nil {
		return ac.DirectEntityID
	}
	// Fall back to AssignedEntityID (for standard/monitoring users)
	return ac.AssignedEntityID
}

// GetEntityIDForOperation returns the entity ID to use for operations like create/update
func (ac *AccessContext) GetEntityIDForOperation() *uint {
	entityID := ac.GetAccessibleEntityID()
	if entityID != nil {
		fmt.Printf("Using entity ID %d for operation (role: %s)\n", *entityID, ac.RoleName)
	} else {
		fmt.Printf("WARNING: No entity ID available for operation (role: %s)\n", ac.RoleName)
	}
	return entityID
}

// CanWrite returns true if the user has write permissions
func (ac *AccessContext) CanWrite() bool {
	return ac.PermissionType == "full"
}

// CanRead returns true if the user has read permissions
func (ac *AccessContext) CanRead() bool {
	return ac.PermissionType == "full" || ac.PermissionType == "readonly"
}

// CanAccessEntity checks if the user can access a specific entity.
// FIXED: standarduser/monitoringuser with a valid TenantID are allowed through.
// The actual entity-to-tenant ownership check happens at the DB/handler level.
func (ac *AccessContext) CanAccessEntity(entityID uint) bool {
	// SuperAdmin can access any entity
	if ac.RoleName == RoleSuperAdmin {
		return true
	}

	// Check direct entity ID match (works for templeadmin, devotee, volunteer)
	accessibleEntityID := ac.GetAccessibleEntityID()
	if accessibleEntityID != nil && *accessibleEntityID == entityID {
		return true
	}

	// For standarduser/monitoringuser:
	// They are assigned to a TENANT (not directly to an entity).
	// If they have a valid TenantID, allow them through here.
	// The handler/DB layer is responsible for verifying the entity belongs to their tenant.
	if (ac.RoleName == RoleStandardUser || ac.RoleName == RoleMonitoringUser) && ac.TenantID > 0 {
		return true
	}

	return false
}

// ExtractTenantIDFromContext extracts tenant ID from the access context
func ExtractTenantIDFromContext(c *gin.Context) *uint {
	if accessCtxVal, exists := c.Get("access_context"); exists {
		if ctx, ok := accessCtxVal.(AccessContext); ok {
			if ctx.TenantID > 0 {
				return &ctx.TenantID
			}
		}
	}
	return nil
}

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
			fmt.Printf("[%s] Entity Resolution - UserID: %d, Role: %s, EntityID: %v, DirectEntityID: %v, TenantID: %v\n",
				operation, ctx.UserID, ctx.RoleName, entityID, ctx.DirectEntityID, ctx.TenantID)
		}
	}
}