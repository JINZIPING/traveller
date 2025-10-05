package model

import (
	"time"
)

// TCPProbeTaskDTO TCP探测任务
type TCPProbeTaskDTO struct {
	TaskID    string    `json:"task_id"`
	IP        string    `json:"ip"`
	Port      string    `json:"port"`
	Timeout   int       `json:"timeout"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TCPProbeResultDTO TCP探测结果
type TCPProbeResultDTO struct {
	TaskID  string        `json:"task_id"`
	IP      string        `json:"ip"`
	Port    string        `json:"port"`
	Success bool          `json:"success"`
	RTT     time.Duration `json:"rtt"`
}
