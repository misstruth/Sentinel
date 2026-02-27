import { AlertTriangle, FileText, Shield, ChevronRight, Clock } from 'lucide-react'
import { cn } from '@/utils'

interface ActionItem {
  id: number
  type: 'event' | 'report'
  title: string
  subtitle?: string
  priority: 'high' | 'medium' | 'low'
  action: string
  link: string
  deadline?: string
}

interface Props {
  criticalCount: number
  pendingReports: number
}

export default function ActionQueue({ criticalCount, pendingReports }: Props) {
  const actions: ActionItem[] = []

  if (criticalCount > 0) {
    actions.push({
      id: 1,
      type: 'event',
      title: `${criticalCount} 个高危事件待处理`,
      subtitle: 'CVE-2024-XXXX 等',
      priority: 'high',
      action: '立即处置',
      link: '/events?severity=critical',
      deadline: '2h',
    })
  }

  if (pendingReports > 0) {
    actions.push({
      id: 2,
      type: 'report',
      title: `${pendingReports} 份报告待审阅`,
      priority: 'medium',
      action: '查看报告',
      link: '/reports',
    })
  }

  actions.push({
    id: 3,
    type: 'event',
    title: '运行Agent分析流水线',
    priority: 'low',
    action: '启动分析',
    link: '/events/analysis',
  })

  if (actions.length === 0) return null

  return (
    <div className="card">
      <div className="card-header flex items-center gap-2">
        <Shield className="w-4 h-4 text-primary-500" />
        待办任务
        <span className="ml-auto text-xs font-normal text-slate-400">{actions.length} 项</span>
      </div>
      <div className="divide-y divide-slate-100">
        {actions.map((item) => (
          <a
            key={item.id}
            href={item.link}
            className="flex items-center gap-3 p-4 hover:bg-slate-50 transition-colors"
          >
            <div className={cn(
              'w-8 h-8 rounded-lg flex items-center justify-center',
              item.priority === 'high' ? 'bg-red-100' : item.priority === 'medium' ? 'bg-amber-100' : 'bg-blue-100'
            )}>
              {item.type === 'event' ? (
                <AlertTriangle className={cn('w-4 h-4', item.priority === 'high' ? 'text-red-600' : 'text-amber-600')} />
              ) : (
                <FileText className="w-4 h-4 text-amber-600" />
              )}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-slate-900 truncate">{item.title}</p>
              {item.subtitle && <p className="text-xs text-slate-500">{item.subtitle}</p>}
            </div>
            {item.deadline && (
              <div className="flex items-center gap-1 text-xs text-red-600">
                <Clock className="w-3 h-3" />
                {item.deadline}
              </div>
            )}
            <span className={cn(
              'text-xs font-medium px-2 py-1 rounded',
              item.priority === 'high' ? 'bg-red-100 text-red-700' : item.priority === 'medium' ? 'bg-amber-100 text-amber-700' : 'bg-blue-100 text-blue-700'
            )}>
              {item.action}
            </span>
            <ChevronRight className="w-4 h-4 text-slate-400" />
          </a>
        ))}
      </div>
    </div>
  )
}
