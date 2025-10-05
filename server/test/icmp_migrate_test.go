package test

import (
	"my_project/server/internal/infra"
	"my_project/server/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMigrateICMPProbeTaskTable 测试建表
func TestMigrateICMPProbeTaskTable(t *testing.T) {
	initTestConfig(t)

	db := infra.InitMySQL()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	err := db.AutoMigrate(&model.ICMPProbeTask{})
	assert.NoError(t, err, "AutoMigrate ICMPProbeTask table failed")
}
