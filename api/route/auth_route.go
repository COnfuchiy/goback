package route

import (
	"github.com/gin-gonic/gin"
	"goback/api/controller"
	"goback/mapper"
	"goback/services"
)

func InitAuthRoute(router *gin.RouterGroup, userService services.IUserService, authService services.IAuthService, userMapper mapper.IUserMapper) {
	authController := controller.NewAuthController(userService, authService, userMapper)
	router.POST("/login", authController.Login)
	router.POST("/register", authController.Register)
	router.POST("/refresh", authController.RefreshToken)
}
