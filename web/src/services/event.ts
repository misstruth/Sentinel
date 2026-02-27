import api from './api'
import {
  SecurityEvent,
  EventFilter,
  EventStatus,
  ApiResponse,
  PaginatedResponse,
} from '@/types'

// 后端返回的事件列表响应格式
interface EventListResponse {
  list: Array<{
    id: number
    title: string
    severity: string
    status: string
    cve_id: string
    event_time: string
  }>
  total: number
}

export const eventService = {
  // 获取事件列表
  async list(filter?: EventFilter): Promise<PaginatedResponse<SecurityEvent>> {
    const params: Record<string, unknown> = {}
    if (filter?.page) params.page = filter.page
    if (filter?.size) params.page_size = filter.size
    if (filter?.severity) params.severity = filter.severity
    if (filter?.status) params.status = filter.status
    if (filter?.keyword) params.keyword = filter.keyword

    const res = await api.get<ApiResponse<EventListResponse>>('/event', { params })
    const data = res.data.data

    // 转换后端返回的数据格式
    const list: SecurityEvent[] = (data.list || []).map(item => ({
      id: item.id,
      subscription_id: 0,
      title: item.title,
      description: '',
      severity: item.severity as SecurityEvent['severity'],
      status: item.status as SecurityEvent['status'],
      source: '',
      source_url: '',
      event_time: item.event_time,
      cve_id: item.cve_id,
      created_at: item.event_time,
    }))

    return {
      list,
      total: data.total || 0,
      page: filter?.page || 1,
      size: filter?.size || 20,
    }
  },

  // 获取单个事件
  async get(id: number): Promise<SecurityEvent> {
    const res = await api.get<ApiResponse<SecurityEvent>>(`/event/${id}`)
    return res.data.data
  },

  // 更新事件状态
  async updateStatus(id: number, status: EventStatus): Promise<void> {
    await api.put(`/event/${id}/status`, { status })
  },

  // 批量更新事件状态
  async batchUpdateStatus(ids: number[], status: EventStatus): Promise<void> {
    await api.post('/event/batch/status', { ids, status })
  },

  // 获取事件统计
  async getStats(filter?: { subscription_id?: number; start_time?: string; end_time?: string }): Promise<{
    total: number
    by_severity: Record<string, number>
    by_status: Record<string, number>
  }> {
    const res = await api.get<ApiResponse<{
      total: number
      by_severity: Record<string, number>
      by_status: Record<string, number>
    }>>('/event/stats', { params: filter })
    return res.data.data
  },

  // 获取事件趋势
  async getTrend(days = 7): Promise<Array<{
    date: string
    total: number
    critical: number
    high: number
    medium: number
    low: number
    info: number
  }>> {
    const res = await api.get<ApiResponse<{
      list: Array<{
        date: string
        total: number
        critical: number
        high: number
        medium: number
        low: number
        info: number
      }>
    }>>('/event/trend', { params: { days } })
    return res.data.data.list || []
  },

  // AI分析事件
  async analyze(id: number): Promise<{ risk_score: number; severity: string; recommendation: string }> {
    const res = await api.post<ApiResponse<{ risk_score: number; severity: string; recommendation: string }>>(`/event/${id}/analyze`)
    return res.data.data
  },

  // 多Agent流水线处理
  async processPipeline(): Promise<{ total_count: number; dedup_count: number; new_count: number; processed_at: string; steps: Array<{agent: string; status: string; message: string; count: number}> }> {
    const res = await api.post<ApiResponse<{ total_count: number; dedup_count: number; new_count: number; processed_at: string; steps: Array<{agent: string; status: string; message: string; count: number}> }>>('/event/pipeline/process')
    return res.data.data
  },
}
