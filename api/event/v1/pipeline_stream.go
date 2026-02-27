package v1

import "github.com/gogf/gf/v2/frame/g"

// PipelineStreamReq SSE流式处理请求
type PipelineStreamReq struct {
	g.Meta   `path:"/event/pipeline/stream" method:"post" tags:"Event" summary:"多Agent流式处理"`
	Mode     string `json:"mode" v:"in:today,latest,specific" dc:"分析模式: today=今日事件, latest=最近10条, specific=指定事件"`
	EventIDs []int  `json:"event_ids" dc:"指定事件ID列表(mode=specific时必填)"`
}
