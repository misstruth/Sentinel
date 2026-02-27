# SuperBizAgent 智能安全运维平台

## 项目概述

SuperBizAgent 是一个基于多 Agent 架构的智能安全运维平台，集成了安全事件管理、威胁情报分析、AI 运维等功能。

### 技术栈

| 层级 | 技术 |
|------|------|
| 后端框架 | GoFrame (gf/v2) |
| 前端框架 | React + TypeScript + Vite |
| AI 框架 | Eino (字节跳动) |
| 数据库 | MySQL + GORM |
| 向量数据库 | Milvus |
| LLM | DeepSeek V3 (火山引擎) |
| 嵌入模型 | Doubao Embedding (火山引擎) |

### 核心功能

1. **安全事件管理** - 实时监控、追踪和处理安全事件
2. **威胁情报分析** - 收集和分析威胁情报
3. **AI 对话助手** - 智能问答和安全分析
4. **多 Agent 协作** - Supervisor 调度多个专业 Agent
5. **Skills 技能系统** - 可扩展的专业技能模块
6. **报告生成** - 自动生成安全分析报告

---

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Milvus 2.0+ (可选，用于 RAG)

### 配置文件

编辑 `manifest/config/config.yaml`:

```yaml
# 数据库配置
database:
  default:
    host: "127.0.0.1"
    port: "3306"
    user: "root"
    pass: "your_password"
    name: "superbiz"

# LLM 配置 (火山引擎)
doubao_model:
  api_key: "your_api_key"
  model: "ep-xxx"  # DeepSeek 端点

# 嵌入模型配置
doubao_embedding_model:
  api_key: "your_api_key"
  model: "ep-xxx"  # Embedding 端点

# Milvus 配置
milvus:
  address: "127.0.0.1:19530"
  collection: "knowledge_base"
```

### 启动服务

```bash
# 后端
GOSUMDB=sum.golang.org go run main.go

# 前端
cd web && npm install && npm run dev
```

访问 http://localhost:3000

---

## 项目结构

```
SuperBizAgent/
├── api/                    # API 定义
│   ├── chat/v1/           # 聊天接口
│   ├── skill/v1/          # Skills 接口
│   ├── event/v1/          # 事件接口
│   └── report/v1/         # 报告接口
├── internal/
│   ├── ai/                # AI 核心
│   │   ├── agent/         # Agent 实现
│   │   ├── skills/        # Skills 系统
│   │   └── tools/         # 工具集
│   ├── controller/        # 控制器
│   ├── database/          # 数据库
│   └── model/             # 数据模型
├── web/                   # 前端项目
├── manifest/config/       # 配置文件
└── main.go               # 入口
```

---

## 多 Agent 架构

### 架构图

```
用户输入
    │
    ▼
┌─────────────────┐
│   Supervisor    │  ← 智能路由决策
└─────────────────┘
    │
    ├──→ ChatAgent     (对话问答)
    ├──→ EventAgent    (事件查询)
    ├──→ ReportAgent   (报告生成)
    ├──→ RiskAgent     (风险评估)
    └──→ PlanAgent     (复杂任务规划)
    │
    ▼
┌─────────────────┐
│  SharedMemory   │  ← 共享上下文
└─────────────────┘
```

### Agent 类型

| Agent | 功能 | 触发场景 |
|-------|------|----------|
| ChatAgent | 通用对话 | 默认/问答 |
| EventAgent | 事件查询 | "查询事件"、"最近告警" |
| ReportAgent | 报告生成 | "生成报告"、"安全总结" |
| RiskAgent | 风险评估 | "评估风险"、"威胁分析" |
| PlanAgent | 复杂规划 | 多步骤任务 |

### 核心特性

1. **智能路由** - Supervisor 自动分析用户意图，选择最合适的 Agent
2. **共享记忆** - 所有 Agent 共享上下文，保持对话连贯性
3. **Agent 互调用** - Agent 可作为 Tool 被其他 Agent 调用
4. **流式输出** - 实时返回处理进度和结果

---

## Skills 系统

Skills 是可扩展的专业技能模块，比 Tools 更高层次的抽象。

### 内置 Skills

| Skill | 功能 | 参数 |
|-------|------|------|
| event_analysis | 安全事件深度分析 | event_id |
| log_diagnosis | 日志诊断 | keyword |
| threat_hunting | 威胁狩猎 | target |

### 自定义 Skill

```go
// internal/ai/skills/builtin/my_skill.go
package builtin

import "SuperBizAgent/internal/ai/skills"

var MySkill = &skills.Skill{
    ID:          "my_skill",
    Name:        "我的技能",
    Description: "技能描述",
    Category:    "custom",
    Enabled:     true,
    Tools:       []string{"query_events"},
    Params: []skills.SkillParam{
        {Name: "param1", Type: "string", Required: true},
    },
    Prompt: `执行任务: {param1}`,
}

func init() { skills.Register(MySkill) }
```

---

## API 接口

### 聊天接口

```
POST /api/chat/v1/stream      # 普通对话 (SSE)
POST /api/chat/v1/supervisor  # 多Agent对话 (SSE)
POST /api/chat/v1/ai-ops      # AI运维分析
```

### Skills 接口

```
GET  /api/skill/v1/list       # 获取技能列表
POST /api/skill/v1/execute    # 执行技能 (SSE)
```

### 事件接口

```
GET  /api/event/v1/list       # 事件列表
POST /api/event/v1/update     # 更新状态
POST /api/event/pipeline/stream  # 事件分析 (SSE)
```

### 报告接口

```
GET  /api/report/list         # 报告列表
POST /api/report/generate     # 生成报告
GET  /api/report/:id          # 获取报告
```

---

## Tools 工具集

内置工具供 Agent 调用：

| Tool | 功能 |
|------|------|
| query_events | 查询安全事件 |
| query_subscriptions | 查询订阅 |
| query_reports | 查询报告 |
| query_internal_docs | 查询内部文档 |
| get_current_time | 获取当前时间 |
| mysql_crud | 数据库操作 |

---

## 前端页面

| 页面 | 路径 | 功能 |
|------|------|------|
| 仪表板 | /dashboard | 安全概览 |
| 事件列表 | /events | 事件管理 |
| 事件分析 | /event-analysis | AI分析 |
| AI助手 | /chat | 对话/Skills |
| 报告 | /reports | 报告管理 |
| 订阅 | /subscriptions | 订阅管理 |

---

## 开发指南

### 添加新 Agent

```go
// internal/ai/agent/supervisor/agent_xxx.go
package supervisor

type XxxAgent struct{}

func (a *XxxAgent) Name() AgentType { return "xxx" }

func (a *XxxAgent) Execute(ctx context.Context, task *Task, cb StreamCallback) (*Result, error) {
    cb("xxx", "[处理中...]\n")
    // 实现逻辑
    return &Result{Content: "结果"}, nil
}

func init() { RegisterAgent(&XxxAgent{}) }
```

### 添加新 Tool

```go
// internal/ai/tools/my_tool.go
type MyTool struct{}

func (t *MyTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
    return &schema.ToolInfo{
        Name: "my_tool",
        Desc: "工具描述",
    }, nil
}

func (t *MyTool) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
    return "结果", nil
}
```

---

## 许可证

MIT License
