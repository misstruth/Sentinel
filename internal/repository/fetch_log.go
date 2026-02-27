package repository

import (
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"

	"gorm.io/gorm"
)

// FetchLogRepository 抓取日志仓储
type FetchLogRepository struct {
	db *gorm.DB
}

// NewFetchLogRepository 创建抓取日志仓储
func NewFetchLogRepository() *FetchLogRepository {
	return &FetchLogRepository{db: database.GetDB()}
}

// Create 创建日志
func (r *FetchLogRepository) Create(log *model.FetchLog) error {
	return r.db.Create(log).Error
}

// ListBySubscription 按订阅源查询日志
func (r *FetchLogRepository) ListBySubscription(subID uint, limit int) ([]model.FetchLog, error) {
	var logs []model.FetchLog
	err := r.db.Where("subscription_id = ?", subID).
		Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

// ListBySubscriptionPaged 按订阅源分页查询日志
func (r *FetchLogRepository) ListBySubscriptionPaged(subID uint, page, pageSize int) ([]model.FetchLog, int64, error) {
	var logs []model.FetchLog
	var total int64

	query := r.db.Model(&model.FetchLog{}).Where("subscription_id = ?", subID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

// GetStats 获取订阅统计
func (r *FetchLogRepository) GetStats(subID uint) (total, success, failed int64, avgDuration int, err error) {
	r.db.Model(&model.FetchLog{}).Where("subscription_id = ?", subID).Count(&total)
	r.db.Model(&model.FetchLog{}).Where("subscription_id = ? AND status = ?", subID, "success").Count(&success)
	r.db.Model(&model.FetchLog{}).Where("subscription_id = ? AND status = ?", subID, "failed").Count(&failed)

	var avg float64
	r.db.Model(&model.FetchLog{}).Where("subscription_id = ?", subID).Select("COALESCE(AVG(duration), 0)").Scan(&avg)
	avgDuration = int(avg)
	return
}
