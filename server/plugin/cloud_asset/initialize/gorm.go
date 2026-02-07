package initialize

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		model.CloudProvider{},
	)
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("[%s] %+v", "cloud_asset", zap.Error(err)))
	}
}
