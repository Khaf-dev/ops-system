package middleware

import (
	"backend/config"
	"backend/internal/app/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "missing authorization")
			c.Abort()
			return
		}
		parts := strings.SplitN(auth, "", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(cfg, parts[1])
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token!")
			c.Abort()
			return
		}
		// set claims to context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
