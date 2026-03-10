package middleware

import (
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// RBACMiddleware checks if the user has one of the allowed roles
func RBACMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Special case for standard users accessing entities endpoint
		if c.Request.URL.Path == "/api/v1/entities" &&
			(c.Request.Method == "GET" || c.Request.Method == "POST") {
			userVal, exists := c.Get("user")
			if exists {
				if user, ok := userVal.(auth.User); ok &&
					(user.Role.RoleName == "standarduser" || user.Role.RoleName == "monitoringuser") {
					c.Next()
					return
				}
			}
		}

		userVal, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

		user, ok := userVal.(auth.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user object"})
			return
		}

		c.Set("user", user)
		c.Set("userID", user.ID)

		for _, role := range allowedRoles {
			if user.Role.RoleName == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
	}
}

// RequireTempleAccess ensures the user has access to a temple/entity.
//
// For templeadmin: checks DirectEntityID != nil (set to user.EntityID in CreateAccessContext).
// The actual per-seva ownership check happens later in canAccessSeva via TenantID.
func RequireTempleAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}

		user, ok := userVal.(auth.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user object"})
			return
		}

		accessContextVal, exists := c.Get("access_context")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "access context missing",
				"user_id":     user.ID,
				"user_role":   user.Role.RoleName,
				"url_path":    c.Request.URL.Path,
				"stack_trace": string(debug.Stack()),
			})
			return
		}

		accessContext, ok := accessContextVal.(AccessContext)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "invalid access context",
				"user_id":     user.ID,
				"user_role":   user.Role.RoleName,
				"url_path":    c.Request.URL.Path,
				"stack_trace": string(debug.Stack()),
			})
			return
		}

		// Allow tenant user management endpoints regardless of temple assignment
		if strings.Contains(c.Request.URL.Path, "/tenants/") && strings.Contains(c.Request.URL.Path, "/user") {
			c.Next()
			return
		}

		switch user.Role.RoleName {
		case RoleSuperAdmin:
			c.Next()
			return

		case RoleTempleAdmin:
			// DirectEntityID is always set to user.EntityID in CreateAccessContext,
			// so this gate will pass as long as the templeadmin has any entity on their record.
			// The per-resource ownership check (which temple they manage) is done in
			// canAccessSeva / similar handlers via TenantID lookup.
			if accessContext.DirectEntityID == nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error":          "templeadmin must have a direct entity assigned",
					"user_id":        user.ID,
					"user_role":      user.Role.RoleName,
					"access_context": accessContext,
					"url_path":       c.Request.URL.Path,
					"stack_trace":    string(debug.Stack()),
				})
				return
			}
			c.Next()
			return

		case RoleStandardUser, RoleMonitoringUser:
			if accessContext.TenantID == 0 && accessContext.DirectEntityID == nil && accessContext.AssignedEntityID == nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error":          "user must have an assigned entity",
					"user_id":        user.ID,
					"user_role":      user.Role.RoleName,
					"access_context": accessContext,
					"url_path":       c.Request.URL.Path,
					"stack_trace":    string(debug.Stack()),
				})
				return
			}
			c.Next()
			return

		case RoleDevotee, RoleVolunteer:
			if accessContext.TenantID == 0 && accessContext.DirectEntityID == nil && accessContext.AssignedEntityID == nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error":          "devotee/volunteer must have an associated entity",
					"user_id":        user.ID,
					"user_role":      user.Role.RoleName,
					"access_context": accessContext,
					"url_path":       c.Request.URL.Path,
					"stack_trace":    string(debug.Stack()),
				})
				return
			}
			accessContext.PermissionType = "readonly"
			c.Set("access_context", accessContext)
			c.Next()
			return

		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":       "unsupported role",
				"user_id":     user.ID,
				"user_role":   user.Role.RoleName,
				"url_path":    c.Request.URL.Path,
				"stack_trace": string(debug.Stack()),
			})
			return
		}
	}
}

// RequireWriteAccess ensures user has write access
func RequireWriteAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if exists {
			user, ok := userVal.(auth.User)
			if ok {
				if user.Role.RoleName == RoleSuperAdmin || user.Role.RoleName == RoleTempleAdmin {
					c.Next()
					return
				}
			}
		}

		accessContext, exists := c.Get("access_context")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
			return
		}

		ctx, ok := accessContext.(AccessContext)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access context"})
			return
		}

		if !ctx.CanWrite() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "write access denied"})
			return
		}

		c.Next()
	}
}