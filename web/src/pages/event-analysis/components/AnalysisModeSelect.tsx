import { useState, useRef, useEffect } from 'react'
import { CalendarDays, Clock, Crosshair, ChevronDown } from 'lucide-react'
import { cn } from '@/utils'

export type AnalysisMode = 'latest' | 'today' | 'specific'

interface Props {
  value: AnalysisMode
  onChange: (mode: AnalysisMode) => void
  disabled?: boolean
  selectedCount?: number
}

const options = [
  { value: 'latest' as const, label: '最近10条', icon: Clock, desc: '按时间倒序' },
  { value: 'today' as const, label: '今日事件', icon: CalendarDays, desc: '今天产生的事件' },
  { value: 'specific' as const, label: '指定事件', icon: Crosshair, desc: '手动选择事件' },
]

export default function AnalysisModeSelect({ value, onChange, disabled, selectedCount }: Props) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [])

  const current = options.find(o => o.value === value) || options[0]
  const Icon = current.icon

  return (
    <div ref={ref} className="relative">
      <button
        onClick={() => !disabled && setOpen(!open)}
        disabled={disabled}
        className={cn(
          'flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-medium transition-all',
          'border border-[#30363D] bg-[#161B22] text-[#C9D1D9] hover:border-[#484F58]',
          disabled && 'opacity-50 cursor-not-allowed',
        )}
      >
        <Icon className="w-3.5 h-3.5 text-[#00F0E0]" />
        <span>{current.label}</span>
        {value === 'specific' && selectedCount !== undefined && selectedCount > 0 && (
          <span className="px-1.5 py-0.5 rounded bg-[#00F0E0]/15 text-[#00F0E0] text-[10px] font-mono">
            {selectedCount}
          </span>
        )}
        <ChevronDown className={cn('w-3 h-3 text-[#484F58] transition-transform', open && 'rotate-180')} />
      </button>

      {open && (
        <div className="absolute top-full left-0 mt-1 w-48 rounded-lg border border-[#30363D] bg-[#161B22] shadow-xl z-50 overflow-hidden">
          {options.map(opt => {
            const OptIcon = opt.icon
            const active = opt.value === value
            return (
              <button
                key={opt.value}
                onClick={() => { onChange(opt.value); setOpen(false) }}
                className={cn(
                  'w-full flex items-center gap-2.5 px-3 py-2.5 text-left transition-colors',
                  active ? 'bg-[#00F0E0]/10 text-[#00F0E0]' : 'text-[#C9D1D9] hover:bg-[#21262D]',
                )}
              >
                <OptIcon className="w-3.5 h-3.5 shrink-0" />
                <div className="flex flex-col">
                  <span className="text-xs font-medium">{opt.label}</span>
                  <span className="text-[10px] text-[#484F58]">{opt.desc}</span>
                </div>
              </button>
            )
          })}
        </div>
      )}
    </div>
  )
}
