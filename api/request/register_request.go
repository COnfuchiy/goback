package request

type RegisterRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}
