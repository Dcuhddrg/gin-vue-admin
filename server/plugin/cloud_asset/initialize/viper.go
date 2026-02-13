package initialize

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/config"
	cloudAssetUtils "github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/utils"
)

func Viper(ctx context.Context) {
	// 尝试从 global.GVA_VP 加载配置
	if global.GVA_VP != nil {
		var cfg config.Config
		// 尝试读取 plug-ins.cloud-asset 配置
		if err := global.GVA_VP.UnmarshalKey("plug-ins.cloud-asset", &cfg); err == nil {
			config.ConfigData = cfg
		}
	}

	// 初始化加密模块
	cloudAssetUtils.InitCrypto(config.ConfigData.EncryptionKey)
}
