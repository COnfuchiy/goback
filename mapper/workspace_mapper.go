package mapper

import (
	"goback/api/request"
	"goback/api/response"
	"goback/domain/entity"
)

type IWorkspaceMapper interface {
	ToWorkspaceResponse(workspace *entity.Workspace) *response.WorkspaceResponse
	FromCreateRequest(createRequest *request.CreateWorkspaceRequest) *entity.Workspace
}

type WorkspaceMapper struct {
	userMapper IUserMapper
}

func NewWorkspaceMapper(userMapper IUserMapper) IWorkspaceMapper {
	return &WorkspaceMapper{userMapper: userMapper}
}

func (m WorkspaceMapper) ToWorkspaceResponse(workspace *entity.Workspace) *response.WorkspaceResponse {
	creatorProfileResponse := m.userMapper.ToProfileResponse(&workspace.Creator)
	var usersProfilesResponses []response.ProfileResponse
	for _, user := range workspace.Users {
		usersProfilesResponses = append(usersProfilesResponses, *m.userMapper.ToProfileResponse(&user))
	}
	return &response.WorkspaceResponse{
		ID:      workspace.ID.String(),
		Name:    workspace.Name,
		Creator: *creatorProfileResponse,
		Users:   usersProfilesResponses,
	}
}

func (m WorkspaceMapper) FromCreateRequest(createRequest *request.CreateWorkspaceRequest) *entity.Workspace {
	return &entity.Workspace{
		Name: createRequest.Name,
	}
}
