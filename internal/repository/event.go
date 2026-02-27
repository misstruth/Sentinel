package repository

import (
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"
	"time"

	"gorm.io/gorm"
)

// EventRepository 安全事件仓储
type EventRepository struct {
	db *gorm.DB
}

// NewEventRepository 创建事件仓储
func NewEventRepository() *EventRepository {
	return &EventRepository{db: database.GetDB()}
}

// Create 创建事件
func (r *EventRepository) Create(event *model.SecurityEvent) error {
	return r.db.Create(event).Error
}

// GetByID 根据ID获取事件
func (r *EventRepository) GetByID(id uint) (*model.SecurityEvent, error) {
	var event model.SecurityEvent
	err := r.db.First(&event, id).Error
	return &event, err
}

// Update 更新事件
func (r *EventRepository) Update(event *model.SecurityEvent) error {
	return r.db.Save(event).Error
}

// Delete 删除事件
func (r *EventRepository) Delete(id uint) error {
	return r.db.Delete(&model.SecurityEvent{}, id).Error
}

// ExistsByHash 检查事件是否存在(去重)
func (r *EventRepository) ExistsByHash(hash string) bool {
	var count int64
	r.db.Model(&model.SecurityEvent{}).Where("unique_hash = ?", hash).Count(&count)
	return count > 0
}

// ListBySubscription 按订阅源查询事件
func (r *EventRepository) ListBySubscription(subID uint, page, size int) ([]model.SecurityEvent, int64, error) {
	var events []model.SecurityEvent
	var total int64

	query := r.db.Model(&model.SecurityEvent{}).Where("subscription_id = ?", subID)
	query.Count(&total)

	offset := (page - 1) * size
	err := query.Order("event_time DESC").Offset(offset).Limit(size).Find(&events).Error
	return events, total, err
}

// CountByStatus 按状态统计事件数
func (r *EventRepository) CountByStatus(status model.EventStatus) int64 {
	var count int64
	r.db.Model(&model.SecurityEvent{}).Where("status = ?", status).Count(&count)
	return count
}

// ListBySeverity 按严重程度查询事件
func (r *EventRepository) ListBySeverity(severity model.SeverityLevel, page, size int) ([]model.SecurityEvent, int64, error) {
	var events []model.SecurityEvent
	var total int64

	query := r.db.Model(&model.SecurityEvent{}).Where("severity = ?", severity)
	query.Count(&total)

	offset := (page - 1) * size
	err := query.Order("event_time DESC").Offset(offset).Limit(size).Find(&events).Error
	return events, total, err
}

// UpdateStatus 更新事件状态
func (r *EventRepository) UpdateStatus(id uint, status model.EventStatus) error {
	return r.db.Model(&model.SecurityEvent{}).Where("id = ?", id).Update("status", status).Error
}

// BatchUpdateStatus 批量更新事件状态
func (r *EventRepository) BatchUpdateStatus(ids []uint, status model.EventStatus) error {
	return r.db.Model(&model.SecurityEvent{}).Where("id IN ?", ids).Update("status", status).Error
}

// GetStats 获取事件统计
func (r *EventRepository) GetStats() (total int64, bySeverity map[string]int64, byStatus map[string]int64, err error) {
	// 总数
	r.db.Model(&model.SecurityEvent{}).Count(&total)

	// 按严重程度统计
	bySeverity = make(map[string]int64)
	var severityResults []struct {
		Severity string
		Count    int64
	}
	r.db.Model(&model.SecurityEvent{}).
		Select("severity, count(*) as count").
		Group("severity").
		Scan(&severityResults)
	for _, r := range severityResults {
		bySeverity[r.Severity] = r.Count
	}

	// 按状态统计
	byStatus = make(map[string]int64)
	var statusResults []struct {
		Status string
		Count  int64
	}
	r.db.Model(&model.SecurityEvent{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusResults)
	for _, r := range statusResults {
		byStatus[r.Status] = r.Count
	}

	return
}

// GetTrend 获取事件趋势
func (r *EventRepository) GetTrend(days int) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, days)

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.AddDate(0, 0, 1)

		var total int64
		r.db.Model(&model.SecurityEvent{}).
			Where("event_time >= ? AND event_time < ?", startOfDay, endOfDay).
			Count(&total)

		// 按严重程度统计
		var severityCounts []struct {
			Severity string
			Count    int64
		}
		r.db.Model(&model.SecurityEvent{}).
			Select("severity, count(*) as count").
			Where("event_time >= ? AND event_time < ?", startOfDay, endOfDay).
			Group("severity").
			Scan(&severityCounts)

		severityMap := map[string]int64{
			"critical": 0,
			"high":     0,
			"medium":   0,
			"low":      0,
			"info":     0,
		}
		for _, sc := range severityCounts {
			severityMap[sc.Severity] = sc.Count
		}

		results[days-1-i] = map[string]interface{}{
			"date":     dateStr,
			"total":    total,
			"critical": severityMap["critical"],
			"high":     severityMap["high"],
			"medium":   severityMap["medium"],
			"low":      severityMap["low"],
			"info":     severityMap["info"],
		}
	}

	return results, nil
}

// GetTodayStats 获取今日统计
func (r *EventRepository) GetTodayStats() (total, critical, high int64) {
	today := time.Now().Format("2006-01-02")
	r.db.Model(&model.SecurityEvent{}).Where("DATE(event_time) = ?", today).Count(&total)
	r.db.Model(&model.SecurityEvent{}).Where("DATE(event_time) = ? AND severity = ?", today, "critical").Count(&critical)
	r.db.Model(&model.SecurityEvent{}).Where("DATE(event_time) = ? AND severity = ?", today, "high").Count(&high)
	return
}
