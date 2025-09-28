package test

import (
	"my_project/pkg/model"
	"my_project/pkg/utils/timeutil"
	"my_project/server/internal/dao"
	"my_project/server/internal/infra"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func initTestConfig(t *testing.T) {
	viper.SetConfigName("server_config.test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../server/config") // relative to server/test

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
}

func TestICMPProbeTaskDAO(t *testing.T) {
	initTestConfig(t)

	db := infra.InitMySQL()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	testIP := "8.8.8.8"

	t.Run("Insert ICMP task", func(t *testing.T) {
		task := &model.ICMPProbeTaskDTO{
			IP:        testIP,
			Count:     4,
			Threshold: 100,
			Timeout:   2,
			CreatedAt: timeutil.NowUTC8(),
			UpdatedAt: timeutil.NowUTC8(),
		}

		err := dao.StoreICMPProbeTask(task)
		assert.NoError(t, err)

		got, err := dao.GetICMPProbeTaskByIP(testIP)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, testIP, got.IP)
	})

	t.Run("Get ICMP task by IP", func(t *testing.T) {
		got, err := dao.GetICMPProbeTaskByIP(testIP)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, testIP, got.IP)
		assert.Equal(t, 4, got.Count)
		assert.Equal(t, 100, got.Threshold)
		assert.Equal(t, 2, got.Timeout)
	})

	t.Run("Delete ICMP task by IP", func(t *testing.T) {
		// Delete the task
		err := dao.DeleteICMPProbeTaskByIP(testIP)
		assert.NoError(t, err)

		// Verify the task is gone
		got, err := dao.GetICMPProbeTaskByIP(testIP)
		assert.Error(t, err)
		assert.Nil(t, got, "expected task to be deleted, but found one")
	})
}
