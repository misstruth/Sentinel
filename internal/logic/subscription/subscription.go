package subscription

import (
	"context"
	"errors"

	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"

	"gorm.io/gorm"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidSubscription  = errors.New("invalid subscription data")
)

// Service 订阅服务
type Service struct {
	db *gorm.DB
}

// NewService 创建订阅服务实例
func NewService() *Service {
	return &Service{
		db: database.GetDB(),
	}
}

// Create 创建订阅
func (s *Service) Create(ctx context.Context, sub *model.Subscription) error {
	if sub.Name == "" || sub.SourceType == "" {
		return ErrInvalidSubscription
	}
	return s.db.WithContext(ctx).Create(sub).Error
}

// GetByID 根据ID获取订阅
func (s *Service) GetByID(ctx context.Context, id uint) (*model.Subscription, error) {
	var sub model.Subscription
	err := s.db.WithContext(ctx).First(&sub, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSubscriptionNotFound
	}
	return &sub, err
}

// Update 更新订阅
func (s *Service) Update(ctx context.Context, sub *model.Subscription) error {
	if sub.ID == 0 {
		return ErrInvalidSubscription
	}
	return s.db.WithContext(ctx).Save(sub).Error
}

// Delete 删除订阅(软删除)
func (s *Service) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&model.Subscription{}, id)
	if result.RowsAffected == 0 {
		return ErrSubscriptionNotFound
	}
	return result.Error
}

// ListQuery 列表查询参数
type ListQuery struct {
	Page       int
	PageSize   int
	SourceType model.SourceType
	Status     model.SubscriptionStatus
	Keyword    string
}

// ListResult 列表查询结果
type ListResult struct {
	Items []*model.Subscription
	Total int64
	Page  int
	Size  int
}

// List 分页查询订阅列表
func (s *Service) List(ctx context.Context, query *ListQuery) (*ListResult, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	db := s.db.WithContext(ctx).Model(&model.Subscription{})

	// 条件过滤
	if query.SourceType != "" {
		db = db.Where("source_type = ?", query.SourceType)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []*model.Subscription
	offset := (query.Page - 1) * query.PageSize
	err := db.Order("created_at DESC").
		Offset(offset).Limit(query.PageSize).
		Find(&items).Error

	return &ListResult{
		Items: items,
		Total: total,
		Page:  query.Page,
		Size:  query.PageSize,
	}, err
}

// GetActiveSubscriptions 获取所有活跃的订阅
func (s *Service) GetActiveSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	var subs []*model.Subscription
	err := s.db.WithContext(ctx).
		Where("status = ?", model.StatusActive).
		Find(&subs).Error
	return subs, err
}

// UpdateStatus 更新订阅状态
func (s *Service) UpdateStatus(ctx context.Context, id uint, status model.SubscriptionStatus) error {
	result := s.db.WithContext(ctx).
		Model(&model.Subscription{}).
		Where("id = ?", id).
		Update("status", status)
	if result.RowsAffected == 0 {
		return ErrSubscriptionNotFound
	}
	return result.Error
}

// Pause 暂停订阅
func (s *Service) Pause(ctx context.Context, id uint) error {
	return s.UpdateStatus(ctx, id, model.StatusPaused)
}

// Resume 恢复订阅
func (s *Service) Resume(ctx context.Context, id uint) error {
	return s.UpdateStatus(ctx, id, model.StatusActive)
}

// Disable 禁用订阅
func (s *Service) Disable(ctx context.Context, id uint) error {
	return s.UpdateStatus(ctx, id, model.StatusDisabled)
}
