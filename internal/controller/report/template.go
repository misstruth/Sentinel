package report

import (
	"context"

	v1 "SuperBizAgent/api/report/v1"
	"SuperBizAgent/internal/logic/report"
	"SuperBizAgent/internal/model"
)

// TemplateController 模板控制器
type TemplateController struct {
	service *report.TemplateService
}

// NewTemplateController 创建模板控制器
func NewTemplateController() *TemplateController {
	return &TemplateController{
		service: report.NewTemplateService(),
	}
}

// TemplateCreate 创建模板
func (c *TemplateController) TemplateCreate(ctx context.Context, req *v1.TemplateCreateReq) (*v1.TemplateCreateRes, error) {
	tpl := &model.ReportTemplate{
		Name:        req.Name,
		Description: req.Description,
		Type:        model.ReportType(req.Type),
		Content:     req.Content,
		IsDefault:   req.IsDefault,
	}

	if err := c.service.Create(ctx, tpl); err != nil {
		return nil, err
	}

	return &v1.TemplateCreateRes{ID: tpl.ID}, nil
}

// TemplateList 获取模板列表
func (c *TemplateController) TemplateList(ctx context.Context, req *v1.TemplateListReq) (*v1.TemplateListRes, error) {
	templates, err := c.service.List(ctx, req.Type)
	if err != nil {
		return nil, err
	}

	list := make([]v1.TemplateItem, len(templates))
	for i, t := range templates {
		list[i] = v1.TemplateItem{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Type:        string(t.Type),
			IsDefault:   t.IsDefault,
		}
	}

	return &v1.TemplateListRes{List: list}, nil
}

// TemplateGet 获取模板详情
func (c *TemplateController) TemplateGet(ctx context.Context, req *v1.TemplateGetReq) (*v1.TemplateGetRes, error) {
	t, err := c.service.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &v1.TemplateGetRes{
		TemplateItem: &v1.TemplateItem{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Type:        string(t.Type),
			IsDefault:   t.IsDefault,
		},
		Content: t.Content,
	}, nil
}

// TemplateUpdate 更新模板
func (c *TemplateController) TemplateUpdate(ctx context.Context, req *v1.TemplateUpdateReq) (*v1.TemplateUpdateRes, error) {
	t, err := c.service.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		t.Name = req.Name
	}
	if req.Description != "" {
		t.Description = req.Description
	}
	if req.Content != "" {
		t.Content = req.Content
	}
	t.IsDefault = req.IsDefault
	return &v1.TemplateUpdateRes{}, c.service.Update(ctx, t)
}

// TemplateDelete 删除模板
func (c *TemplateController) TemplateDelete(ctx context.Context, req *v1.TemplateDeleteReq) (*v1.TemplateDeleteRes, error) {
	return &v1.TemplateDeleteRes{}, c.service.Delete(ctx, req.ID)
}
