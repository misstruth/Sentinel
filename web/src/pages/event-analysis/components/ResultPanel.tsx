import { useState } from 'react'
import ReactMarkdown from 'react-markdown'
import { cn } from '@/utils'
import { ExternalLink, Shield, ChevronRight, AlertTriangle, Scan, X } from 'lucide-react'

interface EventData {
  id: number
  title: string
  desc: string
  cve_id: string
  cvss: number
  severity: string
  vendor: string
  product: string
  source_url: string
  recommendation?: string
  similar_events?: Array<{ title: string; similarity: number }>
}

interface Props {
  data: {
    maxCVSS: number
    count: number
    avgRisk: number
    critical?: number
    highRisk?: number
    events?: EventData[]
  } | null
  isProcessing?: boolean
}

function SkeletonCard() {
  return (
    <div className="px-3 py-2 rounded-lg border border-[#30363D]/30 bg-[#080C13]">
      <div className="flex gap-2">
        <div className="h-3.5 w-14 rounded cyber-skeleton" />
        <div className="h-3.5 w-full rounded cyber-skeleton" />
      </div>
    </div>
  )
}

const severityConfig: Record<string, { color: string; label: string }> = {
  critical: { color: '#F43F5E', label: 'CRIT' },
  high: { color: '#F97316', label: 'HIGH' },
  medium: { color: '#EAB308', label: 'MED' },
  low: { color: '#8B949E', label: 'LOW' },
}

// Detail slide panel overlay
function EventDetail({ event, onClose }: { event: EventData; onClose: () => void }) {
  const cfg = severityConfig[event.severity] || severityConfig.low
  return (
    <div
      className="absolute inset-0 z-10 bg-[#080C13]/95 backdrop-blur-sm flex flex-col"
      style={{ animation: 'cyber-slide-in 0.25s ease-out' }}
    >
      <div className="flex items-center justify-between px-3 py-2.5 border-b border-[#30363D]/60">
        <span className="text-[11px] font-semibold text-[#E6EDF3]">事件详情</span>
        <button onClick={onClose} className="text-[#484F58] hover:text-[#E6EDF3] transition-colors">
          <X className="w-3.5 h-3.5" />
        </button>
      </div>
      <div className="flex-1 overflow-auto p-3 space-y-3">
        <div className="flex items-center gap-2 flex-wrap">
          <span className="px-1.5 py-0.5 rounded text-[10px] font-bold"
            style={{ backgroundColor: cfg.color + '15', color: cfg.color }}>{cfg.label}</span>
          {event.cve_id && (
            <span className="px-1.5 py-0.5 rounded text-[10px] font-mono"
              style={{ backgroundColor: cfg.color + '15', color: cfg.color }}>{event.cve_id}</span>
          )}
          <span className="px-1.5 py-0.5 rounded text-[10px] font-bold"
            style={{ backgroundColor: cfg.color + '15', color: cfg.color }}>CVSS {event.cvss}</span>
        </div>
        <h3 className="text-[13px] text-[#E6EDF3] font-semibold leading-relaxed">{event.title}</h3>
        {event.desc && <p className="text-[11px] text-[#8B949E] leading-relaxed">{event.desc}</p>}
        <div className="space-y-1.5 text-[10px] text-[#484F58]">
          {event.vendor && <div>厂商: <span className="text-[#8B949E]">{event.vendor}</span></div>}
          {event.product && <div>产品: <span className="text-[#8B949E]">{event.product}</span></div>}
        </div>
        {event.source_url && (
          <a href={event.source_url} target="_blank" rel="noopener noreferrer"
            className="inline-flex items-center gap-1 text-[10px] text-[#00F0E0] hover:underline">
            <ExternalLink className="w-3 h-3" /> 查看来源
          </a>
        )}
        {event.recommendation && (
          <div className="p-2.5 rounded-lg bg-[#22C55E]/5 border border-[#22C55E]/20">
            <div className="text-[10px] text-[#22C55E] font-semibold mb-1.5">AI 解决方案</div>
            <div className="prose prose-invert prose-sm max-w-none text-[11px] text-[#C9D1D9] leading-relaxed">
              <ReactMarkdown>{event.recommendation}</ReactMarkdown>
            </div>
          </div>
        )}
        {event.similar_events && event.similar_events.length > 0 && (
          <div className="p-2.5 rounded-lg bg-[#F59E0B]/5 border border-[#F59E0B]/20">
            <div className="text-[10px] text-[#F59E0B] font-semibold mb-1.5">相似历史事件</div>
            <div className="space-y-1">
              {event.similar_events.map((se, i) => (
                <div key={i} className="text-[10px] text-[#8B949E] flex items-center gap-1.5">
                  <span className="w-1 h-1 rounded-full bg-[#F59E0B]/60 shrink-0" />
                  <span className="flex-1 truncate">{se.title}</span>
                  <span className="text-[#F59E0B] shrink-0">{(se.similarity * 100).toFixed(0)}%</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

// Compact event card - severity badge + title only
function CompactCard({ event, onClick }: { event: EventData; onClick: () => void }) {
  const cfg = severityConfig[event.severity] || severityConfig.low
  const isCritical = event.severity === 'critical'

  return (
    <button
      onClick={onClick}
      className={cn(
        'w-full text-left px-3 py-2 rounded-lg border transition-all duration-200 group',
        'bg-[#080C13] hover:bg-[#0D1117]',
        isCritical
          ? 'border-[#F43F5E]/20 hover:border-[#F43F5E]/40'
          : 'border-[#30363D]/30 hover:border-[#30363D]/60',
      )}
      style={
        isCritical
          ? { animation: 'cyber-glow-pulse 3s ease-in-out infinite', '--glow-color': 'rgba(248,63,94,0.12)' } as React.CSSProperties
          : undefined
      }
    >
      <div className="flex items-center gap-2">
        <span
          className="px-1.5 py-0.5 rounded text-[9px] font-bold shrink-0 tracking-wider"
          style={{ backgroundColor: cfg.color + '15', color: cfg.color }}
        >
          {cfg.label}
        </span>
        <span className="text-[11px] text-[#C9D1D9] truncate flex-1 group-hover:text-[#E6EDF3] transition-colors">
          {event.title}
        </span>
        <ChevronRight className="w-3 h-3 text-[#30363D] group-hover:text-[#484F58] shrink-0 transition-colors" />
      </div>
    </button>
  )
}

export default function ResultPanel({ data, isProcessing }: Props) {
  const [selectedEvent, setSelectedEvent] = useState<EventData | null>(null)

  // Skeleton loading state
  if (isProcessing && !data?.events?.length) {
    return (
      <div className="h-full flex flex-col bg-[#080C13]">
        <div className="px-3 py-2.5 border-b border-[#30363D]/60 flex items-center gap-2">
          <Scan className="w-3.5 h-3.5 text-[#00F0E0] animate-pulse" />
          <span className="text-[11px] font-semibold text-[#E6EDF3] tracking-wide">扫描结果</span>
          <span className="ml-auto text-[10px] text-[#00F0E0] font-mono animate-pulse">scanning...</span>
        </div>
        <div className="flex-1 p-2 space-y-1.5 overflow-hidden">
          <SkeletonCard />
          <SkeletonCard />
          <SkeletonCard />
          <SkeletonCard />
          <SkeletonCard />
        </div>
      </div>
    )
  }

  // Empty state
  if (!data?.events?.length) {
    return (
      <div className="h-full flex items-center justify-center text-[#484F58] bg-[#080C13]">
        <div className="text-center">
          <Shield className="w-10 h-10 mx-auto mb-3 opacity-20" />
          <p className="text-xs tracking-wide">启动分析后查看结果</p>
          <p className="text-[10px] text-[#30363D] mt-1 font-mono">AWAITING ANALYSIS</p>
        </div>
      </div>
    )
  }

  // Main view with compact cards
  return (
    <div className="relative h-full flex flex-col bg-[#080C13]">
      <div className="px-3 py-2.5 border-b border-[#30363D]/60 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <AlertTriangle className="w-3.5 h-3.5 text-[#F43F5E]" />
          <span className="text-[11px] font-semibold text-[#E6EDF3] tracking-wide">高危事件</span>
        </div>
        <span className="text-[10px] text-[#484F58] font-mono">{data.events.length} items</span>
      </div>
      <div className="flex-1 overflow-auto p-2 space-y-1 scrollbar-thin">
        {data.events.map((event) => (
          <CompactCard
            key={event.id}
            event={event}
            onClick={() => setSelectedEvent(event)}
          />
        ))}
      </div>

      {/* Detail slide panel overlay */}
      {selectedEvent && (
        <EventDetail
          event={selectedEvent}
          onClose={() => setSelectedEvent(null)}
        />
      )}
    </div>
  )
}