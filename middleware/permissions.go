package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// AccessContext stores user access information
type AccessContext struct {
	UserID             uint
	RoleName           string
	DirectEntityID     *uint   // User's own entity (for templeadmin)
	AssignedEntityID   *uint   // Assigned tenant entity (for standarduser/monitoringuser)
	PermissionType     string  // "full" or "readonly"
}

// GetAccessibleEntityID returns the entity ID the user can access
func (ac AccessContext) GetAccessibleEntityID() *uint {
	if ac.AssignedEntityID != nil {
		return ac.AssignedEntityID
	}
	return ac.DirectEntityID
}

// CanWrite returns true if the user has write permissions
func (ac AccessContext) CanWrite() bool {
	return ac.PermissionType == "full"
}

// CanRead returns true if the user has read permissions
func (ac AccessContext) CanRead() bool {
	return ac.PermissionType == "full" || ac.PermissionType == "readonly"
}

// RequireTempleAccess middleware checks if user can access temple admin functions
func RequireTempleAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessContext, exists := c.Get("access_context")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
			return
		}

		ctx := accessContext.(AccessContext)
		
		// Allow if: templeadmin with entity, OR assigned standarduser/monitoringuser
		if (ctx.RoleName == "templeadmin" && ctx.DirectEntityID != nil) ||
		   (ctx.RoleName == "standarduser" && ctx.AssignedEntityID != nil) ||
		   (ctx.RoleName == "monitoringuser" && ctx.AssignedEntityID != nil) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "temple access required"})
	}
}

// RequireWriteAccess middleware checks if user can perform write operations
func RequireWriteAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessContext, exists := c.Get("access_context")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
			return
		}

		ctx := accessContext.(AccessContext)
		
		if !ctx.CanWrite() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "write access denied"})
			return
		}

		c.Next()
	}
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
		if user.Role.RoleName == "monitoringuser" {
			accessContext.PermissionType = "readonly"
		}
	}
	
	return accessContext
}