import { NavLink, useLocation } from 'react-router-dom'
import {
  LayoutDashboard,
  Rss,
  ShieldAlert,
  FileText,
  Settings,
  ChevronLeft,
  Shield,
  MessageSquare,
  FileSearch,
  Cpu,
} from 'lucide-react'
import { useAppStore } from '@/stores/app'
import { cn } from '@/utils'

const navItems = [
  { path: '/dashboard', icon: LayoutDashboard, label: '仪表盘' },
  { path: '/subscriptions', icon: Rss, label: '订阅管理' },
  { path: '/events', icon: ShieldAlert, label: '安全事件' },
  { path: '/events/analysis', icon: Cpu, label: 'Agent分析' },
  { path: '/reports', icon: FileText, label: '分析报告' },
  { path: '/logs', icon: FileSearch, label: '抓取日志' },
  { path: '/chat', icon: MessageSquare, label: 'AI 助手' },
  { path: '/settings', icon: Settings, label: '系统设置' },
]

export default function Sidebar() {
  const { sidebarCollapsed, toggleSidebar } = useAppStore()
  const location = useLocation()

  return (
    <aside
      className={cn(
        'fixed left-0 top-0 h-screen bg-white z-50 flex flex-col transition-all duration-300 border-r border-slate-200',
        sidebarCollapsed ? 'w-[72px]' : 'w-64'
      )}
    >
      {/* Logo */}
      <div className="h-16 flex items-center px-5 border-b border-slate-100">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-xl bg-slate-900 flex items-center justify-center flex-shrink-0">
            <Shield className="w-5 h-5 text-white" />
          </div>
          {!sidebarCollapsed && (
            <div className="overflow-hidden">
              <h1 className="text-base font-semibold text-slate-900 whitespace-nowrap tracking-tight">
                Sentinel
              </h1>
            </div>
          )}
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 py-4 px-3 overflow-y-auto">
        <div className="space-y-1">
          {navItems.map((item) => {
            const isActive = location.pathname === item.path ||
              (item.path !== '/dashboard' && location.pathname.startsWith(item.path))
            return (
              <NavLink
                key={item.path}
                to={item.path}
                className={cn(
                  'flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-150',
                  isActive
                    ? 'bg-slate-900 text-white'
                    : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100'
                )}
                title={sidebarCollapsed ? item.label : undefined}
              >
                <item.icon className="w-5 h-5 flex-shrink-0" />
                {!sidebarCollapsed && <span>{item.label}</span>}
              </NavLink>
            )
          })}
        </div>
      </nav>

      {/* Collapse Button */}
      <div className="p-3 border-t border-slate-100">
        <button
          onClick={toggleSidebar}
          className="w-full flex items-center justify-center gap-2 px-3 py-2.5 rounded-xl text-sm text-slate-500 hover:text-slate-900 hover:bg-slate-100 transition-all duration-150"
        >
          <ChevronLeft
            className={cn(
              'w-4 h-4 transition-transform duration-300',
              sidebarCollapsed && 'rotate-180'
            )}
          />
          {!sidebarCollapsed && <span>收起</span>}
        </button>
      </div>
    </aside>
  )
}
