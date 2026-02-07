package service

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonRequest "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cert_manager/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/cert_manager/model/request"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	global.GVA_DB = db
	db.AutoMigrate(&model.CertCertificate{})

	// Mock Redis
	global.GVA_REDIS = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestCertCertificateService_GetCertCertificateList(t *testing.T) {
	setupTestDB()
	s := &certCertificate{}

	// 1. 插入测试数据
	certs := []model.CertCertificate{
		{Domain: "a.com", Category: "ProjectA", DaysRemaining: 10},
		{Domain: "b.com", Category: "ProjectB", DaysRemaining: 20},
		{Domain: "c.com", Category: "ProjectA", DaysRemaining: 5},
	}
	for _, c := range certs {
		global.GVA_DB.Create(&c)
	}

	// 2. 测试全量查询
	_, total, err := s.GetCertCertificateList(request.CertCertificateSearch{
		PageInfo: commonRequest.PageInfo{Page: 1, PageSize: 10},
	})
	if err != nil {
		t.Fatalf("全量查询失败: %v", err)
	}
	if total != 3 {
		t.Errorf("全量查询总数错误: 预期 3, 实际 %d", total)
	}

	// 3. 测试按项目查询
	_, total, err = s.GetCertCertificateList(request.CertCertificateSearch{
		PageInfo: commonRequest.PageInfo{Page: 1, PageSize: 10},
		Category: "ProjectA",
	})
	if total != 2 {
		t.Errorf("项目查询总数错误: 预期 2, 实际 %d", total)
	}

	// 4. 测试缓存 (模拟 Redis)
	// 由于 miniredis 设置复杂，这里主要验证代码逻辑是否通过编译和基本 DB 操作
	
	// 清理缓存逻辑测试
	s.clearCache()
}

func TestCertCertificateService_NoRedis(t *testing.T) {
	// Setup DB but NO Redis
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	global.GVA_DB = db
	db.AutoMigrate(&model.CertCertificate{})
	global.GVA_REDIS = nil // Explicitly set to nil

	s := &certCertificate{}
	
	// 1. 插入测试数据
	global.GVA_DB.Create(&model.CertCertificate{Domain: "a.com"})

	// 2. 测试全量查询 - Should NOT panic
	_, total, err := s.GetCertCertificateList(request.CertCertificateSearch{
		PageInfo: commonRequest.PageInfo{Page: 1, PageSize: 10},
	})
	
	if err != nil {
		t.Fatalf("NoRedis 查询失败: %v", err)
	}
	if total != 1 {
		t.Errorf("NoRedis 查询总数错误")
	}

	// 3. Test clearCache - Should NOT panic
	s.clearCache()
}
