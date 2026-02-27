package supervisor

import (
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino/schema"
)

// MessageMeta 消息元数据
type MessageMeta struct {
	Message   *schema.Message
	Timestamp time.Time
	Tokens    int      // 估算的 token 数
	Important bool     // 是否重要（包含工具调用结果）
	Tags      []string // 标签：tool_call, tool_result, query, answer
}

// SharedMemory 共享记忆 - 参考 qmd 的分层上下文设计
type SharedMemory struct {
	mu            sync.RWMutex
	messages      []*MessageMeta
	context       map[string]any
	summary       string // 历史对话摘要
	maxTokens     int    // 最大 token 数
	summaryTokens int    // 摘要占用的 token 数
}

var globalMemory = &SharedMemory{
	messages:  make([]*MessageMeta, 0),
	context:   make(map[string]any),
	maxTokens: 4000, // 预留 4k tokens 给历史
}

func GetSharedMemory() *SharedMemory { return globalMemory }

// estimateTokens 估算消息的 token 数（简单按字符数/3估算）
func estimateTokens(content string) int {
	return len(content) / 3
}

// AddMessage 添加消息，自动标记重要性
func (m *SharedMemory) AddMessage(msg *schema.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()

	meta := &MessageMeta{
		Message:   msg,
		Timestamp: time.Now(),
		Tokens:    estimateTokens(msg.Content),
		Tags:      make([]string, 0),
	}

	// 标记消息类型
	switch msg.Role {
	case schema.User:
		meta.Tags = append(meta.Tags, "query")
	case schema.Assistant:
		meta.Tags = append(meta.Tags, "answer")
	case schema.Tool:
		meta.Tags = append(meta.Tags, "tool_result")
		meta.Important = true // 工具结果通常重要
	}

	// 检查是否包含工具调用
	if len(msg.ToolCalls) > 0 {
		meta.Tags = append(meta.Tags, "tool_call")
		meta.Important = true
	}

	m.messages = append(m.messages, meta)
	m.compressIfNeeded()
}

func (m *SharedMemory) GetMessages() []*schema.Message {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*schema.Message, 0, len(m.messages))
	for _, meta := range m.messages {
		result = append(result, meta.Message)
	}
	return result
}

func (m *SharedMemory) SetContext(key string, val any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.context[key] = val
}

func (m *SharedMemory) GetContext(key string) any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.context[key]
}

func (m *SharedMemory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = make([]*MessageMeta, 0)
	m.context = make(map[string]any)
	m.summary = ""
	m.summaryTokens = 0
}

// compressIfNeeded 当 token 超限时压缩历史
func (m *SharedMemory) compressIfNeeded() {
	totalTokens := m.summaryTokens
	for _, meta := range m.messages {
		totalTokens += meta.Tokens
	}

	if totalTokens <= m.maxTokens {
		return
	}

	// 策略：保留最近的重要消息，压缩旧消息为摘要
	var toSummarize []*MessageMeta
	var toKeep []*MessageMeta

	// 从后往前，保留最近 5 轮对话
	keepCount := 0
	for i := len(m.messages) - 1; i >= 0; i-- {
		if keepCount < 10 || m.messages[i].Important {
			toKeep = append([]*MessageMeta{m.messages[i]}, toKeep...)
			keepCount++
		} else {
			toSummarize = append([]*MessageMeta{m.messages[i]}, toSummarize...)
		}
	}

	// 生成摘要
	if len(toSummarize) > 0 {
		newSummary := m.generateSummary(toSummarize)
		if m.summary != "" {
			m.summary = m.summary + "\n" + newSummary
		} else {
			m.summary = newSummary
		}
		m.summaryTokens = estimateTokens(m.summary)
	}

	m.messages = toKeep
}

// generateSummary 生成对话摘要
func (m *SharedMemory) generateSummary(messages []*MessageMeta) string {
	var sb strings.Builder
	sb.WriteString("[历史摘要] ")

	for _, meta := range messages {
		switch meta.Message.Role {
		case schema.User:
			content := meta.Message.Content
			if len(content) > 50 {
				content = content[:50] + "..."
			}
			sb.WriteString("用户问: " + content + "; ")
		case schema.Assistant:
			content := meta.Message.Content
			if len(content) > 80 {
				content = content[:80] + "..."
			}
			sb.WriteString("回答: " + content + "; ")
		}
	}
	return sb.String()
}

// GetSummary 获取历史摘要
func (m *SharedMemory) GetSummary() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.summary
}

// GetContextString 获取格式化的上下文字符串（用于注入 prompt）
func (m *SharedMemory) GetContextString() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.summary == "" {
		return ""
	}
	return "对话历史摘要: " + m.summary
}
