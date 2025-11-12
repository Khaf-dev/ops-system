package middleware

import (
	"backend/internal/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleAllowed(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "role missing")
			c.Abort()
			return
		}
		role := r.(string)
		for _, rr := range roles {
			if role == rr {
				c.Next()
				return
			}
		}
		utils.ErrorResponse(c, http.StatusForbidden, "Forbidden!")
		c.Abort()
	}
}
