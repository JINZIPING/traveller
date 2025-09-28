package model

import "time"

// ICMPProbeTaskDTO 探测任务
type ICMPProbeTaskDTO struct {
	IP        string    `json:"ip"`         // 探测目标的IP地址
	Count     int       `json:"count"`      // 探测的次数
	Port      int       `json:"port"`       // 探测目标的端口（ICMP可不设置）
	Threshold int       `json:"threshold"`  // 丢包率阈值
	Timeout   int       `json:"timeout"`    // 探测超时时间（秒）
	CreatedAt time.Time `json:"created_at"` // 任务创建时间
	UpdatedAt time.Time `json:"updated_at"` // 任务更新时间
}

// ICMPProbeResultDTO 探测结果
type ICMPProbeResultDTO struct {
	IP         string        `json:"ip"`          // 探测目标的IP地址
	PacketLoss float64       `json:"packet_loss"` // 丢包率
	MinRTT     time.Duration `json:"min_rtt"`     // 最小往返时间
	MaxRTT     time.Duration `json:"max_rtt"`     // 最大往返时间
	AvgRTT     time.Duration `json:"avg_rtt"`     // 平均往返时间
	Threshold  int           `json:"threshold"`   // 丢包率阈值（来自任务）
	Success    bool          `json:"success"`     // 探测是否成功
	Timestamp  time.Time     `json:"timestamp"`   // 探测完成的时间戳
	TaskTime   time.Time     `json:"task_time"`   // 下发任务的时间 (T0)
}
