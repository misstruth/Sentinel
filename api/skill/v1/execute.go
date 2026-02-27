package v1

import "github.com/gogf/gf/v2/frame/g"

type SkillExecuteReq struct {
	g.Meta  `path:"/skill/v1/execute" method:"post"`
	SkillID string         `json:"skill_id" v:"required"`
	Params  map[string]any `json:"params"`
}

type SkillExecuteRes struct {
	g.Meta `mime:"text/event-stream"`
}
