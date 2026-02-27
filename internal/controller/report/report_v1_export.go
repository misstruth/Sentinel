package report

import (
	"context"
	"fmt"

	v1 "SuperBizAgent/api/report/v1"
	reportLogic "SuperBizAgent/internal/logic/report"
)

// Export 导出报告
func (c *Controller) Export(ctx context.Context, req *v1.ExportReq) (*v1.ExportRes, error) {
	// 获取报告
	rpt, err := c.service.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 导出
	exportSvc := reportLogic.NewExportService()
	format := reportLogic.ExportFormat(req.Format)
	data, err := exportSvc.Export(rpt, format)
	if err != nil {
		return nil, err
	}

	// 生成文件名
	ext := req.Format
	if ext == "markdown" {
		ext = "md"
	}
	filename := fmt.Sprintf("report_%d.%s", req.ID, ext)

	return &v1.ExportRes{
		Content:  string(data),
		Filename: filename,
	}, nil
}
