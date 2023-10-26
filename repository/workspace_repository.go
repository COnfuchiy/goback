package repository

import (
	"github.com/google/uuid"
	"goback/domain/entity"
	"gorm.io/gorm"
)

type IWorkspaceRepository interface {
	Create(workspace *entity.Workspace) error
	FindByID(id uuid.UUID) (*entity.Workspace, error)
	FindAllByUserID(userID uuid.UUID, offset int, limit int) ([]entity.Workspace, int64, error)
}

type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) IWorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r WorkspaceRepository) Create(workspace *entity.Workspace) error {
	return r.db.Model(&entity.Workspace{}).Create(&workspace).Error
}

func (r WorkspaceRepository) FindByID(id uuid.UUID) (*entity.Workspace, error) {
	var workspace entity.Workspace
	err := r.db.Model(&entity.Workspace{}).Preload("Creator").Preload("Users").
		Where("id = ?", id).First(&workspace).Error
	return &workspace, err
}

func (r WorkspaceRepository) FindAllByUserID(userID uuid.UUID, offset int, limit int) ([]entity.Workspace, int64, error) {
	var workspaces []entity.Workspace
	var totalCount int64

	// err := r.db.Model(&entity.Workspace{}).Preload("Creator").Preload("Users").
	// 	Where("creator_id = ?", userID.String()).Or("user_workspaces.user_id = ?", userID.String()).Count(&totalCount).
	// 	Limit(limit).Offset(offset).Find(&workspaces).Error
	err := r.db.Model(&entity.Workspace{}).Preload("Creator").Preload("Users").
		Joins("left join users ON workspaces.creator_id = users.id AND users.id = ?", userID.String()).
		Joins("left join user_workspaces ON workspaces.id = user_workspaces.workspace_id").
		Where("creator_id = ? OR user_id = ?", userID.String(), userID.String()).Count(&totalCount).
		Limit(limit).Offset(offset).Find(&workspaces).Error

	if err != nil {
		return nil, 0, err
	}

	return workspaces, totalCount, err
}
