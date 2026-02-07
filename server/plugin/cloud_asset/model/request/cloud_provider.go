package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type CloudProviderSearch struct {
	ProjectId uint   `json:"projectId" form:"projectId"`
	Type      string `json:"type" form:"type"`
	Status    *int   `json:"status" form:"status"`
	request.PageInfo
}

type GetRegionsReq struct {
	Provider  string `json:"provider" binding:"required,oneof=aliyun tencent aws"`
	AccessKey string `json:"accessKey" binding:"required"`
	SecretKey string `json:"secretKey" binding:"required"`
}
