package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{Path: "/cloudProvider/createCloudProvider", Description: "创建云厂商", ApiGroup: "云资产管理", Method: "POST"},
		{Path: "/cloudProvider/deleteCloudProvider", Description: "删除云厂商", ApiGroup: "云资产管理", Method: "DELETE"},
		{Path: "/cloudProvider/updateCloudProvider", Description: "更新云厂商", ApiGroup: "云资产管理", Method: "PUT"},
		{Path: "/cloudProvider/findCloudProvider", Description: "根据ID获取云厂商", ApiGroup: "云资产管理", Method: "GET"},
		{Path: "/cloudProvider/getCloudProviderList", Description: "获取云厂商列表", ApiGroup: "云资产管理", Method: "POST"},
	}
	utils.RegisterApis(entities...)
}
