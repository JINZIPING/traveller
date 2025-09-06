package model

import (
	"time"
)

// ICMPProbeTask ICMP探测结果
type ICMPProbeTask struct {
	IP        string    `json:"ip"`         // 探测目标的IP地址
	Count     int       `json:"count"`      // 探测的次数
	Port      int       `json:"port"`       // 探测目标的端口（ICMP可不设置）
	Threshold int       `json:"threshold"`  // 丢包率阈值
	Timeout   int       `json:"timeout"`    // 探测超时时间（秒）
	CreatedAt time.Time `json:"created_at"` // 任务创建时间
	UpdatedAt time.Time `json:"updated_at"` // 任务更新时间
}

// TCPProbeTask TCP探测结果
type TCPProbeTask struct {
	IP        string    `json:"ip"`
	Port      string    `json:"port"`
	Timeout   int       `json:"timeout"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
