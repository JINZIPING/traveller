package test

import (
	"my_project/server/internal/infra"
	"my_project/server/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMigrateTCPProbeTaskTable 测试创建 TCP 表
func TestMigrateTCPProbeTaskTable(t *testing.T) {
	InitTestConfig(t)

	db := infra.InitMySQL()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	err := db.AutoMigrate(&model.TCPProbeTask{})
	assert.NoError(t, err, "AutoMigrate TCPProbeTask table failed")
}
