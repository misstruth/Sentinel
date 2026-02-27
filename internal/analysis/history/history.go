package history

import (
	"time"
)

// Record 分析历史记录
type Record struct {
	ID         uint
	EventID    uint
	TemplateID uint
	Input      string
	Output     string
	CreatedAt  time.Time
}

// Store 历史记录存储
type Store struct {
	records []*Record
	nextID  uint
}

// NewStore 创建存储
func NewStore() *Store {
	return &Store{
		records: make([]*Record, 0),
		nextID:  1,
	}
}

// Add 添加记录
func (s *Store) Add(r *Record) uint {
	r.ID = s.nextID
	r.CreatedAt = time.Now()
	s.records = append(s.records, r)
	s.nextID++
	return r.ID
}

// Get 获取记录
func (s *Store) Get(id uint) *Record {
	for _, r := range s.records {
		if r.ID == id {
			return r
		}
	}
	return nil
}

// List 列出所有记录
func (s *Store) List() []*Record {
	return s.records
}
