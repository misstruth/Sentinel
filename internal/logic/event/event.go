package event

import (
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"
	"context"
)

// EventService 事件服务
type EventService struct{}

// NewEventService 创建事件服务
func NewEventService() *EventService {
	return &EventService{}
}

// Get 获取事件详情
func (s *EventService) Get(ctx context.Context, id uint) (*model.SecurityEvent, error) {
	db := database.GetDB()
	var event model.SecurityEvent
	err := db.First(&event, id).Error
	return &event, err
}

// List 获取事件列表
func (s *EventService) List(ctx context.Context, page, pageSize int, severity, status, keyword string) ([]model.SecurityEvent, int64, error) {
	db := database.GetDB()
	var events []model.SecurityEvent
	var total int64

	query := db.Model(&model.SecurityEvent{})

	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("event_time DESC").Offset(offset).Limit(pageSize).Find(&events).Error
	return events, total, err
}

// DeleteAll 删除所有事件(硬删除)
func (s *EventService) DeleteAll(ctx context.Context) error {
	db := database.GetDB()
	return db.Unscoped().Where("1=1").Delete(&model.SecurityEvent{}).Error
}
