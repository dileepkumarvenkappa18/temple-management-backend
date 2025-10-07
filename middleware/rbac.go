package middleware

import (
	"net/http"
<<<<<<< HEAD
	"strings"
	
=======
	"strconv"
	"strings"

>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	"github.com/gin-gonic/gin"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// RBACMiddleware checks if the user has one of the allowed roles
func RBACMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
<<<<<<< HEAD
		// Special case for standard users accessing entities endpoint
		if c.Request.URL.Path == "/api/v1/entities" && 
		   (c.Request.Method == "GET" || c.Request.Method == "POST") {
			userVal, exists := c.Get("user")
			if exists {
				if user, ok := userVal.(auth.User); ok && 
				   (user.Role.RoleName == "standarduser" || user.Role.RoleName == "monitoringuser") {
					// Allow standard users to access this endpoint for both GET and POST
					c.Next()
					return
				}
			}
		}

=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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

<<<<<<< HEAD
		// Always set both user and userID in context for downstream handlers
		c.Set("user", user)
		c.Set("userID", user.ID)

=======
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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

<<<<<<< HEAD
// RequireTempleAccess with proper tenant isolation
=======
// RequireTempleAccess ensures user has access to the temple entity
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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
		
<<<<<<< HEAD
		// Get access context from auth middleware
		accessContextVal, exists := c.Get("access_context")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access context missing"})
			return
		}
		
		accessContext, ok := accessContextVal.(AccessContext)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access context"})
			return
		}
		
		// Check if this is a tenant user management endpoint
		// Allow these endpoints even if no temple exists
		if strings.Contains(c.Request.URL.Path, "/tenants/") && strings.Contains(c.Request.URL.Path, "/user") {
			c.Next()
			return
		}
		
		// FIXED: Role-based access control with tenant isolation
		switch user.Role.RoleName {
		case RoleSuperAdmin:
			// Superadmin can access any entity, but should be scoped to requested tenant
			c.Next()
			return
			
		case RoleTempleAdmin:
			// Temple admin can only access their own temple and related entities
			if accessContext.DirectEntityID == nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "templeadmin must have a direct entity assigned",
				})
				return
			}
			c.Next()
			return
			
		case RoleStandardUser, RoleMonitoringUser:
			// Standard/monitoring users can only access their assigned tenant
			if accessContext.AssignedEntityID == nil && accessContext.DirectEntityID == nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "user must have an assigned entity",
				})
				return
			}
			c.Next()
			return
		
		case RoleDevotee, RoleVolunteer:
			// Devotees and volunteers can access their associated temple
			if accessContext.AssignedEntityID == nil && accessContext.DirectEntityID == nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "devotee/volunteer must have an associated entity",
				})
				return
			}
			
			// Set permission type to readonly for regular endpoints
			// This ensures devotees can view but not modify temple data
			accessContext.PermissionType = "readonly"
			c.Set("access_context", accessContext)
			
			c.Next()
			return
			
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unsupported role"})
			return
		}
=======
		// IMPORTANT: SuperAdmin handling - always has full access
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
		
		// For templeadmin - ALWAYS grant full access to their own temple
		if user.Role.RoleName == RoleTempleAdmin {
			// Get the entity ID from the URL parameter
			entityIDStr := c.Param("id")
			
			// Special case for "all" - handled in handler
			if entityIDStr == "all" {
				accessContext := AccessContext{
					UserID:           user.ID,
					RoleName:         RoleTempleAdmin,
					DirectEntityID:   user.EntityID,
					AssignedEntityID: nil,
					PermissionType:   "full", // Always full for templeadmin
				}
				c.Set("access_context", accessContext)
				c.Next()
				return
			}
			
			// Parse the entity ID
			entityID, err := strconv.ParseUint(entityIDStr, 10, 64)
			if err != nil {
				// If not valid ID in URL, default to user's own entity
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
			
			// Create an access context with the entity ID
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
		
		// For standarduser, monitoringuser, check entity access
		if user.Role.RoleName == RoleStandardUser || 
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
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
	}
}

// RequireWriteAccess ensures user has write access
func RequireWriteAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
<<<<<<< HEAD
=======
		// Allow superadmin and templeadmin to bypass write access check
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
		userVal, exists := c.Get("user")
		if exists {
			user, ok := userVal.(auth.User)
			if ok {
<<<<<<< HEAD
				// Superadmin and templeadmin always have write access
=======
				// Both superadmin and templeadmin always have write access
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
				if user.Role.RoleName == RoleSuperAdmin || user.Role.RoleName == RoleTempleAdmin {
					c.Next()
					return
				}
			}
		}
		
<<<<<<< HEAD
=======
		// Continue with regular access context check for other roles
>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
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
<<<<<<< HEAD
}
=======
}

>>>>>>> 94687f1f9b610a9b6c08378c7d37e9a6b831dbf6
