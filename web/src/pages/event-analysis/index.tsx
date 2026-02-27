import { useState } from 'react'
import { FileText, Sparkles, Shield, Database, Filter, Cpu, CheckCircle, AlertTriangle, Zap, Lightbulb } from 'lucide-react'
import { useEventStore } from '@/stores/eventStore'
import AgentFlowGraph from './components/AgentFlowGraph'
import ThinkingConsole from './components/ThinkingConsole'
import StatsBar from './components/StatsBar'
import StartButton from './components/StartButton'
import ResultPanel from './components/ResultPanel'
import ReportModal from './components/ReportModal'
import { cn } from '@/utils'

interface RiskData {
  maxCVSS: number
  count: number
  avgRisk: number
  critical?: number
  highRisk?: number
  events?: Array<{
    id: number
    title: string
    desc: string
    cve_id: string
    cvss: number
    severity: string
    vendor: string
    product: string
    source_url: string
  }>
}

const steps = [
  { id: 1, label: '数据采集', agent: '数据采集Agent', icon: Database, color: '#00F0E0' },
  { id: 2, label: '智能提取', agent: '提取Agent', icon: Filter, color: '#A855F7' },
  { id: 3, label: '去重过滤', agent: '去重Agent', icon: Cpu, color: '#22C55E' },
  { id: 4, label: '风险评估', agent: '风险评估Agent', icon: Shield, color: '#F43F5E' },
  { id: 5, label: '解决方案', agent: '解决方案Agent', icon: Lightbulb, color: '#F59E0B' },
]

export default function EventAnalysis() {
  const { clearLogs, isProcessing, setProcessing, agentLogs, addLog } = useEventStore()
  const [riskData, setRiskData] = useState<RiskData | null>(null)
  const [showReport, setShowReport] = useState(false)
  const [showConclusion, setShowConclusion] = useState(false)

  const getStepStatus = (agent: string) => {
    const logs = agentLogs.filter(l => l.agent === agent)
    if (logs.some(l => l.status === 'success')) return 'completed'
    if (logs.some(l => l.status === 'running')) return 'running'
    return 'pending'
  }

  const startAnalysis = async () => {
    clearLogs()
    setProcessing(true)
    setRiskData(null)
    setShowConclusion(false)
    try {
      const response = await fetch('/api/event/pipeline/stream', { method: 'POST' })
      const reader = response.body?.getReader()
      const decoder = new TextDecoder()
      while (reader) {
        const { done, value } = await reader.read()
        if (done) break
        const lines = decoder.decode(value).split('\n').filter(l => l.startsWith('data: '))
        for (const line of lines) {
          const data = JSON.parse(line.slice(6))
          addLog({
            agent: data.agent,
            status: data.status,
            message: data.message,
            timestamp: data.timestamp,
            data: data.data
          })
          if (data.agent === '风险评估Agent' && data.type === 'agent_complete') {
            setRiskData(data.data)
          }
          if (data.type === 'pipeline_done') {
            setProcessing(false)
            setShowConclusion(true)
            return
          }
        }
      }
    } catch {
      setProcessing(false)
    }
  }

  return (
    <div className="relative h-full flex flex-col bg-[#010409] overflow-hidden">
      {/* 顶部操作栏 */}
      <div className="px-4 py-3 border-b border-[#30363D]/50 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Sparkles className="w-4 h-4 text-[#00F0E0]" />
          <span className="text-sm font-semibold text-[#E6EDF3] tracking-wide">多 Agent 协作研判</span>
          <span className="text-[10px] text-[#484F58] font-mono">MULTI-AGENT ANALYSIS</span>
        </div>
        <div className="flex items-center gap-2">
          {riskData && (
            <button
              onClick={() => setShowReport(true)}
              className="px-3 py-1.5 rounded-lg text-xs font-semibold transition-all flex items-center gap-1.5 border border-transparent"
              style={{
                background: 'linear-gradient(135deg, #D4A017, #F5C842, #D4A017)',
                color: '#1A1200',
              }}
            >
              <FileText className="w-3.5 h-3.5" />
              生成报告
            </button>
          )}
          <StartButton isProcessing={isProcessing} onStart={startAnalysis} hasData={!!riskData} />
        </div>
      </div>

      {/* 步骤条 Stepper */}
      <div className="px-4 py-2.5 border-b border-[#30363D]/30 bg-[#0D1117]/30">
        <div className="flex items-center justify-center gap-1">
          {steps.map((step, idx) => {
            const status = getStepStatus(step.agent)
            const Icon = step.icon
            return (
              <div key={step.id} className="flex items-center">
                <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg" style={{
                  background: status === 'running' ? `${step.color}10` : 'transparent',
                }}>
                  <div className={cn(
                    'w-6 h-6 rounded-md flex items-center justify-center transition-all',
                    status === 'running' && 'animate-pulse',
                  )} style={{
                    backgroundColor: status !== 'pending' ? step.color + '20' : '#30363D20',
                  }}>
                    {status === 'completed' ? (
                      <CheckCircle className="w-3.5 h-3.5" style={{ color: step.color }} />
                    ) : (
                      <Icon className="w-3.5 h-3.5" style={{ color: status === 'running' ? step.color : '#484F58' }} />
                    )}
                  </div>
                  <span className={cn(
                    'text-[11px] font-medium tracking-wide transition-colors',
                    status === 'running' ? 'text-[#E6EDF3]' : status === 'completed' ? 'text-[#8B949E]' : 'text-[#484F58]',
                  )}>{step.label}</span>
                </div>
                {idx < steps.length - 1 && (
                  <div className="w-8 h-[1px] mx-1" style={{
                    background: status === 'completed'
                      ? `linear-gradient(90deg, ${step.color}40, ${steps[idx + 1].color}40)`
                      : '#30363D30',
                  }} />
                )}
              </div>
            )
          })}
        </div>
      </div>

      {/* 统计指标栏 */}
      <div className="px-4 py-3 border-b border-[#30363D]/30">
        <StatsBar data={riskData} isProcessing={isProcessing} />
      </div>

      {/* 主内容区 */}
      <div className="flex-1 flex overflow-hidden">
        {/* 左侧：Agent 协作图 + 思考链路 */}
        <div className="flex-1 flex flex-col p-4 gap-3">
          {/* Agent 协作拓扑图 */}
          <div className="flex-1 rounded-xl border border-[#30363D]/60 bg-[#0D1117]/40 overflow-hidden backdrop-blur-sm">
            <AgentFlowGraph logs={agentLogs} isProcessing={isProcessing} />
          </div>
          {/* 思考链路控制台 */}
          <div className="h-[220px] shrink-0">
            <ThinkingConsole logs={agentLogs} isProcessing={isProcessing} />
          </div>
        </div>

        {/* 右侧：结果面板 */}
        <div className="w-[380px] border-l border-[#30363D]/50 bg-[#0D1117]/20">
          <ResultPanel data={riskData} isProcessing={isProcessing} />
        </div>
      </div>

      {/* AI 研判结论卡片 */}
      {showConclusion && riskData && (
        <div className="absolute inset-0 z-20 flex items-center justify-center bg-black/40 backdrop-blur-sm"
          onClick={() => setShowConclusion(false)}>
          <div
            className="max-w-md w-full mx-4 rounded-xl border overflow-hidden"
            style={{
              animation: 'cyber-pop-in 0.4s ease-out',
              borderColor: riskData.maxCVSS >= 9 ? '#F43F5E40' : '#00F0E040',
              background: 'linear-gradient(135deg, #0D1117, #161B22)',
            }}
            onClick={e => e.stopPropagation()}
          >
            {/* 结论卡片头部 */}
            <div className="px-5 pt-5 pb-3 flex items-center gap-2">
              <Zap className="w-5 h-5 text-[#00F0E0]" />
              <span className="text-sm font-bold text-[#E6EDF3]">AI 研判结论</span>
            </div>

            {/* 结论内容 */}
            <div className="px-5 pb-4">
              <p className="text-[13px] text-[#C9D1D9] leading-relaxed">
                检测到 <strong className="text-[#F43F5E]">{riskData.critical || 0}</strong> 处严重漏洞、
                <strong className="text-[#F97316]">{riskData.highRisk || 0}</strong> 处高危漏洞，
                最高 CVSS <strong className="text-[#F43F5E]">{riskData.maxCVSS}</strong>，
                {riskData.maxCVSS >= 9
                  ? '建议立即启动应急响应，4小时内完成修复。'
                  : riskData.maxCVSS >= 7
                    ? '建议24小时内更新防护规则。'
                    : '风险可控，建议72小时内排查处理。'}
              </p>
            </div>

            {/* 操作按钮 */}
            <div className="px-5 pb-5 flex items-center gap-2">
              <button
                onClick={() => { setShowConclusion(false); setShowReport(true) }}
                className="flex-1 py-2 rounded-lg text-xs font-bold flex items-center justify-center gap-1.5 border border-transparent"
                style={{ background: 'linear-gradient(135deg, #D4A017, #F5C842)', color: '#1A1200' }}
              >
                <FileText className="w-3.5 h-3.5" />
                生成报告
              </button>
              <button
                onClick={() => setShowConclusion(false)}
                className="flex-1 py-2 rounded-lg text-xs font-bold flex items-center justify-center gap-1.5 border border-transparent"
                style={{ background: 'linear-gradient(135deg, #16A34A, #22C55E)', color: '#fff' }}
              >
                <AlertTriangle className="w-3.5 h-3.5" />
                应用规则
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 报告弹窗 */}
      <ReportModal
        visible={showReport}
        onClose={() => setShowReport(false)}
        data={riskData}
        logs={agentLogs}
      />
    </div>
  )
}