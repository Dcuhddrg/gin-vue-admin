package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ProjectApi struct{}

var projectService = service.ServiceGroupApp.ProjectService

// CreateProject 创建项目
// @Tags     Project
// @Summary  创建项目
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.Project true "项目模型"
// @Success  200  {object} response.Response{msg=string} "创建成功"
// @Router   /project/createProject [post]
func (a *ProjectApi) CreateProject(c *gin.Context) {
	var project model.Project
	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := projectService.CreateProject(project); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteProject 删除项目
// @Tags     Project
// @Summary  删除项目
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.Project true "项目模型"
// @Success  200  {object} response.Response{msg=string} "删除成功"
// @Router   /project/deleteProject [delete]
func (a *ProjectApi) DeleteProject(c *gin.Context) {
	var project model.Project
	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := projectService.DeleteProject(project); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateProject 更新项目
// @Tags     Project
// @Summary  更新项目
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data body model.Project true "项目模型"
// @Success  200  {object} response.Response{msg=string} "更新成功"
// @Router   /project/updateProject [put]
func (a *ProjectApi) UpdateProject(c *gin.Context) {
	var project model.Project
	err := c.ShouldBindJSON(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := projectService.UpdateProject(project); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// GetProject 获取项目
// @Tags     Project
// @Summary  获取项目
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query model.Project true "项目ID"
// @Success  200  {object} response.Response{data=model.Project,msg=string} "获取成功"
// @Router   /project/getProject [get]
func (a *ProjectApi) GetProject(c *gin.Context) {
	var project model.Project
	err := c.ShouldBindQuery(&project)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	reproject, err := projectService.GetProject(project.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(reproject, "获取成功", c)
	}
}

// GetProjectList 获取项目列表
// @Tags     Project
// @Summary  获取项目列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data query request.ProjectSearch true "页码, 每页大小, 搜索条件"
// @Success  200  {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router   /project/getProjectList [get]
func (a *ProjectApi) GetProjectList(c *gin.Context) {
	var pageInfo request.ProjectSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := projectService.GetProjectList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
