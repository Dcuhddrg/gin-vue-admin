package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model/response"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"go.uber.org/zap"
)

type CloudProviderService struct{}

var (
	configCache sync.Map // key: projectId_providerType, value: *configCacheItem
)

type configCacheItem struct {
	data      response.ProviderConfigResponse
	expiredAt time.Time
}

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

// GetProviderConfig 获取厂商配置（含区域列表）
func (s *CloudProviderService) GetProviderConfig(projectId uint, providerType string) (response.ProviderConfigResponse, error) {
	// 1. 检查缓存
	cacheKey := fmt.Sprintf("%d_%s", projectId, providerType)
	if cached, ok := configCache.Load(cacheKey); ok {
		item := cached.(*configCacheItem)
		if time.Now().Before(item.expiredAt) {
			return item.data, nil
		}
		configCache.Delete(cacheKey)
	}

	// 2. 查找有效凭证
	var cp model.CloudProvider
	err := global.GVA_DB.Where("project_id = ? AND type = ? AND status = 1", projectId, providerType).First(&cp).Error
	if err != nil {
		return response.ProviderConfigResponse{}, errors.New("该项目未配置有效的云厂商凭证")
	}

	// 3. 获取区域列表
	allRegions, err := s.GetRegions(providerType, cp.AK, cp.SK)
	if err != nil {
		global.GVA_LOG.Error("获取区域列表失败", zap.Error(err))
		return response.ProviderConfigResponse{}, fmt.Errorf("调用云厂商API失败: %v", err)
	}

	// 4. 如果配置了区域白名单，则进行过滤
	var regions []model.RegionVO
	if cp.Region != "" {
		allowedRegions := make(map[string]bool)
		for _, r := range strings.Split(cp.Region, ",") {
			if trimmed := strings.TrimSpace(r); trimmed != "" {
				allowedRegions[trimmed] = true
			}
		}

		for _, r := range allRegions {
			if allowedRegions[r.RegionId] {
				regions = append(regions, r)
			}
		}
	} else {
		regions = allRegions
	}

	// 5. 组装响应
	resp := response.ProviderConfigResponse{
		ProviderInfo: response.ProviderInfo{
			Name:    s.getProviderName(providerType),
			Version: "v1.0",
		},
		Regions: regions,
		Limits: response.ProviderLimits{
			MaxSyncBatch: 5,
		},
	}

	// 6. 写入缓存 (10分钟)
	configCache.Store(cacheKey, &configCacheItem{
		data:      resp,
		expiredAt: time.Now().Add(10 * time.Minute),
	})

	return resp, nil
}

func (s *CloudProviderService) getProviderName(t string) string {
	switch t {
	case "aliyun":
		return "阿里云"
	case "tencent":
		return "腾讯云"
	case "aws":
		return "AWS"
	default:
		return t
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
