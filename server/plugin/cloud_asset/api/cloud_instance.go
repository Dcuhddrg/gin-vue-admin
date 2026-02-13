package api

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CloudInstanceApi struct{}

var instanceService = service.ServiceGroupApp.CloudInstanceService

// CreateCloudInstance 创建云服务器实例
// @Tags     CloudInstance
// @Summary  创建云服务器实例
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.CloudInstance true "创建云服务器实例"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /cloudInstance/createCloudInstance [post]
func (a *CloudInstanceApi) CreateCloudInstance(c *gin.Context) {
	var instance model.CloudInstance
	err := c.ShouldBindJSON(&instance)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = instanceService.CreateCloudInstance(instance)
	if err != nil {
		global.GVA_LOG.Error("创建云服务器实例失败!", zap.Error(err))
		response.FailWithMessage("创建云服务器实例失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteCloudInstance 删除云服务器实例
// @Tags     CloudInstance
// @Summary  删除云服务器实例
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.CloudInstance true "删除云服务器实例"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /cloudInstance/deleteCloudInstance [delete]
func (a *CloudInstanceApi) DeleteCloudInstance(c *gin.Context) {
	var instance model.CloudInstance
	err := c.ShouldBindJSON(&instance)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = instanceService.DeleteCloudInstance(instance)
	if err != nil {
		global.GVA_LOG.Error("删除云服务器实例失败!", zap.Error(err))
		response.FailWithMessage("删除云服务器实例失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// BatchDeleteCloudInstance 批量删除云服务器实例
// @Tags     CloudInstance
// @Summary  批量删除云服务器实例
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.BatchOperationReq true "批量删除请求"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /cloudInstance/batchDeleteCloudInstance [delete]
func (a *CloudInstanceApi) BatchDeleteCloudInstance(c *gin.Context) {
	var req request.BatchOperationReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Operation != "delete" {
		response.FailWithMessage("操作类型错误", c)
		return
	}
	err = instanceService.BatchDeleteCloudInstance(req.InstanceIDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除云服务器实例失败!", zap.Error(err))
		response.FailWithMessage("批量删除云服务器实例失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCloudInstance 更新云服务器实例
// @Tags     CloudInstance
// @Summary  更新云服务器实例
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.CloudInstance true "更新云服务器实例"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /cloudInstance/updateCloudInstance [put]
func (a *CloudInstanceApi) UpdateCloudInstance(c *gin.Context) {
	var instance model.CloudInstance
	err := c.ShouldBindJSON(&instance)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = instanceService.UpdateCloudInstance(instance)
	if err != nil {
		global.GVA_LOG.Error("更新云服务器实例失败!", zap.Error(err))
		response.FailWithMessage("更新云服务器实例失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCloudInstance 用id查询云服务器实例
// @Tags     CloudInstance
// @Summary  用id查询云服务器实例
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query model.CloudInstance true "用id查询云服务器实例"
// @Success  200  {object} response.Response{data=model.CloudInstance,msg=string} "查询成功"
// @Router   /cloudInstance/findCloudInstance [get]
func (a *CloudInstanceApi) FindCloudInstance(c *gin.Context) {
	var req struct {
		ID uint `form:"ID" json:"ID"`
	}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	reinstance, err := instanceService.GetCloudInstance(req.ID)
	if err != nil {
		global.GVA_LOG.Error("查询云服务器实例失败!", zap.Error(err))
		response.FailWithMessage("查询云服务器实例失败", c)
		return
	}
	response.OkWithDetailed(reinstance, "查询成功", c)
}

// GetCloudInstanceList 分页获取云服务器实例列表
// @Tags     CloudInstance
// @Summary  分页获取云服务器实例列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.CloudInstanceSearch true "分页获取云服务器实例列表"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /cloudInstance/getCloudInstanceList [post]
func (a *CloudInstanceApi) GetCloudInstanceList(c *gin.Context) {
	var info request.CloudInstanceSearch
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := instanceService.GetCloudInstanceInfoList(info)
	if err != nil {
		global.GVA_LOG.Error("获取云服务器实例列表失败!", zap.Error(err))
		response.FailWithMessage("获取云服务器实例列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

// GetCloudInstanceVOList 分页获取云服务器实例视图列表
// @Tags     CloudInstance
// @Summary  分页获取云服务器实例视图列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.CloudInstanceSearch true "分页获取云服务器实例视图列表"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /cloudInstance/getCloudInstanceVOList [post]
func (a *CloudInstanceApi) GetCloudInstanceVOList(c *gin.Context) {
	var info request.CloudInstanceSearch
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := instanceService.GetCloudInstanceVOList(info)
	if err != nil {
		global.GVA_LOG.Error("获取云服务器实例列表失败!", zap.Error(err))
		response.FailWithMessage("获取云服务器实例列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

// SyncInstances 同步云服务器实例
// @Tags     CloudInstance
// @Summary  同步云服务器实例
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.SyncInstancesReq true "同步云服务器实例请求"
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "同步成功"
// @Router   /cloudInstance/syncInstances [post]
func (a *CloudInstanceApi) SyncInstances(c *gin.Context) {
	var req request.SyncInstancesReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	ctx := context.Background()

	// 敏感信息脱敏日志
	global.GVA_LOG.Info("开始同步云服务器实例",
		zap.Uint("providerId", req.ProviderID),
		zap.String("region", req.Region),
		zap.Bool("forceSync", req.ForceSync),
	)

	count, err := instanceService.SyncInstances(ctx, req)
	if err != nil {
		global.GVA_LOG.Error("同步云服务器实例失败!",
			zap.Uint("providerId", req.ProviderID),
			zap.String("region", req.Region),
			zap.Error(err),
		)
		response.FailWithMessage("同步云服务器实例失败: "+err.Error(), c)
		return
	}

	global.GVA_LOG.Info("同步云服务器实例成功",
		zap.Uint("providerId", req.ProviderID),
		zap.String("region", req.Region),
		zap.Int("count", count),
	)

	response.OkWithDetailed(gin.H{
		"count":    count,
		"providerId": req.ProviderID,
		"region":     req.Region,
	}, "同步成功", c)
}

// BatchSyncInstances 批量同步云服务器实例
// @Tags     CloudInstance
// @Summary  批量同步云服务器实例
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body struct {
//   ProviderID uint   `json:"providerId" binding:"required"`
//   Regions    []string `json:"regions" binding:"required"`
//   ForceSync  bool   `json:"forceSync"`
// } true "批量同步云服务器实例请求"
// @Success  200  {object} response.Response{data=map[string]interface{},msg=string} "同步成功"
// @Router   /cloudInstance/batchSyncInstances [post]
func (a *CloudInstanceApi) BatchSyncInstances(c *gin.Context) {
	var req struct {
		ProviderID uint   `json:"providerId" binding:"required"`
		Regions    []string `json:"regions" binding:"required"`
		ForceSync  bool   `json:"forceSync"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	ctx := context.Background()

	global.GVA_LOG.Info("开始批量同步云服务器实例",
		zap.Uint("providerId", req.ProviderID),
		zap.Int("regionCount", len(req.Regions)),
		zap.Bool("forceSync", req.ForceSync),
	)

	results, err := instanceService.BatchSyncInstances(ctx, req.ProviderID, req.Regions, req.ForceSync)
	if err != nil {
		global.GVA_LOG.Error("批量同步云服务器实例失败!",
			zap.Uint("providerId", req.ProviderID),
			zap.Error(err),
		)
		response.FailWithMessage("批量同步云服务器实例失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"results":   results,
		"totalSync": len(results),
	}, "批量同步成功", c)
}

// GetInstanceStats 获取实例统计信息
// @Tags     CloudInstance
// @Summary  获取实例统计信息
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200  {object} response.Response{data=model.InstanceStatsVO,msg=string} "获取成功"
// @Router   /cloudInstance/getInstanceStats [get]
func (a *CloudInstanceApi) GetInstanceStats(c *gin.Context) {
	stats, err := instanceService.GetInstanceStats()
	if err != nil {
		global.GVA_LOG.Error("获取实例统计信息失败!", zap.Error(err))
		response.FailWithMessage("获取实例统计信息失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(stats, "获取成功", c)
}

// ClearCache 清除云服务器实例缓存
// @Tags     CloudInstance
// @Summary  清除云服务器实例缓存
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body struct {
//   ProviderID uint   `json:"providerId"`
//   Region     string `json:"region"`
// } true "清除缓存请求"
// @Success  200  {object} response.Response{msg=string} "清除成功"
// @Router   /cloudInstance/clearCache [post]
func (a *CloudInstanceApi) ClearCache(c *gin.Context) {
	var req struct {
		ProviderID uint   `json:"providerId"`
		Region     string `json:"region"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	instanceService.ClearCache(req.ProviderID, req.Region)

	global.GVA_LOG.Info("清除云服务器实例缓存",
		zap.Uint("providerId", req.ProviderID),
		zap.String("region", req.Region),
	)

	response.OkWithMessage("清除缓存成功", c)
}
