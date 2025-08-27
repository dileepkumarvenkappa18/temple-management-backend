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
		c.Set("claims", claims)
		
		// Special handling for "all" entity parameter in reports URLs
		var assignedTenantID *uint
		if user.RoleID == 1 && strings.Contains(c.Request.URL.Path, "/entities/all/reports/") {
			// Superadmin viewing all entities - use their own tenant ID
			id := user.ID
			assignedTenantID = &id
			
			// Check if there's a specific tenant_id in query params that should override
			if tenantQuery := c.Query("tenant_id"); tenantQuery != "" && tenantQuery != "all" {
				if tid, err := strconv.ParseUint(tenantQuery, 10, 32); err == nil {
					id := uint(tid)
					assignedTenantID = &id
				}
			}
		} else {
			// Regular handling for all other cases
			
			// Check URL for tenant_id or id parameter
			if tenantIDParam := c.Param("id"); tenantIDParam != "" && tenantIDParam != "all" {
				if tid, err := strconv.ParseUint(tenantIDParam, 10, 32); err == nil {
					id := uint(tid)
					assignedTenantID = &id
				}
			} else if tenantQuery := c.Query("tenant_id"); tenantQuery != "" && tenantQuery != "all" {
				if tid, err := strconv.ParseUint(tenantQuery, 10, 32); err == nil {
					id := uint(tid)
					assignedTenantID = &id
				}
			} else if tenantsQuery := c.Query("tenants"); tenantsQuery != "" {
				tenantIDs := strings.Split(tenantsQuery, ",")
				if len(tenantIDs) > 0 && tenantIDs[0] != "" && tenantIDs[0] != "all" {
					if tid, err := strconv.ParseUint(tenantIDs[0], 10, 32); err == nil {
						id := uint(tid)
						assignedTenantID = &id
					}
				}
			}
			
			// If no tenant ID from URL, check claims
			if assignedTenantID == nil {
				if assignedTenantIDFloat, exists := claims["assigned_tenant_id"]; exists {
					if tenantID, ok := assignedTenantIDFloat.(float64); ok && tenantID > 0 {
						id := uint(tenantID)
						assignedTenantID = &id
					}
				}
			}
		}
		
		// Create and set access context
		accessContext := ResolveAccessContext(user, assignedTenantID)
		c.Set("access_context", accessContext)

		c.Next()
	}
}