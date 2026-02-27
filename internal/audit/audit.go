package audit

import (
	"time"
)

// Action 操作类型
type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionLogin  Action = "login"
)

// Log 操作日志
type Log struct {
	ID        uint
	UserID    uint
	Action    Action
	Resource  string
	Detail    string
	IP        string
	CreatedAt time.Time
}

// Logger 日志记录器
type Logger struct {
	logs   []*Log
	nextID uint
}

// NewLogger 创建记录器
func NewLogger() *Logger {
	return &Logger{
		logs:   make([]*Log, 0),
		nextID: 1,
	}
}

// Record 记录日志
func (l *Logger) Record(log *Log) {
	log.ID = l.nextID
	log.CreatedAt = time.Now()
	l.logs = append(l.logs, log)
	l.nextID++
}

// List 列出日志
func (l *Logger) List() []*Log {
	return l.logs
}
