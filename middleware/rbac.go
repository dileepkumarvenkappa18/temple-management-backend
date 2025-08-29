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

		// Always set both user and userID in context for downstream handlers
		c.Set("user", user)
		c.Set("userID", user.ID)

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
		
		// IMPORTANT: SuperAdmin handling - always has full access
		if user.Role.RoleName == RoleSuperAdmin {
			// Extract entity ID from multiple possible sources
			var entityIDUint *uint
			
			// 1. URL path parameters
			params := c.Params
			for _, param := range params {
				if strings.Contains(param.Key, "id") || 
				   strings.Contains(param.Key, "entity") || 
				   strings.Contains(param.Key, "tenant") {
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
				AssignedEntityID: entityIDUint,
				PermissionType:   "full",
			}
			c.Set("access_context", accessContext)
			c.Next()
			return
		}
		
		// For templeadmin - ALWAYS grant full access to their own temple
		if user.Role.RoleName == RoleTempleAdmin {
			entityIDStr := c.Param("id")
			
			if entityIDStr == "all" {
				accessContext := AccessContext{
					UserID:           user.ID,
					RoleName:         RoleTempleAdmin,
					DirectEntityID:   user.EntityID,
					AssignedEntityID: nil,
					PermissionType:   "full",
				}
				c.Set("access_context", accessContext)
				c.Next()
				return
			}
			
			entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
			if err != nil {
				if user.EntityID != nil {
					accessContext := AccessContext{
						UserID:           user.ID,
						RoleName:         RoleTempleAdmin,
						DirectEntityID:   user.EntityID,
						AssignedEntityID: nil,
						PermissionType:   "full",
					}
					c.Set("access_context", accessContext)
					c.Next()
					return
				} else {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid entity ID"})
					return
				}
			}
			
			entityIDUint := uint(entityID)
			accessContext := AccessContext{
				UserID:           user.ID,
				RoleName:         RoleTempleAdmin,
				DirectEntityID:   user.EntityID,
				AssignedEntityID: &entityIDUint,
				PermissionType:   "full",
			}
			c.Set("access_context", accessContext)
			c.Next()
			return
		}
		
		// For standarduser, monitoringuser
		if user.Role.RoleName == RoleStandardUser || 
		   user.Role.RoleName == RoleMonitoringUser {
			entityIDStr := c.Param("id")
			
			if entityIDStr == "all" {
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
			
			entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid entity ID"})
				return
			}
			
			permType := "full"
			if user.Role.RoleName == RoleMonitoringUser {
				permType = "readonly"
			}
			
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

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "temple access required"})
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
