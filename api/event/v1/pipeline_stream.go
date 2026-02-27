package v1

import "github.com/gogf/gf/v2/frame/g"

// PipelineStreamReq SSE流式处理请求
type PipelineStreamReq struct {
	g.Meta `path:"/event/pipeline/stream" method:"post" tags:"Event" summary:"多Agent流式处理"`
}
