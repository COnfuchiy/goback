package request

type CreateWorkspaceRequest struct {
	Name string `form:"name" binding:"required"`
}
