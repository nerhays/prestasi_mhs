package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nerhays/prestasi_uas/utils"
)

const (
	ContextUserIDKey      = "userID"
	ContextUsernameKey    = "username"
	ContextRoleKey        = "role"
	ContextPermissionsKey = "permissions"
)

// AuthMiddleware: cek header Authorization: Bearer <token>
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			c.Abort()
			return
		}

		// simpan info user ke context
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextRoleKey, claims.Role)
		c.Set(ContextPermissionsKey, claims.Permissions)

		c.Next()
	}
}

// RequireRole: cek role ("Admin", "Mahasiswa", "Dosen Wali")
func RequireRole(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		roleAny, ok := c.Get(ContextRoleKey)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "role not found in context"})
			c.Abort()
			return
		}

		role, ok := roleAny.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "invalid role type"})
			c.Abort()
			return
		}

		if _, ok := roleSet[role]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "access denied - insufficient role"})
			c.Abort()
			return
		}

		c.Next()
	}
}
