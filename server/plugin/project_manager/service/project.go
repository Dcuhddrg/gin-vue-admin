package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/model/request"
)

type ProjectService struct{}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(project model.Project) error {
	return global.GVA_DB.Create(&project).Error
}

// DeleteProject 删除项目
func (s *ProjectService) DeleteProject(project model.Project) error {
	return global.GVA_DB.Delete(&project).Error
}

// UpdateProject 更新项目
func (s *ProjectService) UpdateProject(project model.Project) error {
	return global.GVA_DB.Save(&project).Error
}

// GetProject 获取项目
func (s *ProjectService) GetProject(id uint) (project model.Project, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&project).Error
	return
}

// GetProjectList 获取项目列表
func (s *ProjectService) GetProjectList(info request.ProjectSearch) (list []model.Project, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&model.Project{})

	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}
