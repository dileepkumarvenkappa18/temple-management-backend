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

		// Extract tenant ID from headers (for standarduser/monitoringuser multi-tenancy)
		assignedTenantID := ExtractTenantIDFromContext(c)

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

			if entityIDStr == "all" || entityIDStr == "" {
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
				// If invalid ID, fall back to user's entity
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
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid entity ID and no assigned entity"})
					return
				}
			}

			entityIDUint := uint(entityID)

			// Check if templeadmin is accessing their own entity
			if user.EntityID != nil && *user.EntityID == entityIDUint {
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
			} else {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied to this entity"})
				return
			}
		}

		// For standarduser, monitoringuser - handle both their entity and assigned tenants
		if user.Role.RoleName == RoleStandardUser || user.Role.RoleName == RoleMonitoringUser {
			entityIDStr := c.Param("id")
			permType := "full"
			if user.Role.RoleName == RoleMonitoringUser {
				permType = "readonly"
			}

			// Handle "all" or empty entity ID - use their assigned entity or direct entity
			if entityIDStr == "all" || entityIDStr == "" {
				var accessibleEntityID *uint

				// Priority: assigned tenant > direct entity
				if assignedTenantID != nil {
					accessibleEntityID = assignedTenantID
				} else if user.EntityID != nil {
					accessibleEntityID = user.EntityID
				}

				if accessibleEntityID == nil {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no accessible entity found"})
					return
				}

				accessContext := AccessContext{
					UserID:           user.ID,
					RoleName:         user.Role.RoleName,
					DirectEntityID:   user.EntityID,
					AssignedEntityID: accessibleEntityID,
					PermissionType:   permType,
				}
				c.Set("access_context", accessContext)
				c.Next()
				return
			}

			// Handle specific entity ID
			entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid entity ID"})
				return
			}

			entityIDUint := uint(entityID)

			// Check if user has access to this specific entity
			hasAccess := false

			// Check against assigned tenant
			if assignedTenantID != nil && *assignedTenantID == entityIDUint {
				hasAccess = true
			}

			// Check against direct entity
			if user.EntityID != nil && *user.EntityID == entityIDUint {
				hasAccess = true
			}

			if !hasAccess {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied to this entity"})
				return
			}

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
				if user.Role.RoleName == RoleSuperAdmin || user.Role.RoleName == RoleTempleAdmin || user.Role.RoleName == RoleStandardUser {
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