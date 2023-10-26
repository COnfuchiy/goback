package middleware

import (
	"github.com/gin-gonic/gin"
	"goback/api/response"
	"goback/services"
	"net/http"
	"strings"
)

type JwtAuthMiddleware struct {
	authService services.IAuthService
	userService services.IUserService
}

func NewJwtAuthMiddleware(authService services.IAuthService, userService services.IUserService) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{authService, userService}
}

func (m JwtAuthMiddleware) Handle() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.Request.Header.Get("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")

		if len(authHeaderParts) == 2 {
			authToken := authHeaderParts[1]
			userID, err := m.authService.ExtractIDFromAccessToken(authToken)
			if err != nil {
				context.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: err.Error()})
				context.Abort()
				return
			}

			user, err := m.userService.GetUserFromContext(userID)
			if err != nil {
				context.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: err.Error()})
				context.Abort()
				return
			}

			context.Set("user", user)
			context.Next()
			return
		}

		context.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "Not authorized"})
		context.Abort()
	}
}
