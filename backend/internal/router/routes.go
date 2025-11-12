package router

import (
	"backend/config"
	"backend/internal/app/handlers"
	"backend/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func Register(
	r *gin.Engine,
	cfg *config.Config,
	authH *handlers.AuthHandler,
	opsH *handlers.OpsRequestHandler,
	approvalH *handlers.ApprovalHandler) {

	api := r.Group("/api")

	// ==== AUTH ==== //
	auth := api.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", authH.Logout)
		auth.GET("/me", middleware.JWTAuth(cfg), authH.Me)
	}

	// ==== OPS ==== //
	ops := api.Group("/ops")
	ops.Use(middleware.JWTAuth(cfg))
	{
		ops.POST("", opsH.CreateOpsRequest)
		ops.GET("", opsH.ListOpsRequests)
		ops.GET("/:id", opsH.GetOpsByRequestByID)
		ops.PUT("/:id", opsH.UpdateOpsRequest)
		ops.DELETE("/:id", opsH.DeleteOpsRequest)
	}

	// ==== APPROVE ==== //
	approve := api.Group("/approve")
	approve.Use(middleware.JWTAuth(cfg), middleware.RoleAllowed("admin"))
	{
		approve.POST("/:id", approvalH.HandleApproval)
		approve.POST("/:request_id/reject", approvalH.HandleApproval)
	}
}
