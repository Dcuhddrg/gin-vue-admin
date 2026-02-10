package service

type ServiceGroup struct {
	CloudProviderService
	CloudInstanceService
}

var ServiceGroupApp = new(ServiceGroup)
