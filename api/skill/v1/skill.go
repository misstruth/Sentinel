package v1

import "github.com/gogf/gf/v2/frame/g"

type SkillListReq struct {
	g.Meta `path:"/skill/v1/list" method:"get"`
}

type SkillListRes struct {
	g.Meta `mime:"application/json"`
	Skills []SkillInfo `json:"skills"`
}

type SkillInfo struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Category    string       `json:"category"`
	Params      []ParamInfo  `json:"params"`
}

type ParamInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}
