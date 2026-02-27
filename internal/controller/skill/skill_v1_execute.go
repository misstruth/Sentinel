package skill

import (
	v1 "SuperBizAgent/api/skill/v1"
	"SuperBizAgent/internal/ai/skills"
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) Execute(ctx context.Context, req *v1.SkillExecuteReq) (*v1.SkillExecuteRes, error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")

	executor, err := skills.NewExecutor(req.SkillID, req.Params)
	if err != nil {
		sendSSE(r, "error", err.Error())
		return nil, nil
	}

	err = executor.Execute(ctx, func(result skills.ExecuteResult) {
		sendSSE(r, result.Type, result.Content)
	})
	if err != nil {
		sendSSE(r, "error", err.Error())
	}

	r.Response.Writef("data: [DONE]\n\n")
	r.Response.Flush()
	return nil, nil
}

func sendSSE(r *ghttp.Request, typ, content string) {
	data, _ := json.Marshal(map[string]string{"type": typ, "content": content})
	r.Response.Writef("data: %s\n\n", data)
	r.Response.Flush()
}
