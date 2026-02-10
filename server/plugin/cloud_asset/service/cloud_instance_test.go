package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonRequest "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model/request"
)

// TestCloudInstanceService_CreateCloudInstance 测试创建云服务器实例
func TestCloudInstanceService_CreateCloudInstance(t *testing.T) {
	service := &CloudInstanceService{}

	// 测试数据
	now := time.Now()
	testInstance := model.CloudInstance{
		InstanceID:     "test-instance-001",
		InstanceName:   "测试实例",
		ProviderID:     1,
		Region:         "cn-hangzhou",
		InstanceType:   "ecs.t6-c1m2.large",
		CPU:            func() *int { i := 2; return &i }(),
		Memory:         func() *int { i := 4; return &i }(),
		OSName:         "CentOS 7.9 64位",
		PublicIP:       "1.2.3.4",
		PrivateIP:      "10.0.0.1",
		Status:         "running",
		CreatedTime:    &now,
		LastSyncAt:     &now,
		LastSyncStatus: "test",
	}

	// 注意：这个测试需要真实的数据库连接
	// 如果没有数据库连接，会跳过
	if global.GVA_DB == nil {
		t.Skip("数据库连接未初始化")
	}

	// 执行创建
	err := service.CreateCloudInstance(testInstance)
	if err != nil {
		t.Errorf("创建云服务器实例失败: %v", err)
	}

	// 验证实例是否存在
	var count int64
	if err := global.GVA_DB.Model(&model.CloudInstance{}).
		Where("instance_id = ? AND provider_id = ?", testInstance.InstanceID, testInstance.ProviderID).
		Count(&count).Error; err != nil {
		t.Errorf("查询实例失败: %v", err)
	}
	if count != 1 {
		t.Errorf("实例创建失败，期望1条记录，实际%d条", count)
	}

	// 清理测试数据
	global.GVA_DB.Where("instance_id = ?", testInstance.InstanceID).Delete(&model.CloudInstance{})
}

// TestCloudInstanceService_CreateDuplicateInstance 测试创建重复实例
func TestCloudInstanceService_CreateDuplicateInstance(t *testing.T) {
	service := &CloudInstanceService{}

	if global.GVA_DB == nil {
		t.Skip("数据库连接未初始化")
	}

	// 测试数据
	testInstance := model.CloudInstance{
		InstanceID:   "test-duplicate-001",
		InstanceName: "测试重复实例",
		ProviderID:   1,
		Region:       "cn-hangzhou",
		Status:       "running",
	}

	// 第一次创建
	err := service.CreateCloudInstance(testInstance)
	if err != nil {
		t.Fatalf("第一次创建失败: %v", err)
	}

	// 第二次创建（应该失败）
	err = service.CreateCloudInstance(testInstance)
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}

	// 清理
	global.GVA_DB.Where("instance_id = ?", testInstance.InstanceID).Delete(&model.CloudInstance{})
}

// TestCloudInstanceService_GetCloudInstanceInfoList 测试分页查询实例列表
func TestCloudInstanceService_GetCloudInstanceInfoList(t *testing.T) {
	service := &CloudInstanceService{}

	if global.GVA_DB == nil {
		t.Skip("数据库连接未初始化")
	}

	// 创建测试数据
	testInstances := []model.CloudInstance{
		{
			InstanceID:   "test-list-001",
			InstanceName: "测试实例1",
			ProviderID:   1,
			Region:       "cn-hangzhou",
			Status:       "running",
		},
		{
			InstanceID:   "test-list-002",
			InstanceName: "测试实例2",
			ProviderID:   1,
			Region:       "cn-shanghai",
			Status:       "stopped",
		},
	}

	for _, inst := range testInstances {
		service.CreateCloudInstance(inst)
	}
	defer func() {
		for _, inst := range testInstances {
			global.GVA_DB.Where("instance_id = ?", inst.InstanceID).Delete(&model.CloudInstance{})
		}
	}()

	// 测试分页查询
	search := request.CloudInstanceSearch{
		PageInfo: commonRequest.PageInfo{
			Page:     1,
			PageSize: 10,
		},
	}

	list, total, err := service.GetCloudInstanceInfoList(search)
	if err != nil {
		t.Errorf("查询实例列表失败: %v", err)
	}

	if total < 2 {
		t.Errorf("期望至少2条记录，实际%d条", total)
	}

	if len(list) < 2 {
		t.Errorf("期望返回至少2条记录，实际%d条", len(list))
	}
}

// TestCloudInstanceService_GetCloudInstanceInfoList_WithFilter 测试带过滤条件的查询
func TestCloudInstanceService_GetCloudInstanceInfoList_WithFilter(t *testing.T) {
	service := &CloudInstanceService{}

	if global.GVA_DB == nil {
		t.Skip("数据库连接未初始化")
	}

	// 创建测试数据
	testInstance := model.CloudInstance{
		InstanceID:   "test-filter-001",
		InstanceName: "测试过滤实例",
		ProviderID:   1,
		Region:       "cn-hangzhou",
		Status:       "running",
	}
	service.CreateCloudInstance(testInstance)
	defer global.GVA_DB.Where("instance_id = ?", testInstance.InstanceID).Delete(&model.CloudInstance{})

	// 测试按状态过滤
	search := request.CloudInstanceSearch{
		PageInfo: commonRequest.PageInfo{
			Page:     1,
			PageSize: 10,
		},
		Status: "running",
	}

	list, total, err := service.GetCloudInstanceInfoList(search)
	if err != nil {
		t.Errorf("查询实例列表失败: %v", err)
	}

	if total < 1 {
		t.Errorf("期望至少1条running状态的记录，实际%d条", total)
	}

	for _, inst := range list {
		if inst.Status != "running" {
			t.Errorf("期望status=running，实际status=%s", inst.Status)
		}
	}
}

// TestCloudInstanceService_UpdateCloudInstance 测试更新实例
func TestCloudInstanceService_UpdateCloudInstance(t *testing.T) {
	service := &CloudInstanceService{}

	if global.GVA_DB == nil {
		t.Skip("数据库连接未初始化")
	}

	// 创建测试实例
	testInstance := model.CloudInstance{
		InstanceID:   "test-update-001",
		InstanceName: "测试更新实例",
		ProviderID:   1,
		Region:       "cn-hangzhou",
		Status:       "running",
	}
	service.CreateCloudInstance(testInstance)
	defer global.GVA_DB.Where("instance_id = ?", testInstance.InstanceID).Delete(&model.CloudInstance{})

	// 获取创建的实例
	created, _ := service.GetCloudInstance(testInstance.ID)

	// 更新实例
	created.InstanceName = "更新后的实例名"
	created.Remark = "更新备注"

	err := service.UpdateCloudInstance(created)
	if err != nil {
		t.Errorf("更新实例失败: %v", err)
	}

	// 验证更新
	updated, err := service.GetCloudInstance(testInstance.ID)
	if err != nil {
		t.Errorf("获取更新后的实例失败: %v", err)
	}

	if updated.InstanceName != "更新后的实例名" {
		t.Errorf("期望InstanceName=更新后的实例名，实际=%s", updated.InstanceName)
	}

	if updated.Remark != "更新备注" {
		t.Errorf("期望Remark=更新备注，实际=%s", updated.Remark)
	}
}

// TestCloudInstanceService_DeleteCloudInstance 测试删除实例
func TestCloudInstanceService_DeleteCloudInstance(t *testing.T) {
	service := &CloudInstanceService{}

	if global.GVA_DB == nil {
		t.Skip("数据库连接未初始化")
	}

	// 创建测试实例
	testInstance := model.CloudInstance{
		InstanceID:   "test-delete-001",
		InstanceName: "测试删除实例",
		ProviderID:   1,
		Region:       "cn-hangzhou",
		Status:       "running",
	}
	service.CreateCloudInstance(testInstance)

	// 获取创建的实例
	created, _ := service.GetCloudInstance(testInstance.ID)

	// 删除实例
	err := service.DeleteCloudInstance(created)
	if err != nil {
		t.Errorf("删除实例失败: %v", err)
	}

	// 验证删除
	var count int64
	if err := global.GVA_DB.Model(&model.CloudInstance{}).
		Where("instance_id = ?", testInstance.InstanceID).
		Count(&count).Error; err != nil {
		t.Errorf("查询实例失败: %v", err)
	}

	if count != 0 {
		t.Errorf("期望0条记录，实际%d条", count)
	}
}

// TestCloudInstanceService_BatchDeleteCloudInstance 测试批量删除
func TestCloudInstanceService_BatchDeleteCloudInstance(t *testing.T) {
	service := &CloudInstanceService{}

	if global.GVA_DB == nil {
		t.Skip("数据库连接未初始化")
	}

	// 创建测试数据
	ids := []uint{}
	for i := 0; i < 3; i++ {
		testInstance := model.CloudInstance{
			InstanceID:   testID("test-batch-delete", i),
			InstanceName: "测试批量删除实例",
			ProviderID:   1,
			Region:       "cn-hangzhou",
			Status:       "running",
		}
		service.CreateCloudInstance(testInstance)
		created, _ := service.GetCloudInstance(testInstance.ID)
		ids = append(ids, created.ID)
	}

	// 批量删除
	err := service.BatchDeleteCloudInstance(ids)
	if err != nil {
		t.Errorf("批量删除实例失败: %v", err)
	}

	// 验证删除
	var count int64
	global.GVA_DB.Model(&model.CloudInstance{}).Where("id IN ?", ids).Count(&count)

	if count != 0 {
		t.Errorf("期望0条记录，实际%d条", count)
	}
}

// TestCloudInstanceService_BatchDeleteCloudInstance_Empty 测试批量删除空数组
func TestCloudInstanceService_BatchDeleteCloudInstance_Empty(t *testing.T) {
	service := &CloudInstanceService{}

	err := service.BatchDeleteCloudInstance([]uint{})
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}

	if err.Error() != "ids不能为空" {
		t.Errorf("期望错误消息='ids不能为空'，实际='%s'", err.Error())
	}
}

// TestCloudInstanceService_GetInstanceStats 测试获取实例统计信息
func TestCloudInstanceService_GetInstanceStats(t *testing.T) {
	service := &CloudInstanceService{}

	if global.GVA_DB == nil {
		t.Skip("数据库连接未初始化")
	}

	// 创建测试数据
	testInstances := []model.CloudInstance{
		{
			InstanceID:   "test-stats-running-001",
			InstanceName: "测试运行实例",
			ProviderID:   1,
			Region:       "cn-hangzhou",
			Status:       "running",
		},
		{
			InstanceID:   "test-stats-stopped-001",
			InstanceName: "测试停止实例",
			ProviderID:   1,
			Region:       "cn-hangzhou",
			Status:       "stopped",
		},
	}

	for _, inst := range testInstances {
		service.CreateCloudInstance(inst)
	}
	defer func() {
		for _, inst := range testInstances {
			global.GVA_DB.Where("instance_id = ?", inst.InstanceID).Delete(&model.CloudInstance{})
		}
	}()

	// 获取统计信息
	stats, err := service.GetInstanceStats()
	if err != nil {
		t.Errorf("获取实例统计信息失败: %v", err)
	}

	if stats.TotalCount < 2 {
		t.Errorf("期望总数至少2，实际%d", stats.TotalCount)
	}

	if stats.RunningCount < 1 {
		t.Errorf("期望运行中至少1，实际%d", stats.RunningCount)
	}

	if stats.StoppedCount < 1 {
		t.Errorf("期望停止中至少1，实际%d", stats.StoppedCount)
	}
}

// TestCloudInstanceService_ClearCache 测试清除缓存
func TestCloudInstanceService_ClearCache(t *testing.T) {
	service := &CloudInstanceService{}

	// 测试清除特定缓存（不会报错）
	service.ClearCache(1, "cn-hangzhou")

	// 测试清除所有缓存（不会报错）
	service.ClearCache(0, "")
}

// TestCloudInstanceService_SyncInstances_InvalidProvider 测试同步无效云厂商
func TestCloudInstanceService_SyncInstances_InvalidProvider(t *testing.T) {
	service := &CloudInstanceService{}

	ctx := context.Background()
	req := request.SyncInstancesReq{
		ProviderID: 99999, // 不存在的云厂商ID
		Region:     "cn-hangzhou",
	}

	_, err := service.SyncInstances(ctx, req)
	if err == nil {
		t.Error("期望返回错误，但没有返回")
	}

	if err.Error() != "云厂商不存在" {
		t.Errorf("期望错误消息='云厂商不存在'，实际='%s'", err.Error())
	}
}

// TestCloudInstanceService_ValidateInstanceOwnership 测试验证实例所有权
func TestCloudInstanceService_ValidateInstanceOwnership(t *testing.T) {
	service := &CloudInstanceService{}

	if global.GVA_DB == nil {
		t.Skip("数据库连接未初始化")
	}

	// 创建测试实例
	testInstance := model.CloudInstance{
		InstanceID:   "test-ownership-001",
		InstanceName: "测试所有权实例",
		ProviderID:   1,
		Region:       "cn-hangzhou",
		Status:       "running",
	}
	service.CreateCloudInstance(testInstance)
	defer global.GVA_DB.Where("instance_id = ?", testInstance.InstanceID).Delete(&model.CloudInstance{})

	// 获取创建的实例
	created, _ := service.GetCloudInstance(testInstance.ID)

	// 验证正确的所有权
	valid, err := service.ValidateInstanceOwnership(created.ID, created.ProviderID)
	if err != nil {
		t.Errorf("验证所有权失败: %v", err)
	}
	if !valid {
		t.Error("期望所有权验证通过，但失败")
	}

	// 验证错误的所有权
	invalid, err := service.ValidateInstanceOwnership(created.ID, 99999)
	if err != nil {
		t.Errorf("验证所有权失败: %v", err)
	}
	if invalid {
		t.Error("期望所有权验证失败，但通过")
	}
}

// 辅助函数：生成测试ID
func testID(prefix string, index int) string {
	return fmt.Sprintf("%s-%03d", prefix, index)
}
