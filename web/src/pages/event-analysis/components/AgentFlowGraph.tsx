import { useEffect, useRef, useState } from 'react'
import { cn } from '@/utils'
import { Database, Filter, Shield, Cpu, CheckCircle, Circle, Lightbulb } from 'lucide-react'
import { AgentLog } from '@/types/agent'

interface Props {
  logs: AgentLog[]
  isProcessing: boolean
}

const agents = [
  { id: '数据采集Agent', icon: Database, color: '#00F0E0', x: 8, y: 50, label: '数据采集', desc: 'Data Collection' },
  { id: '提取Agent', icon: Filter, color: '#A855F7', x: 28, y: 28, label: '智能提取', desc: 'Extraction' },
  { id: '去重Agent', icon: Cpu, color: '#22C55E', x: 48, y: 50, label: '去重过滤', desc: 'Deduplication' },
  { id: '风险评估Agent', icon: Shield, color: '#F43F5E', x: 68, y: 28, label: '风险评估', desc: 'Risk Assessment' },
  { id: '解决方案Agent', icon: Lightbulb, color: '#F59E0B', x: 92, y: 50, label: '解决方案', desc: 'Solution' },
]

const connections = [
  { from: 0, to: 1 },
  { from: 1, to: 2 },
  { from: 2, to: 3 },
  { from: 3, to: 4 },
]

function getBezier(from: { x: number; y: number }, to: { x: number; y: number }, w: number, h: number) {
  const x1 = (from.x / 100) * w
  const y1 = (from.y / 100) * h
  const x2 = (to.x / 100) * w
  const y2 = (to.y / 100) * h
  const mx = (x1 + x2) / 2
  const cy = Math.min(y1, y2) - 30
  return { x1, y1, x2, y2, cx1: mx, cy1: cy, cx2: mx, cy2: cy }
}

function bezierPoint(t: number, p: ReturnType<typeof getBezier>) {
  const u = 1 - t
  const x = u * u * u * p.x1 + 3 * u * u * t * p.cx1 + 3 * u * t * t * p.cx2 + t * t * t * p.x2
  const y = u * u * u * p.y1 + 3 * u * u * t * p.cy1 + 3 * u * t * t * p.cy2 + t * t * t * p.y2
  return { x, y }
}

// Parse hex color to rgba
function hexToRgba(hex: string, alpha: number) {
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  return `rgba(${r},${g},${b},${alpha})`
}

// Draw rotating arc on canvas
function drawRotatingArc(
  ctx: CanvasRenderingContext2D,
  cx: number, cy: number, radius: number,
  color: string, tick: number
) {
  const startAngle = (tick * 0.04) % (Math.PI * 2)
  const arcLen = Math.PI * 1.2

  // Outer glow arc
  ctx.beginPath()
  ctx.arc(cx, cy, radius + 2, startAngle, startAngle + arcLen)
  ctx.strokeStyle = hexToRgba(color, 0.1)
  ctx.lineWidth = 6
  ctx.lineCap = 'round'
  ctx.stroke()

  // Main arc
  ctx.beginPath()
  ctx.arc(cx, cy, radius, startAngle, startAngle + arcLen)
  ctx.strokeStyle = hexToRgba(color, 0.7)
  ctx.lineWidth = 2.5
  ctx.lineCap = 'round'
  ctx.shadowColor = color
  ctx.shadowBlur = 10
  ctx.stroke()
  ctx.shadowBlur = 0

  // Arc head dot
  const headAngle = startAngle + arcLen
  const hx = cx + Math.cos(headAngle) * radius
  const hy = cy + Math.sin(headAngle) * radius
  ctx.beginPath()
  ctx.arc(hx, hy, 3, 0, Math.PI * 2)
  ctx.fillStyle = color
  ctx.fill()
}

// Draw directional arrow on bezier curve
function drawArrowHead(
  ctx: CanvasRenderingContext2D,
  p: ReturnType<typeof getBezier>,
  t: number, color: string, size: number
) {
  const dt = 0.01
  const p1 = bezierPoint(Math.max(0, t - dt), p)
  const p2 = bezierPoint(Math.min(1, t + dt), p)
  const angle = Math.atan2(p2.y - p1.y, p2.x - p1.x)
  const pt = bezierPoint(t, p)

  ctx.save()
  ctx.translate(pt.x, pt.y)
  ctx.rotate(angle)
  ctx.beginPath()
  ctx.moveTo(size, 0)
  ctx.lineTo(-size, -size * 0.6)
  ctx.lineTo(-size, size * 0.6)
  ctx.closePath()
  ctx.fillStyle = color
  ctx.fill()
  ctx.restore()
}

export default function AgentFlowGraph({ logs, isProcessing }: Props) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const [elapsed, setElapsed] = useState(0)

  const getStatus = (id: string) => {
    const agentLogs = logs.filter(l => l.agent === id)
    if (agentLogs.some(l => l.status === 'success')) return 'completed'
    if (agentLogs.some(l => l.status === 'running')) return 'running'
    return 'pending'
  }

  // Timer
  useEffect(() => {
    if (!isProcessing) return
    const t = setInterval(() => setElapsed(e => e + 1), 1000)
    return () => clearInterval(t)
  }, [isProcessing])

  useEffect(() => {
    if (!isProcessing && logs.length === 0) setElapsed(0)
  }, [isProcessing, logs.length])

  // Canvas animation: particles, rotating arcs, directional dashed lines
  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const resize = () => {
      canvas.width = canvas.offsetWidth * 2
      canvas.height = canvas.offsetHeight * 2
      ctx.scale(2, 2)
    }
    resize()

    const particles: { progress: number; conn: number; speed: number; size: number }[] = []
    let animationId: number
    let tick = 0

    const animate = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height)
      const w = canvas.offsetWidth
      const h = canvas.offsetHeight
      tick++

      // --- Draw bezier connections with flowing dashed lines ---
      connections.forEach((conn) => {
        const from = agents[conn.from]
        const to = agents[conn.to]
        const bp = getBezier(from, to, w, h)
        const fromStatus = getStatus(from.id)
        const toStatus = getStatus(to.id)
        const isActive = fromStatus === 'completed' || fromStatus === 'running'

        // Base dashed line (always visible)
        ctx.beginPath()
        ctx.setLineDash([6, 4])
        ctx.moveTo(bp.x1, bp.y1)
        ctx.bezierCurveTo(bp.cx1, bp.cy1, bp.cx2, bp.cy2, bp.x2, bp.y2)
        ctx.strokeStyle = isActive ? `${from.color}30` : 'rgba(48,54,61,0.4)'
        ctx.lineWidth = 1.5
        ctx.stroke()
        ctx.setLineDash([])

        // Active: animated flowing dashed line with dash offset
        if (isActive) {
          ctx.beginPath()
          ctx.setLineDash([8, 6])
          ctx.lineDashOffset = -tick * 0.8
          ctx.moveTo(bp.x1, bp.y1)
          ctx.bezierCurveTo(bp.cx1, bp.cy1, bp.cx2, bp.cy2, bp.x2, bp.y2)
          const grad = ctx.createLinearGradient(bp.x1, bp.y1, bp.x2, bp.y2)
          grad.addColorStop(0, from.color + '80')
          grad.addColorStop(1, to.color + '80')
          ctx.strokeStyle = grad
          ctx.lineWidth = 2
          ctx.shadowColor = from.color
          ctx.shadowBlur = 6
          ctx.stroke()
          ctx.shadowBlur = 0
          ctx.setLineDash([])
          ctx.lineDashOffset = 0

          // Directional arrow at midpoint
          const arrowColor = hexToRgba(from.color, 0.8)
          drawArrowHead(ctx, bp, 0.5, arrowColor, 5)

          // Second arrow at 75% if target is also active
          if (toStatus === 'running' || toStatus === 'completed') {
            drawArrowHead(ctx, bp, 0.75, hexToRgba(to.color, 0.6), 4)
          }
        }
      })

      // --- Draw rotating progress arcs on running nodes ---
      agents.forEach((agent) => {
        const status = getStatus(agent.id)
        if (status === 'running') {
          const cx = (agent.x / 100) * w
          const cy = (agent.y / 100) * h
          drawRotatingArc(ctx, cx, cy, 38, agent.color, tick)
        }
      })

      // --- Spawn particles ---
      if (isProcessing && tick % 3 === 0) {
        connections.forEach((conn, idx) => {
          const status = getStatus(agents[conn.from].id)
          if (status === 'running' || status === 'completed') {
            if (Math.random() > 0.5) {
              particles.push({
                progress: 0,
                conn: idx,
                speed: 0.008 + Math.random() * 0.008,
                size: 2 + Math.random() * 2,
              })
            }
          }
        })
      }

      // --- Update and draw particles ---
      for (let i = particles.length - 1; i >= 0; i--) {
        const p = particles[i]
        p.progress += p.speed
        if (p.progress >= 1) { particles.splice(i, 1); continue }

        const conn = connections[p.conn]
        const from = agents[conn.from]
        const to = agents[conn.to]
        const bp = getBezier(from, to, w, h)
        const pt = bezierPoint(p.progress, bp)

        // Particle trail
        const tailLen = 5
        for (let t = 1; t <= tailLen; t++) {
          const tp = Math.max(0, p.progress - t * 0.02)
          const tpt = bezierPoint(tp, bp)
          const alpha = (1 - t / tailLen) * 0.3
          ctx.fillStyle = hexToRgba(from.color, alpha)
          ctx.beginPath()
          ctx.arc(tpt.x, tpt.y, p.size * (1 - t / tailLen * 0.5), 0, Math.PI * 2)
          ctx.fill()
        }

        // Particle glow
        const glow = ctx.createRadialGradient(pt.x, pt.y, 0, pt.x, pt.y, 16)
        glow.addColorStop(0, from.color + '80')
        glow.addColorStop(0.4, from.color + '20')
        glow.addColorStop(1, 'transparent')
        ctx.fillStyle = glow
        ctx.beginPath()
        ctx.arc(pt.x, pt.y, 16, 0, Math.PI * 2)
        ctx.fill()

        // Particle core
        ctx.fillStyle = '#fff'
        ctx.beginPath()
        ctx.arc(pt.x, pt.y, p.size, 0, Math.PI * 2)
        ctx.fill()
      }

      animationId = requestAnimationFrame(animate)
    }

    animate()
    window.addEventListener('resize', resize)
    return () => {
      cancelAnimationFrame(animationId)
      window.removeEventListener('resize', resize)
    }
  }, [logs, isProcessing])

  const formatTime = (s: number) => `${Math.floor(s / 60).toString().padStart(2, '0')}:${(s % 60).toString().padStart(2, '0')}`

  return (
    <div className="relative w-full h-full min-h-[220px]">
      {/* 背景：点阵网格 */}
      <div className="absolute inset-0 opacity-[0.07]" style={{
        backgroundImage: 'radial-gradient(circle, #00F0E0 1px, transparent 1px)',
        backgroundSize: '32px 32px',
      }} />

      {/* 扫描线动效 */}
      {isProcessing && (
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div
            className="absolute left-0 right-0 h-[2px]"
            style={{
              background: 'linear-gradient(90deg, transparent, #00F0E040, transparent)',
              animation: 'cyber-scan 3s linear infinite',
            }}
          />
        </div>
      )}

      {/* 粒子画布 */}
      <canvas ref={canvasRef} className="absolute inset-0 w-full h-full" />

      {/* Agent 节点 */}
      {agents.map((agent) => {
        const status = getStatus(agent.id)
        const Icon = agent.icon
        return (
          <div
            key={agent.id}
            className="absolute transform -translate-x-1/2 -translate-y-1/2 flex flex-col items-center"
            style={{ left: `${agent.x}%`, top: `${agent.y}%` }}
          >
            {/* 节点光晕 */}
            <div
              className={cn(
                'absolute w-16 h-16 rounded-full blur-xl transition-opacity duration-500',
                status === 'running' ? 'opacity-50' : status === 'completed' ? 'opacity-30' : 'opacity-0'
              )}
              style={{ backgroundColor: agent.color }}
            />

            {/* 节点主体 */}
            <div
              className={cn(
                'relative w-14 h-14 rounded-xl flex items-center justify-center transition-all duration-300',
                'border backdrop-blur-sm',
                status === 'running' && 'scale-110',
              )}
              style={{
                borderColor: status !== 'pending' ? agent.color : '#30363D',
                background: status === 'running'
                  ? `linear-gradient(135deg, ${agent.color}15, ${agent.color}08)`
                  : 'rgba(13,17,23,0.7)',
                boxShadow: status !== 'pending' ? `0 0 20px ${agent.color}20` : 'none',
              }}
            >
              {status === 'completed' ? (
                <CheckCircle className="w-6 h-6" style={{ color: agent.color }} />
              ) : (
                <Icon className="w-6 h-6" style={{ color: status === 'running' ? agent.color : '#484F58' }} />
              )}
            </div>

            {/* 标签 */}
            <div className="mt-2.5 text-center">
              <div className={cn(
                'text-xs font-semibold tracking-wide transition-colors',
                status !== 'pending' ? 'text-[#E6EDF3]' : 'text-[#484F58]'
              )}>
                {agent.label}
              </div>
              <div className={cn(
                'text-[10px] mt-0.5 tracking-wider transition-colors',
                status !== 'pending' ? 'text-[#8B949E]' : 'text-[#30363D]'
              )}>
                {agent.desc}
              </div>
              {status === 'running' && (
                <div
                  className="text-[10px] mt-1 font-mono"
                  style={{ color: agent.color, animation: 'cyber-cursor 1s step-end infinite' }}
                >
                  processing...
                </div>
              )}
            </div>
          </div>
        )
      })}

      {/* 中心状态指示器 */}
      <div className="absolute left-1/2 top-[75%] transform -translate-x-1/2 flex items-center gap-2">
        {isProcessing ? (
          <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-[#00F0E0]/10 border border-[#00F0E0]/20">
            <div className="w-2 h-2 rounded-full bg-[#00F0E0] animate-ping" />
            <span className="text-[11px] text-[#00F0E0] font-mono tracking-wider">
              ANALYZING · {formatTime(elapsed)}
            </span>
          </div>
        ) : logs.length > 0 ? (
          <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-[#22C55E]/10 border border-[#22C55E]/20">
            <CheckCircle className="w-3.5 h-3.5 text-[#22C55E]" />
            <span className="text-[11px] text-[#22C55E] font-mono tracking-wider">
              COMPLETE · {formatTime(elapsed)}
            </span>
          </div>
        ) : (
          <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-[#30363D]/30 border border-[#30363D]/50">
            <Circle className="w-3 h-3 text-[#484F58]" />
            <span className="text-[11px] text-[#484F58] font-mono tracking-wider">STANDBY</span>
          </div>
        )}
      </div>
    </div>
  )
}