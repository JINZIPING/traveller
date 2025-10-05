package mark

import (
	"context"
	"sync"
	"time"
)

type store struct {
	mu sync.RWMutex
	m  map[string]time.Time
}

var Default = &store{m: make(map[string]time.Time)}

type ctxKey string

const keyTaskID ctxKey = "mark.task_id"

// Start 记录任务开始时间
func Start(ctx context.Context, taskID string, t ...time.Time) context.Context {
	start := time.Now().UTC()
	if len(t) > 0 {
		start = t[0].UTC()
	}
	Default.mu.Lock()
	Default.m[taskID] = start
	Default.mu.Unlock()

	return context.WithValue(ctx, keyTaskID, taskID)
}

// Finish 结束并计算耗时
func Finish(taskID string) (time.Duration, bool) {
	now := time.Now().UTC()
	Default.mu.Lock()
	start, ok := Default.m[taskID]
	if ok {
		delete(Default.m, taskID)
	}
	Default.mu.Unlock()
	if !ok {
		return 0, false
	}
	return now.Sub(start), true
}
