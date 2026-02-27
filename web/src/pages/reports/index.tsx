import { useState, useEffect, useCallback } from 'react'
import {
  Search,
  FileText,
  Download,
  Trash2,
  Eye,
  Sparkles,
  MoreHorizontal,
  RefreshCw,
  Loader2,
} from 'lucide-react'
import { cn, formatDate } from '@/utils'
import ReportDetailModal from './components/ReportDetailModal'
import GenerateReportModal from './components/GenerateReportModal'
import { reportService } from '@/services/report'
import type { Report } from '@/types'
import toast from 'react-hot-toast'

const reportTypeConfig: Record<string, { label: string; class: string }> = {
  daily: { label: '日报', class: 'tag-primary' },
  weekly: { label: '周报', class: 'tag-success' },
  monthly: { label: '月报', class: 'bg-purple-100 text-purple-800' },
  vuln_alert: { label: '漏洞告警', class: 'tag-danger' },
  threat_brief: { label: '威胁简报', class: 'tag-warning' },
  custom: { label: '决策简报', class: 'tag-default' },
}

const statusConfig: Record<string, { label: string; class: string }> = {
  pending: { label: '待生成', class: 'tag-default' },
  generating: { label: '生成中', class: 'tag-warning' },
  completed: { label: '已完成', class: 'tag-success' },
  failed: { label: '失败', class: 'tag-danger' },
}

export default function Reports() {
  const [searchQuery, setSearchQuery] = useState('')
  const [typeFilter, setTypeFilter] = useState('all')
  const [selectedReport, setSelectedReport] = useState<Report | null>(null)
  const [showGenerateModal, setShowGenerateModal] = useState(false)
  const [activeDropdown, setActiveDropdown] = useState<number | null>(null)
  const [dropdownPos, setDropdownPos] = useState<{ top: number; left: number; openUp: boolean }>({ top: 0, left: 0, openUp: false })

  const toggleDropdown = useCallback((id: number, e: React.MouseEvent<HTMLButtonElement>) => {
    if (activeDropdown === id) {
      setActiveDropdown(null)
      return
    }
    const rect = e.currentTarget.getBoundingClientRect()
    const menuHeight = 200
    const openUp = rect.bottom + menuHeight > window.innerHeight
    setDropdownPos({
      top: openUp ? rect.top : rect.bottom + 4,
      left: rect.right - 180,
      openUp,
    })
    setActiveDropdown(id)
  }, [activeDropdown])

  const [reports, setReports] = useState<Report[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(20)

  const fetchReports = async () => {
    try {
      setLoading(true)
      const res = await reportService.list(page, pageSize, typeFilter)
      setReports(res.list || [])
      setTotal(res.total || 0)
    } catch (error) {
      toast.error('获取报告列表失败')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchReports()
  }, [page, typeFilter])

  // 当列表中存在 generating 状态的报告时，每 5s 自动刷新
  useEffect(() => {
    const hasGenerating = reports.some(r => r.status === 'generating')
    if (!hasGenerating) return

    const interval = setInterval(() => {
      fetchReports()
    }, 5000)

    return () => clearInterval(interval)
  }, [reports])

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这个报告吗？')) return
    try {
      await reportService.delete(id)
      toast.success('已删除报告')
      fetchReports()
    } catch (error) {
      toast.error('删除失败')
    }
    setActiveDropdown(null)
  }

  const handleExport = async (id: number, format: 'markdown' | 'html' | 'json') => {
    try {
      const blob = await reportService.export(id, format)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `report-${id}.${format === 'markdown' ? 'md' : format}`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
      toast.success('导出成功')
    } catch (error) {
      toast.error('导出失败')
    }
    setActiveDropdown(null)
  }

  const handleViewReport = async (report: Report) => {
    try {
      const fullReport = await reportService.get(report.id)
      setSelectedReport(fullReport)
    } catch (error) {
      toast.error('获取报告详情失败')
    }
  }

  const filteredReports = reports.filter((report) => {
    return report.title.toLowerCase().includes(searchQuery.toLowerCase())
  })

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900 tracking-tight">分析报告</h1>
          <p className="text-sm text-slate-500 mt-1">AI 生成的安全分析报告</p>
        </div>
        <button
          onClick={() => setShowGenerateModal(true)}
          className="btn-primary"
        >
          <Sparkles className="w-4 h-4" />
          生成报告
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
              placeholder="搜索报告标题..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="input pl-9"
            />
          </div>

          {/* Type Filter */}
          <select
            value={typeFilter}
            onChange={(e) => setTypeFilter(e.target.value)}
            className="select w-36"
          >
            <option value="all">全部类型</option>
            <option value="daily">日报</option>
            <option value="weekly">周报</option>
            <option value="monthly">月报</option>
            <option value="vuln_alert">漏洞告警</option>
            <option value="threat_brief">威胁简报</option>
            <option value="custom">决策简报</option>
          </select>

          {/* Refresh Button */}
          <button
            onClick={fetchReports}
            disabled={loading}
            className="btn-default"
          >
            <RefreshCw className={cn('w-4 h-4', loading && 'animate-spin')} />
            刷新
          </button>
        </div>
      </div>

      {/* Report Table */}
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
                  <th>报告标题</th>
                  <th className="w-24">类型</th>
                  <th className="w-20">状态</th>
                  <th className="w-20">事件数</th>
                  <th className="w-32">创建时间</th>
                  <th className="w-24">操作</th>
                </tr>
              </thead>
              <tbody>
                {filteredReports.map((report) => {
                  const typeConfig = reportTypeConfig[report.type] || reportTypeConfig.custom
                  const stConfig = statusConfig[report.status] || statusConfig.pending
                  return (
                    <tr
                      key={report.id}
                      className="cursor-pointer"
                      onClick={() => handleViewReport(report)}
                    >
                      <td className="max-w-[300px]">
                        <div className="min-w-0">
                          <p className="font-medium text-slate-900 truncate">{report.title}</p>
                          <p className="text-xs text-slate-500 truncate mt-0.5">{report.summary || '-'}</p>
                        </div>
                      </td>
                      <td>
                        <span className={cn('tag', typeConfig.class)}>{typeConfig.label}</span>
                      </td>
                      <td>
                        <span className={cn('tag', stConfig.class)}>
                          {report.status === 'generating' && <Loader2 className="w-3 h-3 animate-spin mr-1" />}
                          {stConfig.label}
                        </span>
                      </td>
                      <td className="text-slate-700">{report.event_count}</td>
                      <td className="text-slate-600 text-sm">
                        {formatDate(report.created_at, 'YYYY-MM-DD')}
                      </td>
                      <td onClick={(e) => e.stopPropagation()}>
                        <div className="flex items-center gap-1">
                          <button
                            onClick={() => handleViewReport(report)}
                            className="p-1.5 rounded text-slate-500 hover:text-slate-900 hover:bg-slate-100 transition-colors"
                            title="查看"
                          >
                            <Eye className="w-4 h-4" />
                          </button>
                          <button
                            className="p-1.5 rounded text-slate-500 hover:text-slate-900 hover:bg-slate-100 transition-colors"
                            title="下载"
                            onClick={() => handleExport(report.id, 'markdown')}
                          >
                            <Download className="w-4 h-4" />
                          </button>
                          <button
                            onClick={(e) => toggleDropdown(report.id, e)}
                            className="p-1.5 rounded text-slate-500 hover:text-slate-900 hover:bg-slate-100 transition-colors"
                          >
                            <MoreHorizontal className="w-4 h-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        )}

        {!loading && filteredReports.length === 0 && (
          <div className="empty">
            <FileText className="empty-icon" />
            <p className="empty-text">没有找到匹配的报告</p>
            <button
              onClick={() => setShowGenerateModal(true)}
              className="btn-primary mt-4"
            >
              <Sparkles className="w-4 h-4" />
              生成报告
            </button>
          </div>
        )}

        {/* Pagination */}
        {!loading && filteredReports.length > 0 && (
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

      {/* Dropdown Portal - fixed position to avoid overflow clipping */}
      {activeDropdown !== null && (
        <>
          <div className="fixed inset-0 z-40" onClick={() => setActiveDropdown(null)} />
          <div
            className="fixed z-50 min-w-[180px] py-2 bg-white rounded-xl border border-slate-200"
            style={{
              top: dropdownPos.openUp ? undefined : dropdownPos.top,
              bottom: dropdownPos.openUp ? window.innerHeight - dropdownPos.top : undefined,
              left: dropdownPos.left,
              boxShadow: '0 10px 40px rgba(0, 0, 0, 0.15)',
            }}
          >
            <button
              className="dropdown-item w-full"
              onClick={() => handleExport(activeDropdown, 'markdown')}
            >
              <Download className="w-4 h-4" />
              导出 Markdown
            </button>
            <button
              className="dropdown-item w-full"
              onClick={() => handleExport(activeDropdown, 'html')}
            >
              <Download className="w-4 h-4" />
              导出 HTML
            </button>
            <button
              className="dropdown-item w-full"
              onClick={() => handleExport(activeDropdown, 'json')}
            >
              <Download className="w-4 h-4" />
              导出 JSON
            </button>
            <div className="dropdown-divider" />
            <button
              className="dropdown-item w-full text-red-600"
              onClick={() => handleDelete(activeDropdown)}
            >
              <Trash2 className="w-4 h-4" />
              删除
            </button>
          </div>
        </>
      )}

      {/* Report Detail Modal */}
      <ReportDetailModal
        report={selectedReport}
        onClose={() => setSelectedReport(null)}
      />

      {/* Generate Report Modal */}
      <GenerateReportModal
        isOpen={showGenerateModal}
        onClose={() => setShowGenerateModal(false)}
        onSuccess={fetchReports}
      />
    </div>
  )
}
