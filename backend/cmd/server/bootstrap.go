package server

import (
	"backend/config"
	"backend/internal/app/handlers"
	"backend/internal/app/middleware"
	"backend/internal/app/repository"
	"backend/internal/app/services"
	"backend/internal/app/utils"
	"backend/internal/database"
	"backend/internal/router"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config *config.Config
	Router *gin.Engine
}

// SEtupApp => bootstrap full dependency
func SetupApp() *App {
	cfg := config.Load()
	database.Connect(cfg)
	db := database.InitDB(cfg)

	limiter := utils.NewLoginLimiter(5, time.Minute*5, time.Minute*15)
	limiter.StartCleanup(time.Minute)

	r := gin.Default()
	r.Use(middleware.CORS())

	authSvc := services.NewAuthService(db, cfg)
	authH := handlers.NewAuthHandler(db, authSvc, cfg)

	opsRepo := repository.NewOpsRequestRepository(db)
	opsSvc := services.NewOpsRequestService(opsRepo)
	opsH := handlers.NewOpsRequestHandler(opsSvc)

	approvalSvc := services.NewApprovalService(db)
	approvalH := handlers.NewApprovalHandler(approvalSvc)

	router.Register(r, cfg, authH, opsH, approvalH)

	return &App{
		Config: cfg,
		Router: r,
	}
}
