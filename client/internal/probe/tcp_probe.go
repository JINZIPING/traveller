package probe

import (
	"fmt"
	"log"
	"my_project/pkg/model"
	"net"
	"strconv"
	"time"
)

// tcpProbeFailure
func tcpProbeFailure(task *model.TCPProbeTask, reason string) *model.TCPProbeResult {
	log.Printf("TCP probe failed: ip=%s, port=%s, reason=%s", task.IP, task.Port, reason)
	return &model.TCPProbeResult{
		IP:        task.IP,
		Port:      task.Port,
		Timestamp: time.Now(),
		Success:   false,
		RTT:       0,
	}
}

// ExecuteTCPProbeTask 执行 TCP 探测
func ExecuteTCPProbeTask(task *model.TCPProbeTask) *model.TCPProbeResult {
	// 校验 IP
	if net.ParseIP(task.IP) == nil {
		return tcpProbeFailure(task, "invalid IP address")
	}

	// 校验端口
	portNum, err := strconv.Atoi(task.Port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return tcpProbeFailure(task, "invalid port number")
	}

	address := net.JoinHostPort(task.IP, task.Port)
	timeout := time.Duration(task.Timeout) * time.Second

	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return tcpProbeFailure(task, fmt.Sprintf("dial error: %v", err))
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
