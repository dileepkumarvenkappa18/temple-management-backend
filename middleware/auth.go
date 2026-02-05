package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sharath018/temple-management-backend/config"
	"github.com/sharath018/temple-management-backend/internal/auth"
)

// AuthMiddleware handles JWT authentication and sets up access context
/*
func AuthMiddleware(cfg *config.Config, authSvc auth.Service, opt ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		optional := false
		if len(opt) > 0 {
			optional = opt[0] // if caller passes a value, use it
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
				return
			} else {
				c.Next()
				return
			}
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
				return
			} else {
				c.Next()
				return
			}
		} else if optional {
			c.Next()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTAccessSecret), nil
		})
		if err != nil || !token.Valid {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			} else {
				c.Next()
				return
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
				return
			} else {
				c.Next()
				return
			}
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id missing in token"})
				return
			} else {
				c.Next()
				return
			}
		}

		userID := uint(userIDFloat)
		user, err := authSvc.GetUserByID(userID)
		if err != nil {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				return
			} else {
				c.Next()
				return
			}
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("claims", claims)

		// Determine correct entity ID
		entityID := ResolveEntityIDForOperation(c, user, claims)

		// Create access context (now includes TenantID)
		accessContext := CreateAccessContext(c, user, claims, entityID)
		c.Set("access_context", accessContext)

		// Set resolved entity ID for quick access
		if entityID != nil {
			c.Set("entity_id", *entityID)
		}

		c.Next()
	}
}
*/

func AuthMiddleware(cfg *config.Config, authSvc auth.Service, opt ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		optional := false
		if len(opt) > 0 {
			optional = opt[0]
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
				return
			}
			fmt.Println("1Coming out here")
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
				return
			}
			fmt.Println("2Coming out here")
			c.Next()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTAccessSecret), nil
		})
		if err != nil || !token.Valid {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}
			fmt.Println("3Coming out here")
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
				return
			}
			fmt.Println("4Coming out here")
			c.Next()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id missing in token"})
				return
			}
			fmt.Println("5Coming out here")
			c.Next()
			return
		}

		userID := uint(userIDFloat)
		user, err := authSvc.GetUserByID(userID)
		if err != nil {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				return
			}
			fmt.Println("6Coming out here")
			c.Next()
			return
		}

		fmt.Println("No chance here otherwise it will panic")
		// At this point, user exists. No need for nil check for structs.
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("claims", claims)

		// Only resolve entity & access context if a valid user exists
		entityID := ResolveEntityIDForOperation(c, user, claims)
		accessContext := CreateAccessContext(c, user, claims, entityID)
		c.Set("access_context", accessContext)
		if entityID != nil {
			c.Set("entity_id", *entityID)
		}

		c.Next()
	}
}

// ResolveEntityIDForOperation determines the correct entity ID for the current operation
func ResolveEntityIDForOperation(c *gin.Context, user auth.User, claims jwt.MapClaims) *uint {
	// Priority 1: X-Entity-ID header
	if entityHeader := c.GetHeader("X-Entity-ID"); entityHeader != "" && entityHeader != "all" {
		if eid, err := strconv.ParseUint(entityHeader, 10, 32); err == nil {
			id := uint(eid)
			fmt.Printf("%s using entity ID from X-Entity-ID header: %d\n", user.Role.RoleName, id)
			return &id
		}
	}

	// Priority 2: Entity ID from URL path (/entity/123/...)
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

	// Priority 4: Role-specific fallback logic
	switch user.Role.RoleName {
	case RoleSuperAdmin:
		if tenantID := ResolveTenantIDFromRequest(c, claims); tenantID != nil {
			fmt.Printf("SuperAdmin using tenant ID as entity ID: %d\n", *tenantID)
			return tenantID
		}
		fmt.Println("SuperAdmin with global access (no specific entity)")
		return nil

	case RoleTempleAdmin:
		if user.EntityID != nil {
			fmt.Printf("TempleAdmin fallback to assigned entity ID: %d\n", *user.EntityID)
			return user.EntityID
		}

	case RoleStandardUser, RoleMonitoringUser:
		if assignedTenantIDFloat, exists := claims["assigned_tenant_id"]; exists {
			if tenantID, ok := assignedTenantIDFloat.(float64); ok && tenantID > 0 {
				id := uint(tenantID)
				fmt.Printf("%s using assigned tenant ID: %d\n", user.Role.RoleName, id)
				return &id
			}
		}
		if user.EntityID != nil {
			fmt.Printf("%s fallback to own entity ID: %d\n", user.Role.RoleName, *user.EntityID)
			return user.EntityID
		}

	case RoleDevotee, RoleVolunteer:
		if user.EntityID != nil {
			fmt.Printf("%s using own entity ID: %d\n", user.Role.RoleName, *user.EntityID)
			return user.EntityID
		}
	}

	fmt.Printf("⚠️ Could not resolve entity ID for user %d (role: %s)\n", user.ID, user.Role.RoleName)
	return nil
}

// ResolveTenantIDFromRequest extracts tenant ID from request
func ResolveTenantIDFromRequest(c *gin.Context, claims jwt.MapClaims) *uint {
	if tenantIDParam := c.Param("id"); tenantIDParam != "" && tenantIDParam != "all" {
		if tid, err := strconv.ParseUint(tenantIDParam, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}
	if tenantQuery := c.Query("tenant_id"); tenantQuery != "" && tenantQuery != "all" {
		if tid, err := strconv.ParseUint(tenantQuery, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}
	if tenantHeader := c.GetHeader("X-Tenant-ID"); tenantHeader != "" && tenantHeader != "all" {
		if tid, err := strconv.ParseUint(tenantHeader, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}
	return nil
}

// ExtractEntityIDFromPath extracts entity ID from URL
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

// CreateAccessContext creates the access context with proper entity + tenant resolution
func CreateAccessContext(c *gin.Context, user auth.User, claims jwt.MapClaims, entityID *uint) AccessContext {
	accessContext := AccessContext{
		UserID:   user.ID,
		RoleName: user.Role.RoleName,
	}

	switch user.Role.RoleName {
	case RoleSuperAdmin:
		accessContext.PermissionType = "full"
		accessContext.AssignedEntityID = entityID

	case RoleTempleAdmin:
		accessContext.PermissionType = "full"
		accessContext.DirectEntityID = user.EntityID
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

	// ✅ Extract TenantID (important for standard & monitoring users)
	// ✅ Extract TenantID (important for standard & monitoring users)
	if tenantIDFloat, ok := claims["tenant_id"].(float64); ok {
		accessContext.TenantID = uint(tenantIDFloat)
	} else if assignedTenantIDFloat, ok := claims["assigned_tenant_id"].(float64); ok {
		accessContext.TenantID = uint(assignedTenantIDFloat)
	}

	fmt.Printf("✅ AccessContext initialized: Role=%s, TenantID=%d, EntityID=%v\n",
		accessContext.RoleName, accessContext.TenantID, accessContext.AssignedEntityID)

	return accessContext
}
