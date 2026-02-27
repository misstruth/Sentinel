import api from './api'
import {
  Report,
  ReportTemplate,
  GenerateReportRequest,
  CreateTemplateRequest,
  ExportFormat,
  ApiResponse,
  PaginatedResponse,
} from '@/types'

// 后端返回的报告列表响应格式
interface ReportListResponse {
  list: Array<{
    id: number
    title: string
    type: string
    status: string
    summary: string
    event_count: number
    critical_count: number
    high_count: number
    generated_by: string
    error_msg: string
    created_at: string
  }>
  total: number
}

export const reportService = {
  // 生成报告
  async generate(data: GenerateReportRequest): Promise<Report> {
    const res = await api.post<ApiResponse<{
      report_id: number
      title: string
      summary: string
      status: string
    }>>('/report/generate', data)
    const d = res.data.data
    return {
      id: d.report_id,
      title: d.title,
      type: data.type,
      status: d.status as Report['status'],
      start_time: data.start_time || '',
      end_time: data.end_time || '',
      content: '',
      summary: d.summary,
      event_ids: '',
      subscription_ids: '',
      event_count: 0,
      critical_count: 0,
      high_count: 0,
      generated_by: 'manual',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }
  },

  // 获取报告列表
  async list(page = 1, pageSize = 20, type?: string): Promise<PaginatedResponse<Report>> {
    const params: Record<string, unknown> = { page, page_size: pageSize }
    if (type && type !== 'all') {
      params.type = type
    }
    const res = await api.get<ApiResponse<ReportListResponse>>('/report', { params })
    const data = res.data.data

    const list: Report[] = (data.list || []).map(item => ({
      id: item.id,
      title: item.title,
      type: item.type as Report['type'],
      status: item.status as Report['status'],
      start_time: '',
      end_time: '',
      content: '',
      summary: item.summary || '',
      event_ids: '',
      subscription_ids: '',
      event_count: item.event_count || 0,
      critical_count: item.critical_count || 0,
      high_count: item.high_count || 0,
      generated_by: (item.generated_by || 'manual') as Report['generated_by'],
      error_msg: item.error_msg || '',
      created_at: item.created_at,
      updated_at: item.created_at,
    }))

    return {
      list,
      total: data.total || 0,
      page,
      size: pageSize,
    }
  },

  // 获取单个报告
  async get(id: number): Promise<Report> {
    const res = await api.get<ApiResponse<{
      id: number
      title: string
      type: string
      status: string
      content: string
      summary: string
      event_count: number
      critical_count: number
      high_count: number
      start_time: string
      end_time: string
      generated_by: string
      error_msg: string
      created_at: string
    }>>(`/report/${id}`)
    const d = res.data.data
    return {
      id: d.id,
      title: d.title,
      type: d.type as Report['type'],
      status: d.status as Report['status'],
      start_time: d.start_time || '',
      end_time: d.end_time || '',
      content: d.content,
      summary: d.summary || '',
      event_ids: '',
      subscription_ids: '',
      event_count: d.event_count || 0,
      critical_count: d.critical_count || 0,
      high_count: d.high_count || 0,
      generated_by: (d.generated_by || 'manual') as Report['generated_by'],
      error_msg: d.error_msg || '',
      created_at: d.created_at,
      updated_at: d.created_at,
    }
  },

  // 删除报告
  async delete(id: number): Promise<void> {
    await api.delete(`/report/${id}`)
  },

  // 导出报告
  async export(id: number, format: ExportFormat): Promise<Blob> {
    const res = await api.get<ApiResponse<{ content: string; filename: string }>>(`/report/${id}/export`, {
      params: { format },
    })
    const { content } = res.data.data
    // 将内容转换为 Blob
    const mimeTypes: Record<string, string> = {
      markdown: 'text/markdown',
      html: 'text/html',
      json: 'application/json',
    }
    return new Blob([content], { type: mimeTypes[format] || 'text/plain' })
  },

  // 获取模板列表
  async listTemplates(): Promise<ReportTemplate[]> {
    const res = await api.get<ApiResponse<{ list: Array<{
      id: number
      name: string
      description: string
      type: string
      is_default: boolean
    }> }>>('/report/template')
    return (res.data.data.list || []).map(item => ({
      id: item.id,
      name: item.name,
      description: item.description,
      type: item.type as ReportTemplate['type'],
      content: '',
      is_default: item.is_default,
      created_at: '',
      updated_at: '',
    }))
  },

  // 创建模板
  async createTemplate(data: CreateTemplateRequest): Promise<ReportTemplate> {
    const res = await api.post<ApiResponse<{ id: number }>>('/report/template', data)
    return {
      id: res.data.data.id,
      name: data.name,
      description: data.description || '',
      type: data.type,
      content: data.content,
      is_default: data.is_default || false,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }
  },

  // 获取单个模板
  async getTemplate(id: number): Promise<ReportTemplate> {
    const res = await api.get<ApiResponse<ReportTemplate>>(`/report/template/${id}`)
    return res.data.data
  },

  // 更新模板
  async updateTemplate(id: number, data: Partial<CreateTemplateRequest>): Promise<void> {
    await api.put(`/report/template/${id}`, data)
  },

  // 删除模板
  async deleteTemplate(id: number): Promise<void> {
    await api.delete(`/report/template/${id}`)
  },
}
