package event

import (
	"SuperBizAgent/internal/logic/event"
)

// Controller 事件控制器
type Controller struct {
	service *event.EventService
}

// New 创建事件控制器
func New() *Controller {
	return &Controller{
		service: event.NewEventService(),
	}
}
