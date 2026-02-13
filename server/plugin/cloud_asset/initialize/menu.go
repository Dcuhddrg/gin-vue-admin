package initialize

import (
	"context"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	parent := model.SysBaseMenu{
		ParentId:  0,
		Path:      "cloudAssetManagement",
		Name:      "cloudAssetManagement",
		Hidden:    false,
		Component: "view/routerHolder.vue",
		Sort:      10,
		Meta:      model.Meta{Title: "云资产管理", Icon: "cloudy"},
	}
	// 云厂商菜单
	cloudProvider := model.SysBaseMenu{
		Path:      "cloudProvider",
		Name:      "cloudProvider",
		Hidden:    false,
		Component: "plugin/cloud_asset/view/cloudProvider.vue",
		Sort:      1,
		Meta:      model.Meta{Title: "云厂商", Icon: "box"},
	}

	// 云服务器实例菜单
	cloudInstance := model.SysBaseMenu{
		Path:      "cloudInstance",
		Name:      "cloudInstance",
		Hidden:    false,
		Component: "plugin/cloud_asset/view/cloudInstance.vue",
		Sort:      2,
		Meta:      model.Meta{Title: "云服务器实例", Icon: "server"},
	}

	utils.RegisterMenus(parent, cloudProvider, cloudInstance)
}
