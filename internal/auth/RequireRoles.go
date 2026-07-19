package auth

import (
	"net/http"

	"github.com/andrebarone77/cardiaflow-api/internal/domain"
	"github.com/gin-gonic/gin"
)

func RequireRoles(roles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": domain.ErrCouldNotProcessRequest.Error()})
			return
		}
		userRole, ok := role.(domain.Role)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": domain.ErrCouldNotProcessRequest.Error()})
			return
		}
		for _, requiredRole := range roles {
			if userRole == requiredRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": domain.ErrForbidden.Error()})

	}
}
