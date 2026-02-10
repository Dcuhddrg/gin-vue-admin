package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	projectManagerModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/model"
)

type CloudProvider struct {
	global.GVA_MODEL
	AK        string                       `json:"ak" gorm:"column:ak;comment:AccessKey"`
	SK        string                       `json:"sk" gorm:"column:sk;comment:SecretKey"`
	Type      string                       `json:"type" gorm:"column:type;comment:厂商类型"`
	Region    string                       `json:"region" gorm:"column:region;comment:区域"`
	Status    *int                         `json:"status" gorm:"column:status;comment:状态:1启用,2禁用;default:1"`
	Remark    string                       `json:"remark" gorm:"column:remark;comment:备注"`
	ProjectID uint                         `json:"projectId" gorm:"column:project_id;comment:关联项目ID"`
	Project   projectManagerModel.Project `json:"project" gorm:"foreignKey:ProjectID"`
}

func (CloudProvider) TableName() string {
	return "gva_cloud_providers"
}
