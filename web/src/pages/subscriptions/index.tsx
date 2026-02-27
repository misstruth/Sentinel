import { useState, useEffect } from 'react'
import {
  Plus,
  Search,
  MoreHorizontal,
  Play,
  Pause,
  Trash2,
  RefreshCw,
  Github,
  Rss,
  Database,
  CheckCircle,
  PauseCircle,
  Edit,
  FileSearch,
  Ban,
  Loader2,
} from 'lucide-react'
import { cn, formatRelativeTime, getSourceTypeLabel } from '@/utils'
import AddSubscriptionModal from './components/AddSubscriptionModal'
import { subscriptionService } from '@/services/subscription'
import type { Subscription } from '@/types'
import toast from 'react-hot-toast'

const sourceTypeIcons: Record<string, typeof Github> = {
  github_repo: Github,
  rss: Rss,
  nvd: Database,
  cve: Database,
  vulnerability: Database,
  threat_intel: Database,
  vendor_advisory: Database,
  attack_activity: Database,
  webhook: Database,
}

const statusConfig = {
  active: {
    icon: CheckCircle,
    class: 'status-active',
    label: '运行中',
  },
  paused: {
    icon: PauseCircle,
    class: 'status-paused',
    label: '已暂停',
  },
  disabled: {
    icon: Ban,
    class: 'status-disabled',
    label: '已禁用',
  },
}

export default function Subscriptions() {
  const [showAddModal, setShowAddModal] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const [filterType, setFilterType] = useState<string>('all')
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [activeDropdown, setActiveDropdown] = useState<number | null>(null)

  const [subscriptions, setSubscriptions] = useState<Subscription[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(20)

  const fetchSubscriptions = async () => {
    try {
      setLoading(true)
      const res = await subscriptionService.list(page, pageSize)
      setSubscriptions(res.list || [])
      setTotal(res.total || 0)
    } catch (error) {
      toast.error('获取订阅列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSubscriptions()
  }, [page])

  const handlePause = async (id: number) => {
    try {
      await subscriptionService.pause(id)
      toast.success('已暂停订阅')
      fetchSubscriptions()
    } catch (error) {
      toast.error('操作失败')
    }
    setActiveDropdown(null)
  }

  const handleResume = async (id: number) => {
    try {
      await subscriptionService.resume(id)
      toast.success('已恢复订阅')
      fetchSubscriptions()
    } catch (error) {
      toast.error('操作失败')
    }
    setActiveDropdown(null)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这个订阅吗？')) return
    try {
      await subscriptionService.delete(id)
      toast.success('已删除订阅')
      fetchSubscriptions()
    } catch (error) {
      toast.error('删除失败')
    }
    setActiveDropdown(null)
  }

  const handleFetch = async (id: number) => {
    try {
      toast.loading('正在抓取...', { id: 'fetch' })
      const res = await subscriptionService.fetch(id)
      toast.success(
        `抓取完成: 获取${res.fetched_count}条, 新增${res.new_count}条, 总计${res.total_events}条, 耗时${res.duration_ms}ms`,
        { id: 'fetch', duration: 4000 }
      )
      fetchSubscriptions()
    } catch (error) {
      toast.error('抓取失败', { id: 'fetch' })
    }
    setActiveDropdown(null)
  }

  const filteredSubscriptions = subscriptions.filter((sub) => {
    const matchesSearch = sub.name.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesType = filterType === 'all' || sub.source_type === filterType
    const matchesStatus = filterStatus === 'all' || sub.status === filterStatus
    return matchesSearch && matchesType && matchesStatus
  })

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900 tracking-tight">订阅管理</h1>
          <p className="text-sm text-slate-500 mt-1">管理安全事件数据源订阅</p>
        </div>
        <button onClick={() => setShowAddModal(true)} className="btn-primary">
          <Plus className="w-4 h-4" />
          添加订阅
        </button>
      </div>

      {/* Filters */}
      <div className="card card-body">
        <div className="flex flex-col sm:flex-row gap-4">
          {/* Search */}
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
            <input
              type="text"
              placeholder="搜索订阅名称..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input pl-9"
            />
          </div>

          {/* Type Filter */}
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
            className="select w-40"
          >
            <option value="all">全部类型</option>
            <option value="github_repo">GitHub</option>
            <option value="rss">RSS</option>
            <option value="nvd">NVD</option>
            <option value="cve">CVE</option>
            <option value="vulnerability">漏洞</option>
            <option value="threat_intel">威胁情报</option>
            <option value="vendor_advisory">厂商公告</option>
            <option value="attack_activity">攻击活动</option>
            <option value="webhook">Webhook</option>
          </select>

          {/* Status Filter */}
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className="select w-32"
          >
            <option value="all">全部状态</option>
            <option value="active">运行中</option>
            <option value="paused">已暂停</option>
            <option value="disabled">已禁用</option>
          </select>

          {/* Refresh Button */}
          <button
            onClick={fetchSubscriptions}
            disabled={loading}
            className="btn-default"
          >
            <RefreshCw className={cn('w-4 h-4', loading && 'animate-spin')} />
            刷新
          </button>
        </div>
      </div>

      {/* Subscription Table */}
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
                  <th>名称</th>
                  <th>类型</th>
                  <th>状态</th>
                  <th>Cron 表达式</th>
                  <th>上次抓取</th>
                  <th className="text-right">事件数</th>
                  <th className="w-20">操作</th>
                </tr>
              </thead>
              <tbody>
                {filteredSubscriptions.map((sub) => {
                  const SourceIcon = sourceTypeIcons[sub.source_type] || Database
                  const config = statusConfig[sub.status as keyof typeof statusConfig]
                  const StatusIcon = config?.icon || CheckCircle

                  return (
                    <tr key={sub.id}>
                      <td>
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 rounded bg-slate-100 flex items-center justify-center flex-shrink-0">
                            <SourceIcon className="w-4 h-4 text-slate-500" />
                          </div>
                          <div className="min-w-0">
                            <p className="font-medium text-slate-900 truncate">{sub.name}</p>
                            <a
                              href={sub.source_url}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="text-xs text-slate-500 hover:text-blue-600 truncate block max-w-[200px]"
                            >
                              {sub.source_url}
                            </a>
                          </div>
                        </div>
                      </td>
                      <td>
                        <span className="tag tag-default">
                          {getSourceTypeLabel(sub.source_type)}
                        </span>
                      </td>
                      <td>
                        <span className={cn('flex items-center gap-1.5', config?.class)}>
                          <StatusIcon className="w-3.5 h-3.5" />
                          {config?.label || sub.status}
                        </span>
                      </td>
                      <td className="text-slate-700 font-mono text-xs">
                        {sub.cron_expr || '-'}
                      </td>
                      <td className="text-slate-600">
                        {sub.last_fetch_at ? formatRelativeTime(sub.last_fetch_at) : '-'}
                      </td>
                      <td className="text-right text-slate-900 font-medium">
                        {(sub.total_events || 0).toLocaleString()}
                      </td>
                      <td>
                        <div className="relative">
                          <button
                            onClick={() => setActiveDropdown(activeDropdown === sub.id ? null : sub.id)}
                            className="p-1.5 rounded text-slate-500 hover:text-slate-900 hover:bg-slate-100 transition-colors"
                          >
                            <MoreHorizontal className="w-4 h-4" />
                          </button>
                          {activeDropdown === sub.id && (
                            <>
                              <div
                                className="fixed inset-0 z-40"
                                onClick={() => setActiveDropdown(null)}
                              />
                              <div className="dropdown right-0 top-full mt-1 z-50">
                                <button className="dropdown-item w-full" onClick={() => handleFetch(sub.id)}>
                                  <RefreshCw className="w-4 h-4" />
                                  立即抓取
                                </button>
                                <button className="dropdown-item w-full">
                                  <FileSearch className="w-4 h-4" />
                                  查看日志
                                </button>
                                <button className="dropdown-item w-full">
                                  <Edit className="w-4 h-4" />
                                  编辑
                                </button>
                                <div className="dropdown-divider" />
                                {sub.status === 'active' ? (
                                  <button
                                    className="dropdown-item w-full text-warning-500"
                                    onClick={() => handlePause(sub.id)}
                                  >
                                    <Pause className="w-4 h-4" />
                                    暂停
                                  </button>
                                ) : (
                                  <button
                                    className="dropdown-item w-full text-success-500"
                                    onClick={() => handleResume(sub.id)}
                                  >
                                    <Play className="w-4 h-4" />
                                    启动
                                  </button>
                                )}
                                <button
                                  className="dropdown-item w-full text-danger-500"
                                  onClick={() => handleDelete(sub.id)}
                                >
                                  <Trash2 className="w-4 h-4" />
                                  删除
                                </button>
                              </div>
                            </>
                          )}
                        </div>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        )}

        {!loading && filteredSubscriptions.length === 0 && (
          <div className="empty">
            <Database className="empty-icon" />
            <p className="empty-text">没有找到匹配的订阅</p>
            <button
              onClick={() => setShowAddModal(true)}
              className="btn-primary mt-4"
            >
              <Plus className="w-4 h-4" />
              添加订阅
            </button>
          </div>
        )}

        {/* Pagination */}
        {!loading && filteredSubscriptions.length > 0 && (
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

      {/* Add Modal */}
      <AddSubscriptionModal
        isOpen={showAddModal}
        onClose={() => setShowAddModal(false)}
        onSuccess={fetchSubscriptions}
      />
    </div>
  )
}
