package test

import (
	"my_project/client/internal/probe"
	"my_project/pkg/model"
	"my_project/pkg/utils/timeutil"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecuteTCPProbeTask(t *testing.T) {
	// Invalid IP
	t.Run("invalid ip", func(t *testing.T) {
		task := &model.TCPProbeTaskDTO{
			IP:        "999.999.999.999",
			Port:      "9999",
			Timeout:   1,
			CreatedAt: timeutil.NowUTC8(),
		}
		result := probe.ExecuteTCPProbeTask(task)
		assert.False(t, result.Success)
	})

	// Invalid Port
	t.Run("invalid port", func(t *testing.T) {
		task := &model.TCPProbeTaskDTO{
			IP:        "127.0.0.1",
			Port:      "abc",
			Timeout:   1,
			CreatedAt: timeutil.NowUTC8(),
		}
		result := probe.ExecuteTCPProbeTask(task)
		assert.False(t, result.Success)
	})

	t.Run("timeout zero", func(t *testing.T) {
		task := &model.TCPProbeTaskDTO{
			IP:        "127.0.0.1",
			Port:      "80",
			Timeout:   0,
			CreatedAt: timeutil.NowUTC8(),
		}
		result := probe.ExecuteTCPProbeTask(task)
		assert.False(t, result.Success)
	})

	// close port
	t.Run("closed port", func(t *testing.T) {
		task := &model.TCPProbeTaskDTO{
			IP:        "127.0.0.1",
			Port:      "65000",
			Timeout:   1,
			CreatedAt: timeutil.NowUTC8(),
		}
		result := probe.ExecuteTCPProbeTask(task)
		assert.False(t, result.Success)
	})

	t.Run("open port", func(t *testing.T) {
		// 启动一个本地 TCP server
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		assert.NoError(t, err)
		defer ln.Close()

		// 提取端口
		_, port, _ := net.SplitHostPort(ln.Addr().String())

		go func() {
			conn, _ := ln.Accept()
			if conn != nil {
				defer conn.Close()
				time.Sleep(1 * time.Millisecond)
			}
		}()

		task := &model.TCPProbeTaskDTO{
			IP:        "127.0.0.1",
			Port:      port,
			Timeout:   2,
			CreatedAt: timeutil.NowUTC8(),
		}
		result := probe.ExecuteTCPProbeTask(task)
		assert.True(t, result.Success)
		assert.Greater(t, result.RTT.Nanoseconds(), int64(0))
	})
}
