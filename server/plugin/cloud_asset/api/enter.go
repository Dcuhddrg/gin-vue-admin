package api

type ApiGroup struct {
	CloudProviderApi
	CloudInstanceApi
}

var ApiGroupApp = new(ApiGroup)
