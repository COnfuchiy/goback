package route

import (
	"github.com/gin-gonic/gin"
	"goback/api/middleware"
	"goback/bootstrap"
	"goback/mapper"
	"goback/repository"
	"goback/services"
	"gorm.io/gorm"
)

const (
	ApiRoute         = "/api/v1"
	BaseRoute        = "/"
	UserRoute        = "/user"
	WorkspaceRoute   = "/workspace"
	FileHistoryRoute = "/file-history"
	FileRoute        = "/file"
)

func Init(db *gorm.DB, gin *gin.Engine, env *bootstrap.Env) {

	userRepository := repository.NewUserRepository(db)
	workspaceRepository := repository.NewWorkspaceRepository(db)
	fileHistoryRepository := repository.NewFileHistoryRepository(db)
	fileRepository := repository.NewFileRepository(db)
	paginateRepository := repository.NewPaginateRepository(env.EntitiesPerPage)

	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(env.AccessTokenSecret, env.RefreshTokenSecret, env.AccessTokenExpiryHour, env.RefreshTokenExpiryHour)
	workspaceService := services.NewWorkspaceService(workspaceRepository, paginateRepository)
	fileHistoryService := services.NewFileHistoryService(fileHistoryRepository, paginateRepository)
	fileService := services.NewFileService(fileRepository)
	fileStorageService := services.NewFileStorageService()

	userMapper := mapper.NewUserMapper()
	workspaceMapper := mapper.NewWorkspaceMapper(userMapper)
	fileMapper := mapper.NewFileMapper(userMapper)
	fileHistoryMapper := mapper.NewFileHistoryMapper(fileMapper)
	paginateMapper := mapper.NewPaginationMapper(env.EntitiesPerPage)

	authMiddlewareHandler := middleware.NewJwtAuthMiddleware(authService, userService).Handle()
	workspaceMiddlewareHandler := middleware.NewWorkSpaceMiddleware(workspaceService).Handle()

	apiRoute := gin.Group(ApiRoute)
	publicRoute := apiRoute.Group(BaseRoute)
	userRoute := apiRoute.Group(UserRoute, authMiddlewareHandler)
	workspaceRoute := apiRoute.Group(WorkspaceRoute, authMiddlewareHandler, workspaceMiddlewareHandler)
	fileHistoryRoute := workspaceRoute.Group("/:workspace_id" + FileHistoryRoute)
	fileRoute := workspaceRoute.Group("/:workspace_id" + FileRoute)

	InitAuthRoute(publicRoute, userService, authService, userMapper)
	InitUserRoute(userRoute, userService, workspaceService, userMapper, workspaceMapper, paginateMapper)
	InitWorkspaceRoute(workspaceRoute, workspaceService, fileHistoryService, workspaceMapper, fileHistoryMapper, paginateMapper)
	InitFileHistoryRoute(fileHistoryRoute, fileHistoryService, fileHistoryMapper)
	InitFileRoute(fileRoute, userService, fileHistoryService, fileService, fileStorageService, workspaceService, fileMapper)
}
