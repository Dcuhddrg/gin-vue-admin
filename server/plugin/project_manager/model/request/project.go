package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// ProjectSearch 项目搜索参数
type ProjectSearch struct {
	Name string `json:"name" form:"name"`
	request.PageInfo
}
