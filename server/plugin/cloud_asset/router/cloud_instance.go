package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/api"
	"github.com/gin-gonic/gin"
)

type CloudInstanceRouter struct{}

func (s *CloudInstanceRouter) InitCloudInstanceRouter(Router *gin.RouterGroup) {
	instanceRouter := Router.Group("cloudInstance").Use(middleware.OperationRecord())
	instanceRouterWithoutRecord := Router.Group("cloudInstance")
	instanceApi := api.ApiGroupApp.CloudInstanceApi
	{
		// 需要记录操作的路由
		instanceRouter.POST("createCloudInstance", instanceApi.CreateCloudInstance)
		instanceRouter.DELETE("deleteCloudInstance", instanceApi.DeleteCloudInstance)
		instanceRouter.DELETE("batchDeleteCloudInstance", instanceApi.BatchDeleteCloudInstance)
		instanceRouter.PUT("updateCloudInstance", instanceApi.UpdateCloudInstance)
		instanceRouter.POST("syncInstances", instanceApi.SyncInstances)
		instanceRouter.POST("batchSyncInstances", instanceApi.BatchSyncInstances)
		instanceRouter.POST("clearCache", instanceApi.ClearCache)
	}
	{
		// 不记录操作的路由（查询类接口）
		instanceRouterWithoutRecord.GET("findCloudInstance", instanceApi.FindCloudInstance)
		instanceRouterWithoutRecord.POST("getCloudInstanceList", instanceApi.GetCloudInstanceList)
		instanceRouterWithoutRecord.POST("getCloudInstanceVOList", instanceApi.GetCloudInstanceVOList)
		instanceRouterWithoutRecord.GET("getInstanceStats", instanceApi.GetInstanceStats)
	}
}
