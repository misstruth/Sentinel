package skill

import (
	v1 "SuperBizAgent/api/skill/v1"
	"SuperBizAgent/internal/ai/skills"
	"context"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) List(ctx context.Context, req *v1.SkillListReq) (*v1.SkillListRes, error) {
	list := skills.List()
	res := &v1.SkillListRes{Skills: make([]v1.SkillInfo, 0, len(list))}
	for _, s := range list {
		info := v1.SkillInfo{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			Category:    s.Category,
			Params:      make([]v1.ParamInfo, 0, len(s.Params)),
		}
		for _, p := range s.Params {
			info.Params = append(info.Params, v1.ParamInfo{
				Name:        p.Name,
				Type:        p.Type,
				Description: p.Description,
				Required:    p.Required,
			})
		}
		res.Skills = append(res.Skills, info)
	}
	return res, nil
}
