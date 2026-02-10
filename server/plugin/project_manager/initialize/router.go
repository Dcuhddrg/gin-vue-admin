package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/router"
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	group := engine.Group("")
	router.RouterGroupApp.InitProjectRouter(group)
}
