package middleware

import (
	"net/http"
	"strconv"
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

		// FIXED: Proper tenant resolution
		assignedTenantID := ResolveTenantID(c, user, claims)
		
		// Create access context
		accessContext := ResolveAccessContext(user, assignedTenantID)
		c.Set("access_context", accessContext)

		c.Next()
	}
}

// FIXED: Simplified tenant ID resolution
func ResolveTenantID(c *gin.Context, user auth.User, claims jwt.MapClaims) *uint {
	// Priority 1: URL parameter (for specific tenant access)
	if tenantIDParam := c.Param("id"); tenantIDParam != "" && tenantIDParam != "all" {
		if tid, err := strconv.ParseUint(tenantIDParam, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}

	// Priority 2: Query parameter
	if tenantQuery := c.Query("tenant_id"); tenantQuery != "" && tenantQuery != "all" {
		if tid, err := strconv.ParseUint(tenantQuery, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}

	// Priority 3: Header (X-Tenant-ID)
	if tenantHeader := c.GetHeader("X-Tenant-ID"); tenantHeader != "" && tenantHeader != "all" {
		if tid, err := strconv.ParseUint(tenantHeader, 10, 32); err == nil {
			id := uint(tid)
			return &id
		}
	}

	// Priority 4: JWT token assigned_tenant_id
	if assignedTenantIDFloat, exists := claims["assigned_tenant_id"]; exists {
		if tenantID, ok := assignedTenantIDFloat.(float64); ok && tenantID > 0 {
			id := uint(tenantID)
			return &id
		}
	}

	// For standarduser and monitoringuser, they must have an assigned tenant
	if user.Role.RoleName == "standarduser" || user.Role.RoleName == "monitoringuser" {
		// Return their own entity ID as the assigned tenant
		return user.EntityID
	}

	return nil
}