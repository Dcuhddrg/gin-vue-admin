package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CloudProviderApi struct{}

var cpService = service.ServiceGroupApp.CloudProviderService

// CreateCloudProvider 创建云厂商
// @Tags     CloudProvider
// @Summary  创建云厂商
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.CloudProvider true "创建云厂商"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /cloudProvider/createCloudProvider [post]
func (a *CloudProviderApi) CreateCloudProvider(c *gin.Context) {
	var cp model.CloudProvider
	err := c.ShouldBindJSON(&cp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = cpService.CreateCloudProvider(cp)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteCloudProvider 删除云厂商
// @Tags     CloudProvider
// @Summary  删除云厂商
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.CloudProvider true "删除云厂商"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /cloudProvider/deleteCloudProvider [delete]
func (a *CloudProviderApi) DeleteCloudProvider(c *gin.Context) {
	var cp model.CloudProvider
	err := c.ShouldBindJSON(&cp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = cpService.DeleteCloudProvider(cp)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateCloudProvider 更新云厂商
// @Tags     CloudProvider
// @Summary  更新云厂商
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.CloudProvider true "更新云厂商"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /cloudProvider/updateCloudProvider [put]
func (a *CloudProviderApi) UpdateCloudProvider(c *gin.Context) {
	var cp model.CloudProvider
	err := c.ShouldBindJSON(&cp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = cpService.UpdateCloudProvider(cp)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCloudProvider 用id查询云厂商
// @Tags     CloudProvider
// @Summary  用id查询云厂商
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query model.CloudProvider true "用id查询云厂商"
// @Success  200  {object} response.Response{data=model.CloudProvider,msg=string} "查询成功"
// @Router   /cloudProvider/findCloudProvider [get]
func (a *CloudProviderApi) FindCloudProvider(c *gin.Context) {
	var cp model.CloudProvider
	err := c.ShouldBindQuery(&cp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	recp, err := cpService.GetCloudProvider(cp.ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(recp, "查询成功", c)
}

// GetCloudProviderList 分页获取云厂商列表
// @Tags     CloudProvider
// @Summary  分页获取云厂商列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.CloudProviderSearch true "分页获取云厂商列表"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /cloudProvider/getCloudProviderList [post]
func (a *CloudProviderApi) GetCloudProviderList(c *gin.Context) {
	var info request.CloudProviderSearch
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := cpService.GetCloudProviderInfoList(info)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

// GetRegions 获取云厂商可用区
// @Tags     CloudProvider
// @Summary  获取云厂商可用区
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body request.GetRegionsReq true "获取可用区参数"
// @Success  200  {object} response.Response{data=[]model.RegionVO,msg=string} "获取成功"
// @Router   /cloudProvider/getRegions [post]
func (a *CloudProviderApi) GetRegions(c *gin.Context) {
	var req request.GetRegionsReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 敏感信息脱敏日志
	global.GVA_LOG.Info("获取可用区请求",
		zap.String("provider", req.Provider),
		zap.String("accessKey", "***"),
		zap.String("secretKey", "***"),
	)

	regions, err := cpService.GetRegions(req.Provider, req.AccessKey, req.SecretKey)
	if err != nil {
		// 记录错误但不打印完整密钥
		global.GVA_LOG.Error("获取可用区失败",
			zap.String("provider", req.Provider),
			zap.Error(err),
		)
		response.FailWithMessage("获取可用区失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(regions, "获取成功", c)
}
