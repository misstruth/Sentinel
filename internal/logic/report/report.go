package report

import (
	"SuperBizAgent/internal/ai/agent/report_generator"
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"
	"context"
	"fmt"
)

// ReportService 报告服务
type ReportService struct{}

// NewReportService 创建报告服务
func NewReportService() *ReportService {
	return &ReportService{}
}

// Generate 生成报告
func (s *ReportService) Generate(ctx context.Context, req *report_generator.ReportRequest) (*report_generator.ReportResponse, error) {
	generator := report_generator.NewReportGenerator(ctx)
	return generator.Generate(req)
}

// Get 获取报告
func (s *ReportService) Get(ctx context.Context, id uint) (*model.Report, error) {
	db := database.GetDB()
	var report model.Report
	if err := db.First(&report, id).Error; err != nil {
		return nil, fmt.Errorf("报告不存在: %w", err)
	}
	return &report, nil
}

// List 获取报告列表
func (s *ReportService) List(ctx context.Context, page, pageSize int, reportType string) ([]model.Report, int64, error) {
	db := database.GetDB()
	var reports []model.Report
	var total int64

	query := db.Model(&model.Report{})
	if reportType != "" {
		query = query.Where("type = ?", reportType)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// Delete 删除报告
func (s *ReportService) Delete(ctx context.Context, id uint) error {
	db := database.GetDB()
	return db.Delete(&model.Report{}, id).Error
}
