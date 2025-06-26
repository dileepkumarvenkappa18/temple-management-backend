package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RBACMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

		tokenClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		roleIDFloat, ok := tokenClaims["role_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role missing in token"})
			return
		}

		roleMap := map[int]string{
			1: "superadmin",
			2: "templeadmin",
			3: "devotee",
			4: "volunteer",
		}
		roleName := roleMap[int(roleIDFloat)]

		for _, allowed := range allowedRoles {
			if strings.ToLower(allowed) == roleName {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
	}
}
