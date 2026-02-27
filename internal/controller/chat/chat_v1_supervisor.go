package chat

import (
	v1 "SuperBizAgent/api/chat/v1"
	"SuperBizAgent/internal/ai/agent/supervisor"
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Supervisor(ctx context.Context, req *v1.SupervisorChatReq) (*v1.SupervisorChatRes, error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")

	sup := supervisor.NewSupervisor(ctx)
	sup.Execute(req.Query, func(agent supervisor.AgentType, chunk string) {
		data, _ := json.Marshal(map[string]string{"agent": string(agent), "content": chunk})
		r.Response.Writef("data: %s\n\n", data)
		r.Response.Flush()
	})

	r.Response.Writef("data: [DONE]\n\n")
	r.Response.Flush()
	return nil, nil
}
