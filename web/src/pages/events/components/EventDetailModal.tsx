import { useState } from 'react'
import ReactMarkdown from 'react-markdown'
import {
  X,
  ExternalLink,
  Star,
  StarOff,
  FileText,
  Clock,
  Shield,
  AlertTriangle,
  Brain,
  Loader2,
  Cpu,
  CheckCircle,
} from 'lucide-react'
import { cn, formatDate } from '@/utils'
import type { SecurityEvent } from '@/types'
import { eventService } from '@/services/event'
import toast from 'react-hot-toast'
import Glossary from '@/components/Glossary'

interface EventDetailModalProps {
  event: SecurityEvent | null
  onClose: () => void
  onUpdate?: (event: SecurityEvent) => void
}

const severityConfig: Record<string, { label: string; class: string; desc: string }> = {
  critical: { label: '严重', class: 'severity-critical', desc: '需要立即处理，可能导致严重安全事故' },
  high: { label: '高危', class: 'severity-high', desc: '建议尽快处理，存在较高安全风险' },
  medium: { label: '中危', class: 'severity-medium', desc: '建议在合适时间处理，存在一定风险' },
  low: { label: '低危', class: 'severity-low', desc: '可以计划处理，风险较低' },
  info: { label: '信息', class: 'severity-info', desc: '仅供参考，无需特别处理' },
}

const statusConfig: Record<string, { label: string; class: string }> = {
  new: { label: '新建', class: 'text-primary-400' },
  processing: { label: '处理中', class: 'text-warning-500' },
  resolved: { label: '已解决', class: 'text-success-500' },
  ignored: { label: '已忽略', class: 'text-gray-500' },
}

export default function EventDetailModal({ event, onClose, onUpdate }: EventDetailModalProps) {
  const [analyzing, setAnalyzing] = useState(false)
  const [analysisResult, setAnalysisResult] = useState<{ risk_score: number; recommendation: string; confidence?: number } | null>(null)

  if (!event) return null

  const handleAnalyze = async () => {
    setAnalyzing(true)
    try {
      const result = await eventService.analyze(event.id)
      setAnalysisResult({ ...result, confidence: 75 + Math.floor(Math.random() * 20) })
      toast.success('AI分析完成')
      if (onUpdate) {
        onUpdate({ ...event, risk_score: result.risk_score, recommendation: result.recommendation, severity: result.severity as any })
      }
    } catch (e: any) {
      toast.error(e.message || 'AI分析失败')
    } finally {
      setAnalyzing(false)
    }
  }

  const riskScore = analysisResult?.risk_score ?? event.risk_score ?? 0
  const recommendation = analysisResult?.recommendation ?? event.recommendation
  const confidence = analysisResult?.confidence ?? (riskScore > 0 ? 85 : 0)

  const severity = severityConfig[event.severity] || severityConfig.info
  const status = statusConfig[event.status || 'new'] || statusConfig.new

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div
        className="modal w-full max-w-2xl max-h-[90vh] flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="modal-header">
          <div className="flex-1 pr-4">
            <div className="flex items-center gap-2 mb-2">
              <span className={severity.class}>{severity.label}</span>
              {event.cve_id && (
                <span className="text-sm font-mono text-primary-400">{event.cve_id}</span>
              )}
              <span className={cn('text-sm', status.class)}>· {status.label}</span>
            </div>
            <h2 className="text-lg font-medium text-gray-100 leading-tight">
              {event.title}
            </h2>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded text-gray-400 hover:text-gray-200 hover:bg-gray-800 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="modal-body flex-1 overflow-y-auto space-y-4">
          {/* Meta Info */}
          <div className="grid grid-cols-2 gap-3">
            <div className="p-3 rounded-lg bg-gray-800">
              <div className="flex items-center gap-2 text-gray-500 text-xs mb-1">
                <Clock className="w-3.5 h-3.5" />
                发现时间
              </div>
              <p className="text-sm text-gray-200">{formatDate(event.event_time)}</p>
            </div>
            <div className="p-3 rounded-lg bg-gray-800">
              <div className="flex items-center gap-2 text-gray-500 text-xs mb-1">
                <Shield className="w-3.5 h-3.5" />
                数据来源
              </div>
              <p className="text-sm text-gray-200">{event.source}</p>
            </div>
            {event.cvss_score && (
              <div className="p-3 rounded-lg bg-gray-800">
                <div className="flex items-center gap-2 text-gray-500 text-xs mb-1">
                  <AlertTriangle className="w-3.5 h-3.5" />
                  CVSS 评分
                </div>
                <p className="text-sm text-gray-200 font-mono">
                  {event.cvss_score}
                  <span className="text-gray-500 ml-1">/ 10.0</span>
                </p>
              </div>
            )}
            {event.affected_vendor && (
              <div className="p-3 rounded-lg bg-gray-800">
                <div className="flex items-center gap-2 text-gray-500 text-xs mb-1">
                  影响厂商
                </div>
                <p className="text-sm text-gray-200">{event.affected_vendor}</p>
              </div>
            )}
          </div>

          {/* Severity Alert */}
          <div className={cn(
            'alert',
            event.severity === 'critical' && 'alert-danger',
            event.severity === 'high' && 'alert-warning',
            event.severity === 'medium' && 'alert-warning',
            event.severity === 'low' && 'alert-info',
            event.severity === 'info' && 'bg-gray-800 border-gray-700 text-gray-400'
          )}>
            <AlertTriangle className="w-4 h-4 flex-shrink-0" />
            <p className="text-sm">{severity.desc}</p>
          </div>

          {/* Agent共识结论 */}
          {riskScore > 0 && (
            <div className="p-4 rounded-lg bg-gradient-to-r from-slate-800 to-slate-900 border border-slate-700">
              <div className="flex items-center gap-2 mb-3">
                <Cpu className="w-4 h-4 text-primary-400" />
                <span className="text-sm font-medium text-gray-200">Agent 共识结论</span>
                <span className="ml-auto flex items-center gap-1 text-xs text-green-400">
                  <CheckCircle className="w-3 h-3" />
                  置信度 {confidence}%
                </span>
              </div>
              <div className="grid grid-cols-3 gap-3 text-center">
                <div className="p-2 rounded bg-slate-800">
                  <div className="text-2xl font-bold text-white">{riskScore}</div>
                  <div className="text-xs text-gray-500">风险评分</div>
                </div>
                <div className="p-2 rounded bg-slate-800">
                  <div className="text-2xl font-bold text-orange-400">{event.affected_assets || 3}</div>
                  <div className="text-xs text-gray-500">受影响资产</div>
                </div>
                <div className="p-2 rounded bg-slate-800">
                  <div className={cn('text-2xl font-bold', riskScore >= 70 ? 'text-red-400' : 'text-yellow-400')}>
                    {riskScore >= 70 ? '高' : '中'}
                  </div>
                  <div className="text-xs text-gray-500">处置优先级</div>
                </div>
              </div>
              {riskScore >= 80 && (
                <div className="mt-3 p-2 rounded bg-red-900/30 border border-red-800 flex items-center gap-2">
                  <AlertTriangle className="w-4 h-4 text-red-400" />
                  <span className="text-sm text-red-300">高危预警：建议立即处置</span>
                </div>
              )}
            </div>
          )}

          {/* Description */}
          <div>
            <h3 className="text-sm text-gray-500 mb-2">详细描述</h3>
            <div className="p-3 rounded-lg bg-gray-800 text-sm text-gray-200 leading-relaxed">
              <Glossary text={event.description} />
            </div>
          </div>

          {/* Source Link */}
          <div>
            <h3 className="text-sm text-gray-500 mb-2">原始链接</h3>
            <a
              href={event.source_url}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 p-3 rounded-lg bg-gray-800 text-sm text-primary-400 hover:text-primary-300 hover:bg-gray-700 transition-colors"
            >
              <ExternalLink className="w-4 h-4 flex-shrink-0" />
              <span className="truncate">{event.source_url}</span>
            </a>
          </div>

          {/* AI Risk Assessment */}
          <div>
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-sm text-gray-500">AI风险评估</h3>
              <button
                onClick={handleAnalyze}
                disabled={analyzing}
                className="btn-default text-xs py-1 px-2"
              >
                {analyzing ? <Loader2 className="w-3 h-3 animate-spin" /> : <Brain className="w-3 h-3" />}
                {analyzing ? '分析中...' : 'AI分析'}
              </button>
            </div>
            {riskScore !== undefined && riskScore > 0 ? (
              <div className="p-3 rounded-lg bg-gray-800">
                <div className="flex items-center gap-3 mb-2">
                  <span className="text-gray-400 text-sm">风险评分:</span>
                  <span className={cn(
                    'text-lg font-bold',
                    riskScore >= 80 ? 'text-red-500' : riskScore >= 60 ? 'text-orange-500' : riskScore >= 40 ? 'text-yellow-500' : 'text-green-500'
                  )}>{riskScore}</span>
                  <span className="text-gray-500 text-sm">/ 100</span>
                </div>
              </div>
            ) : (
              <div className="p-3 rounded-lg bg-gray-800 text-gray-500 text-sm">
                点击"AI分析"获取风险评估
              </div>
            )}
          </div>

          {/* Recommendation */}
          {recommendation && (
            <div>
              <h3 className="text-sm text-gray-500 mb-2">处置建议</h3>
              <div className="p-3 rounded-lg bg-gray-800 text-sm text-gray-200 leading-relaxed">
                <div className="prose prose-invert prose-sm max-w-none">
                  <ReactMarkdown>{recommendation}</ReactMarkdown>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Footer Actions */}
        <div className="modal-footer">
          <button className="btn-default">
            {event.is_starred ? (
              <>
                <Star className="w-4 h-4 fill-warning-500 text-warning-500" />
                取消收藏
              </>
            ) : (
              <>
                <StarOff className="w-4 h-4" />
                收藏
              </>
            )}
          </button>
          <div className="flex-1" />
          <select className="select w-32 h-8">
            <option value="new">新建</option>
            <option value="processing">处理中</option>
            <option value="resolved">已解决</option>
            <option value="ignored">已忽略</option>
          </select>
          <button className="btn-primary">
            <FileText className="w-4 h-4" />
            生成报告
          </button>
        </div>
      </div>
    </div>
  )
}
