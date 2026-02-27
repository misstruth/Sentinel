import { useState, useEffect } from 'react'
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
  monthly: { label: '月报', class: 'bg-purple-500/15 text-purple-400 border border-purple-500/30' },
  vuln_alert: { label: '漏洞告警', class: 'tag-danger' },
  threat_brief: { label: '威胁简报', class: 'tag-warning' },
  custom: { label: '自定义', class: 'tag-default' },
}

export default function Reports() {
  const [searchQuery, setSearchQuery] = useState('')
  const [typeFilter, setTypeFilter] = useState('all')
  const [selectedReport, setSelectedReport] = useState<Report | null>(null)
  const [showGenerateModal, setShowGenerateModal] = useState(false)
  const [activeDropdown, setActiveDropdown] = useState<number | null>(null)

  const [reports, setReports] = useState<Report[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(20)

  const fetchReports = async () => {
    try {
      setLoading(true)
      const res = await reportService.list(page, pageSize)
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
  }, [page])

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

  const filteredReports = reports.filter((report) => {
    const matchesSearch = report.title.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesType = typeFilter === 'all' || report.type === typeFilter
    return matchesSearch && matchesType
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
            <option value="custom">自定义</option>
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
                  <th className="w-20">事件数</th>
                  <th className="w-32">创建时间</th>
                  <th className="w-24">操作</th>
                </tr>
              </thead>
              <tbody>
                {filteredReports.map((report) => {
                  const typeConfig = reportTypeConfig[report.type] || reportTypeConfig.custom
                  return (
                    <tr
                      key={report.id}
                      className="cursor-pointer"
                      onClick={() => setSelectedReport(report)}
                    >
                      <td>
                        <div className="min-w-0">
                          <p className="font-medium text-slate-900 truncate">{report.title}</p>
                          <p className="text-xs text-slate-500 truncate mt-0.5">{report.summary || '-'}</p>
                        </div>
                      </td>
                      <td>
                        <span className={cn('tag', typeConfig.class)}>{typeConfig.label}</span>
                      </td>
                      <td className="text-slate-700">{report.event_count}</td>
                      <td className="text-slate-600 text-sm">
                        {formatDate(report.created_at, 'YYYY-MM-DD')}
                      </td>
                      <td onClick={(e) => e.stopPropagation()}>
                        <div className="flex items-center gap-1">
                          <button
                            onClick={() => setSelectedReport(report)}
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
                          <div className="relative">
                            <button
                              onClick={() => setActiveDropdown(activeDropdown === report.id ? null : report.id)}
                              className="p-1.5 rounded text-slate-500 hover:text-slate-900 hover:bg-slate-100 transition-colors"
                            >
                              <MoreHorizontal className="w-4 h-4" />
                            </button>
                            {activeDropdown === report.id && (
                              <>
                                <div
                                  className="fixed inset-0 z-40"
                                  onClick={() => setActiveDropdown(null)}
                                />
                                <div className="dropdown right-0 top-full mt-1 z-50">
                                  <button
                                    className="dropdown-item w-full"
                                    onClick={() => handleExport(report.id, 'markdown')}
                                  >
                                    <Download className="w-4 h-4" />
                                    导出 Markdown
                                  </button>
                                  <button
                                    className="dropdown-item w-full"
                                    onClick={() => handleExport(report.id, 'html')}
                                  >
                                    <Download className="w-4 h-4" />
                                    导出 HTML
                                  </button>
                                  <button
                                    className="dropdown-item w-full"
                                    onClick={() => handleExport(report.id, 'json')}
                                  >
                                    <Download className="w-4 h-4" />
                                    导出 JSON
                                  </button>
                                  <div className="dropdown-divider" />
                                  <button
                                    className="dropdown-item w-full text-danger-500"
                                    onClick={() => handleDelete(report.id)}
                                  >
                                    <Trash2 className="w-4 h-4" />
                                    删除
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
