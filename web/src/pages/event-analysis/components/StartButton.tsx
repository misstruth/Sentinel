import { Play, Loader2, RotateCcw } from 'lucide-react'
import { cn } from '@/utils'

interface Props {
  isProcessing: boolean
  onStart: () => void
  hasData: boolean
}

export default function StartButton({ isProcessing, onStart, hasData }: Props) {
  return (
    <button
      onClick={onStart}
      disabled={isProcessing}
      className={cn(
        'relative group px-8 py-3.5 rounded-xl font-semibold text-sm transition-all duration-300',
        'flex items-center gap-2.5 overflow-hidden tracking-wide',
        isProcessing
          ? 'bg-[#00F0E0]/10 text-[#00F0E0] cursor-wait border border-[#00F0E0]/30'
          : hasData
            ? 'bg-[#161B22] text-[#8B949E] hover:text-[#E6EDF3] border border-[#30363D] hover:border-[#484F58]'
            : 'bg-gradient-to-r from-[#00F0E0] to-[#00D9FF] text-[#010409] hover:shadow-lg hover:shadow-[#00F0E0]/20 border border-transparent'
      )}
    >
      {/* 呼吸灯外圈 */}
      {!isProcessing && !hasData && (
        <div
          className="absolute -inset-[2px] rounded-xl opacity-40"
          style={{
            background: 'linear-gradient(135deg, #00F0E0, #00D9FF, #A855F7)',
            animation: 'cyber-breathe 2.5s ease-in-out infinite',
            filter: 'blur(8px)',
          }}
        />
      )}

      {/* 扫描线 */}
      {!isProcessing && !hasData && (
        <div className="absolute inset-0 overflow-hidden rounded-xl">
          <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/25 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-700" />
        </div>
      )}

      {/* 处理中的扫描动效 */}
      {isProcessing && (
        <div className="absolute inset-0 overflow-hidden rounded-xl">
          <div
            className="absolute left-0 right-0 h-full"
            style={{
              background: 'linear-gradient(180deg, transparent, #00F0E010, transparent)',
              animation: 'cyber-scan 2s linear infinite',
            }}
          />
        </div>
      )}

      <span className="relative flex items-center gap-2.5">
        {isProcessing ? (
          <>
            <Loader2 className="w-4 h-4 animate-spin" />
            <span className="font-mono">Agent 分析中</span>
            <span className="flex gap-0.5">
              <span className="w-1 h-1 rounded-full bg-current animate-bounce" style={{ animationDelay: '0ms' }} />
              <span className="w-1 h-1 rounded-full bg-current animate-bounce" style={{ animationDelay: '150ms' }} />
              <span className="w-1 h-1 rounded-full bg-current animate-bounce" style={{ animationDelay: '300ms' }} />
            </span>
          </>
        ) : hasData ? (
          <>
            <RotateCcw className="w-4 h-4" />
            <span>重新分析</span>
          </>
        ) : (
          <>
            <Play className="w-4 h-4" />
            <span>启动 AI 研判</span>
          </>
        )}
      </span>
    </button>
  )
}
