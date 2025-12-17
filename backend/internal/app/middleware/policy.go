package middleware

import (
	"backend/internal/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PolicyContext struct {
	UserID uuid.UUID
	Role   string
	Group  []string
}

// RequirePolicy checks permission based on policy string
func RequirePolicy(policy string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ========= ambil context auth ========= //
		uidRaw, ok := c.Get("user_id")
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "user_id missing")
			c.Abort()
			return
		}

		roleRaw, ok := c.Get("role")
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "role missing")
			c.Abort()
			return
		}

		ctx := PolicyContext{
			UserID: uidRaw.(uuid.UUID),
			Role:   roleRaw.(string),
		}

		// ========= evaluasi policy ========= //
		if !evaluatePolicy(policy, &ctx) {
			utils.ErrorResponse(c, http.StatusForbidden, "policy denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

func evaluatePolicy(policy string, ctx *PolicyContext) bool {
	switch policy {

	// =====================================
	// Approval Policies
	// =====================================
	case "approval.view":
		return ctx.Role == "admin" ||
			ctx.Role == "finance" ||
			ctx.Role == "manager"

	case "approval.start":
		return ctx.Role == "user" ||
			ctx.Role == "admin"

	case "approval.action": // approve / reject
		return ctx.Role == "finance" ||
			ctx.Role == "manager"

		// =====================================
		// Admin Policies
		// =====================================
	case "admin.manage":
		return ctx.Role == "admin"

	default:
		return false
	}
}
