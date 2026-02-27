# Sentinel — 多Agent安全事件智能研判平台

基于多Agent协作架构的企业安全事件分析平台，结合大语言模型（DeepSeek V3）、向量数据库（Milvus）和 RAG 检索增强生成技术，实现安全事件的自动采集、智能提取、语义去重、AI风险评估和解决方案生成。

## 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                    React Frontend                        │
│         Vite + TypeScript + Tailwind CSS                │
│    Canvas动画拓扑图 / SSE实时流 / Zustand状态管理         │
└──────────────────────┬──────────────────────────────────┘
                       │ SSE / REST API
┌──────────────────────▼──────────────────────────────────┐
│                  Go Backend (GoFrame)                     │
│                    Port: 6872                            │
├─────────────────────────────────────────────────────────┤
│  多Agent SSE Pipeline (/api/event/pipeline/stream)       │
│                                                          │
│  ┌──────────┐  ┌──────────┐  ┌────────┐  ┌──────────┐  │
│  │数据采集   │→│智能提取   │→│去重过滤 │→│风险评估   │  │
│  │Agent     │  │Agent     │  │Agent   │  │Agent     │  │
│  │(DB查询)  │  │(DeepSeek)│  │(MD5)   │  │(DeepSeek)│  │
│  └──────────┘  └──────────┘  └────────┘  └────┬─────┘  │
│                                                │         │
│                                    ┌───────────▼───────┐ │
│                                    │ 解决方案Agent      │ │
│                                    │ (ReAct + RAG)     │ │
│                                    │ ┌───────────────┐ │ │
│                                    │ │ 6个AI工具:     │ │ │
│                                    │ │ • query_events │ │ │
│                                    │ │ • event_detail │ │ │
│                                    │ │ • similar_events│ │ │
│                                    │ │ • internal_docs│ │ │
│                                    │ │ • prom_alerts  │ │ │
│                                    │ │ • current_time │ │ │
│                                    │ └───────────────┘ │ │
│                                    └───────────────────┘ │
├─────────────────────────────────────────────────────────┤
│  Milvus向量库          MySQL关系库         Prometheus    │
│  (语义检索/事件索引)    (事件/订阅/报告)    (监控告警)    │
└─────────────────────────────────────────────────────────┘
```

## 多Agent协作流程

5个Agent通过SSE流式管道串联，前端实时展示每个Agent的执行状态和思考链路：

| 阶段 | Agent | 技术实现 | 说明 |
|------|-------|---------|------|
| 1 | 数据采集Agent | MySQL查询 | 扫描 `status=new` 的安全事件，最多10条 |
| 2 | 智能提取Agent | DeepSeek V3 | 对原始事件做AI结构化提取（标题/CVE/严重程度/影响产品） |
| 3 | 去重过滤Agent | MD5哈希 | 基于标题+CVE的语义去重 |
| 4 | 风险评估Agent | DeepSeek V3 | AI评分（0-100），输出severity/recommendation/factors |
| 5 | 解决方案Agent | ReAct + Milvus RAG | 对高危事件调用6个工具，搜索历史相似事件和知识库，生成处置方案 |

解决方案Agent是核心，基于 Cloudwego Eino 框架的 ReAct 推理循环，最多25步工具调用，自主决定使用哪些工具来分析事件。

## 向量数据库 (Milvus RAG)

系统使用 Milvus 向量数据库实现语义检索和事件记忆：

```
事件处理完成 → Doubao Embedding → 1024维向量 → Milvus集合(biz)
                                                    ↓
用户查询/新事件 → Doubao Embedding → 向量相似搜索 → TopK=3 相似结果
```

**核心能力：**

- **事件索引**：处理后的安全事件自动写入 Milvus，内容包含标题+描述+CVE+厂商+AI分析结论（上限7000字符）
- **相似搜索**：解决方案Agent通过 `search_similar_events` 工具检索历史相似事件，辅助生成处置方案
- **知识库检索**：`internal_docs` 工具检索预导入的安全知识文档，提供修复参考
- **内存管理**：集合容量上限10000条，通过 `EnsureCapacity` 监控并 Flush 持久化数据

**Embedding模型：** 豆包（Doubao）文本向量模型，1024维，通过字节跳动 Ark API 调用。

## 技术栈

### 后端

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.24 |
| Web框架 | GoFrame (gf/v2) | v2 |
| ORM | GORM | v1.25 |
| Agent框架 | Cloudwego Eino | v0.7.13 |
| LLM | DeepSeek V3 (via Ark API) | — |
| Embedding | 豆包 Doubao (via Ark API) | 1024维 |
| 向量数据库 | Milvus | v2.5.10 |
| 关系数据库 | MySQL | 8.0 |
| 实时通信 | SSE (Server-Sent Events) | — |

### 前端

| 组件 | 技术 | 版本 |
|------|------|------|
| 框架 | React + TypeScript | 18.3 / 5.6 |
| 构建工具 | Vite | 6.0 |
| 样式 | Tailwind CSS | 3.4 |
| 状态管理 | Zustand | 4.5 |
| 图表 | ECharts | 5.4 |
| 动画 | Framer Motion | 11.0 |
| Markdown渲染 | react-markdown | 9.0 |
| 路由 | React Router | 6.22 |
| 请求 | Axios + React Query | — |

## Docker 容器化部署

项目依赖以下容器化服务，分为两组 Docker Compose 配置：

### 1. Milvus 向量数据库集群

路径：`manifest/docker/docker-compose.yml`

| 服务 | 镜像 | 端口 | 说明 |
|------|------|------|------|
| etcd | quay.io/coreos/etcd:v3.5.18 | — | Milvus 元数据存储 |
| minio | minio/minio:RELEASE.2023-03-20 | 9000/9001 | Milvus 对象存储后端 |
| milvus-standalone | milvusdb/milvus:v2.5.10 | 19530/9091 | 向量数据库主服务 |
| attu | zilliz/attu:v2.6 | 8000→3000 | Milvus 可视化管理界面 |

### 2. MySQL + 应用服务

路径：`docker-compose.yml`（项目根目录）

| 服务 | 镜像 | 端口 | 说明 |
|------|------|------|------|
| mysql | mysql:8.0 | 3307→3306 | 关系数据库（事件/订阅/报告） |
| app | 本地构建 (Dockerfile) | 8000 | Go 后端服务 |

## 快速启动

### 前置条件

- Docker & Docker Compose
- Go 1.24+
- Node.js 18+
- 字节跳动 Ark API Key（用于 DeepSeek V3 和 Doubao Embedding）

### Step 1: 启动 Milvus 向量数据库

```bash
cd manifest/docker
docker-compose up -d
```

等待 Milvus 健康检查通过（约90秒），可通过 http://localhost:8000 访问 Attu 管理界面。

### Step 2: 启动 MySQL

```bash
# 项目根目录
docker-compose up -d mysql
```

MySQL 映射到本地 3307 端口，数据库名 `fo_sentinel`。

### Step 3: 配置后端

```bash
# 复制并编辑配置文件
cp manifest/config/config.yaml manifest/config/config.local.yaml
```

修改 `config.local.yaml` 中的关键配置：
- `ds_think_chat_model.api_key` — Ark API Key
- `doubao_embedding_model.api_key` — 同上
- `database.host/port/pass` — MySQL 连接信息
- `file_dir` — 知识库文档目录的绝对路径

### Step 4: 启动后端

```bash
go run main.go
```

后端服务启动在 `localhost:6872`，SSE 管道端点：`/api/event/pipeline/stream`。

### Step 5: 启动前端

```bash
cd web
npm install
npm run dev
```

前端开发服务器默认在 `localhost:5173`。

## 项目结构

```
SuperBizAgent/
├── main.go                      # 应用入口
├── go.mod / go.sum              # Go 依赖管理
├── Dockerfile                   # 多阶段构建镜像
├── docker-compose.yml           # MySQL + App 编排
├── api/                         # API 请求/响应模型定义
│   ├── chat/v1/                 #   聊天接口模型
│   ├── event/v1/                #   事件接口模型
│   ├── report/v1/               #   报告接口模型
│   ├── skill/v1/                #   技能接口模型
│   └── subscription/v1/         #   订阅接口模型
├── internal/                    # 核心业务逻辑
│   ├── ai/                      #   AI Agent 管道
│   │   ├── agent/               #     多Agent实现
│   │   ├── models/              #     LLM模型集成
│   │   └── skills/              #     AI工具定义
│   ├── controller/              #   HTTP 控制器
│   ├── service/                 #   业务逻辑层
│   ├── repository/              #   数据访问层
│   ├── database/                #   数据库初始化与迁移
│   ├── model/                   #   GORM 数据模型
│   ├── config/                  #   配置管理
│   ├── middleware/              #   HTTP 中间件
│   ├── auth/                    #   认证与 RBAC
│   ├── cache/                   #   缓存层
│   ├── notify/                  #   通知系统
│   └── health/                  #   健康检查
├── web/                         # React 前端
│   └── src/
│       ├── pages/               #     页面组件
│       ├── components/          #     公共组件
│       ├── services/            #     API 客户端
│       ├── stores/              #     Zustand 状态管理
│       ├── types/               #     TypeScript 类型
│       └── utils/               #     工具函数
├── cmd/                         # CLI 工具
├── test/                        # 测试文件
├── manifest/                    # 配置与 Docker 文件
│   ├── config/config.yaml       #   后端配置模板
│   └── docker/docker-compose.yml#   Milvus 集群编排
└── docs/                        # 文档
```

## API 接口

### 事件管理 `/api/event`

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/event` | 事件列表（分页、按严重程度/状态/关键词过滤） |
| GET | `/event/{id}` | 事件详情 |
| PUT | `/event/{id}/status` | 更新单个事件状态 |
| POST | `/event/batch/status` | 批量更新事件状态 |
| GET | `/event/stats` | 事件统计（总数/今日/严重数/按等级分布） |
| GET | `/event/trend` | 事件趋势数据（默认7天） |
| DELETE | `/event/all` | 清空所有事件 |
| POST | `/event/{id}/analyze` | AI 分析单个事件 |
| POST | `/event/pipeline/process` | 多Agent管道处理 |
| POST | `/event/pipeline/stream` | SSE 流式管道（核心接口） |

### 订阅管理 `/api/subscriptions`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/subscriptions` | 创建订阅 |
| GET | `/subscriptions` | 订阅列表（支持过滤） |
| GET | `/subscriptions/{id}` | 订阅详情 |
| PUT | `/subscriptions/{id}` | 更新订阅 |
| DELETE | `/subscriptions/{id}` | 删除订阅 |
| POST | `/subscriptions/{id}/pause` | 暂停订阅 |
| POST | `/subscriptions/{id}/resume` | 恢复订阅 |
| POST | `/subscriptions/{id}/disable` | 禁用订阅 |
| GET | `/subscriptions/{id}/logs` | 获取抓取日志 |
| GET | `/subscriptions/{id}/stats` | 订阅统计 |
| POST | `/subscriptions/{id}/fetch` | 手动触发抓取 |

### 聊天与AI `/api/chat`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/chat` | 发送聊天消息，返回完整响应 |
| POST | `/chat_stream` | SSE 流式聊天（实时逐字输出） |
| POST | `/chat/v1/supervisor` | Supervisor 多Agent路由聊天（SSE流式） |
| POST | `/upload` | 上传文件到知识库（PDF/TXT/MD/CSV/DOC） |
| POST | `/ai_ops` | AI 运维分析（生成告警分析报告） |

### 报告管理 `/api/report`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/report/generate` | 生成报告（daily/weekly/monthly） |
| GET | `/report` | 报告列表（分页） |
| GET | `/report/:id` | 报告详情 |
| DELETE | `/report/:id` | 删除报告 |
| GET | `/report/:id/export` | 导出报告（markdown/html/json） |
| POST | `/report/template` | 创建报告模板 |
| GET | `/report/template` | 模板列表 |
| GET | `/report/template/{id}` | 模板详情 |
| PUT | `/report/template/{id}` | 更新模板 |
| DELETE | `/report/template/{id}` | 删除模板 |

### 技能执行 `/api/skill/v1`

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/skill/v1/list` | 获取可用技能列表及参数定义 |
| POST | `/skill/v1/execute` | 执行技能（SSE流式返回结果） |

### 抓取日志 `/api/fetchlog`

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/fetchlog/list` | 抓取日志列表（按订阅ID分页查询） |

## Supervisor 多Agent路由架构

除了事件处理的5个Agent管道外，系统还实现了基于 Supervisor 模式的智能路由架构，用于处理用户的自由对话请求：

```
用户查询 → Supervisor
              │
              ├─ Router（DeepSeek V3 意图分类）
              │
              ├──→ Chat Agent      （通用对话/安全咨询/知识问答）
              ├──→ Event Agent     （事件查询/事件分析/告警关联）
              ├──→ Report Agent    （报告生成/报告分析）
              ├──→ Risk Agent      （风险评估/威胁分析/CVE评估）
              └──→ Plan Agent      （复杂多步任务，Plan-Execute-Replan）
```

| Agent | 职责 | 推理模式 | 最大步数 |
|-------|------|---------|---------|
| Chat Agent | 通用对话、安全咨询、知识问答、日志查询 | ReAct 循环 | 25 |
| Event Agent | 安全事件查询、事件分析、告警关联 | 事件分析管道 | — |
| Report Agent | 报告生成与分析 | 直接执行 | — |
| Risk Agent | 风险评估、威胁分析、漏洞评分 | ReAct 循环 | — |
| Plan Agent | 复杂多步任务规划与执行 | Plan-Execute-Replan | 20 |

Router 使用 DeepSeek V3 Quick 模型对用户意图进行分类，返回 JSON 格式的路由决策，将请求分发到对应的专业 Agent。

## AI 工具集

Agent 在 ReAct 推理循环中可自主调用以下工具：

| 工具 | 说明 | 参数 |
|------|------|------|
| `query_events` | 查询安全事件数据库 | severity, status, limit(1-5) |
| `get_event_detail` | 获取事件详细信息 | event_id |
| `search_similar_events` | Milvus 向量相似搜索历史事件 | 事件特征 |
| `query_internal_docs` | 检索知识库文档（RAG） | 搜索关键词 |
| `query_prometheus_alerts` | 查询 Prometheus 活跃告警 | — |
| `get_current_time` | 获取当前系统时间 | — |
| `query_subscriptions` | 查询订阅和通知配置 | — |
| `query_reports` | 查询已生成的报告 | — |
| `mysql_crud` | 数据库 CRUD 操作 | 表名、操作类型、条件 |
| `query_log` | 通过 MCP 协议查询系统日志 | 查询条件 |

## 预定义技能

系统内置3个可通过 `/skill/v1/execute` 调用的预定义技能，每个技能绑定特定工具组合：

| 技能ID | 名称 | 分类 | 绑定工具 | 参数 |
|--------|------|------|---------|------|
| `event_analysis` | 安全事件分析 | security | query_events, query_internal_docs | event_id (number) |
| `log_diagnosis` | 日志诊断 | ops | query_internal_docs, get_current_time | keyword (string) |
| `threat_hunting` | 威胁狩猎 | security | query_events, query_subscriptions | target (string) |

## 认证与权限

### JWT 认证

- 算法：HMAC-SHA256
- Token 有效期：24小时
- 请求头：`Authorization: Bearer <token>`
- Claims：UserID、Username、过期时间

### RBAC 角色权限

| 角色 | Read | Write | Delete |
|------|------|-------|--------|
| Admin | ✓ | ✓ | ✓ |
| User | ✓ | ✓ | ✗ |
| Viewer | ✓ | ✗ | ✗ |

## 中间件

| 中间件 | 说明 |
|--------|------|
| CORS | 跨域请求支持，基于 GoFrame CORSDefault |
| Auth | Bearer Token 验证，集成 JWT 解析，401 拒绝无效请求 |
| RateLimit | 令牌桶限流，100次/分钟/IP，超限返回 429 |
| Logger | 记录 HTTP 方法、路径、状态码、耗时 |

## 通知系统

支持 Webhook 通知渠道，当安全事件触发告警时自动推送：

- 协议：HTTP POST，JSON 格式
- 超时：10秒
- 认证：可选 `X-Webhook-Secret` 头
- Payload 结构：`{event, data, time}`

## 数据模型

### SecurityEvent（安全事件）

| 字段 | 类型 | 说明 |
|------|------|------|
| Title | string | 事件标题 |
| Description | text | 事件描述 |
| Severity | enum | critical / high / medium / low / info |
| Status | enum | new / processing / resolved / ignored |
| SourceURL | string | 来源链接 |
| CVEID | string | CVE 编号 |
| CVSSScore | float | CVSS 评分 |
| AffectedVendor | string | 受影响厂商 |
| AffectedProduct | string | 受影响产品 |
| RiskScore | int | AI 风险评分（0-100） |
| Recommendation | text | AI 处置建议 |
| UniqueHash | string | 去重哈希 |
| Tags | json | 标签 |
| Starred | bool | 是否标星 |

### Subscription（数据订阅）

| 字段 | 类型 | 说明 |
|------|------|------|
| Name | string | 订阅名称 |
| SourceType | enum | vulnerability / threat_intel / vendor_advisory / attack_activity / github_repo / rss / webhook / nvd / cve |
| SourceURL | string | 数据源地址 |
| Status | enum | active / paused / disabled |
| CronExpr | string | 定时抓取 Cron 表达式 |
| FetchTimeout | int | 抓取超时（秒） |
| AuthType | enum | none / api_key / oauth / basic |
| Keywords | json | 过滤关键词 |
| MinSeverity | string | 最低严重程度过滤 |
| TotalEvents | int | 累计抓取事件数 |
| FailedFetches | int | 失败次数 |

### Report（安全报告）

| 字段 | 类型 | 说明 |
|------|------|------|
| Title | string | 报告标题 |
| Type | enum | daily / weekly / monthly / custom / vuln_alert / threat_brief |
| Status | enum | pending / generating / completed / failed |
| Summary | text | 报告摘要 |
| Content | longtext | 报告正文 |
| TemplateID | uint | 关联模板ID |
| EventIDs | json | 关联事件ID列表 |
| EventCount | int | 事件总数 |
| CriticalCount | int | 严重事件数 |
| HighCount | int | 高危事件数 |
| GenerateMethod | enum | manual / scheduled / api |

## 前端页面

React 前端包含以下页面，通过 React Router v6 管理路由：

| 路由 | 页面 | 说明 |
|------|------|------|
| `/dashboard` | Dashboard | 安全态势总览：统计卡片、事件趋势图、严重程度分布饼图、安全漏斗 |
| `/events` | Events | 事件列表：搜索过滤、批量操作、状态管理、事件详情弹窗 |
| `/events/analysis` | EventAnalysis | 多Agent协作分析：5步管道可视化、SSE实时流、思考链展示、风险评分仪表盘 |
| `/chat` | Chat | AI助手：通用对话/AI运维/技能执行/多Agent四种模式 |
| `/subscriptions` | Subscriptions | 订阅管理：数据源配置、状态控制、手动抓取、统计信息 |
| `/reports` | Reports | 报告管理：生成/查看/导出（markdown/html/json）、模板管理 |
| `/logs` | Logs | 抓取日志：按订阅查看抓取历史和状态 |
| `/settings` | Settings | 系统设置 |

### 前端状态管理 (Zustand)

| Store | 说明 | 持久化 |
|-------|------|--------|
| `appStore` | 侧边栏折叠状态、主题（dark/light） | localStorage |
| `eventStore` | Agent日志列表、管道处理状态 | 否 |
| `contextStore` | 当前页面上下文、选中事件ID/标题 | 否 |

### SSE 流式通信

前端通过 Fetch API + ReadableStream 实现 SSE 流式接收，支持以下4个流式端点：

| 端点 | 用途 |
|------|------|
| `/api/event/pipeline/stream` | 多Agent事件分析管道 |
| `/api/chat_stream` | 标准聊天流式响应 |
| `/api/chat/v1/supervisor` | Supervisor 多Agent路由聊天 |
| `/api/skill/v1/execute` | 技能执行流式返回 |

协议格式：
```
data: {"agent":"AgentName","content":"message","type":"event_type"}
data: [DONE]
```

## 配置说明

配置文件路径：`manifest/config/config.yaml`，建议复制为 `config.local.yaml` 后修改。

### 服务配置

```yaml
server:
  address: ":8000"           # 监听端口
  openapiPath: "/api.json"   # OpenAPI 文档端点
  swaggerPath: "/swagger"    # Swagger UI 端点
```

### AI 模型配置

```yaml
# DeepSeek V3 Think 模型（用于 Plan Agent 复杂推理）
ds_think_chat_model:
  api_key: "your-ark-api-key"
  base_url: "https://ark.cn-beijing.volces.com/api/v3"
  model: "deepseek-v3-1-terminus"

# DeepSeek V3 Quick 模型（用于 Router/Chat Agent 快速推理）
ds_quick_chat_model:
  api_key: "your-ark-api-key"
  base_url: "https://ark.cn-beijing.volces.com/api/v3"
  model: "deepseek-v3-1-terminus"

# 豆包 Embedding 模型（用于 Milvus 向量化）
doubao_embedding_model:
  api_key: "your-ark-api-key"
  model: "ep-20260222183953-gv8q5"
```

### 数据库配置

```yaml
database:
  host: "127.0.0.1"
  port: "3307"
  user: "root"
  pass: "sentinel123"
  name: "fo_sentinel"
```

### 知识库目录

```yaml
file_dir: "/your/path/to/knowledge_docs"
```

该目录存放 RAG 检索用的安全知识文档，上传的文件会自动向量化写入 Milvus。

## Legacy 前端

`SuperBizAgentFrontend/` 目录包含一个轻量级的原生 JavaScript 前端，适合快速体验核心聊天功能：

- 纯 Vanilla JS（ES6+），无需构建工具
- 聊天界面 + SSE 流式响应
- 文件拖拽上传到知识库（支持 PDF/TXT/MD/CSV/DOC/DOCX，上限 50MB）

启动方式：

```bash
cd SuperBizAgentFrontend
python3 -m http.server 8080
```

## CLI 命令

项目提供 `Fo-Sentinel` CLI 工具（`cmd/cli/`），支持以下命令：

| 命令 | 说明 |
|------|------|
| `help` | 显示使用帮助 |
| `version` | 显示版本号（v1.0.0） |
| `batch` | 批量操作框架 |

参数：
- `-c` 指定配置文件路径
- `-v` 开启详细日志

## 缓存系统

内置内存缓存（`internal/cache/`），支持 TTL 自动过期：

- 线程安全（sync.RWMutex）
- 可配置 TTL 过期时间
- 后台 goroutine 自动清理过期条目
- 用于减少数据库查询和 API 调用频率

## 健康检查

可插拔的健康检查框架（`internal/health/`），支持注册自定义检查函数：

- 状态值：`healthy` / `unhealthy`
- 可扩展：注册数据库、Milvus、外部服务等检查项
- 统一执行所有检查并返回结果汇总

## Docker 构建

Dockerfile 采用多阶段构建，优化镜像体积：

```
阶段1: golang:1.21-alpine（编译）
  → go mod download
  → CGO_ENABLED=0 go build → fo-sentinel 二进制

阶段2: alpine:3.19（运行）
  → 仅包含二进制 + manifest 配置
  → 时区: Asia/Shanghai
  → 暴露端口: 8000
```

构建命令：

```bash
docker build -t superbizagent .
docker run -p 8000:8000 superbizagent
```
<img width="3023" height="1714" alt="image" src="https://github.com/user-attachments/assets/41317efd-5fc1-497a-a75b-b5069cacf628" />
<img width="3023" height="1708" alt="image" src="https://github.com/user-attachments/assets/dd442ab5-7b86-49fd-9267-681a9547c357" />
<img width="3023" height="1713" alt="image" src="https://github.com/user-attachments/assets/ecee3168-8146-4c35-b390-72a30dd522cf" />
