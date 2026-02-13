package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		// CloudProvider APIs
		{Path: "/cloudProvider/createCloudProvider", Description: "创建云厂商", ApiGroup: "云资产管理", Method: "POST"},
		{Path: "/cloudProvider/deleteCloudProvider", Description: "删除云厂商", ApiGroup: "云资产管理", Method: "DELETE"},
		{Path: "/cloudProvider/updateCloudProvider", Description: "更新云厂商", ApiGroup: "云资产管理", Method: "PUT"},
		{Path: "/cloudProvider/findCloudProvider", Description: "根据ID获取云厂商", ApiGroup: "云资产管理", Method: "GET"},
		{Path: "/cloudProvider/getCloudProviderList", Description: "获取云厂商列表", ApiGroup: "云资产管理", Method: "POST"},

		// CloudInstance APIs
		{Path: "/cloudInstance/createCloudInstance", Description: "创建云服务器实例", ApiGroup: "云资产管理", Method: "POST"},
		{Path: "/cloudInstance/deleteCloudInstance", Description: "删除云服务器实例", ApiGroup: "云资产管理", Method: "DELETE"},
		{Path: "/cloudInstance/batchDeleteCloudInstance", Description: "批量删除云服务器实例", ApiGroup: "云资产管理", Method: "DELETE"},
		{Path: "/cloudInstance/updateCloudInstance", Description: "更新云服务器实例", ApiGroup: "云资产管理", Method: "PUT"},
		{Path: "/cloudInstance/findCloudInstance", Description: "根据ID获取云服务器实例", ApiGroup: "云资产管理", Method: "GET"},
		{Path: "/cloudInstance/getCloudInstanceList", Description: "获取云服务器实例列表", ApiGroup: "云资产管理", Method: "POST"},
		{Path: "/cloudInstance/getCloudInstanceVOList", Description: "获取云服务器实例视图列表", ApiGroup: "云资产管理", Method: "POST"},
		{Path: "/cloudInstance/syncInstances", Description: "同步云服务器实例", ApiGroup: "云资产管理", Method: "POST"},
		{Path: "/cloudInstance/batchSyncInstances", Description: "批量同步云服务器实例", ApiGroup: "云资产管理", Method: "POST"},
		{Path: "/cloudInstance/getInstanceStats", Description: "获取实例统计信息", ApiGroup: "云资产管理", Method: "GET"},
		{Path: "/cloudInstance/clearCache", Description: "清除云服务器实例缓存", ApiGroup: "云资产管理", Method: "POST"},
	}
	utils.RegisterApis(entities...)
}
