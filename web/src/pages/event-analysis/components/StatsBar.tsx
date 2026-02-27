import { useEffect, useState } from 'react'
import { Shield, AlertTriangle, Activity, TrendingUp } from 'lucide-react'
import { cn } from '@/utils'

interface Props {
  data: {
    maxCVSS: number
    count: number
    avgRisk: number
    critical?: number
    highRisk?: number
  } | null
  isProcessing: boolean
}

function AnimatedNumber({ value, decimals = 0, color }: { value: number; decimals?: number; color: string }) {
  const [display, setDisplay] = useState(0)

  useEffect(() => {
    if (value === 0) { setDisplay(0); return }
    const duration = 800
    const start = performance.now()
    const from = 0

    const tick = (now: number) => {
      const elapsed = now - start
      const progress = Math.min(elapsed / duration, 1)
      const eased = 1 - Math.pow(1 - progress, 3)
      setDisplay(from + (value - from) * eased)
      if (progress < 1) requestAnimationFrame(tick)
    }
    requestAnimationFrame(tick)
  }, [value])

  return (
    <span className="tabular-nums font-mono" style={{ color }}>
      {decimals > 0 ? display.toFixed(decimals) : Math.round(display)}
    </span>
  )
}

export default function StatsBar({ data, isProcessing }: Props) {
  const stats = [
    {
      label: '最高CVSS',
      value: data?.maxCVSS ?? 0,
      decimals: 1,
      icon: AlertTriangle,
      color: '#F43F5E',
      glow: data != null && data.maxCVSS >= 9,
    },
    {
      label: '严重事件',
      value: data?.critical ?? 0,
      icon: Shield,
      color: '#F43F5E',
    },
    {
      label: '高危事件',
      value: data?.highRisk ?? 0,
      icon: Activity,
      color: '#F97316',
    },
    {
      label: '分析总数',
      value: data?.count ?? 0,
      icon: TrendingUp,
      color: '#00F0E0',
    },
  ]

  return (
    <div className="grid grid-cols-4 gap-3">
      {stats.map((stat, idx) => {
        const Icon = stat.icon
        const hasValue = data != null
        return (
          <div
            key={idx}
            className={cn(
              'relative p-3 rounded-xl border transition-all duration-300',
              'bg-[#0D1117]/80 border-[#30363D]/60 backdrop-blur-sm',
              stat.glow && 'border-[#F43F5E]/40'
            )}
          >
            {/* 严重时的脉冲背景 */}
            {stat.glow && (
              <div className="absolute inset-0 rounded-xl bg-[#F43F5E]/5 animate-pulse" />
            )}

            <div className="relative flex items-center gap-3">
              <div
                className="w-9 h-9 rounded-lg flex items-center justify-center shrink-0"
                style={{ backgroundColor: stat.color + '15' }}
              >
                <Icon className="w-4 h-4" style={{ color: stat.color }} />
              </div>
              <div>
                <div className="text-lg font-bold leading-tight">
                  {hasValue ? (
                    <AnimatedNumber
                      value={stat.value}
                      decimals={stat.decimals || 0}
                      color={stat.color}
                    />
                  ) : (
                    <span className={cn(
                      'text-[#30363D]',
                      isProcessing && 'animate-pulse'
                    )}>--</span>
                  )}
                </div>
                <div className="text-[10px] text-[#484F58] tracking-wide mt-0.5">
                  {stat.label}
                </div>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
