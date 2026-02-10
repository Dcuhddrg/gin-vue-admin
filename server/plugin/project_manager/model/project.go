package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Project 项目模型
type Project struct {
	global.GVA_MODEL
	Name string `json:"name" gorm:"column:name;comment:项目名称;not null;uniqueIndex:idx_project_name"`
	Desc string `json:"desc" gorm:"column:desc;comment:项目描述"`
}

func (Project) TableName() string {
	return "gva_projects"
}
