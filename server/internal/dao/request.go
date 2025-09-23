package dao

type TCPProbeTaskReq struct {
	IP          string `json:"ip" binding:"required"`
	Port        string `json:"port" binding:"required"`
	Timeout     int    `json:"timeout" binding:"required,min=1"`
	IntervalSec int    `json:"interval_sec" binding:"omitempty,min=1"`
}

type ICMPProbeTaskReq struct {
	IP          string `json:"ip" binding:"required"`
	Count       int    `json:"count" binding:"required,min=1"`
	Threshold   int    `json:"threshold" binding:"required,min=0,max=100"`
	Timeout     int    `json:"timeout" binding:"required,min=1"`
	IntervalSec int    `json:"interval_sec" binding:"omitempty,min=1"`
}
