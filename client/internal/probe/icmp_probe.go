package probe

import (
	"fmt"
	"log"
	"my_project/pkg/model"
	"net"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// icmpProbeFailure 封装失败返回逻辑
func icmpProbeFailure(task *model.ICMPProbeTask, reason string) *model.ICMPProbeResult {
	log.Printf("ICMP probe failed: ip=%s, reason=%s", task.IP, reason)
	return &model.ICMPProbeResult{
		IP:         task.IP,
		Timestamp:  time.Now(),
		PacketLoss: 100.0, // 失败时认为完全丢包
		MinRTT:     0,
		MaxRTT:     0,
		AvgRTT:     0,
		Threshold:  task.Threshold,
		Success:    false,
	}
}

func ExecuteICMPProbeTask(task *model.ICMPProbeTask) *model.ICMPProbeResult {
	// 校验 IP
	if net.ParseIP(task.IP) == nil {
		return icmpProbeFailure(task, "invalid IP address")
	}

	pinger, err := probing.NewPinger(task.IP)
	if err != nil {
		return icmpProbeFailure(task, fmt.Sprintf("failed to create pinger: %v", err))
	}

	pinger.Count = task.Count
	pinger.Timeout = time.Duration(task.Timeout) * time.Second

	var packetLoss float64
	var minRTT, maxRTT, avgRTT time.Duration

	pinger.OnRecv = func(pkt *probing.Packet) {
		log.Printf("Received packet from %s: time=%v", pkt.IPAddr, pkt.Rtt)
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		log.Printf("Probe finished. Packet loss: %v%%, Min RTT: %v, Max RTT: %v, Avg RTT: %v",
			stats.PacketLoss, stats.MinRtt, stats.MaxRtt, stats.AvgRtt)

		packetLoss = stats.PacketLoss
		minRTT = stats.MinRtt
		maxRTT = stats.MaxRtt
		avgRTT = stats.AvgRtt
	}

	log.Printf("Starting ICMP probe to %s", task.IP)
	pinger.Run()

	// 创建 ProbeResult 结构体并返回
	result := &model.ICMPProbeResult{
		IP:         task.IP,
		Timestamp:  time.Now(),
		PacketLoss: packetLoss,
		MinRTT:     minRTT,
		MaxRTT:     maxRTT,
		AvgRTT:     avgRTT,
		Threshold:  task.Threshold,
		Success:    packetLoss <= float64(task.Threshold),
	}

	return result
}
