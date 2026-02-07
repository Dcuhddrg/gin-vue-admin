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
	child := model.SysBaseMenu{
		Path:      "cloudProvider",
		Name:      "cloudProvider",
		Hidden:    false,
		Component: "plugin/cloud_asset/view/cloudProvider.vue",
		Sort:      1,
		Meta:      model.Meta{Title: "云厂商", Icon: "box"},
	}
	utils.RegisterMenus(parent, child)
}
