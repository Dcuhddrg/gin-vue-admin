package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CloudInstanceSearch 云服务器查询请求
type CloudInstanceSearch struct {
	ProviderID uint   `json:"providerId" form:"providerId"`
	ProjectID  uint   `json:"projectId" form:"projectId"`
	ProviderType string `json:"providerType" form:"providerType"`
	Region     string `json:"region" form:"region"`
	InstanceID string `json:"instanceId" form:"instanceId"`
	InstanceName string `json:"instanceName" form:"instanceName"`
	Status     string `json:"status" form:"status"`
	PublicIP   string `json:"publicIp" form:"publicIp"`
	PrivateIP  string `json:"privateIp" form:"privateIp"`
	request.PageInfo
}

// SyncInstancesReq 同步云服务器实例请求
type SyncInstancesReq struct {
	ProviderID  uint   `json:"providerId" binding:"required"`
	Region      string `json:"region" binding:"required"`
	ForceSync   bool   `json:"forceSync"` // 是否强制同步（忽略缓存）
}

// GetInstanceDetailReq 获取实例详情请求
type GetInstanceDetailReq struct {
	ProviderID  uint   `json:"providerId" binding:"required"`
	InstanceID  string `json:"instanceId" binding:"required"`
	Region      string `json:"region" binding:"required"`
}

// StartInstanceReq 启动实例请求
type StartInstanceReq struct {
	ProviderID  uint   `json:"providerId" binding:"required"`
	InstanceID  string `json:"instanceId" binding:"required"`
	Region      string `json:"region" binding:"required"`
}

// StopInstanceReq 停止实例请求
type StopInstanceReq struct {
	ProviderID  uint   `json:"providerId" binding:"required"`
	InstanceID  string `json:"instanceId" binding:"required"`
	Region      string `json:"region" binding:"required"`
	ForceStop   bool   `json:"forceStop"` // 是否强制停止
}

// RestartInstanceReq 重启实例请求
type RestartInstanceReq struct {
	ProviderID  uint   `json:"providerId" binding:"required"`
	InstanceID  string `json:"instanceId" binding:"required"`
	Region      string `json:"region" binding:"required"`
}

// BatchOperationReq 批量操作请求
type BatchOperationReq struct {
	InstanceIDs []uint `json:"instanceIds" binding:"required,min=1,max=100"`
	Operation   string `json:"operation" binding:"required,oneof=start stop restart delete"`
}
