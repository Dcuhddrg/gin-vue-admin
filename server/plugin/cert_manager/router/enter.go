package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/cert_manager/api"

var (
	RouterGroupApp     = new(routerGroup)
	apiCertCertificate = api.ApiGroupApp.CertCertificate
	apiCertAdvanced    = api.ApiGroupApp.CertAdvanced
)

type routerGroup struct {
	CertCertificate    certCertificate
	CertAdvancedRouter certAdvancedRouter
}
