package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{
			Path:        "/certCertificate/createCertCertificate",
			Description: "新建证书监控",
			ApiGroup:    "证书监控",
			Method:      "POST",
		},
		{
			Path:        "/certCertificate/deleteCertCertificate",
			Description: "删除证书监控",
			ApiGroup:    "证书监控",
			Method:      "DELETE",
		},
		{
			Path:        "/certCertificate/deleteCertCertificateByIds",
			Description: "批量删除证书监控",
			ApiGroup:    "证书监控",
			Method:      "DELETE",
		},
		{
			Path:        "/certCertificate/updateCertCertificate",
			Description: "更新证书监控",
			ApiGroup:    "证书监控",
			Method:      "PUT",
		},
		{
			Path:        "/certCertificate/findCertCertificate",
			Description: "根据ID获取证书监控",
			ApiGroup:    "证书监控",
			Method:      "GET",
		},
		{
			Path:        "/certCertificate/getCertCertificateList",
			Description: "获取证书监控列表",
			ApiGroup:    "证书监控",
			Method:      "GET",
		},
		{
			Path:        "/certCertificate/probeCertificate",
			Description: "探测证书",
			ApiGroup:    "证书监控",
			Method:      "POST",
		},
		{
			Path:        "/certCertificate/updateAllCertificates",
			Description: "批量更新证书",
			ApiGroup:    "证书监控",
			Method:      "POST",
		},
		{
			Path:        "/certAdvanced/discoverSubdomains",
			Description: "子域名发现",
			ApiGroup:    "域名监控",
			Method:      "POST",
		},
		{
			Path:        "/certAdvanced/getSubdomainTree",
			Description: "获取域名树",
			ApiGroup:    "域名监控",
			Method:      "GET",
		},
		{
			Path:        "/certAdvanced/batchReprobe",
			Description: "批量重新探测",
			ApiGroup:    "域名监控",
			Method:      "POST",
		},
		{
			Path:        "/certAdvanced/batchIgnore",
			Description: "批量忽略",
			ApiGroup:    "域名监控",
			Method:      "POST",
		},
		{
			Path:        "/certAdvanced/getStatistics",
			Description: "获取统计",
			ApiGroup:    "域名监控",
			Method:      "GET",
		},
	}
	utils.RegisterApis(entities...)
}
