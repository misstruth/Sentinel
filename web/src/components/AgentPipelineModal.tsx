import { useState } from 'react'
import { X, Cpu, CheckCircle, Loader2, Database, Filter, Brain, Save } from 'lucide-react'
import { eventService } from '@/services/event'

interface PipelineStep {
  agent: string
  status: string
  message: string
  count: number
}

interface Props {
  onClose: () => void
  onComplete: () => void
}

const agentIcons: Record<string, typeof Cpu> = {
  '数据采集': Database,
  '提取Agent': Cpu,
  '去重Agent': Filter,
  '风险评估Agent': Brain,
  '数据持久化': Save,
}

export default function AgentPipelineModal({ onClose, onComplete }: Props) {
  const [running, setRunning] = useState(false)
  const [steps, setSteps] = useState<PipelineStep[]>([])
  const [result, setResult] = useState<{ total: number; new: number } | null>(null)

  const runPipeline = async () => {
    setRunning(true)
    setSteps([
      { agent: '数据采集', status: 'running', message: '正在获取待处理事件...', count: 0 },
    ])

    try {
      const res = await eventService.processPipeline()
      setSteps(res.steps || [])
      setResult({ total: res.total_count, new: res.new_count })
      onComplete()
    } catch {
      setSteps([{ agent: '错误', status: 'error', message: '处理失败', count: 0 }])
    } finally {
      setRunning(false)
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-slate-900 rounded-lg w-[500px] max-h-[80vh] overflow-hidden">
        <div className="flex items-center justify-between p-4 border-b border-slate-700">
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <Cpu className="w-5 h-5 text-primary-400" />
            多Agent协作处理
          </h2>
          <button onClick={onClose} className="text-slate-400 hover:text-white">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-4 space-y-4">
          {/* Agent流程图 */}
          <div className="bg-slate-800 rounded-lg p-4">
            <div className="text-sm text-slate-400 mb-3">处理流程</div>
            <div className="flex items-center justify-between text-xs">
              {['采集', '提取', '去重', '评估', '存储'].map((s, i) => (
                <div key={s} className="flex items-center">
                  <div className="w-12 h-12 rounded-full bg-slate-700 flex items-center justify-center text-primary-400">
                    {s}
                  </div>
                  {i < 4 && <div className="w-8 h-0.5 bg-slate-600" />}
                </div>
              ))}
            </div>
          </div>

          {/* 执行步骤 */}
          {steps.length > 0 && (
            <div className="space-y-2">
              {steps.map((step, i) => {
                const Icon = agentIcons[step.agent] || Cpu
                return (
                  <div key={i} className="flex items-center gap-3 p-3 bg-slate-800 rounded-lg">
                    <Icon className="w-5 h-5 text-primary-400" />
                    <div className="flex-1">
                      <div className="font-medium">{step.agent}</div>
                      <div className="text-sm text-slate-400">{step.message}</div>
                    </div>
                    {step.status === 'completed' ? (
                      <CheckCircle className="w-5 h-5 text-green-500" />
                    ) : step.status === 'running' ? (
                      <Loader2 className="w-5 h-5 text-primary-400 animate-spin" />
                    ) : null}
                  </div>
                )
              })}
            </div>
          )}

          {/* 结果 */}
          {result && (
            <div className="bg-green-900/30 border border-green-700 rounded-lg p-4 text-center">
              <div className="text-green-400 font-semibold">处理完成</div>
              <div className="text-sm text-slate-300 mt-1">
                共处理 {result.total} 个事件，更新 {result.new} 个
              </div>
            </div>
          )}

          {/* 操作按钮 */}
          <button
            onClick={runPipeline}
            disabled={running}
            className="w-full btn-primary py-3"
          >
            {running ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                Agent处理中...
              </>
            ) : (
              <>
                <Cpu className="w-4 h-4" />
                启动多Agent处理
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  )
}
