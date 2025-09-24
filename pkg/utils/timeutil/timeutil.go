package timeutil

import "time"

// NowUTC8 返回当前时间（UTC+8 / Asia/Shanghai）
func NowUTC8() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc)
}
