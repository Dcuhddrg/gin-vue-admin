package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/api"
	"github.com/gin-gonic/gin"
)

type CloudProviderRouter struct{}

func (s *CloudProviderRouter) InitCloudProviderRouter(Router *gin.RouterGroup) {
	cpRouter := Router.Group("cloudProvider").Use(middleware.OperationRecord())
	cpRouterWithoutRecord := Router.Group("cloudProvider")
	cpApi := api.ApiGroupApp.CloudProviderApi
	{
		cpRouter.POST("createCloudProvider", cpApi.CreateCloudProvider)   // 新建云厂商
		cpRouter.DELETE("deleteCloudProvider", cpApi.DeleteCloudProvider) // 删除云厂商
		cpRouter.PUT("updateCloudProvider", cpApi.UpdateCloudProvider)    // 更新云厂商
		cpRouter.POST("getRegions", cpApi.GetRegions)                     // 获取可用区
	}
	{
		cpRouterWithoutRecord.GET("findCloudProvider", cpApi.FindCloudProvider)        // 根据ID获取云厂商
		cpRouterWithoutRecord.POST("getCloudProviderList", cpApi.GetCloudProviderList) // 获取云厂商列表
	}
}
