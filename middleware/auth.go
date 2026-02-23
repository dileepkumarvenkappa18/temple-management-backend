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
func AuthMiddleware(cfg *config.Config, authSvc auth.Service, authRepo auth.Repository, opt ...bool) gin.HandlerFunc {
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
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
				return
			}
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
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
				return
			}
			c.Next()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			if !optional {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user_id missing in token"})
				return
			}
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
			c.Next()
			return
		}

		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("claims", claims)

		// Resolve entity & access context
		entityID := ResolveEntityIDForOperation(c, user, claims)
		accessContext := CreateAccessContext(c, user, claims, entityID, authRepo)
		c.Set("access_context", accessContext)
		if entityID != nil {
			c.Set("entity_id", *entityID)
		}

		c.Next()
	}
}

// ResolveEntityIDForOperation determines the correct entity ID for the current operation
func ResolveEntityIDForOperation(c *gin.Context, user auth.User, claims jwt.MapClaims) *uint {

	// Priority 1: Entity ID from request header X-Entity-ID
	if entityHeader := c.GetHeader("X-Entity-ID"); entityHeader != "" && entityHeader != "all" {
		if eid, err := strconv.ParseUint(entityHeader, 10, 32); err == nil {
			id := uint(eid)
			return &id
		}
	}

	// Priority 2: Entity ID from URL path (/entity/123/... OR /entities/123/...)
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
		if entityQuery := c.Query("entity_id"); entityQuery != "" {
			if eid, err := strconv.ParseUint(entityQuery, 10, 32); err == nil {
				id := uint(eid)
				return &id
			}
		}
		fmt.Println("SuperAdmin with global access (no specific entity)")
		return nil

	case RoleTempleAdmin:
		if user.EntityID != nil {
			fmt.Printf("TempleAdmin fallback to assigned entity ID: %d\n", *user.EntityID)
			return user.EntityID
		}

	case RoleStandardUser, RoleMonitoringUser:
		// DO NOT use assigned_tenant_id as entity ID.
		// assigned_tenant_id is stored separately in accessContext.TenantID.
		// Fall back to user's own EntityID only if nothing found in URL/header/query.
		if user.EntityID != nil {
			fmt.Printf("%s fallback to own entity ID: %d\n", user.Role.RoleName, *user.EntityID)
			return user.EntityID
		}
		// No entity assigned - return nil, access will be checked via TenantID in accessContext
		return nil

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
		fmt.Println("Hope we don't hit this.. A problem for SuperAdmin. Don't like headers.")
		if tid, err := strconv.ParseUint(tenantHeader, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}
	return nil
}

// ExtractEntityIDFromPath extracts entity ID from URL path.
// FIXED: Now handles both /entity/123 and /entities/123 patterns.
func ExtractEntityIDFromPath(c *gin.Context) *uint {
	path := c.Request.URL.Path
	parts := strings.Split(path, "/")
	for i, part := range parts {
		// Match both singular "/entity/123" and plural "/entities/123"
		if (part == "entity" || part == "entities") && i+1 < len(parts) {
			entityIDFromPath, err := strconv.ParseUint(parts[i+1], 10, 32)
			if err == nil {
				uintID := uint(entityIDFromPath)
				fmt.Printf("Extracted entity ID %d from URL path: %s\n", uintID, path)
				return &uintID
			}
		}
	}
	return nil
}

// CreateAccessContext creates the access context with proper entity + tenant resolution
func CreateAccessContext(c *gin.Context, user auth.User, claims jwt.MapClaims, entityID *uint, authRepo auth.Repository) AccessContext {
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
		accessContext.TenantID = user.ID

	case RoleStandardUser:
		accessContext.PermissionType = "full"
		accessContext.DirectEntityID = user.EntityID
		// Use URL entity ID if provided, otherwise fall back to user's own entity
		if entityID != nil {
			accessContext.AssignedEntityID = entityID
		} else {
			accessContext.AssignedEntityID = user.EntityID
		}
		// Always resolve TenantID from DB for standarduser
		if authRepo != nil {
			if assignedTenantID, err := authRepo.GetAssignedTenantID(user.ID); err == nil && assignedTenantID != nil {
				accessContext.TenantID = *assignedTenantID
				fmt.Printf("✅ StandardUser %d assigned to tenant: %d\n", user.ID, accessContext.TenantID)
			} else {
				fmt.Printf("⚠️ No tenant assigned to StandardUser %d\n", user.ID)
			}
		}

	case RoleMonitoringUser:
		accessContext.PermissionType = "readonly"
		accessContext.DirectEntityID = user.EntityID
		if entityID != nil {
			accessContext.AssignedEntityID = entityID
		} else {
			accessContext.AssignedEntityID = user.EntityID
		}
		// Always resolve TenantID from DB for monitoringuser
		if authRepo != nil {
			if assignedTenantID, err := authRepo.GetAssignedTenantID(user.ID); err == nil && assignedTenantID != nil {
				accessContext.TenantID = *assignedTenantID
				fmt.Printf("✅ MonitoringUser %d assigned to tenant: %d\n", user.ID, accessContext.TenantID)
			} else {
				fmt.Printf("⚠️ No tenant assigned to MonitoringUser %d\n", user.ID)
			}
		}

	case RoleDevotee, RoleVolunteer:
		accessContext.PermissionType = "readonly"
		if entityID != nil {
			accessContext.AssignedEntityID = entityID
		} else {
			accessContext.DirectEntityID = user.EntityID
		}
	}

	fmt.Printf("✅ AccessContext initialized: Role=%s, TenantID=%d, EntityID=%v\n",
		accessContext.RoleName, accessContext.TenantID, accessContext.AssignedEntityID)

	return accessContext
}

// GetTenantIDFromAccessContext is a helper to get tenant ID from the access context
func GetTenantIDFromAccessContext(c *gin.Context) uint {
	if accessCtxVal, exists := c.Get("access_context"); exists {
		if accessCtx, ok := accessCtxVal.(AccessContext); ok {
			return accessCtx.TenantID
		}
	}
	return 0
}