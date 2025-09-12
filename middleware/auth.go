
package middleware

import (
	"net/http"
	"strconv"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// AuthMiddleware handles JWT authentication and sets up access context
func AuthMiddleware(cfg *config.Config, authSvc auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTAccessSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id missing in token"})
			return
		}

		userID := uint(userIDFloat)
		user, err := authSvc.GetUserByID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("claims", claims)

		// FIXED: Comprehensive entity ID resolution - TempleAdmin can access requested entities
		entityID := ResolveEntityIDForOperation(c, user, claims)
		
		// Create access context with proper entity ID
		accessContext := CreateAccessContext(c, user, claims, entityID)
		c.Set("access_context", accessContext)

		// Also set the resolved entity_id for easy access
		if entityID != nil {
			c.Set("entity_id", *entityID)
		}

		c.Next()
	}
}

// ResolveEntityIDForOperation determines the correct entity ID for the current operation
func ResolveEntityIDForOperation(c *gin.Context, user auth.User, claims jwt.MapClaims) *uint {
	// Universal priority order for all roles (except specific restrictions)
	
	// Priority 1: X-Entity-ID header (most explicit)
	if entityHeader := c.GetHeader("X-Entity-ID"); entityHeader != "" && entityHeader != "all" {
		if eid, err := strconv.ParseUint(entityHeader, 10, 32); err == nil {
			id := uint(eid)
			fmt.Printf("%s using entity ID from X-Entity-ID header: %d\n", user.Role.RoleName, id)
			return &id
		}
	}

	// Priority 2: Entity ID from URL path (e.g., /entity/123/...)
	if entityIDFromPath := ExtractEntityIDFromPath(c); entityIDFromPath != nil {
		fmt.Printf("%s using entity ID from URL path: %d\n", user.Role.RoleName, *entityIDFromPath)
		return entityIDFromPath
	}

	// Priority 3: Query parameter entity_id
	if entityQuery := c.Query("entity_id"); entityQuery != "" && entityQuery != "all" {
		if eid, err := strconv.ParseUint(entityQuery, 10, 32); err == nil {
			id := uint(eid)
			fmt.Printf("%s using entity ID from query parameter: %d\n", user.Role.RoleName, id)
			return &id
		}
	}

	// Priority 4: Role-specific fallbacks
	switch user.Role.RoleName {
	case RoleSuperAdmin:
		// Try X-Tenant-ID header or tenant param for multi-tenant operations
		if tenantID := ResolveTenantIDFromRequest(c, claims); tenantID != nil {
			fmt.Printf("SuperAdmin using tenant ID as entity ID: %d\n", *tenantID)
			return tenantID
		}
		// SuperAdmin without specific entity - return nil for global access
		fmt.Printf("SuperAdmin with global access (no specific entity)\n")
		return nil

	case RoleTempleAdmin:
		// TempleAdmin: Fallback to their assigned entity if no specific request
		if user.EntityID != nil {
			fmt.Printf("TempleAdmin fallback to assigned entity ID: %d\n", *user.EntityID)
			return user.EntityID
		}

	case RoleStandardUser, RoleMonitoringUser:
		// Use assigned tenant ID from claims
		if assignedTenantIDFloat, exists := claims["assigned_tenant_id"]; exists {
			if tenantID, ok := assignedTenantIDFloat.(float64); ok && tenantID > 0 {
				id := uint(tenantID)
				fmt.Printf("%s using assigned tenant ID: %d\n", user.Role.RoleName, id)
				return &id
			}
		}
		// Fallback to their own entity ID if no assigned tenant
		if user.EntityID != nil {
			fmt.Printf("%s fallback to own entity ID: %d\n", user.Role.RoleName, *user.EntityID)
			return user.EntityID
		}

	case RoleDevotee, RoleVolunteer:
		// Fallback to user's own entity
		if user.EntityID != nil {
			fmt.Printf("%s using own entity ID: %d\n", user.Role.RoleName, *user.EntityID)
			return user.EntityID
		}
	}

	fmt.Printf("WARNING: Could not resolve entity ID for user %d (role: %s)\n", user.ID, user.Role.RoleName)
	return nil
}

// ResolveTenantIDFromRequest extracts tenant ID from request for multi-tenant operations
func ResolveTenantIDFromRequest(c *gin.Context, claims jwt.MapClaims) *uint {
	// Try URL parameter
	if tenantIDParam := c.Param("id"); tenantIDParam != "" && tenantIDParam != "all" {
		if tid, err := strconv.ParseUint(tenantIDParam, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}

	// Try query parameter
	if tenantQuery := c.Query("tenant_id"); tenantQuery != "" && tenantQuery != "all" {
		if tid, err := strconv.ParseUint(tenantQuery, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}

	// Try header
	if tenantHeader := c.GetHeader("X-Tenant-ID"); tenantHeader != "" && tenantHeader != "all" {
		if tid, err := strconv.ParseUint(tenantHeader, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}

	return nil
}

// ExtractEntityIDFromPath attempts to extract an entity ID from URL paths like /entity/123/...
func ExtractEntityIDFromPath(c *gin.Context) *uint {
	path := c.Request.URL.Path
	if strings.Contains(path, "/entity/") {
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if part == "entity" && i+1 < len(parts) {
				entityIDFromPath, err := strconv.ParseUint(parts[i+1], 10, 32)
				if err == nil {
					uintID := uint(entityIDFromPath)
					return &uintID
				}
			}
		}
	}
	return nil
}

// CreateAccessContext creates the access context with proper entity resolution
func CreateAccessContext(c *gin.Context, user auth.User, claims jwt.MapClaims, entityID *uint) AccessContext {
	accessContext := AccessContext{
		UserID:   user.ID,
		RoleName: user.Role.RoleName,
	}

	switch user.Role.RoleName {
	case RoleSuperAdmin:
		accessContext.PermissionType = "full"
		accessContext.AssignedEntityID = entityID // Could be nil for global access

	case RoleTempleAdmin:
		accessContext.PermissionType = "full"
		accessContext.DirectEntityID = user.EntityID
		// Use the resolved entity ID (can be requested entity or assigned entity)
		accessContext.AssignedEntityID = entityID

	case RoleStandardUser:
		accessContext.PermissionType = "full"
		accessContext.AssignedEntityID = entityID

	case RoleMonitoringUser:
		accessContext.PermissionType = "readonly"
		accessContext.AssignedEntityID = entityID

	case RoleDevotee, RoleVolunteer:
		accessContext.PermissionType = "readonly"
		if entityID != nil {
			accessContext.AssignedEntityID = entityID
		} else {
			accessContext.DirectEntityID = user.EntityID
		}
	}

	return accessContext
}


