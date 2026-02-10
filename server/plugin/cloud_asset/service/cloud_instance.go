package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model/request"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CloudInstanceService struct{}

// 缓存配置
const (
	// 缓存过期时间：5分钟
	cacheExpireTime = 5 * time.Minute
	// 重试次数
	maxRetryAttempts = 3
	// 重试延迟
	retryDelay = 1 * time.Second
	// 最大并发数
	maxConcurrentSync = 5
)

// 缓存结构
type syncCache struct {
	data      []model.CloudInstance
	expiredAt time.Time
}

var (
	instanceCache sync.Map // key: providerID_region, value: *syncCache
	cacheLock     sync.RWMutex
)

// CreateCloudInstance 创建云服务器实例（手动录入）
func (s *CloudInstanceService) CreateCloudInstance(instance model.CloudInstance) error {
	// 检查实例是否已存在
	var count int64
	if err := global.GVA_DB.Model(&model.CloudInstance{}).
		Where("instance_id = ? AND provider_id = ?", instance.InstanceID, instance.ProviderID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("实例已存在")
	}

	now := time.Now()
	instance.LastSyncAt = &now
	instance.LastSyncStatus = "manual"

	return global.GVA_DB.Create(&instance).Error
}

// DeleteCloudInstance 删除云服务器实例
func (s *CloudInstanceService) DeleteCloudInstance(instance model.CloudInstance) error {
	return global.GVA_DB.Delete(&instance).Error
}

// BatchDeleteCloudInstance 批量删除云服务器实例
func (s *CloudInstanceService) BatchDeleteCloudInstance(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("ids不能为空")
	}
	return global.GVA_DB.Delete(&model.CloudInstance{}, ids).Error
}

// UpdateCloudInstance 更新云服务器实例
func (s *CloudInstanceService) UpdateCloudInstance(instance model.CloudInstance) error {
	return global.GVA_DB.Save(&instance).Error
}

// GetCloudInstance 获取云服务器实例
func (s *CloudInstanceService) GetCloudInstance(id uint) (instance model.CloudInstance, err error) {
	err = global.GVA_DB.Where("id = ?", id).Preload("Provider").First(&instance).Error
	return
}

// GetCloudInstanceInfoList 分页获取云服务器实例列表
func (s *CloudInstanceService) GetCloudInstanceInfoList(info request.CloudInstanceSearch) (list []model.CloudInstance, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&model.CloudInstance{})

	// 构建查询条件
	if info.ProviderID > 0 {
		db = db.Where("provider_id = ?", info.ProviderID)
	}
	if info.ProviderType != "" {
		db = db.Joins("JOIN gva_cloud_providers ON gva_cloud_providers.id = gva_cloud_instances.provider_id").
			Where("gva_cloud_providers.type = ?", info.ProviderType)
	}
	if info.Region != "" {
		db = db.Where("region = ?", info.Region)
	}
	if info.InstanceID != "" {
		db = db.Where("instance_id LIKE ?", "%"+info.InstanceID+"%")
	}
	if info.InstanceName != "" {
		db = db.Where("instance_name LIKE ?", "%"+info.InstanceName+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.PublicIP != "" {
		db = db.Where("public_ip LIKE ?", "%"+info.PublicIP+"%")
	}
	if info.PrivateIP != "" {
		db = db.Where("private_ip LIKE ?", "%"+info.PrivateIP+"%")
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 获取列表
	var instances []model.CloudInstance
	err = db.Limit(limit).Offset(offset).
		Preload("Provider").
		Order("created_at DESC").
		Find(&instances).Error

	return instances, total, err
}

// GetCloudInstanceVOList 分页获取云服务器实例视图列表
func (s *CloudInstanceService) GetCloudInstanceVOList(info request.CloudInstanceSearch) (list []model.InstanceVO, total int64, err error) {
	instances, total, err := s.GetCloudInstanceInfoList(info)
	if err != nil {
		return nil, 0, err
	}

	// 转换为视图对象
	voList := make([]model.InstanceVO, 0, len(instances))
	for _, inst := range instances {
		vo := model.InstanceVO{
			ID:             inst.ID,
			InstanceID:     inst.InstanceID,
			InstanceName:   inst.InstanceName,
			ProviderName:   inst.Provider.Remark,
			ProviderType:   inst.Provider.Type,
			Region:         inst.Region,
			InstanceType:   inst.InstanceType,
			CPU:            inst.CPU,
			Memory:         inst.Memory,
			OSName:         inst.OSName,
			PublicIP:       inst.PublicIP,
			PrivateIP:      inst.PrivateIP,
			Status:         inst.Status,
			CreatedTime:    inst.CreatedTime,
			ExpiredTime:    inst.ExpiredTime,
			ChargeType:     inst.ChargeType,
			LastSyncAt:     inst.LastSyncAt,
		}
		voList = append(voList, vo)
	}

	return voList, total, nil
}

// SyncInstances 同步云服务器实例
func (s *CloudInstanceService) SyncInstances(ctx context.Context, req request.SyncInstancesReq) (int, error) {
	// 获取云厂商信息
	var provider model.CloudProvider
	if err := global.GVA_DB.Where("id = ?", req.ProviderID).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("云厂商不存在")
		}
		return 0, err
	}

	// 检查云厂商状态
	if provider.Status == nil || *provider.Status != 1 {
		return 0, errors.New("云厂商已禁用")
	}

	cacheKey := fmt.Sprintf("%d_%s", req.ProviderID, req.Region)

	// 检查缓存（非强制同步时）
	if !req.ForceSync {
		cacheLock.RLock()
		if cached, ok := instanceCache.Load(cacheKey); ok {
			cache := cached.(*syncCache)
			cacheLock.RUnlock()
			if time.Now().Before(cache.expiredAt) {
				// 缓存有效，直接返回缓存数据
				return s.updateInstancesFromCache(cache.data, req.ProviderID, req.Region)
			}
		} else {
			cacheLock.RUnlock()
		}
	}

	// 调用云厂商API获取实例列表
	var instances []model.CloudInstance
	var err error

	switch provider.Type {
	case "aliyun":
		instances, err = s.syncAliyunInstances(provider, req.Region)
	case "tencent":
		instances, err = s.syncTencentInstances(provider, req.Region)
	case "aws":
		return 0, errors.New("AWS 暂未支持")
	default:
		return 0, fmt.Errorf("不支持的云厂商类型: %s", provider.Type)
	}

	if err != nil {
		// 记录同步失败状态
		_ = s.updateSyncStatus(req.ProviderID, req.Region, "failed")
		return 0, fmt.Errorf("同步云服务器实例失败: %w", err)
	}

	// 更新数据库
	count, err := s.updateInstancesFromCache(instances, req.ProviderID, req.Region)
	if err != nil {
		return 0, fmt.Errorf("更新数据库失败: %w", err)
	}

	// 更新缓存
	cacheLock.Lock()
	instanceCache.Store(cacheKey, &syncCache{
		data:      instances,
		expiredAt: time.Now().Add(cacheExpireTime),
	})
	cacheLock.Unlock()

	// 记录同步成功状态
	_ = s.updateSyncStatus(req.ProviderID, req.Region, "success")

	return count, nil
}

// updateInstancesFromCache 根据缓存数据更新数据库
func (s *CloudInstanceService) updateInstancesFromCache(instances []model.CloudInstance, providerID uint, region string) (int, error) {
	if len(instances) == 0 {
		return 0, nil
	}

	// 使用事务更新
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Upsert Instances
	count, err := s.upsertInstances(tx, instances, providerID, region)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// 2. Delete Stale Instances
	if err := s.deleteStaleInstances(tx, instances, providerID, region); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *CloudInstanceService) upsertInstances(tx *gorm.DB, instances []model.CloudInstance, providerID uint, region string) (int, error) {
	now := time.Now()
	count := 0

	for _, instance := range instances {
		instance.ProviderID = providerID
		instance.Region = region
		instance.LastSyncAt = &now
		instance.LastSyncStatus = "synced"

		// 使用 Upsert 逻辑：存在则更新，不存在则创建
		var existing model.CloudInstance
		err := tx.Where("instance_id = ? AND provider_id = ?", instance.InstanceID, providerID).
			First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 新实例，插入
			if err := tx.Create(&instance).Error; err != nil {
				return 0, err
			}
			count++
		} else if err != nil {
			return 0, err
		} else {
			// 已存在，更新
			instance.ID = existing.ID
			if err := tx.Model(&existing).Updates(&instance).Error; err != nil {
				return 0, err
			}
			count++
		}
	}
	return count, nil
}

func (s *CloudInstanceService) deleteStaleInstances(tx *gorm.DB, instances []model.CloudInstance, providerID uint, region string) error {
	var existingIDs []string
	for _, inst := range instances {
		existingIDs = append(existingIDs, inst.InstanceID)
	}

	return tx.Where("provider_id = ? AND region = ?", providerID, region).
		Not("instance_id", existingIDs).
		Delete(&model.CloudInstance{}).Error
}

// updateSyncStatus 更新同步状态
func (s *CloudInstanceService) updateSyncStatus(providerID uint, region string, status string) error {
	now := time.Now()
	return global.GVA_DB.Model(&model.CloudInstance{}).
		Where("provider_id = ? AND region = ?", providerID, region).
		Updates(map[string]interface{}{
			"last_sync_status": status,
			"last_sync_at":     now,
		}).Error
}

// syncAliyunInstances 同步阿里云实例
func (s *CloudInstanceService) syncAliyunInstances(provider model.CloudProvider, region string) ([]model.CloudInstance, error) {
	var instances []model.CloudInstance

	// 重试机制
	for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
		client, err := ecs.NewClientWithAccessKey(region, provider.AK, provider.SK)
		if err != nil {
			if attempt < maxRetryAttempts {
				time.Sleep(retryDelay * time.Duration(attempt))
				continue
			}
			return nil, fmt.Errorf("创建阿里云客户端失败: %w", err)
		}

		request := ecs.CreateDescribeInstancesRequest()
		request.Scheme = "https"
		request.PageSize = "100"
		request.RegionId = region

		response, err := client.DescribeInstances(request)
		if err != nil {
			if attempt < maxRetryAttempts {
				global.GVA_LOG.Warn("同步阿里云实例失败，重试中",
					zap.String("attempt", fmt.Sprintf("%d/%d", attempt, maxRetryAttempts)),
					zap.String("error", err.Error()),
				)
				time.Sleep(retryDelay * time.Duration(attempt))
				continue
			}
			return nil, fmt.Errorf("获取阿里云实例失败: %w", err)
		}

		// 解析实例列表
		for _, inst := range response.Instances.Instance {
			instances = append(instances, s.convertAliyunInstance(inst))
		}

		break
	}

	return instances, nil
}

// convertAliyunInstance 转换阿里云实例数据
func (s *CloudInstanceService) convertAliyunInstance(aliInst ecs.Instance) model.CloudInstance {
	instance := model.CloudInstance{
		InstanceID:     aliInst.InstanceId,
		InstanceName:   aliInst.InstanceName,
		InstanceType:   aliInst.InstanceType,
		Region:         aliInst.RegionId,
		OSName:         aliInst.OSName,
		PublicIP:       extractPublicIP(aliInst),
		PrivateIP:      extractPrivateIP(aliInst),
		Status:         s.mapAliyunStatus(aliInst.Status),
		InstanceStatus: aliInst.Status,
	}

	// 处理时间
	if aliInst.CreationTime != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", aliInst.CreationTime); err == nil {
			instance.CreatedTime = &t
		}
	}
	if aliInst.ExpiredTime != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", aliInst.ExpiredTime); err == nil {
			instance.ExpiredTime = &t
		}
	}

	// CPU和内存（aliInst.Cpu 和 Memory 是 int 类型）
	if aliInst.Cpu > 0 {
		cpu := aliInst.Cpu
		instance.CPU = &cpu
	}
	if aliInst.Memory > 0 {
		memory := aliInst.Memory / 1024 / 1024 / 1024 // 转换为GB
		instance.Memory = &memory
	}

	// 付费类型
	instance.ChargeType = aliInst.InstanceChargeType

	return instance
}

// mapAliyunStatus 映射阿里云状态
func (s *CloudInstanceService) mapAliyunStatus(status string) string {
	statusMap := map[string]string{
		"Running":   "running",
		"Starting":  "starting",
		"Stopping":  "stopping",
		"Stopped":   "stopped",
		"Pending":   "pending",
		"Deleting":  "deleting",
	}
	if mapped, ok := statusMap[status]; ok {
		return mapped
	}
	return "unknown"
}

// extractPublicIP 提取公网IP
func extractPublicIP(aliInst ecs.Instance) string {
	if len(aliInst.PublicIpAddress.IpAddress) > 0 {
		return aliInst.PublicIpAddress.IpAddress[0]
	}
	if aliInst.EipAddress.IpAddress != "" {
		return aliInst.EipAddress.IpAddress
	}
	return ""
}

// extractPrivateIP 提取私网IP
func extractPrivateIP(aliInst ecs.Instance) string {
	if len(aliInst.InnerIpAddress.IpAddress) > 0 {
		return aliInst.InnerIpAddress.IpAddress[0]
	}
	if len(aliInst.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		return aliInst.VpcAttributes.PrivateIpAddress.IpAddress[0]
	}
	return ""
}

// syncTencentInstances 同步腾讯云实例
func (s *CloudInstanceService) syncTencentInstances(provider model.CloudProvider, region string) ([]model.CloudInstance, error) {
	var instances []model.CloudInstance

	// 重试机制
	for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
		credential := common.NewCredential(provider.AK, provider.SK)
		cpf := profile.NewClientProfile()
		client, err := cvm.NewClient(credential, region, cpf)
		if err != nil {
			if attempt < maxRetryAttempts {
				time.Sleep(retryDelay * time.Duration(attempt))
				continue
			}
			return nil, fmt.Errorf("创建腾讯云客户端失败: %w", err)
		}

		request := cvm.NewDescribeInstancesRequest()
		request.Limit = common.Int64Ptr(100)

		response, err := client.DescribeInstances(request)
		if err != nil {
			if attempt < maxRetryAttempts {
				global.GVA_LOG.Warn("同步腾讯云实例失败，重试中",
					zap.String("attempt", fmt.Sprintf("%d/%d", attempt, maxRetryAttempts)),
					zap.String("error", err.Error()),
				)
				time.Sleep(retryDelay * time.Duration(attempt))
				continue
			}
			return nil, fmt.Errorf("获取腾讯云实例失败: %w", err)
		}

		// 解析实例列表
		for _, inst := range response.Response.InstanceSet {
			instances = append(instances, s.convertTencentInstance(inst, region))
		}

		break
	}

	return instances, nil
}

// convertTencentInstance 转换腾讯云实例数据
func (s *CloudInstanceService) convertTencentInstance(tencentInst *cvm.Instance, region string) model.CloudInstance {
	instance := model.CloudInstance{
		InstanceID:     *tencentInst.InstanceId,
		InstanceName:   *tencentInst.InstanceName,
		InstanceType:   *tencentInst.InstanceType,
		Region:         region,
		OSName:         *tencentInst.OsName,
		Status:         s.mapTencentStatus(*tencentInst.InstanceState),
		InstanceStatus: *tencentInst.InstanceState,
	}

	// CPU和内存
	if tencentInst.CPU != nil {
		cpu := int(*tencentInst.CPU)
		instance.CPU = &cpu
	}
	if tencentInst.Memory != nil {
		memory := int(*tencentInst.Memory)
		instance.Memory = &memory
	}

	// 公网IP和私网IP
	if len(tencentInst.PublicIpAddresses) > 0 {
		instance.PublicIP = *tencentInst.PublicIpAddresses[0]
	}
	if len(tencentInst.PrivateIpAddresses) > 0 {
		instance.PrivateIP = *tencentInst.PrivateIpAddresses[0]
	}

	// 处理时间
	if tencentInst.CreatedTime != nil && *tencentInst.CreatedTime != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", *tencentInst.CreatedTime); err == nil {
			instance.CreatedTime = &t
		}
	}

	// 付费类型
	if tencentInst.InstanceChargeType != nil {
		instance.ChargeType = *tencentInst.InstanceChargeType
	}

	return instance
}

// mapTencentStatus 映射腾讯云状态
func (s *CloudInstanceService) mapTencentStatus(status string) string {
	statusMap := map[string]string{
		"RUNNING":     "running",
		"LAUNCHING":   "starting",
		"STOPPING":    "stopping",
		"STOPPED":     "stopped",
		"TERMINATING": "deleting",
	}
	if mapped, ok := statusMap[status]; ok {
		return mapped
	}
	return "unknown"
}

// BatchSyncInstances 批量同步云服务器实例
func (s *CloudInstanceService) BatchSyncInstances(ctx context.Context, providerID uint, regions []string, forceSync bool) (map[string]int, error) {
	results := make(map[string]int)

	// 使用信号量控制并发
	sem := make(chan struct{}, maxConcurrentSync)
	errChan := make(chan error, len(regions))
	var wg sync.WaitGroup

	for _, region := range regions {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			req := request.SyncInstancesReq{
				ProviderID: providerID,
				Region:     r,
				ForceSync:  forceSync,
			}

			count, err := s.SyncInstances(ctx, req)
			if err != nil {
				errChan <- fmt.Errorf("同步区域 %s 失败: %w", r, err)
				return
			}
			results[r] = count
		}(region)
	}

	// 等待所有任务完成
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 收集错误
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return results, fmt.Errorf("批量同步部分失败: %v", errs)
	}

	return results, nil
}

// GetInstanceStats 获取实例统计信息
func (s *CloudInstanceService) GetInstanceStats() (model.InstanceStatsVO, error) {
	var stats model.InstanceStatsVO

	// 总数统计
	global.GVA_DB.Model(&model.CloudInstance{}).Count(&stats.TotalCount)

	// 状态统计
	global.GVA_DB.Model(&model.CloudInstance{}).Where("status = ?", "running").Count(&stats.RunningCount)
	global.GVA_DB.Model(&model.CloudInstance{}).Where("status = ?", "stopped").Count(&stats.StoppedCount)

	// 按云厂商统计
	type providerResult struct {
		ProviderName string
		ProviderType string
		Count        int64
	}
	var providerResults []providerResult
	global.GVA_DB.Table("gva_cloud_instances").
		Select("gva_cloud_providers.remark as provider_name, gva_cloud_providers.type as provider_type, COUNT(*) as count").
		Joins("JOIN gva_cloud_providers ON gva_cloud_providers.id = gva_cloud_instances.provider_id").
		Group("gva_cloud_providers.id, gva_cloud_providers.remark, gva_cloud_providers.type").
		Scan(&providerResults)

	for _, r := range providerResults {
		stats.ProviderStats = append(stats.ProviderStats, model.ProviderStatVO{
			ProviderName: r.ProviderName,
			ProviderType: r.ProviderType,
			Count:        r.Count,
		})
	}

	// 按区域统计
	type regionResult struct {
		Region string
		Count  int64
	}
	global.GVA_DB.Model(&model.CloudInstance{}).
		Select("region, COUNT(*) as count").
		Group("region").
		Scan(&stats.RegionStats)

	return stats, nil
}

// ClearCache 清除缓存
func (s *CloudInstanceService) ClearCache(providerID uint, region string) {
	if providerID > 0 && region != "" {
		cacheKey := fmt.Sprintf("%d_%s", providerID, region)
		instanceCache.Delete(cacheKey)
	} else {
		// 清除所有缓存
		instanceCache.Range(func(key, value interface{}) bool {
			instanceCache.Delete(key)
			return true
		})
	}
}

// ValidateInstanceOwnership 验证实例所有权
func (s *CloudInstanceService) ValidateInstanceOwnership(instanceID uint, providerID uint) (bool, error) {
	var count int64
	err := global.GVA_DB.Model(&model.CloudInstance{}).
		Where("id = ? AND provider_id = ?", instanceID, providerID).
		Count(&count).Error
	return count > 0, err
}
