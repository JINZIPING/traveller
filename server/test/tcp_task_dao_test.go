package test

import (
	"my_project/pkg/model"
	"my_project/pkg/utils/timeutil"
	"my_project/server/internal/dao"
	"my_project/server/internal/infra"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPProbeTaskDAO(t *testing.T) {
	initTestConfig(t)
	db := infra.InitMySQL()
	defer db.Close()

	testIP := "1.1.1.1"
	_, _ = db.Exec("DELETE FROM tcp_probe_tasks WHERE ip = ?", testIP)

	t.Run("Insert TCP task", func(t *testing.T) {
		task := &model.TCPProbeTask{
			IP:        testIP,
			Port:      "80", // 注意这里是 string
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
		assert.NoError(t, err)
		assert.Nil(t, got, "expected task to be deleted, but found one")
	})
}
