package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/api"
	"github.com/gin-gonic/gin"
)

type ProjectRouter struct{}

func (r *ProjectRouter) InitProjectRouter(Router *gin.RouterGroup) {
	projectRouter := Router.Group("project").Use(middleware.OperationRecord())
	projectRouterWithoutRecord := Router.Group("project")
	projectApi := api.ApiGroupApp.Project
	{
		// 需要记录操作的路由
		projectRouter.POST("createProject", projectApi.CreateProject)
		projectRouter.DELETE("deleteProject", projectApi.DeleteProject)
		projectRouter.PUT("updateProject", projectApi.UpdateProject)
	}
	{
		// 不记录操作的路由（查询类接口）
		projectRouterWithoutRecord.GET("getProjectList", projectApi.GetProjectList)
		projectRouterWithoutRecord.GET("getProject", projectApi.GetProject)
	}
}
