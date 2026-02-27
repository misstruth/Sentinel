package service

import (
	"context"
	"fmt"
	"strings"

	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// AnalyzeEvent AI分析事件并生成处置建议
func AnalyzeEvent(ctx context.Context, eventID uint) error {
	db := database.GetDB()
	var event model.SecurityEvent
	if err := db.First(&event, eventID).Error; err != nil {
		return err
	}

	// 构建分析提示词
	prompt := fmt.Sprintf(`作为安全分析专家，请分析以下安全事件并给出评估：

事件标题: %s
事件描述: %s
来源: %s
CVE编号: %s

请按以下格式输出：
风险评分: (0-100的数字，100为最高风险)
严重级别: (critical/high/medium/low/info)
处置建议: (具体的处置步骤，每条建议一行)`,
		event.Title, event.Description, event.SourceURL, event.CVEID)

	// 调用LLM
	content, err := callAnalysisLLM(ctx, prompt)
	if err != nil {
		return err
	}

	// 解析结果
	riskScore, severity, recommendation := parseAnalysisResult(content)

	// 更新事件
	updates := map[string]interface{}{
		"risk_score":     riskScore,
		"recommendation": recommendation,
	}
	if severity != "" {
		updates["severity"] = severity
	}

	return db.Model(&event).Updates(updates).Error
}

func callAnalysisLLM(ctx context.Context, prompt string) (string, error) {
	gctx := gctx.New()
	apiKey, _ := g.Cfg().Get(gctx, "ds_quick_chat_model.api_key")
	baseURL, _ := g.Cfg().Get(gctx, "ds_quick_chat_model.base_url")
	modelName, _ := g.Cfg().Get(gctx, "ds_quick_chat_model.model")

	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  apiKey.String(),
		BaseURL: baseURL.String(),
		Model:   modelName.String(),
	})
	if err != nil {
		return "", err
	}

	resp, err := chatModel.Generate(ctx, []*schema.Message{schema.UserMessage(prompt)})
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

func parseAnalysisResult(content string) (int, string, string) {
	riskScore := 50
	severity := ""
	var recommendations []string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "风险评分:") {
			fmt.Sscanf(strings.TrimPrefix(line, "风险评分:"), "%d", &riskScore)
		} else if strings.HasPrefix(line, "严重级别:") {
			s := strings.TrimSpace(strings.TrimPrefix(line, "严重级别:"))
			if s == "critical" || s == "high" || s == "medium" || s == "low" || s == "info" {
				severity = s
			}
		} else if strings.HasPrefix(line, "处置建议:") {
			recommendations = append(recommendations, strings.TrimPrefix(line, "处置建议:"))
		} else if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "•") {
			recommendations = append(recommendations, strings.TrimLeft(line, "-• "))
		}
	}

	return riskScore, severity, strings.Join(recommendations, "\n")
}
