package probe

import (
	"log"
	"my_project/pkg/model"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func ExecuteICMPProbeTask(task *model.ICMPProbeTask) *model.ICMPProbeResult {
	pinger, err := probing.NewPinger(task.IP)
	if err != nil {
		log.Printf("Failed to create pinger: %v", err)
		return nil
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
