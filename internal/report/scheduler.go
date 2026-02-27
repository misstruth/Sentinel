package report

import (
	"context"
	"sync"
	"time"
)

// Schedule 报告调度配置
type Schedule struct {
	ID       uint
	Name     string
	Cron     string
	Type     string
	Enabled  bool
}

// Scheduler 报告调度器
type Scheduler struct {
	schedules map[uint]*Schedule
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		schedules: make(map[uint]*Schedule),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Add 添加调度
func (s *Scheduler) Add(sch *Schedule) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.schedules[sch.ID] = sch
}

// Remove 移除调度
func (s *Scheduler) Remove(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.schedules, id)
}

// Start 启动调度器
func (s *Scheduler) Start() {
	go s.run()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cancel()
}

// run 运行调度循环
func (s *Scheduler) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkAndRun()
		}
	}
}

// checkAndRun 检查并执行调度
func (s *Scheduler) checkAndRun() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sch := range s.schedules {
		if sch.Enabled {
			go s.generate(sch)
		}
	}
}

// generate 生成报告
func (s *Scheduler) generate(sch *Schedule) {
	// 根据类型生成报告
	_ = sch
}
