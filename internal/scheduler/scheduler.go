package scheduler

import (
	"context"
	"sync"
	"time"

	"SuperBizAgent/internal/model"
)

// Job 调度任务
type Job struct {
	ID           uint
	Subscription *model.Subscription
	NextRun      time.Time
}

// Handler 任务处理函数
type Handler func(ctx context.Context, sub *model.Subscription) error

// Scheduler 调度器
type Scheduler struct {
	jobs    map[uint]*Job
	handler Handler
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// New 创建调度器
func New(handler Handler) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		jobs:    make(map[uint]*Job),
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Add 添加任务
func (s *Scheduler) Add(sub *model.Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[sub.ID] = &Job{
		ID:           sub.ID,
		Subscription: sub,
		NextRun:      time.Now(),
	}
}

// Remove 移除任务
func (s *Scheduler) Remove(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, id)
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
			s.tick()
		}
	}
}

// tick 执行一次调度检查
func (s *Scheduler) tick() {
	s.mu.RLock()
	jobs := make([]*Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		jobs = append(jobs, j)
	}
	s.mu.RUnlock()

	now := time.Now()
	for _, j := range jobs {
		if j.NextRun.Before(now) {
			go s.execute(j)
		}
	}
}

// execute 执行任务
func (s *Scheduler) execute(j *Job) {
	sub := j.Subscription
	if sub.Status != model.StatusActive {
		return
	}

	_ = s.handler(s.ctx, sub)

	// 更新下次执行时间
	s.mu.Lock()
	if job, ok := s.jobs[j.ID]; ok {
		job.NextRun = time.Now().Add(time.Hour)
	}
	s.mu.Unlock()
}
