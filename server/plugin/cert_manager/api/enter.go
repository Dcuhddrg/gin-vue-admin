package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/cert_manager/service"

var (
	ApiGroupApp            = new(apiGroup)
	serviceCertCertificate = service.ServiceGroupApp.CertCertificate
)

type apiGroup struct {
	CertCertificate certCertificate
	CertAdvanced    certAdvanced
}
