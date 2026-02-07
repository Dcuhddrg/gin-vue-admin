package service

import (
	"errors"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model/request"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type CloudProviderService struct{}

// CreateCloudProvider 创建云厂商
func (s *CloudProviderService) CreateCloudProvider(cp model.CloudProvider) error {
	return global.GVA_DB.Create(&cp).Error
}

// DeleteCloudProvider 删除云厂商
func (s *CloudProviderService) DeleteCloudProvider(cp model.CloudProvider) error {
	return global.GVA_DB.Delete(&cp).Error
}

// UpdateCloudProvider 更新云厂商
func (s *CloudProviderService) UpdateCloudProvider(cp model.CloudProvider) error {
	return global.GVA_DB.Save(&cp).Error
}

// GetCloudProvider 获取云厂商
func (s *CloudProviderService) GetCloudProvider(id uint) (cp model.CloudProvider, err error) {
	err = global.GVA_DB.Where("id = ?", id).Preload("Project").First(&cp).Error
	return
}

// GetCloudProviderInfoList 分页获取云厂商列表
func (s *CloudProviderService) GetCloudProviderInfoList(info request.CloudProviderSearch) (list []model.CloudProvider, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&model.CloudProvider{})
	var cps []model.CloudProvider

	if info.ProjectId > 0 {
		db = db.Where("project_id = ?", info.ProjectId)
	}
	if info.Type != "" {
		db = db.Where("type = ?", info.Type)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Preload("Project").Find(&cps).Error
	return cps, total, err
}

// GetRegions 获取云厂商可用区
func (s *CloudProviderService) GetRegions(provider, accessKey, secretKey string) ([]model.RegionVO, error) {
	switch provider {
	case "aliyun":
		return s.getAliyunRegions(accessKey, secretKey)
	case "tencent":
		return s.getTencentRegions(accessKey, secretKey)
	case "aws":
		return nil, errors.New("AWS 暂未支持")
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *CloudProviderService) getAliyunRegions(ak, sk string) ([]model.RegionVO, error) {
	client, err := ecs.NewClientWithAccessKey("cn-hangzhou", ak, sk)
	if err != nil {
		return nil, err
	}

	request := ecs.CreateDescribeRegionsRequest()
	response, err := client.DescribeRegions(request)
	if err != nil {
		return nil, err
	}

	var regions []model.RegionVO
	for _, r := range response.Regions.Region {
		regions = append(regions, model.RegionVO{
			RegionId:  r.RegionId,
			LocalName: r.LocalName,
		})
	}
	return regions, nil
}

func (s *CloudProviderService) getTencentRegions(ak, sk string) ([]model.RegionVO, error) {
	credential := common.NewCredential(ak, sk)
	cpf := profile.NewClientProfile()
	client, err := cvm.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return nil, err
	}

	request := cvm.NewDescribeRegionsRequest()
	response, err := client.DescribeRegions(request)
	if err != nil {
		return nil, err
	}

	var regions []model.RegionVO
	for _, r := range response.Response.RegionSet {
		if r.RegionState != nil && *r.RegionState == "AVAILABLE" {
			regions = append(regions, model.RegionVO{
				RegionId:  *r.Region,
				LocalName: *r.RegionName,
			})
		}
	}
	return regions, nil
}
