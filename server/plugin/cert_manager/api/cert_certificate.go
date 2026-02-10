package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cert_manager/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cert_manager/model/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type certCertificate struct{}

// CreateCertCertificate 创建证书
// @Tags     CertCertificate
// @Summary  创建证书
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      model.CertCertificate  true  "证书模型"
// @Success  200   {object}  response.Response{msg=string}  "创建成功"
// @Router   /certCertificate/createCertCertificate [post]
func (a *certCertificate) CreateCertCertificate(c *gin.Context) {
	var cert model.CertCertificate
	err := c.ShouldBindJSON(&cert)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = serviceCertCertificate.CreateCertCertificate(&cert)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteCertCertificate 删除证书
// @Tags     CertCertificate
// @Summary  删除证书
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    ID    query     string  true  "ID"
// @Success  200   {object}  response.Response{msg=string}  "删除成功"
// @Router   /certCertificate/deleteCertCertificate [delete]
func (a *certCertificate) DeleteCertCertificate(c *gin.Context) {
	ID := c.Query("ID")
	err := serviceCertCertificate.DeleteCertCertificate(ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteCertCertificateByIds 批量删除证书
// @Tags     CertCertificate
// @Summary  批量删除证书
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    IDs   query     []string  true  "IDs"
// @Success  200   {object}  response.Response{msg=string}  "批量删除成功"
// @Router   /certCertificate/deleteCertCertificateByIds [delete]
func (a *certCertificate) DeleteCertCertificateByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := serviceCertCertificate.DeleteCertCertificateByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateCertCertificate 更新证书
// @Tags     CertCertificate
// @Summary  更新证书
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      model.CertCertificate  true  "证书模型"
// @Success  200   {object}  response.Response{msg=string}  "更新成功"
// @Router   /certCertificate/updateCertCertificate [put]
func (a *certCertificate) UpdateCertCertificate(c *gin.Context) {
	var cert model.CertCertificate
	err := c.ShouldBindJSON(&cert)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = serviceCertCertificate.UpdateCertCertificate(cert)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindCertCertificate 用id查询证书
// @Tags     CertCertificate
// @Summary  用id查询证书
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    ID    query     string  true  "ID"
// @Success  200   {object}  response.Response{data=model.CertCertificate,msg=string}  "查询成功"
// @Router   /certCertificate/findCertCertificate [get]
func (a *certCertificate) FindCertCertificate(c *gin.Context) {
	ID := c.Query("ID")
	reCert, err := serviceCertCertificate.GetCertCertificate(ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(reCert, c)
}

// GetCertCertificateList 分页获取证书列表
// @Tags     CertCertificate
// @Summary  分页获取证书列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  query     request.CertCertificateSearch  true  "页码, 每页大小, 搜索条件"
// @Success  200   {object}  response.Response{data=response.PageResult,msg=string}  "获取成功"
// @Router   /certCertificate/getCertCertificateList [get]
func (a *certCertificate) GetCertCertificateList(c *gin.Context) {
	var pageInfo request.CertCertificateSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceCertCertificate.GetCertCertificateList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

type ProbeCertificateRequest struct {
	Domain string `json:"domain" binding:"required"`
}

// ProbeCertificate 探测证书
// @Tags     CertCertificate
// @Summary  探测证书
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      ProbeCertificateRequest  true  "域名"
// @Success  200   {object}  response.Response{msg=string}  "探测成功"
// @Router   /certCertificate/probeCertificate [post]
func (a *certCertificate) ProbeCertificate(c *gin.Context) {
	var req ProbeCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := serviceCertCertificate.ProbeAndUpdateCertificate(req.Domain)
	if err != nil {
		global.GVA_LOG.Error("探测证书失败!", zap.Error(err))
		response.FailWithMessage("探测证书失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("探测成功", c)
}

// UpdateAllCertificates 批量更新证书
// @Tags     CertCertificate
// @Summary  批量更新证书
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Success  200   {object}  response.Response{msg=string}  "批量更新成功"
// @Router   /certCertificate/updateAllCertificates [post]
func (a *certCertificate) UpdateAllCertificates(c *gin.Context) {
	err := serviceCertCertificate.UpdateAllCertificates()
	if err != nil {
		global.GVA_LOG.Error("批量更新证书失败!", zap.Error(err))
		response.FailWithMessage("批量更新证书失败", c)
		return
	}
	response.OkWithMessage("批量更新成功", c)
}
