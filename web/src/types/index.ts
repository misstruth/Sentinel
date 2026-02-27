// ==================== 订阅相关类型 ====================
export interface Subscription {
  id: number
  name: string
  description: string
  source_type: SourceType
  source_url: string
  status: SubscriptionStatus
  config: string // JSON配置
  cron_expr: string
  last_fetch_at?: string
  next_fetch_at?: string
  fetch_timeout: number
  auth_type: AuthType
  auth_token?: string
  keywords: string // JSON数组
  min_severity: SeverityLevel
  tags: string
  total_events: number
  failed_fetches: number
  created_at: string
  updated_at: string
}

export type SourceType =
  | 'vulnerability'
  | 'threat_intel'
  | 'vendor_advisory'
  | 'attack_activity'
  | 'github_repo'
  | 'rss'
  | 'webhook'
  | 'nvd'
  | 'cve'

export type SubscriptionStatus = 'active' | 'paused' | 'disabled'
export type AuthType = 'none' | 'api_key' | 'oauth' | 'basic'

export interface CreateSubscriptionRequest {
  name: string
  description?: string
  source_type: SourceType
  source_url: string
  cron_expr?: string
  fetch_timeout?: number
  auth_type?: AuthType
  auth_token?: string
  keywords?: string[]
  min_severity?: SeverityLevel
  tags?: string[]
  config?: string
}

export interface UpdateSubscriptionRequest {
  name?: string
  description?: string
  source_url?: string
  cron_expr?: string
  fetch_timeout?: number
  auth_type?: AuthType
  auth_token?: string
  keywords?: string[]
  min_severity?: SeverityLevel
  tags?: string[]
}

// ==================== 安全事件相关类型 ====================
export interface SecurityEvent {
  id: number
  subscription_id: number
  title: string
  description: string
  severity: SeverityLevel
  status: EventStatus
  source: string
  source_url: string
  event_time: string
  raw_data?: string
  cve_id?: string
  cvss_score?: number
  affected_vendor?: string
  affected_product?: string
  tags?: string
  unique_hash?: string
  is_starred?: boolean
  risk_score?: number
  recommendation?: string
  affected_assets?: number
  created_at: string
  updated_at?: string
}

export type SeverityLevel = 'critical' | 'high' | 'medium' | 'low' | 'info'
export type EventStatus = 'new' | 'processing' | 'resolved' | 'ignored'

export interface EventFilter {
  subscription_id?: number
  severity?: SeverityLevel
  status?: EventStatus
  start_time?: string
  end_time?: string
  keyword?: string
  cve_id?: string
  page?: number
  size?: number
}

export interface EventStats {
  total: number
  by_severity: Record<SeverityLevel, number>
  by_status: Record<EventStatus, number>
}

// ==================== 报告相关类型 ====================
export interface Report {
  id: number
  title: string
  type: ReportType
  status: ReportStatus
  start_time: string
  end_time: string
  content: string
  summary?: string
  template_id?: number
  template_data?: string
  event_ids: string // JSON数组
  subscription_ids: string // JSON数组
  event_count: number
  critical_count: number
  high_count: number
  generated_by: 'manual' | 'scheduled' | 'api'
  error_msg?: string
  created_at: string
  updated_at: string
}

export type ReportType = 'daily' | 'weekly' | 'monthly' | 'custom' | 'vuln_alert' | 'threat_brief'
export type ReportStatus = 'pending' | 'generating' | 'completed' | 'failed'
export type ExportFormat = 'markdown' | 'html' | 'json'

export interface ReportTemplate {
  id: number
  name: string
  description: string
  type: ReportType
  content: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface GenerateReportRequest {
  title: string
  type: ReportType
  start_time?: string
  end_time?: string
  event_ids?: number[]
  subscription_ids?: number[]
  template_id?: number
}

export interface CreateTemplateRequest {
  name: string
  description?: string
  type: ReportType
  content: string
  is_default?: boolean
}

// ==================== 抓取日志类型 ====================
export interface FetchLog {
  id: number
  subscription_id: number
  status: FetchStatus
  event_count: number
  error_msg?: string
  duration: number // 毫秒
  created_at: string
}

export type FetchStatus = 'success' | 'failed' | 'timeout'

export interface FetchStats {
  total_fetches: number
  success_count: number
  failed_count: number
  total_events: number
  avg_duration_ms: number
}

// ==================== AI 对话类型 ====================
export interface ChatMessage {
  role: 'user' | 'assistant' | 'system'
  content: string
}

export interface ChatRequest {
  message: string
  history?: ChatMessage[]
  context?: string
}

export interface ChatResponse {
  reply: string
  tokens_used?: number
}

export interface ChatStreamRequest {
  message: string
  history?: ChatMessage[]
}

export interface AIOpsRequest {
  type: 'log_analysis' | 'anomaly_detection' | 'root_cause'
  data: string
  context?: string
}

export interface AIOpsResponse {
  result: string
  suggestions?: string[]
  related_events?: SecurityEvent[]
}

// ==================== 通知配置类型 ====================
export interface NotifyConfig {
  email_enabled: boolean
  email_addresses: string[]
  email_levels: SeverityLevel[]
  webhook_enabled: boolean
  webhook_url: string
  webhook_secret?: string
}

// ==================== 用户认证类型 ====================
export interface User {
  id: number
  username: string
  email: string
  role: UserRole
  created_at: string
}

export type UserRole = 'admin' | 'analyst' | 'viewer'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

// ==================== 通用响应类型 ====================
export interface ApiResponse<T> {
  code?: number
  message: string
  data: T
}

export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  size: number
}

// ==================== 仪表盘统计类型 ====================
export interface DashboardStats {
  total_subscriptions: number
  active_subscriptions: number
  paused_subscriptions: number
  total_events: number
  new_events: number
  critical_events: number
  high_events: number
  total_reports: number
  events_today: number
  events_this_week: number
}

export interface TrendData {
  date: string
  total: number
  critical: number
  high: number
  medium: number
  low: number
  info: number
}

// ==================== 审计日志类型 ====================
export interface AuditLog {
  id: number
  user_id: number
  action: string
  resource_type: string
  resource_id: number
  details: string
  ip_address: string
  created_at: string
}
