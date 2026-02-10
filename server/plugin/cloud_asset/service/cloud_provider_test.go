package service

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	projectManagerModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/project_manager/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model"
	cpRequest "github.com/flipped-aurora/gin-vue-admin/server/plugin/cloud_asset/model/request"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	global.GVA_DB = db
	db.AutoMigrate(&model.CloudProvider{}, &projectManagerModel.Project{})
}

func TestCloudProviderService(t *testing.T) {
	setupTestDB()
	s := &CloudProviderService{}

	// 1. 创建关联项目
	project := projectManagerModel.Project{
		Name: "测试项目",
	}
	global.GVA_DB.Create(&project)

	// 2. 创建云厂商并关联项目
	cp := model.CloudProvider{
		ProjectID: project.ID,
		Type:      "aliyun",
		AK:        "test-ak",
		SK:        "test-sk",
	}
	err := s.CreateCloudProvider(cp)
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}

	// 3. 验证关联查询和数据一致性
	list, total, err := s.GetCloudProviderInfoList(cpRequest.CloudProviderSearch{
		PageInfo: request.PageInfo{Page: 1, PageSize: 10},
	})
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if total != 1 {
		t.Errorf("预期总数 1，实际 %d", total)
	}
	if list[0].Project.Name != "测试项目" {
		t.Errorf("预期项目名称 '测试项目'，实际 '%s'", list[0].Project.Name)
	}

	// 4. 验证级联更新同步 (模拟项目更名)
	global.GVA_DB.Model(&project).Update("name", "更新后的项目")

	listUpdated, _, _ := s.GetCloudProviderInfoList(cpRequest.CloudProviderSearch{
		PageInfo: request.PageInfo{Page: 1, PageSize: 10},
	})
	if listUpdated[0].Project.Name != "更新后的项目" {
		t.Errorf("同步更新失败: 预期 '更新后的项目'，实际 '%s'", listUpdated[0].Project.Name)
	}
}
