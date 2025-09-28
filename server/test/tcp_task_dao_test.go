package test

import (
	pkgModel "my_project/pkg/model"
	"my_project/pkg/utils/timeutil"
	"my_project/server/internal/dao"
	"my_project/server/internal/infra"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTCPProbeTaskDAO(t *testing.T) {
	initTestConfig(t)

	db := infra.InitMySQL()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	testIP := "1.1.1.1"

	t.Run("Insert TCP task", func(t *testing.T) {
		task := &pkgModel.TCPProbeTaskDTO{
			IP:        testIP,
			Port:      "80",
			Timeout:   5,
			CreatedAt: timeutil.NowUTC8(),
			UpdatedAt: timeutil.NowUTC8(),
		}
		err := dao.StoreTCPProbeTask(task)
		assert.NoError(t, err)

		got, err := dao.GetTCPProbeTaskByIP(testIP)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, testIP, got.IP)
	})

	t.Run("Get TCP task by IP", func(t *testing.T) {
		got, err := dao.GetTCPProbeTaskByIP(testIP)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, testIP, got.IP)
		assert.Equal(t, "80", got.Port)
		assert.Equal(t, 5, got.Timeout)
	})

	t.Run("Delete TCP task by IP", func(t *testing.T) {
		err := dao.DeleteTCPProbeTaskByIP(testIP)
		assert.NoError(t, err)

		got, err := dao.GetTCPProbeTaskByIP(testIP)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, got, "expected task to be deleted, but found one")
	})
}
