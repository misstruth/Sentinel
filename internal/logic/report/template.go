package report

import (
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"
	"context"
	"fmt"
)

// TemplateService 模板服务
type TemplateService struct{}

// NewTemplateService 创建模板服务
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// Create 创建模板
func (s *TemplateService) Create(ctx context.Context, tpl *model.ReportTemplate) error {
	db := database.GetDB()
	return db.Create(tpl).Error
}

// Get 获取模板
func (s *TemplateService) Get(ctx context.Context, id uint) (*model.ReportTemplate, error) {
	db := database.GetDB()
	var tpl model.ReportTemplate
	if err := db.First(&tpl, id).Error; err != nil {
		return nil, fmt.Errorf("模板不存在: %w", err)
	}
	return &tpl, nil
}

// Update 更新模板
func (s *TemplateService) Update(ctx context.Context, tpl *model.ReportTemplate) error {
	db := database.GetDB()
	return db.Save(tpl).Error
}

// Delete 删除模板
func (s *TemplateService) Delete(ctx context.Context, id uint) error {
	db := database.GetDB()
	return db.Delete(&model.ReportTemplate{}, id).Error
}

// List 获取模板列表
func (s *TemplateService) List(ctx context.Context, reportType string) ([]model.ReportTemplate, error) {
	db := database.GetDB()
	var templates []model.ReportTemplate

	query := db.Model(&model.ReportTemplate{})
	if reportType != "" {
		query = query.Where("type = ?", reportType)
	}

	if err := query.Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}

// GetDefault 获取默认模板
func (s *TemplateService) GetDefault(ctx context.Context, reportType model.ReportType) (*model.ReportTemplate, error) {
	db := database.GetDB()
	var tpl model.ReportTemplate
	err := db.Where("type = ? AND is_default = ?", reportType, true).First(&tpl).Error
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}
