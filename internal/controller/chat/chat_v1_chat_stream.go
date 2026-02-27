package chat

import (
	"SuperBizAgent/api/chat/v1"
	"SuperBizAgent/internal/ai/agent/chat_pipeline"
	"SuperBizAgent/utility/log_call_back"
	"SuperBizAgent/utility/mem"
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) ChatStream(ctx context.Context, req *v1.ChatStreamReq) (res *v1.ChatStreamRes, err error) {
	id := req.Id
	msg := req.Question

	ctx = context.WithValue(ctx, "client_id", req.Id)
	client, err := c.service.Create(ctx, g.RequestFromCtx(ctx))
	if err != nil {
		return nil, err
	}

	userMessage := &chat_pipeline.UserMessage{
		ID:      id,
		Query:   msg,
		History: mem.GetSimpleMemory(id).GetMessages(),
	}

	runner, err := chat_pipeline.BuildChatAgent(ctx)
	if err != nil {
		client.SendToClient("error", err.Error())
		return nil, err
	}

	// 使用 Invoke 替代 Stream，确保工具调用后能获取完整回复
	out, err := runner.Invoke(ctx, userMessage, compose.WithCallbacks(log_call_back.LogCallback(nil)))
	if err != nil {
		client.SendToClient("error", err.Error())
		return nil, err
	}

	// 调试日志
	g.Log().Infof(ctx, "AI Response Content: %s", out.Content)

	content := out.Content
	// 分段发送内容，每次发送一小段
	runes := []rune(content)
	chunkSize := 20 // 每次发送20个字符
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunk := string(runes[i:end])
		if chunk != "" {
			client.SendToClient("message", chunk)
		}
	}

	client.SendToClient("done", "[DONE]")

	if content != "" {
		mem.GetSimpleMemory(id).SetMessages(schema.UserMessage(msg))
		mem.GetSimpleMemory(id).SetMessages(schema.SystemMessage(content))
	}

	return &v1.ChatStreamRes{}, nil
}
