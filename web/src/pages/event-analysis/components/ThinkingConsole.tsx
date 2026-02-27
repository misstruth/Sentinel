import { useEffect, useRef, useState } from 'react'
import { cn } from '@/utils'
import { Terminal, ChevronRight, ChevronDown, Cpu, CheckCircle, AlertCircle, Database, Filter, Shield, Lightbulb } from 'lucide-react'
import { AgentLog } from '@/types/agent'

interface Props {
  logs: AgentLog[]
  isProcessing: boolean
}

const agentColors: Record<string, string> = {
  '数据采集Agent': '#00F0E0',
  '提取Agent': '#A855F7',
  '去重Agent': '#22C55E',
  '风险评估Agent': '#F43F5E',
  '解决方案Agent': '#F59E0B',
}

const agentIcons: Record<string, typeof Database> = {
  '数据采集Agent': Database,
  '提取Agent': Filter,
  '去重Agent': Cpu,
  '风险评估Agent': Shield,
  '解决方案Agent': Lightbulb,
}

// Typewriter text component
function TypewriterText({ text, speed = 20 }: { text: string; speed?: number }) {
  const [displayed, setDisplayed] = useState('')
  const [done, setDone] = useState(false)

  useEffect(() => {
    if (!text) return
    let i = 0
    setDisplayed('')
    setDone(false)
    const timer = setInterval(() => {
      i++
      setDisplayed(text.slice(0, i))
      if (i >= text.length) {
        clearInterval(timer)
        setDone(true)
      }
    }, speed)
    return () => clearInterval(timer)
  }, [text, speed])

  return (
    <span className="text-[#C9D1D9] break-all text-[11px] leading-relaxed">
      {displayed}
      {!done && (
        <span className="inline-block w-1 h-3 bg-[#00F0E0] ml-0.5 align-middle"
          style={{ animation: 'cyber-cursor 0.6s step-end infinite' }} />
      )}
    </span>
  )
}

export default function ThinkingConsole({ logs, isProcessing }: Props) {
  const scrollRef = useRef<HTMLDivElement>(null)
  const [collapsed, setCollapsed] = useState(true)
  const prevLogsLen = useRef(0)

  // Auto-expand when processing starts, track latest log for typewriter
  useEffect(() => {
    if (isProcessing && collapsed) setCollapsed(false)
  }, [isProcessing])

  // Auto-scroll on new logs
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
    prevLogsLen.current = logs.length
  }, [logs])

  // Only apply typewriter to the latest log entry
  const latestIdx = logs.length - 1

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'running':
        return <Cpu className="w-3 h-3 text-[#00F0E0] animate-pulse" />
      case 'success':
        return <CheckCircle className="w-3 h-3 text-[#22C55E]" />
      case 'error':
        return <AlertCircle className="w-3 h-3 text-[#F43F5E]" />
      default:
        return <ChevronRight className="w-3 h-3 text-[#484F58]" />
    }
  }

  return (
    <div className={cn(
      'flex flex-col bg-[#0D1117] border border-[#30363D]/60 rounded-xl overflow-hidden transition-all duration-300',
      collapsed ? 'h-[38px]' : 'h-full'
    )}>
      {/* Header - clickable to toggle */}
      <div
        className="flex items-center gap-2 px-3 py-2 border-b border-[#30363D]/60 bg-[#161B22]/80 cursor-pointer select-none shrink-0"
        onClick={() => setCollapsed(c => !c)}
      >
        {collapsed ? (
          <ChevronRight className="w-3 h-3 text-[#484F58] transition-transform" />
        ) : (
          <ChevronDown className="w-3 h-3 text-[#484F58] transition-transform" />
        )}
        <Terminal className="w-3.5 h-3.5 text-[#00F0E0]" />
        <span className="text-[11px] font-semibold text-[#E6EDF3] tracking-wide">思考链路</span>
        <span className="text-[10px] text-[#484F58] font-mono">REASONING TRACE</span>
        {isProcessing && (
          <span className="ml-auto flex items-center gap-1.5 text-[10px] text-[#00F0E0] font-mono">
            <span className="w-1.5 h-1.5 rounded-full bg-[#00F0E0] animate-ping" />
            LIVE
          </span>
        )}
        {!isProcessing && logs.length > 0 && (
          <span className="ml-auto text-[10px] text-[#484F58] font-mono">{logs.length} entries</span>
        )}
      </div>

      {/* Log Stream */}
      {!collapsed && (
        <div ref={scrollRef} className="flex-1 overflow-auto p-2 space-y-0.5 font-mono text-xs scrollbar-thin">
          {logs.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-[#484F58] gap-2">
              <Terminal className="w-5 h-5 opacity-40" />
              <span className="text-[11px]">等待 Agent 启动...</span>
            </div>
          ) : (
            logs.map((log, idx) => {
              const AgentIcon = agentIcons[log.agent]
              const isLatest = idx === latestIdx && isProcessing
              return (
                <div
                  key={idx}
                  className={cn(
                    'flex items-start gap-2 px-2 py-1.5 rounded-md transition-all',
                    log.status === 'running' && 'bg-[#00F0E0]/[0.03]',
                    log.status === 'success' && 'bg-[#22C55E]/[0.02]',
                  )}
                >
                  <span className="mt-0.5 shrink-0">{getStatusIcon(log.status)}</span>
                  {AgentIcon && (
                    <AgentIcon
                      className="w-3 h-3 mt-0.5 shrink-0"
                      style={{ color: agentColors[log.agent] || '#484F58' }}
                    />
                  )}
                  <span
                    className="font-semibold shrink-0 text-[11px]"
                    style={{ color: agentColors[log.agent] || '#484F58' }}
                  >
                    [{log.agent.replace('Agent', '')}]
                  </span>
                  {isLatest ? (
                    <TypewriterText text={log.message} speed={18} />
                  ) : (
                    <span className="text-[#C9D1D9] break-all text-[11px] leading-relaxed">{log.message}</span>
                  )}
                  {log.timestamp && (
                    <span className="ml-auto text-[#30363D] shrink-0 text-[10px] tabular-nums">
                      {new Date(log.timestamp).toLocaleTimeString()}
                    </span>
                  )}
                </div>
              )
            })
          )}

          {/* 打字机光标 */}
          {isProcessing && (
            <div className="flex items-center gap-1.5 px-2 py-1 text-[#00F0E0]">
              <span
                className="w-1.5 h-3.5 bg-[#00F0E0]"
                style={{ animation: 'cyber-cursor 1s step-end infinite' }}
              />
              <span className="text-[10px] text-[#30363D] font-mono">awaiting next instruction...</span>
            </div>
          )}
        </div>
      )}
    </div>
  )
}