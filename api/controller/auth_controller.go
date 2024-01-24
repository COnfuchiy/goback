package controller

import (
	"github.com/gin-gonic/gin"
	"goback/api/request"
	"goback/api/response"
	"goback/mapper"
	"goback/services"
	"net/http"
)

type AuthController struct {
	userService services.IUserService
	authService services.IAuthService
	userMapper  mapper.IUserMapper
}

func NewAuthController(userService services.IUserService, authService services.IAuthService, userMapper mapper.IUserMapper) *AuthController {
	return &AuthController{userService, authService, userMapper}
}

// Login godoc
// @Summary	login user
// @Description	log in user in system
// @Tags auth
// @Accept x-www-form-urlencoded
// @Param email formData string true "email" Format(email)
// @Param password formData string true "password"
// @Produce json
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /login [post]
func (c AuthController) Login(context *gin.Context) {

	var req request.LoginRequest

	err := context.ShouldBind(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := c.userService.GetByEmail(req.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid credentials"})
		return
	}

	if c.authService.ComparePassword(*user.Password, req.Password) != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := c.authService.CreateAccessAndRefreshTokens(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response.LoginResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Register godoc
// @Summary	register user
// @Description	register user in system
// @Tags auth
// @Accept x-www-form-urlencoded
// @Param email formData string true "email" Format(email)
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Produce json
// @Success 200 {object} response.RegisterResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /register [post]
func (c AuthController) Register(context *gin.Context) {
	var req request.RegisterRequest

	err := context.ShouldBind(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = c.userService.GetByEmail(req.Email)
	if err == nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User already exists with the given email"})
		return
	}

	req.Password, err = c.authService.HashPassword(req.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	user := c.userMapper.FromRegisterRequest(&req)

	err = c.userService.Create(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	accessToken, refreshToken, err := c.authService.CreateAccessAndRefreshTokens(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response.RegisterResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshToken godoc
// @Summary	refresh user
// @Description	get refreshed user refresh token
// @Tags auth
// @Accept x-www-form-urlencoded
// @Param refresh_token formData string true "refresh token"
// @Produce json
// @Success 200 {object} response.RegisterResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /refresh [post]
func (c AuthController) RefreshToken(context *gin.Context) {
	var req request.RefreshTokenRequest

	err := context.ShouldBind(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	userID, err := c.authService.ExtractIDFromRefreshToken(req.RefreshToken)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := c.userService.GetUserFromContext(userID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User not found"})
		return
	}

	accessToken, refreshToken, err := c.authService.CreateAccessAndRefreshTokens(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	context.JSON(http.StatusOK, response.RegisterResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
