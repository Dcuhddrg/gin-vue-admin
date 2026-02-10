package api

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cert_manager/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cert_manager/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type certAdvanced struct{}

var CertAdvancedApi = new(certAdvanced)

type DiscoverSubdomainsRequest struct {
	RootDomain string `json:"rootDomain" binding:"required"`
	Category   string `json:"category"`
	DeepScan   bool   `json:"deepScan"`
}

type BatchDiscoverRequest struct {
	RootDomains []string `json:"rootDomains" binding:"required"`
	Category    string   `json:"category"`
	DeepScan    bool     `json:"deepScan"`
}

type ReProbeRequest struct {
	RootDomain string `json:"rootDomain" binding:"required"`
}

type IgnoreDomainRequest struct {
	DomainID uint `json:"domainId" binding:"required"`
	Ignore   bool `json:"ignore"`
}

// DiscoverSubdomains 子域名发现
// @Tags     CertAdvanced
// @Summary  子域名发现
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      DiscoverSubdomainsRequest  true  "发现参数"
// @Success  200   {object}  response.Response{msg=string}  "子域名发现成功"
// @Router   /certAdvanced/discoverSubdomains [post]
func (a *certAdvanced) DiscoverSubdomains(c *gin.Context) {
	var req DiscoverSubdomainsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := service.ServiceGroupApp.SubdomainDiscovery.DiscoverAndStoreSubdomains(req.RootDomain, req.Category, req.DeepScan); err != nil {
		global.GVA_LOG.Error("子域名发现失败", zap.Error(err))
		response.FailWithMessage("子域名发现失败", c)
		return
	}

	response.OkWithMessage("子域名发现成功", c)
}

// BatchDiscoverSubdomains 批量子域名发现
// @Tags     CertAdvanced
// @Summary  批量子域名发现
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      BatchDiscoverRequest  true  "批量发现参数"
// @Success  200   {object}  response.Response{msg=string}  "批量子域名发现成功"
// @Router   /certAdvanced/batchDiscoverSubdomains [post]
func (a *certAdvanced) BatchDiscoverSubdomains(c *gin.Context) {
	var req BatchDiscoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := service.ServiceGroupApp.SubdomainDiscovery.BatchDiscoverSubdomains(req.RootDomains, req.Category, req.DeepScan); err != nil {
		global.GVA_LOG.Error("批量子域名发现失败", zap.Error(err))
		response.FailWithMessage("批量子域名发现失败", c)
		return
	}

	response.OkWithMessage("批量子域名发现成功", c)
}

// GetDomainTree 获取域名树
// @Tags     CertAdvanced
// @Summary  获取域名树
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    rootDomain  query     string  true  "根域名"
// @Success  200         {object}  response.Response{data=map[string]interface{},msg=string}  "获取成功"
// @Router   /certAdvanced/getDomainTree [get]
func (a *certAdvanced) GetDomainTree(c *gin.Context) {
	rootDomain := c.Query("rootDomain")
	if rootDomain == "" {
		response.FailWithMessage("rootDomain 参数必填", c)
		return
	}

	tree, err := service.ServiceGroupApp.SubdomainDiscovery.GetDomainTree(rootDomain)
	if err != nil {
		global.GVA_LOG.Error("获取域名树失败", zap.Error(err))
		response.FailWithMessage("获取域名树失败", c)
		return
	}

	response.OkWithData(tree, c)
}

// ExportSubdomainReport 导出子域名报告
// @Tags     CertAdvanced
// @Summary  导出子域名报告
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/octet-stream
// @Param    rootDomain  query     string  true  "根域名"
// @Success  200         {file}    file    "文件流"
// @Router   /certAdvanced/exportSubdomainReport [get]
func (a *certAdvanced) ExportSubdomainReport(c *gin.Context) {
	rootDomain := c.Query("rootDomain")
	if rootDomain == "" {
		response.FailWithMessage("rootDomain 参数必填", c)
		return
	}

	filePath, fileName, err := service.ServiceGroupApp.SubdomainDiscovery.ExportSubdomainReport(rootDomain)
	if err != nil {
		global.GVA_LOG.Error("导出报告失败", zap.Error(err))
		response.FailWithMessage("导出报告失败", c)
		return
	}
	defer os.Remove(filePath)

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
}

// ReProbeDomainTree 重新探测域名树
// @Tags     CertAdvanced
// @Summary  重新探测域名树
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      ReProbeRequest  true  "探测参数"
// @Success  200   {object}  response.Response{msg=string}  "重新探测成功"
// @Router   /certAdvanced/reProbeDomainTree [post]
func (a *certAdvanced) ReProbeDomainTree(c *gin.Context) {
	var req ReProbeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	assets, err := service.ServiceGroupApp.DomainAssetService.GetDomainAssetsByRootDomain(req.RootDomain)
	if err != nil {
		global.GVA_LOG.Error("查询域名资产失败", zap.Error(err))
		response.FailWithMessage("查询域名资产失败", c)
		return
	}

	for _, asset := range assets {
		if asset.IsIgnored {
			continue
		}
		if err := service.ServiceGroupApp.CertCertificate.ProbeAndUpdateCertificate(asset.Domain); err != nil {
			global.GVA_LOG.Error("探测域名失败", zap.String("domain", asset.Domain), zap.Error(err))
		}
	}

	response.OkWithMessage("重新探测成功", c)
}

// IgnoreDomain 忽略/取消忽略域名
// @Tags     CertAdvanced
// @Summary  忽略/取消忽略域名
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      IgnoreDomainRequest  true  "忽略参数"
// @Success  200   {object}  response.Response{msg=string}  "更新成功"
// @Router   /certAdvanced/ignoreDomain [post]
func (a *certAdvanced) IgnoreDomain(c *gin.Context) {
	var req IgnoreDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := service.ServiceGroupApp.DomainAssetService.UpdateIgnoreStatus(req.DomainID, req.Ignore); err != nil {
		global.GVA_LOG.Error("更新忽略状态失败", zap.Error(err))
		response.FailWithMessage("更新忽略状态失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// GetDomainAssetList 获取域名资产列表
// @Tags     CertAdvanced
// @Summary  获取域名资产列表
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  query     request.DomainAssetSearch  true  "搜索条件"
// @Success  200   {object}  response.Response{data=response.PageResult,msg=string}  "获取成功"
// @Router   /certAdvanced/getDomainAssetList [get]
func (a *certAdvanced) GetDomainAssetList(c *gin.Context) {
	var pageInfo request.DomainAssetSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := service.ServiceGroupApp.DomainAssetService.GetDomainAssetList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取域名资产列表失败", zap.Error(err))
		response.FailWithMessage("获取域名资产列表失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

type BatchIdsRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

// BatchReprobe 批量重新探测
// @Tags     CertAdvanced
// @Summary  批量重新探测
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      BatchIdsRequest  true  "IDs"
// @Success  200   {object}  response.Response{msg=string}  "批量探测完成"
// @Router   /certAdvanced/batchReprobe [post]
func (a *certAdvanced) BatchReprobe(c *gin.Context) {
	var req BatchIdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	assets, err := service.ServiceGroupApp.DomainAssetService.GetDomainAssetsByIds(req.IDs)
	if err != nil {
		global.GVA_LOG.Error("查询域名记录失败", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}

	if len(assets) == 0 {
		response.OkWithMessage("无匹配记录", c)
		return
	}

	var wg sync.WaitGroup
	errChan := make(chan string, len(assets))

	for _, asset := range assets {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			if err := service.ServiceGroupApp.CertCertificate.ProbeAndUpdateCertificate(domain); err != nil {
				global.GVA_LOG.Error("重新探测失败", zap.String("domain", domain), zap.Error(err))
				errChan <- fmt.Sprintf("%s: %v", domain, err)
			}
		}(asset.Domain)
	}

	wg.Wait()
	close(errChan)

	var errMsgs []string
	for msg := range errChan {
		errMsgs = append(errMsgs, msg)
	}

	if len(errMsgs) > 0 {
		response.OkWithMessage(fmt.Sprintf("部分探测失败: %s", strings.Join(errMsgs, "; ")), c)
		return
	}

	response.OkWithMessage("批量探测完成", c)
}

// BatchIgnore 批量忽略
// @Tags     CertAdvanced
// @Summary  批量忽略
// @Security ApiKeyAuth
// @accept   application/json
// @Produce  application/json
// @Param    data  body      BatchIdsRequest  true  "IDs"
// @Success  200   {object}  response.Response{msg=string}  "批量忽略成功"
// @Router   /certAdvanced/batchIgnore [post]
func (a *certAdvanced) BatchIgnore(c *gin.Context) {
	var req BatchIdsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := service.ServiceGroupApp.DomainAssetService.BatchUpdateIgnoreStatus(req.IDs, true); err != nil {
		global.GVA_LOG.Error("批量忽略失败", zap.Error(err))
		response.FailWithMessage("操作失败", c)
		return
	}

	response.OkWithMessage("批量忽略成功", c)
}
