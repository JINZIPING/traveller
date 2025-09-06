package probe

import (
	"log"
	"my_project/pkg/model"
	"net"
	"time"
)

// ExecuteTCPProbeTask 执行 TCP 探测
func ExecuteTCPProbeTask(task *model.TCPProbeTask) *model.TCPProbeResult {
	address := net.JoinHostPort(task.IP, task.Port)
	timeout := time.Duration(task.Timeout) * time.Second

	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Printf("TCP probe failed: %s, error: %v", address, err)
		return &model.TCPProbeResult{
			IP:        task.IP,
			Port:      task.Port,
			Timestamp: time.Now(),
			Success:   false,
			RTT:       0,
		}
	}
	defer conn.Close()

	rtt := time.Since(start)
	log.Printf("TCP probe success: %s, RTT=%v", address, rtt)

	return &model.TCPProbeResult{
		IP:        task.IP,
		Port:      task.Port,
		Timestamp: time.Now(),
		RTT:       rtt,
		Success:   true,
	}
}
