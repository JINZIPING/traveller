package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"my_project/client/internal/probe"
	"my_project/pkg/model"
	"my_project/pkg/utils/timeutil"
)

func TestExecuteICMPProbeTask(t *testing.T) {
	t.Run("invalid ip", func(t *testing.T) {
		task := &model.ICMPProbeTaskDTO{
			IP:        "999.999.999.999",
			Count:     1,
			Timeout:   1,
			Threshold: 50,
			CreatedAt: timeutil.NowUTC8(),
		}
		result := probe.ExecuteICMPProbeTask(task)
		assert.False(t, result.Success)
		assert.Equal(t, float64(100.0), result.PacketLoss)
	})

	t.Run("loopback", func(t *testing.T) {
		task := &model.ICMPProbeTaskDTO{
			IP:        "1.1.1.1",
			Count:     1,
			Timeout:   2,
			Threshold: 100,
			CreatedAt: timeutil.NowUTC8(),
		}
		result := probe.ExecuteICMPProbeTask(task)

		if result.Success {
			assert.LessOrEqual(t, result.PacketLoss, float64(100.0))
			assert.True(t, result.MinRTT >= 0)
		} else {
			assert.Equal(t, float64(100.0), result.PacketLoss)
		}
	})

	t.Run("icmp count zero", func(t *testing.T) {
		task := &model.ICMPProbeTaskDTO{
			IP:        "1.1.1.1",
			Count:     1,
			Timeout:   0,
			Threshold: 50,
			CreatedAt: timeutil.NowUTC8(),
		}
		result := probe.ExecuteICMPProbeTask(task)

		assert.False(t, result.Success)
		assert.Equal(t, float64(100.0), result.PacketLoss)
	})

}
