# SuperBizAgent Skills 功能技术方案

## 1. 概述

Skills 是一种可扩展的能力模块，允许 AI 助手执行特定领域的复杂任务。与 Tools 不同，Skills 是更高层次的抽象，可以组合多个 Tools 和自定义逻辑来完成复杂工作流。

## 2. Skills vs Tools 对比

| 特性 | Tools | Skills |
|------|-------|--------|
| 粒度 | 单一操作 | 复杂工作流 |
| 触发方式 | LLM 自动选择 | 用户显式调用或 LLM 推荐 |
| 配置 | 代码硬编码 | 支持动态配置 |
| 示例 | query_events | 安全事件分析报告生成 |

## 3. 核心 Skills 设计

### 3.1 预置 Skills

1. **event_analysis** - 安全事件深度分析
2. **threat_hunting** - 威胁狩猎
3. **log_diagnosis** - 日志诊断
4. **report_generation** - 报告生成
5. **alert_triage** - 告警分类处理

### 3.2 Skill 数据结构

```go
type Skill struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Category    string            `json:"category"`
    Enabled     bool              `json:"enabled"`
    Config      map[string]any    `json:"config"`
    Prompt      string            `json:"prompt"`      // Skill 专用提示词
    Tools       []string          `json:"tools"`       // 依赖的 Tools
}
```

## 4. 架构设计

```
┌─────────────────────────────────────────────────────┐
│                    Frontend                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │
│  │ Skill 列表  │  │ Skill 配置  │  │ Skill 触发  │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  │
└─────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────┐
│                   API Layer                          │
│  POST /api/skill/v1/list                            │
│  POST /api/skill/v1/execute                         │
│  POST /api/skill/v1/config                          │
└─────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────┐
│                 Skill Engine                         │
│  ┌─────────────────────────────────────────────┐   │
│  │            SkillRegistry                     │   │
│  │  - RegisterSkill()                          │   │
│  │  - GetSkill()                               │   │
│  │  - ListSkills()                             │   │
│  └─────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────┐   │
│  │            SkillExecutor                     │   │
│  │  - Execute(skillID, params)                 │   │
│  │  - BuildPrompt()                            │   │
│  │  - SelectTools()                            │   │
│  └─────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────┐
│              Existing Components                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐         │
│  │  Tools   │  │   LLM    │  │ Retriever│         │
│  └──────────┘  └──────────┘  └──────────┘         │
└─────────────────────────────────────────────────────┘
```

## 5. 目录结构

```
internal/
├── ai/
│   ├── skills/                    # 新增
│   │   ├── registry.go           # Skill 注册中心
│   │   ├── executor.go           # Skill 执行器
│   │   ├── types.go              # 类型定义
│   │   └── builtin/              # 内置 Skills
│   │       ├── event_analysis.go
│   │       ├── threat_hunting.go
│   │       ├── log_diagnosis.go
│   │       └── report_gen.go
│   └── ...
├── controller/
│   └── skill/                     # 新增
│       └── skill_v1_*.go
└── ...

api/
└── skill/
    └── v1/
        └── skill.go               # API 定义
```

## 6. 核心接口设计

### 6.1 API 接口

```
GET  /api/skill/v1/list           # 获取 Skill 列表
POST /api/skill/v1/execute        # 执行 Skill (SSE 流式)
PUT  /api/skill/v1/config/:id     # 更新 Skill 配置
```

### 6.2 执行请求

```json
{
  "skill_id": "event_analysis",
  "params": {
    "event_id": 123
  },
  "stream": true
}
```

### 6.3 执行响应 (SSE)

```
data: {"type": "step", "content": "正在获取事件详情..."}
data: {"type": "tool_call", "tool": "query_events", "args": {...}}
data: {"type": "step", "content": "分析威胁指标..."}
data: {"type": "result", "content": "分析完成..."}
data: [DONE]
```

## 7. 实现步骤

### Phase 1: 基础框架
1. 定义 Skill 类型和接口
2. 实现 SkillRegistry
3. 实现 SkillExecutor
4. 添加 API 路由

### Phase 2: 内置 Skills
1. 实现 event_analysis skill
2. 实现 log_diagnosis skill
3. 实现 report_generation skill

### Phase 3: 前端集成
1. Skill 列表展示
2. Skill 执行界面
3. 执行结果流式展示

## 8. 示例: event_analysis Skill

```go
var EventAnalysisSkill = &Skill{
    ID:          "event_analysis",
    Name:        "安全事件分析",
    Description: "对安全事件进行深度分析，包括威胁评估、影响范围、处置建议",
    Category:    "security",
    Enabled:     true,
    Tools:       []string{"query_events", "query_log", "query_internal_docs"},
    Prompt: `你是安全分析专家。请对事件进行深度分析：
1. 事件概述
2. 威胁评估
3. 影响范围
4. 处置建议

事件ID: {event_id}`,
}
```

## 9. 与现有系统集成

Skills 复用现有组件：
- **LLM**: 使用 chat_pipeline 中的 model
- **Tools**: 复用 internal/ai/tools 中的工具
- **Retriever**: 复用 RAG 检索能力
- **SSE**: 复用现有 SSE 服务

## 10. 前端交互设计

在 Chat 页面添加 Skills 入口：
- 输入框上方显示可用 Skills 快捷按钮
- 点击后弹出参数配置面板
- 执行结果以流式消息展示在对话中
