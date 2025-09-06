package model

import "time"

// ICMPProbeResult 定义探测结果的结构体
type ICMPProbeResult struct {
	IP         string        `json:"ip"`          // 探测目标的IP地址
	Timestamp  time.Time     `json:"timestamp"`   // 探测完成的时间戳
	PacketLoss float64       `json:"packet_loss"` // 丢包率
	MinRTT     time.Duration `json:"min_rtt"`     // 最小往返时间
	MaxRTT     time.Duration `json:"max_rtt"`     // 最大往返时间
	AvgRTT     time.Duration `json:"avg_rtt"`     // 平均往返时间
	Threshold  int           `json:"threshold"`   // 丢包率阈值（来自任务）
	Success    bool          `json:"success"`     // 探测是否成功
}

// TCPProbeResult 定义TCP探测结果
type TCPProbeResult struct {
	IP        string        `json:"ip"`
	Port      string        `json:"port"`
	Timestamp time.Time     `json:"timestamp"`
	RTT       time.Duration `json:"rtt"`
	Success   bool          `json:"success"`
}
