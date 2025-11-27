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
	db := database.Connect(cfg)

	limiter := utils.NewLoginLimiter(5, time.Minute*5, time.Minute*15)
	limiter.StartCleanup(time.Minute)

	r := gin.Default()
	r.Use(middleware.CORS())

	// repositori
	opsRepo := repository.NewOpsRequestRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	userRepo := repository.NewUserRepository(db)
	approvalRepo := repository.NewApprovalRepository(db)
	ReqTypeRepo := repository.NewRequestTypeRepository(db)
	levelRepo := repository.NewLevelRepository(db)

	// service
	authSvc := services.NewAuthService(db, cfg)
	opsSvc := services.NewOpsRequestService(opsRepo)
	attachmentSvc := services.NewAttachmentService(attachmentRepo)
	approvalSvc := services.NewApprovalService(db, opsRepo, approvalRepo, userRepo, ReqTypeRepo)
	adminSvc := services.NewAdminService(userRepo, levelRepo, ReqTypeRepo)
	reqTypeSvc := services.NewRequestTypeService(ReqTypeRepo)
	levelSvc := services.NewLevelService(levelRepo)
	userSvc := services.NewUserService(userRepo)

	// handler
	authH := handlers.NewAuthHandler(db, authSvc, cfg)
	opsH := handlers.NewOpsRequestHandler(opsSvc, approvalSvc)
	attachmentHandler := handlers.NewAttachmentHandler(attachmentSvc)
	approvalH := handlers.NewApprovalHandler(approvalSvc)
	adminH := handlers.NewAdminHandler(adminSvc, levelSvc, userSvc)
	adminReqTypeHandler := handlers.NewAdminRequestTypeHandler(reqTypeSvc)

	router.Register(r, cfg, authH, opsH, approvalH, adminH, attachmentHandler, adminReqTypeHandler)

	return &App{
		Config: cfg,
		Router: r,
	}
}
