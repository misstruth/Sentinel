package report

import (
	"SuperBizAgent/internal/logic/report"
)

// Controller 报告控制器
type Controller struct {
	service *report.ReportService
}

// New 创建报告控制器
func New() *Controller {
	return &Controller{
		service: report.NewReportService(),
	}
}
