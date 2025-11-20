package auth

import (
	"strings"

	"github.com/TBuckholz5/workouttracker/internal/util/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		}
		authHeader = strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
		if err := jwtService.ValidateJwt(c, authHeader); err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		}

		c.Next()
	}
}
