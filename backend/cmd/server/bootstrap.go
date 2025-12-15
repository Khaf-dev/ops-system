package server

import (
	"backend/config"
	"backend/internal/app/handlers"
	"backend/internal/app/logic"
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

	// ==============================================
	// Repository Layer
	// ==============================================
	opsRepo := repository.NewOpsRequestRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)

	userRepo := repository.NewUserRepository(db)
	levelRepo := repository.NewLevelRepository(db)
	reqTypeRepo := repository.NewRequestTypeRepository(db)

	approverCfgRepo := repository.NewApproverConfigRepository(db)
	approvalFlowRepo := repository.NewApprovalFlowRepository(db)
	approvalStepRepo := repository.NewApprovalStepRepository(db)
	approvalLogRepo := repository.NewApprovalLogRepository(db)

	// ==============================================
	// Logic Layer (pure)
	// ==============================================
	approvalLogic := logic.NewApprovalLogic()

	// ==============================================
	// Service Layer
	// ==============================================
	authSvc := services.NewAuthService(db, cfg)

	opsSvc := services.NewOpsRequestService(
		opsRepo,
	)

	attachmentSvc := services.NewAttachmentService(
		attachmentRepo,
	)
	approvalFlowSvc := services.NewApprovalFlowService(
		db,
		opsRepo,
		approvalFlowRepo,
		approvalStepRepo,
		approverCfgRepo,
		approvalLogic,
	)

	approvalActionSvc := services.NewApprovalActionService(
		db,
		approvalFlowRepo,
		approvalStepRepo,
		opsRepo,
		approvalLogRepo,
		userRepo,
		approvalLogic,
	)

	approvalSvc := services.NewApprovalService(
		db,
		opsRepo,
		approvalFlowRepo,
		approverCfgRepo,
		approvalStepRepo,
		approvalLogRepo,
		userRepo,
		approvalLogic,
	)

	adminSvc := services.NewAdminService(
		userRepo,
		levelRepo,
		reqTypeRepo,
	)

	reqTypeSvc := services.NewRequestTypeService(reqTypeRepo)
	levelSvc := services.NewLevelService(levelRepo)
	userSvc := services.NewUserService(userRepo)

	// ==============================================
	// Handler Layer
	// ==============================================
	authH := handlers.NewAuthHandler(db, authSvc, cfg)
	opsH := handlers.NewOpsRequestHandler(opsSvc)
	attachmentH := handlers.NewAttachmentHandler(attachmentSvc)
	approvalH := handlers.NewApprovalHandler(approvalFlowSvc, approvalActionSvc)
	adminH := handlers.NewAdminHandler(adminSvc, levelSvc, userSvc)
	adminReqTypeH := handlers.NewAdminRequestTypeHandler(reqTypeSvc)

	// ==============================================
	// Router
	// ==============================================
	router.Register(
		r,
		cfg,
		authH,
		opsH,
		approvalH,
		attachmentH,
		adminH,
		adminReqTypeH,
	)

	return &App{
		Config: cfg,
		Router: r,
	}
}
