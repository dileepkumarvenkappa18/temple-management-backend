
package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// RBACMiddleware checks if the user has one of the allowed roles
func RBACMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Check if the user has one of the allowed roles
		for _, role := range allowedRoles {
			if user.Role.RoleName == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
	}
}

// RequireTempleAccess ensures user has access to the temple entity
func RequireTempleAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
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
		
		// IMPORTANT: Enhanced SuperAdmin handling
		if user.Role.RoleName == RoleSuperAdmin {
			// Extract entity ID from multiple possible sources
			var entityIDUint *uint
			
			// 1. URL path parameters
			// Look for any path parameter that might contain an entity ID
			params := c.Params
			for _, param := range params {
				// Check if param name contains 'id', 'entity', or 'tenant'
				if strings.Contains(param.Key, "id") || 
				   strings.Contains(param.Key, "entity") || 
				   strings.Contains(param.Key, "tenant") {
					// Try to parse it as a number
					if param.Value != "all" && param.Value != "" {
						id, err := strconv.ParseUint(param.Value, 10, 64)
						if err == nil {
							idUint := uint(id)
							entityIDUint = &idUint
							break
						}
					}
				}
			}
			
			// 2. Query parameters
			if entityIDUint == nil {
				// Check common query param names for entity IDs
				queryParams := []string{"entity_id", "entityId", "tenant_id", "tenantId", "id"}
				for _, paramName := range queryParams {
					if idStr := c.Query(paramName); idStr != "" {
						id, err := strconv.ParseUint(idStr, 10, 64)
						if err == nil {
							idUint := uint(id)
							entityIDUint = &idUint
							break
						}
					}
				}
			}
			
			// 3. Headers
			if entityIDUint == nil {
				headerParams := []string{"X-Entity-ID", "X-Tenant-ID"}
				for _, headerName := range headerParams {
					if idStr := c.GetHeader(headerName); idStr != "" {
						id, err := strconv.ParseUint(idStr, 10, 64)
						if err == nil {
							idUint := uint(id)
							entityIDUint = &idUint
							break
						}
					}
				}
			}
			
			// Create superadmin access context with the target entity ID
			accessContext := AccessContext{
				UserID:           user.ID,
				RoleName:         RoleSuperAdmin,
				DirectEntityID:   nil,
				AssignedEntityID: entityIDUint, // Give superadmin access to the requested entity
				PermissionType:   "full",
			}
			c.Set("access_context", accessContext)
			c.Next()
			return
		}
		
		// For templeadmin, standarduser, monitoringuser, check entity access
		if user.Role.RoleName == RoleTempleAdmin || 
		   user.Role.RoleName == RoleStandardUser || 
		   user.Role.RoleName == RoleMonitoringUser {
			// Get the entity ID from the URL parameter
			entityIDStr := c.Param("id")
			
			// Special case for "all" - handled in handler
			if entityIDStr == "all" {
				// Determine permission type based on role
				permType := "full"
				if user.Role.RoleName == RoleMonitoringUser {
					permType = "readonly"
				}
				
				accessContext := AccessContext{
					UserID:           user.ID,
					RoleName:         user.Role.RoleName,
					DirectEntityID:   user.EntityID,
					AssignedEntityID: nil,
					PermissionType:   permType,
				}
				c.Set("access_context", accessContext)
				c.Next()
				return
			}
			
			// Parse the entity ID
			entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid entity ID"})
				return
			}
			
			// TODO: Check if user has access to this entity
			// This would involve a database lookup to verify the user's access rights
			
			// Determine permission type based on role
			permType := "full"
			if user.Role.RoleName == RoleMonitoringUser {
				permType = "readonly"
			}
			
			// Create an access context with the entity ID
			entityIDUint := uint(entityID)
			accessContext := AccessContext{
				UserID:           user.ID,
				RoleName:         user.Role.RoleName,
				DirectEntityID:   user.EntityID,
				AssignedEntityID: &entityIDUint,
				PermissionType:   permType,
			}
			c.Set("access_context", accessContext)
			c.Next()
			return
		}

		// If we get here, the user doesn't have temple access
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "temple access required"})
	}
}

// RequireWriteAccess ensures user has write access (templeadmin or standarduser)
func RequireWriteAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow superadmin to bypass write access check
		userVal, exists := c.Get("user")
		if exists {
			user, ok := userVal.(auth.User)
			if ok && user.Role.RoleName == RoleSuperAdmin {
				c.Next()
				return
			}
		}
		
		// Continue with regular access context check
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