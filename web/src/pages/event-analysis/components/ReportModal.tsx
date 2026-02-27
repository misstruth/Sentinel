import { X, Download, AlertTriangle, ExternalLink } from 'lucide-react'

interface EventDetail {
  id: number; title: string; desc: string; cve_id: string; cvss: number
  severity: string; vendor: string; product: string; source_url: string
}

interface Props {
  visible: boolean
  onClose: () => void
  data: { count: number; maxCVSS: number; avgRisk: number; critical?: number; highRisk?: number; events?: EventDetail[] } | null
  logs: { agent: string; message: string }[]
}

const severityColor: Record<string, string> = {
  critical: 'text-[#F85149] bg-[#F85149]/10',
  high: 'text-[#F0883E] bg-[#F0883E]/10',
  medium: 'text-[#E3B341] bg-[#E3B341]/10',
  low: 'text-[#8B949E] bg-[#8B949E]/10',
}

export default function ReportModal({ visible, onClose, data, logs }: Props) {
  if (!visible || !data) return null

  const events = data.events || []
  const urgency = data.maxCVSS >= 9 ? '4小时内' : data.maxCVSS >= 7 ? '24小时内' : '72小时内'

  const generateMarkdown = () => {
    const eventsMd = events.map(e =>
      `### ${e.cve_id || '未知CVE'} - ${e.title}\n- 严重程度: ${e.severity} | CVSS: ${e.cvss}\n- 厂商/产品: ${e.vendor || '未知'} / ${e.product || '未知'}\n- 描述: ${e.desc || '无描述'}\n- 来源: ${e.source_url || '无'}\n`
    ).join('\n')

    return `# 安全事件分析报告
生成时间: ${new Date().toLocaleString()}

## 摘要
- 分析事件: ${data.count} 个
- 最高CVSS: ${data.maxCVSS}
- 严重漏洞: ${data.critical || 0} 个
- 高危漏洞: ${data.highRisk || 0} 个
- 建议响应时间: ${urgency}

## 事件详情
${eventsMd}
## Agent分析轨迹
${logs.filter(l => l.agent).map(l => `- [${l.agent}] ${l.message}`).join('\n')}

---
*本报告由 Sentinel-Agent 自动生成*`
  }

  const download = () => {
    const blob = new Blob([generateMarkdown()], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `安全分析报告_${new Date().toISOString().slice(0,10)}.md`
    a.click()
    URL.revokeObjectURL(url)
  }

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-[#0D1117] border border-[#30363D] rounded-lg w-[800px] max-h-[85vh] flex flex-col" onClick={e => e.stopPropagation()}>
        {/* 头部 */}
        <div className="flex items-center justify-between p-4 border-b border-[#30363D]">
          <h2 className="text-sm font-bold text-[#E6EDF3]">安全分析报告</h2>
          <button onClick={onClose} className="text-[#8B949E] hover:text-white"><X className="w-4 h-4" /></button>
        </div>

        {/* 内容 */}
        <div className="flex-1 overflow-auto p-4 space-y-4">
          {/* 摘要 */}
          <div className={`p-3 rounded border ${data.maxCVSS >= 9 ? 'bg-[#F85149]/10 border-[#F85149]/40' : 'bg-[#F0883E]/10 border-[#F0883E]/40'}`}>
            <div className="flex items-center gap-2 mb-2">
              <AlertTriangle className={`w-4 h-4 ${data.maxCVSS >= 9 ? 'text-[#F85149]' : 'text-[#F0883E]'}`} />
              <span className="text-xs font-bold text-[#E6EDF3]">摘要</span>
            </div>
            <p className="text-xs text-[#E6EDF3]">
              发现 <strong className="text-[#F85149]">{data.critical || 0}</strong> 个严重漏洞，
              <strong className="text-[#F0883E]">{data.highRisk || 0}</strong> 个高危漏洞，
              最高CVSS <strong>{data.maxCVSS}</strong>，建议 <strong className="text-[#F0883E]">{urgency}</strong> 响应。
            </p>
          </div>

          {/* 事件列表 */}
          <div>
            <div className="text-xs text-[#8B949E] mb-2">事件详情 ({events.length})</div>
            <div className="space-y-2 max-h-[400px] overflow-auto">
              {events.map(e => (
                <div key={e.id} className="p-3 bg-[#161B22] border border-[#30363D] rounded">
                  <div className="flex items-start justify-between gap-2 mb-2">
                    <div className="flex-1">
                      <div className="text-xs font-bold text-[#E6EDF3] mb-1">{e.title}</div>
                      <div className="flex items-center gap-2 text-[10px]">
                        {e.cve_id && <span className="text-[#58A6FF]">{e.cve_id}</span>}
                        <span className={`px-1.5 py-0.5 rounded ${severityColor[e.severity] || severityColor.low}`}>{e.severity}</span>
                        <span className="text-[#8B949E]">CVSS: {e.cvss}</span>
                      </div>
                    </div>
                    {e.source_url && (
                      <a href={e.source_url} target="_blank" rel="noopener noreferrer" className="text-[#58A6FF] hover:underline">
                        <ExternalLink className="w-3 h-3" />
                      </a>
                    )}
                  </div>
                  <div className="text-[11px] text-[#8B949E] mb-2 line-clamp-2">{e.desc || '无描述'}</div>
                  <div className="flex gap-2 text-[10px] text-[#8B949E]">
                    {e.vendor && <span>厂商: {e.vendor}</span>}
                    {e.product && <span>产品: {e.product}</span>}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* 底部 */}
        <div className="p-4 border-t border-[#30363D] flex justify-end gap-2">
          <button onClick={onClose} className="px-3 py-1.5 border border-[#30363D] rounded text-[#8B949E] text-xs">关闭</button>
          <button onClick={download} className="px-3 py-1.5 bg-[#58A6FF] text-white text-xs rounded flex items-center gap-1">
            <Download className="w-3 h-3" /> 下载报告
          </button>
        </div>
      </div>
    </div>
  )
}
