package request

type RefreshTokenRequest struct {
	RefreshToken string `form:"refresh_token" binding:"required"`
}
