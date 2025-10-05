package service

import (
	"sync"
	"time"
)

type tickJob struct {
	stop chan struct{}
}

type Scheduler struct {
	mu   sync.Mutex // 因为可能有多个goroutine同时增删任务
	jobs map[string]*tickJob
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: make(map[string]*tickJob),
	}
}

func (s *Scheduler) AddTCPJob(id string, every time.Duration, fn func())  { s.addJob(id, every, fn) }
func (s *Scheduler) AddICMPJob(id string, every time.Duration, fn func()) { s.addJob(id, every, fn) }

func (s *Scheduler) addJob(id string, every time.Duration, fn func()) {

	/*
		1. 停掉旧的同名任务
		2. 新建任务
		3. 启动一个goroutine定时执行
	*/

	s.mu.Lock()
	defer s.mu.Unlock()

	if old, ok := s.jobs[id]; ok {
		close(old.stop)
	}
	tj := &tickJob{
		stop: make(chan struct{}),
	}
	s.jobs[id] = tj

	go func() {
		ticker := time.NewTicker(every) // 间隔时间
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C: // 定时传入任务函数
				fn()
			case <-tj.stop:
				return
			}
		}
	}()
}

func (s *Scheduler) Remove(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if j, ok := s.jobs[id]; ok {
		close(j.stop)
		delete(s.jobs, id)
	}
}
