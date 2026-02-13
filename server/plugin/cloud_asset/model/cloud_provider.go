package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/utils"
	projectManagerModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/model"
	"gorm.io/gorm"
)

type CloudProvider struct {
	global.GVA_MODEL
	AK        string                      `json:"ak" gorm:"column:ak;comment:AccessKey"`
	SK        string                      `json:"sk" gorm:"column:sk;comment:SecretKey"`
	Type      string                      `json:"type" gorm:"column:type;comment:厂商类型"`
	Region    string                      `json:"region" gorm:"column:region;comment:区域"`
	Status    *int                        `json:"status" gorm:"column:status;comment:状态:1启用,2禁用;default:1"`
	Remark    string                      `json:"remark" gorm:"column:remark;comment:备注"`
	ProjectID uint                        `json:"projectId" gorm:"column:project_id;comment:关联项目ID"`
	Project   projectManagerModel.Project `json:"project" gorm:"foreignKey:ProjectID"`
}

func (CloudProvider) TableName() string {
	return "gva_cloud_providers"
}

// BeforeSave 保存前加密敏感信息
func (cp *CloudProvider) BeforeSave(tx *gorm.DB) error {
	if cp.AK != "" {
		encrypted, err := utils.Encrypt(cp.AK)
		if err != nil {
			return err
		}
		cp.AK = encrypted
	}
	if cp.SK != "" {
		encrypted, err := utils.Encrypt(cp.SK)
		if err != nil {
			return err
		}
		cp.SK = encrypted
	}
	return nil
}

// AfterFind 查询后解密敏感信息
func (cp *CloudProvider) AfterFind(tx *gorm.DB) error {
	if cp.AK != "" {
		decrypted, err := utils.Decrypt(cp.AK)
		if err != nil {
			// 解密失败可能是因为数据未加密（如旧数据），记录错误但不中断
			// 这里假设解密失败返回空字符串或者保留原文视业务需求而定
			// 也可以选择返回 error
			return err
		}
		cp.AK = decrypted
	}
	if cp.SK != "" {
		decrypted, err := utils.Decrypt(cp.SK)
		if err != nil {
			return err
		}
		cp.SK = decrypted
	}
	return nil
}
