import { useState, useEffect } from 'react'
import {
  ShieldAlert,
  Rss,
  FileText,
  TrendingUp,
  AlertTriangle,
  ArrowUp,
  ArrowDown,
  RefreshCw,
  Loader2,
} from 'lucide-react'
import { cn } from '@/utils'
import EventTrendChart from './components/EventTrendChart'
import SeverityDistribution from './components/SeverityDistribution'
import RecentEvents from './components/RecentEvents'
import SubscriptionStatus from './components/SubscriptionStatus'
import ActionQueue from './components/ActionQueue'
import SecurityFunnel from './components/SecurityFunnel'
import { eventService } from '@/services/event'
import { subscriptionService } from '@/services/subscription'
import { reportService } from '@/services/report'

// 迷你折线图组件
function MiniChart({ data, color }: { data: number[]; color: string }) {
  const max = Math.max(...data)
  const min = Math.min(...data)
  const range = max - min || 1
  const height = 40
  const width = 80
  const points = data.map((v, i) => {
    const x = (i / (data.length - 1)) * width
    const y = height - ((v - min) / range) * height
    return `${x},${y}`
  }).join(' ')

  return (
    <svg width={width} height={height} className="overflow-visible">
      <defs>
        <linearGradient id={`gradient-${color}`} x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor={color} stopOpacity="0.3" />
          <stop offset="100%" stopColor={color} stopOpacity="0" />
        </linearGradient>
      </defs>
      <polygon
        points={`0,${height} ${points} ${width},${height}`}
        fill={`url(#gradient-${color})`}
      />
      <polyline
        points={points}
        fill="none"
        stroke={color}
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

interface DashboardStats {
  eventCount: number
  subscriptionCount: number
  reportCount: number
  todayCount: number
  criticalCount: number
}

export default function Dashboard() {
  const [loading, setLoading] = useState(true)
  const [refreshing, setRefreshing] = useState(false)
  const [stats, setStats] = useState<DashboardStats>({
    eventCount: 0,
    subscriptionCount: 0,
    reportCount: 0,
    todayCount: 0,
    criticalCount: 0,
  })

  const fetchStats = async () => {
    try {
      const [eventsRes, subscriptionsRes, reportsRes, eventStats] = await Promise.all([
        eventService.list({ page: 1, size: 1 }),
        subscriptionService.list(1, 1),
        reportService.list(1, 1),
        eventService.getStats(),
      ])

      setStats({
        eventCount: eventsRes.total || 0,
        subscriptionCount: subscriptionsRes.total || 0,
        reportCount: reportsRes.total || 0,
        todayCount: (eventStats as any).today_count || 0,
        criticalCount: (eventStats as any).critical_count || 0,
      })
    } catch (error) {
      console.error('获取统计数据失败:', error)
    } finally {
      setLoading(false)
      setRefreshing(false)
    }
  }

  useEffect(() => {
    fetchStats()
  }, [])

  const handleRefresh = () => {
    setRefreshing(true)
    fetchStats()
  }

  const statsConfig = [
    {
      label: '安全事件',
      value: stats.eventCount.toLocaleString(),
      change: 0,
      icon: ShieldAlert,
      color: '#ef4444',
      bgColor: 'bg-red-50',
      textColor: 'text-red-600',
      data: [45, 52, 38, 65, 48, 72, stats.eventCount % 100 || 58],
    },
    {
      label: '活跃订阅',
      value: stats.subscriptionCount.toString(),
      change: 0,
      icon: Rss,
      color: '#3b82f6',
      bgColor: 'bg-blue-50',
      textColor: 'text-blue-600',
      data: [18, 20, 19, 22, 21, 23, stats.subscriptionCount || 24],
    },
    {
      label: '分析报告',
      value: stats.reportCount.toString(),
      change: 0,
      icon: FileText,
      color: '#10b981',
      bgColor: 'bg-emerald-50',
      textColor: 'text-emerald-600',
      data: [120, 135, 128, 142, 138, 150, stats.reportCount || 156],
    },
    {
      label: '今日新增',
      value: stats.todayCount.toString(),
      change: 0,
      icon: TrendingUp,
      color: '#f59e0b',
      bgColor: 'bg-amber-50',
      textColor: 'text-amber-600',
      data: [62, 55, 68, 52, 58, 50, stats.todayCount || 47],
    },
  ]

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-primary-500" />
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-slate-900 tracking-tight">仪表盘</h1>
          <p className="text-sm text-slate-500 mt-1">安全态势总览与实时监控</p>
        </div>
        <button
          onClick={handleRefresh}
          className="btn-default"
          disabled={refreshing}
        >
          <RefreshCw className={cn('w-4 h-4', refreshing && 'animate-spin')} />
          刷新
        </button>
      </div>

      {/* Security Funnel */}
      <SecurityFunnel
        collected={stats.eventCount}
        deduplicated={Math.floor(stats.eventCount * 0.3)}
        critical={stats.criticalCount}
      />

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {statsConfig.map((stat) => (
          <div key={stat.label} className="card card-body">
            <div className="flex items-start justify-between mb-4">
              <div className={cn('w-10 h-10 rounded-xl flex items-center justify-center', stat.bgColor)}>
                <stat.icon className={cn('w-5 h-5', stat.textColor)} />
              </div>
              {stat.change !== 0 && (
                <div
                  className={cn(
                    'flex items-center gap-1 text-xs font-medium',
                    stat.change >= 0 ? 'text-emerald-600' : 'text-red-600'
                  )}
                >
                  {stat.change >= 0 ? (
                    <ArrowUp className="w-3 h-3" />
                  ) : (
                    <ArrowDown className="w-3 h-3" />
                  )}
                  {Math.abs(stat.change)}%
                </div>
              )}
            </div>

            <div className="flex items-end justify-between">
              <div>
                <p className="text-3xl font-bold text-slate-900 tracking-tight">{stat.value}</p>
                <p className="text-sm text-slate-500 mt-1">{stat.label}</p>
              </div>
              <MiniChart data={stat.data} color={stat.color} />
            </div>
          </div>
        ))}
      </div>

      {/* Action Queue - 待办任务流 */}
      <ActionQueue criticalCount={stats.criticalCount} pendingReports={stats.reportCount > 0 ? 2 : 0} />

      {/* Alert Banner */}
      {stats.criticalCount > 0 && (
        <div className="flex items-start gap-4 p-5 bg-red-50 rounded-2xl border border-red-100">
          <div className="w-10 h-10 rounded-xl bg-red-100 flex items-center justify-center flex-shrink-0">
            <AlertTriangle className="w-5 h-5 text-red-600" />
          </div>
          <div className="flex-1">
            <p className="font-semibold text-red-900">发现 {stats.criticalCount} 个高危漏洞需要关注</p>
            <p className="text-sm text-red-700 mt-1">
              建议立即处理
            </p>
          </div>
          <a href="/events?severity=critical" className="btn-sm bg-red-600 hover:bg-red-700 text-white">
            立即查看
          </a>
        </div>
      )}

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 card">
          <div className="card-header flex items-center justify-between">
            <span>事件趋势</span>
            <span className="text-xs text-slate-400 font-normal">最近 7 天</span>
          </div>
          <div className="card-body">
            <EventTrendChart />
          </div>
        </div>

        <div className="card">
          <div className="card-header">严重级别分布</div>
          <div className="card-body">
            <SeverityDistribution />
          </div>
        </div>
      </div>

      {/* Bottom Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <div className="card-header flex items-center justify-between">
            <span>最新事件</span>
            <a href="/events" className="text-sm text-slate-500 hover:text-slate-900 font-normal transition-colors">
              查看全部 →
            </a>
          </div>
          <div className="p-0">
            <RecentEvents />
          </div>
        </div>

        <div className="card">
          <div className="card-header flex items-center justify-between">
            <span>订阅状态</span>
            <a href="/subscriptions" className="text-sm text-slate-500 hover:text-slate-900 font-normal transition-colors">
              管理订阅 →
            </a>
          </div>
          <div className="p-0">
            <SubscriptionStatus />
          </div>
        </div>
      </div>
    </div>
  )
}
