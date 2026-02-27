import api from './api'
import {
  Subscription,
  CreateSubscriptionRequest,
  UpdateSubscriptionRequest,
  ApiResponse,
  PaginatedResponse,
  FetchLog,
  FetchStats,
} from '@/types'

// 后端返回的订阅列表响应格式
interface SubscriptionListResponse {
  items: Subscription[]
  total: number
  page: number
  size: number
}

export const subscriptionService = {
  // 获取订阅列表
  async list(page = 1, pageSize = 20): Promise<PaginatedResponse<Subscription>> {
    const res = await api.get<ApiResponse<SubscriptionListResponse>>('/subscriptions', {
      params: { page, page_size: pageSize },
    })
    const data = res.data.data
    return {
      list: data.items || [],
      total: data.total || 0,
      page: data.page || page,
      size: data.size || pageSize,
    }
  },

  // 获取单个订阅
  async get(id: number): Promise<Subscription> {
    const res = await api.get<ApiResponse<Subscription>>(`/subscriptions/${id}`)
    return res.data.data
  },

  // 创建订阅
  async create(data: CreateSubscriptionRequest): Promise<{ id: number }> {
    const res = await api.post<ApiResponse<{ id: number }>>('/subscriptions', data)
    return res.data.data
  },

  // 更新订阅
  async update(id: number, data: UpdateSubscriptionRequest): Promise<void> {
    await api.put(`/subscriptions/${id}`, data)
  },

  // 删除订阅
  async delete(id: number): Promise<void> {
    await api.delete(`/subscriptions/${id}`)
  },

  // 暂停订阅
  async pause(id: number): Promise<void> {
    await api.post(`/subscriptions/${id}/pause`)
  },

  // 恢复订阅
  async resume(id: number): Promise<void> {
    await api.post(`/subscriptions/${id}/resume`)
  },

  // 禁用订阅
  async disable(id: number): Promise<void> {
    await api.post(`/subscriptions/${id}/disable`)
  },

  // 获取订阅的抓取日志
  async getFetchLogs(subscriptionId: number, page = 1, pageSize = 20): Promise<PaginatedResponse<FetchLog>> {
    const res = await api.get<ApiResponse<PaginatedResponse<FetchLog>>>(`/subscriptions/${subscriptionId}/logs`, {
      params: { page, page_size: pageSize },
    })
    return res.data.data
  },

  // 获取订阅的抓取统计
  async getFetchStats(subscriptionId: number): Promise<FetchStats> {
    try {
      const res = await api.get<ApiResponse<FetchStats>>(`/subscriptions/${subscriptionId}/stats`)
      return res.data.data
    } catch (error) {
      console.warn('获取抓取统计失败:', error)
      return {
        total_fetches: 0,
        success_count: 0,
        failed_count: 0,
        total_events: 0,
        avg_duration_ms: 0,
      }
    }
  },

  // 手动触发抓取
  async fetch(id: number): Promise<{
    fetched_count: number
    new_count: number
    total_events: number
    duration_ms: number
    message: string
  }> {
    const res = await api.post<ApiResponse<any>>(`/subscriptions/${id}/fetch`)
    return res.data.data
  },
}
