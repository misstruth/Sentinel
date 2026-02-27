import { useState, useEffect } from 'react'
import {
  Search,
  Star,
  StarOff,
  ExternalLink,
  CheckCircle,
  Circle,
  FileText,
  Shield,
  MoreHorizontal,
  Check,
  Clock,
  XCircle,
  Eye,
  RefreshCw,
  Loader2,
  Cpu,
} from 'lucide-react'
import { cn, formatRelativeTime } from '@/utils'
import { useContextStore } from '@/stores/contextStore'
import EventDetailModal from './components/EventDetailModal'
import AgentPipelineModal from '@/components/AgentPipelineModal'
import type { SecurityEvent, EventStatus } from '@/types'
import { eventService } from '@/services/event'
import toast from 'react-hot-toast'

const severityConfig: Record<string, { label: string; class: string }> = {
  critical: { label: '严重', class: 'severity-critical' },
  high: { label: '高危', class: 'severity-high' },
  medium: { label: '中危', class: 'severity-medium' },
  low: { label: '低危', class: 'severity-low' },
  info: { label: '信息', class: 'severity-info' },
}

const statusConfig: Record<string, { label: string; icon: typeof Check; class: string }> = {
  new: { label: '新建', icon: Circle, class: 'text-primary-400' },
  processing: { label: '处理中', icon: Clock, class: 'text-warning-500' },
  resolved: { label: '已解决', icon: Check, class: 'text-success-500' },
  ignored: { label: '已忽略', icon: XCircle, class: 'text-gray-500' },
}

export default function Events() {
  const [searchQuery, setSearchQuery] = useState('')
  const [severityFilter, setSeverityFilter] = useState('all')
  const [statusFilter, setStatusFilter] = useState('all')
  const [selectedEvents, setSelectedEvents] = useState<number[]>([])
  const [selectedEvent, setSelectedEvent] = useState<SecurityEvent | null>(null)
  const [activeDropdown, setActiveDropdown] = useState<number | null>(null)

  const [events, setEvents] = useState<SecurityEvent[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(20)
  const [processing, _setProcessing] = useState(false)
  const [showAgentModal, setShowAgentModal] = useState(false)
  const setContext = useContextStore((s) => s.setContext)

  const fetchEvents = async () => {
    try {
      setLoading(true)
      const filter: Record<string, unknown> = { page, size: pageSize }
      if (searchQuery) filter.keyword = searchQuery
      if (severityFilter !== 'all') filter.severity = severityFilter
      if (statusFilter !== 'all') filter.status = statusFilter

      const res = await eventService.list(filter)
      setEvents(res.list || [])
      setTotal(res.total || 0)
    } catch (error) {
      toast.error('获取事件列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchEvents()
  }, [page, severityFilter, statusFilter])

  const handleSearch = () => {
    setPage(1)
    fetchEvents()
  }

  const handleUpdateStatus = async (id: number, status: EventStatus) => {
    try {
      await eventService.updateStatus(id, status)
      toast.success('状态已更新')
      fetchEvents()
    } catch (error) {
      toast.error('更新失败')
    }
    setActiveDropdown(null)
  }

  const handleProcessPipeline = async () => {
    setShowAgentModal(true)
  }

  const handleBatchUpdateStatus = async (status: EventStatus) => {
    if (selectedEvents.length === 0) return
    try {
      await eventService.batchUpdateStatus(selectedEvents, status)
      toast.success('批量更新成功')
      setSelectedEvents([])
      fetchEvents()
    } catch (error) {
      toast.error('批量更新失败')
    }
  }

  const toggleSelectEvent = (id: number) => {
    setSelectedEvents((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    )
  }

  const toggleSelectAll = () => {
    if (selectedEvents.length === events.length) {
      setSelectedEvents([])
    } else {
      setSelectedEvents(events.map((e) => e.id))
    }
  }

  const newCount = events.filter((e) => e.status === 'new').length
  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900 tracking-tight">安全事件</h1>
          <p className="text-sm text-slate-500 mt-1">
            共 {total} 个事件，{newCount} 个待处理
          </p>
        </div>
        <button
          onClick={handleProcessPipeline}
          disabled={processing}
          className="btn-primary"
        >
          {processing ? <Loader2 className="w-4 h-4 animate-spin" /> : <Cpu className="w-4 h-4" />}
          {processing ? 'Agent处理中...' : 'AI Agent分析'}
        </button>
      </div>
      {selectedEvents.length > 0 && (
        <div className="flex items-center gap-2 mt-4">
          <span className="text-sm text-slate-600">已选择 {selectedEvents.length} 项</span>
          <select
              className="select w-32 h-8"
              onChange={(e) => {
                if (e.target.value) {
                  handleBatchUpdateStatus(e.target.value as EventStatus)
                  e.target.value = ''
                }
              }}
            >
              <option value="">批量操作</option>
              <option value="processing">标记处理中</option>
              <option value="resolved">标记已解决</option>
              <option value="ignored">标记已忽略</option>
            </select>
            <button className="btn-primary h-8">
              <FileText className="w-4 h-4" />
              生成报告
            </button>
        </div>
      )}

      {/* Filters */}
      <div className="card card-body">
        <div className="flex flex-col lg:flex-row gap-4">
          {/* Search */}
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
            <input
              type="text"
              placeholder="搜索事件标题、描述、CVE ID..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              className="input pl-9"
            />
          </div>

          {/* Severity Filter */}
          <select
            value={severityFilter}
            onChange={(e) => {
              setSeverityFilter(e.target.value)
              setPage(1)
            }}
            className="select w-32"
          >
            <option value="all">全部级别</option>
            <option value="critical">严重</option>
            <option value="high">高危</option>
            <option value="medium">中危</option>
            <option value="low">低危</option>
            <option value="info">信息</option>
          </select>

          {/* Status Filter */}
          <select
            value={statusFilter}
            onChange={(e) => {
              setStatusFilter(e.target.value)
              setPage(1)
            }}
            className="select w-32"
          >
            <option value="all">全部状态</option>
            <option value="new">新建</option>
            <option value="processing">处理中</option>
            <option value="resolved">已解决</option>
            <option value="ignored">已忽略</option>
          </select>

          {/* Refresh Button */}
          <button
            onClick={fetchEvents}
            disabled={loading}
            className="btn-default"
          >
            <RefreshCw className={cn('w-4 h-4', loading && 'animate-spin')} />
            刷新
          </button>
        </div>
      </div>

      {/* Event Table */}
      <div className="card">
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="w-8 h-8 animate-spin text-primary-500" />
          </div>
        ) : (
          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th className="w-10">
                    <button
                      onClick={toggleSelectAll}
                      className="p-1 rounded hover:bg-slate-100 transition-colors"
                    >
                      {selectedEvents.length === events.length && events.length > 0 ? (
                        <CheckCircle className="w-4 h-4 text-primary-400" />
                      ) : (
                        <Circle className="w-4 h-4 text-gray-500" />
                      )}
                    </button>
                  </th>
                  <th>事件</th>
                  <th className="w-20">级别</th>
                  <th className="w-24">状态</th>
                  <th className="w-20">来源</th>
                  <th className="w-28">时间</th>
                  <th className="w-20">操作</th>
                </tr>
              </thead>
              <tbody>
                {events.map((event) => {
                  const severity = severityConfig[event.severity] || severityConfig.info
                  const status = statusConfig[event.status] || statusConfig.new
                  const StatusIcon = status.icon

                  return (
                    <tr
                      key={event.id}
                      className={cn(
                        'cursor-pointer',
                        event.status === 'new' && 'bg-blue-50/50'
                      )}
                      onClick={() => {
                      setSelectedEvent(event)
                      setContext('events', event.id, event.title)
                    }}
                    >
                      <td onClick={(e) => e.stopPropagation()}>
                        <button
                          onClick={() => toggleSelectEvent(event.id)}
                          className="p-1 rounded hover:bg-slate-100 transition-colors"
                        >
                          {selectedEvents.includes(event.id) ? (
                            <CheckCircle className="w-4 h-4 text-primary-400" />
                          ) : (
                            <Circle className="w-4 h-4 text-gray-500" />
                          )}
                        </button>
                      </td>
                      <td>
                        <div className="flex items-start gap-2">
                          {event.status === 'new' && (
                            <span className="w-1.5 h-1.5 rounded-full bg-primary-500 mt-2 flex-shrink-0" />
                          )}
                          <div className="min-w-0">
                            <div className="flex items-center gap-2">
                              <p className={cn(
                                'text-sm truncate max-w-[350px]',
                                event.status === 'new' ? 'font-medium text-slate-900' : 'text-slate-700'
                              )}>
                                {event.title}
                              </p>
                              {event.source_url && (
                                <a
                                  href={event.source_url}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  onClick={(e) => e.stopPropagation()}
                                  className="text-slate-400 hover:text-primary-500 transition-colors"
                                  title="查看原文"
                                >
                                  <ExternalLink className="w-3.5 h-3.5" />
                                </a>
                              )}
                            </div>
                            {event.cve_id && (
                              <p className="text-xs text-slate-500 font-mono mt-0.5">
                                {event.cve_id}
                                {event.cvss_score && ` · CVSS ${event.cvss_score}`}
                              </p>
                            )}
                          </div>
                        </div>
                      </td>
                      <td>
                        <span className={severity.class}>{severity.label}</span>
                      </td>
                      <td>
                        <span className={cn('flex items-center gap-1.5 text-sm', status.class)}>
                          <StatusIcon className="w-3.5 h-3.5" />
                          {status.label}
                        </span>
                      </td>
                      <td>
                        <span className="tag tag-default">{event.source}</span>
                      </td>
                      <td className="text-slate-600 text-sm whitespace-nowrap">
                        {formatRelativeTime(event.event_time)}
                      </td>
                      <td onClick={(e) => e.stopPropagation()}>
                        <div className="flex items-center gap-1">
                          <button
                            onClick={() => {}}
                            className={cn(
                              'p-1.5 rounded transition-colors',
                              event.is_starred
                                ? 'text-amber-500 hover:bg-slate-100'
                                : 'text-slate-400 hover:text-amber-500 hover:bg-slate-100'
                            )}
                          >
                            {event.is_starred ? (
                              <Star className="w-4 h-4 fill-current" />
                            ) : (
                              <StarOff className="w-4 h-4" />
                            )}
                          </button>
                          <div className="relative">
                            <button
                              onClick={() => setActiveDropdown(activeDropdown === event.id ? null : event.id)}
                              className="p-1.5 rounded text-slate-500 hover:text-slate-900 hover:bg-slate-100 transition-colors"
                            >
                              <MoreHorizontal className="w-4 h-4" />
                            </button>
                            {activeDropdown === event.id && (
                              <>
                                <div
                                  className="fixed inset-0 z-40"
                                  onClick={() => setActiveDropdown(null)}
                                />
                                <div className="dropdown right-0 top-full mt-1 z-50">
                                  <button
                                    className="dropdown-item w-full"
                                    onClick={() => {
                                      setSelectedEvent(event)
                                      setActiveDropdown(null)
                                    }}
                                  >
                                    <Eye className="w-4 h-4" />
                                    查看详情
                                  </button>
                                  {event.source_url && (
                                    <a
                                      href={event.source_url}
                                      target="_blank"
                                      rel="noopener noreferrer"
                                      className="dropdown-item w-full"
                                      onClick={() => setActiveDropdown(null)}
                                    >
                                      <ExternalLink className="w-4 h-4" />
                                      查看原文
                                    </a>
                                  )}
                                  <div className="dropdown-divider" />
                                  <button
                                    className="dropdown-item w-full"
                                    onClick={() => handleUpdateStatus(event.id, 'processing')}
                                  >
                                    <Clock className="w-4 h-4" />
                                    标记处理中
                                  </button>
                                  <button
                                    className="dropdown-item w-full text-success-500"
                                    onClick={() => handleUpdateStatus(event.id, 'resolved')}
                                  >
                                    <Check className="w-4 h-4" />
                                    标记已解决
                                  </button>
                                  <button
                                    className="dropdown-item w-full text-gray-500"
                                    onClick={() => handleUpdateStatus(event.id, 'ignored')}
                                  >
                                    <XCircle className="w-4 h-4" />
                                    标记已忽略
                                  </button>
                                </div>
                              </>
                            )}
                          </div>
                        </div>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        )}

        {!loading && events.length === 0 && (
          <div className="empty">
            <Shield className="empty-icon" />
            <p className="empty-text">没有找到匹配的安全事件</p>
          </div>
        )}

        {/* Pagination */}
        {!loading && events.length > 0 && (
          <div className="px-4 py-3 border-t border-slate-200 flex items-center justify-between">
            <span className="text-sm text-slate-600">
              共 {total} 条记录
            </span>
            <div className="pagination">
              <button
                className={cn('pagination-item', page <= 1 && 'disabled')}
                onClick={() => setPage(p => Math.max(1, p - 1))}
                disabled={page <= 1}
              >
                上一页
              </button>
              <button className="pagination-item active">{page}</button>
              <button
                className={cn('pagination-item', page >= totalPages && 'disabled')}
                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                disabled={page >= totalPages}
              >
                下一页
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Event Detail Modal */}
      <EventDetailModal
        event={selectedEvent}
        onClose={() => setSelectedEvent(null)}
      />

      {/* Agent Pipeline Modal */}
      {showAgentModal && (
        <AgentPipelineModal
          onClose={() => setShowAgentModal(false)}
          onComplete={() => fetchEvents()}
        />
      )}
    </div>
  )
}
