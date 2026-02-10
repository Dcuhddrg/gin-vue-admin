package model

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// CloudInstance 云服务器实例模型
type CloudInstance struct {
	global.GVA_MODEL
	// 基础信息
	InstanceID     string `json:"instanceId" gorm:"column:instance_id;type:varchar(128);uniqueIndex:idx_provider_region_instance;comment:实例ID"`
	InstanceName   string `json:"instanceName" gorm:"column:instance_name;type:varchar(255);comment:实例名称"`
	ProviderID     uint   `json:"providerId" gorm:"column:provider_id;index;comment:云厂商ID"`
	Provider       CloudProvider `json:"provider" gorm:"foreignKey:ProviderID"`
	Region         string `json:"region" gorm:"column:region;type:varchar(64);index;comment:所属区域"`

	// 实例配置
	InstanceType   string `json:"instanceType" gorm:"column:instance_type;type:varchar(64);comment:实例规格"`
	CPU            *int   `json:"cpu" gorm:"column:cpu;comment:CPU核数"`
	Memory         *int   `json:"memory" gorm:"column:memory;comment:内存GB"`
	OSName         string `json:"osName" gorm:"column:os_name;type:varchar(255);comment:操作系统名称"`
	OSVersion      string `json:"osVersion" gorm:"column:os_version;type:varchar(64);comment:操作系统版本"`
	DiskSize       *int   `json:"diskSize" gorm:"column:disk_size;comment:磁盘大小GB"`
	PublicIP       string `json:"publicIp" gorm:"column:public_ip;type:varchar(64);comment:公网IP"`
	PrivateIP      string `json:"privateIp" gorm:"column:private_ip;type:varchar(64);comment:内网IP"`

	// 状态信息
	Status         string `json:"status" gorm:"column:status;type:varchar(32);index;comment:实例状态"`
	InstanceStatus string `json:"instanceStatus" gorm:"column:instance_status;type:varchar(64);comment:详细状态描述"`
	CreatedTime    *time.Time `json:"createdTime" gorm:"column:created_time;comment:实例创建时间"`
	ExpiredTime    *time.Time `json:"expiredTime" gorm:"column:expired_time;comment:实例到期时间"`

	// 扩展信息
	ChargeType     string `json:"chargeType" gorm:"column:charge_type;type:varchar(32);comment:付费类型"`
	Remark         string `json:"remark" gorm:"column:remark;type:text;comment:备注"`

	// 同步信息
	LastSyncAt     *time.Time `json:"lastSyncAt" gorm:"column:last_sync_at;comment:最后同步时间"`
	LastSyncStatus string `json:"lastSyncStatus" gorm:"column:last_sync_status;type:varchar(32);comment:最后同步状态"`
}

func (CloudInstance) TableName() string {
	return "gva_cloud_instances"
}

// InstanceVO 实例视图对象，用于API返回
type InstanceVO struct {
	ID             uint        `json:"id"`
	InstanceID     string      `json:"instanceId"`
	InstanceName   string      `json:"instanceName"`
	ProviderName   string      `json:"providerName"`
	ProviderType   string      `json:"providerType"`
	Region         string      `json:"region"`
	InstanceType   string      `json:"instanceType"`
	CPU            *int        `json:"cpu"`
	Memory         *int        `json:"memory"`
	OSName         string      `json:"osName"`
	PublicIP       string      `json:"publicIp"`
	PrivateIP      string      `json:"privateIp"`
	Status         string      `json:"status"`
	CreatedTime    *time.Time  `json:"createdTime"`
	ExpiredTime    *time.Time  `json:"expiredTime"`
	ChargeType     string      `json:"chargeType"`
	LastSyncAt     *time.Time  `json:"lastSyncAt"`
}

// InstanceStatsVO 实例统计视图对象
type InstanceStatsVO struct {
	TotalCount       int64            `json:"totalCount"`
	RunningCount     int64            `json:"runningCount"`
	StoppedCount     int64            `json:"stoppedCount"`
	ProviderStats    []ProviderStatVO `json:"providerStats"`
	RegionStats      []RegionStatVO   `json:"regionStats"`
}

// ProviderStatVO 云厂商统计
type ProviderStatVO struct {
	ProviderName string `json:"providerName"`
	ProviderType string `json:"providerType"`
	Count        int64  `json:"count"`
}

// RegionStatVO 区域统计
type RegionStatVO struct {
	Region string `json:"region"`
	Count  int64  `json:"count"`
}
