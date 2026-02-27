import { useState } from 'react'
import {
  Settings as SettingsIcon,
  Bell,
  Mail,
  Webhook,
  Key,
  Database,
  Shield,
  Save,
  TestTube,
  Loader2,
  Check,
  Trash2,
  Download,
  Upload,
} from 'lucide-react'
import { cn } from '@/utils'
import toast from 'react-hot-toast'

const tabs = [
  { id: 'general', label: '通用设置', icon: SettingsIcon },
  { id: 'notifications', label: '通知配置', icon: Bell },
  { id: 'api', label: 'API 密钥', icon: Key },
  { id: 'database', label: '数据管理', icon: Database },
]

export default function Settings() {
  const [activeTab, setActiveTab] = useState('general')
  const [isSaving, setIsSaving] = useState(false)

  // General settings
  const [generalSettings, setGeneralSettings] = useState({
    siteName: 'Security Sentinel',
    defaultFetchInterval: 60,
    maxEventsPerPage: 50,
    autoMarkRead: true,
    enableAnalytics: true,
  })

  // Notification settings
  const [notifySettings, setNotifySettings] = useState({
    emailEnabled: true,
    emailAddresses: 'admin@example.com',
    emailOnCritical: true,
    emailOnHigh: true,
    emailOnMedium: false,
    webhookEnabled: false,
    webhookUrl: '',
    webhookSecret: '',
  })

  // API settings
  const [apiSettings, setApiSettings] = useState({
    githubToken: '',
    nvdApiKey: '',
    openaiApiKey: '',
  })

  const handleSave = async () => {
    setIsSaving(true)
    await new Promise((resolve) => setTimeout(resolve, 1000))
    setIsSaving(false)
    toast.success('设置已保存')
  }

  const handleTestEmail = async () => {
    toast.promise(
      new Promise((resolve) => setTimeout(resolve, 2000)),
      {
        loading: '发送测试邮件...',
        success: '测试邮件已发送',
        error: '发送失败',
      }
    )
  }

  const handleTestWebhook = async () => {
    toast.promise(
      new Promise((resolve) => setTimeout(resolve, 1500)),
      {
        loading: '测试 Webhook...',
        success: 'Webhook 连接成功',
        error: '连接失败',
      }
    )
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-medium text-gray-100">系统设置</h1>
          <p className="text-sm text-gray-500 mt-1">配置系统参数和集成服务</p>
        </div>
        <button onClick={handleSave} disabled={isSaving} className="btn-primary">
          {isSaving ? (
            <>
              <Loader2 className="w-4 h-4 animate-spin" />
              保存中...
            </>
          ) : (
            <>
              <Save className="w-4 h-4" />
              保存设置
            </>
          )}
        </button>
      </div>

      <div className="flex gap-6">
        {/* Sidebar */}
        <div className="w-48 flex-shrink-0">
          <nav className="card p-2 space-y-1">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={cn(
                  'w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg transition-colors text-left text-sm',
                  activeTab === tab.id
                    ? 'bg-primary-500/10 text-primary-400'
                    : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800'
                )}
              >
                <tab.icon className="w-4 h-4" />
                <span>{tab.label}</span>
              </button>
            ))}
          </nav>
        </div>

        {/* Content */}
        <div className="flex-1">
          <div className="card">
            {/* General Settings */}
            {activeTab === 'general' && (
              <div className="card-body space-y-6">
                <div>
                  <h3 className="text-base font-medium text-gray-100 mb-4">通用设置</h3>
                  <div className="space-y-4">
                    <div className="form-item">
                      <label className="label">系统名称</label>
                      <input
                        type="text"
                        value={generalSettings.siteName}
                        onChange={(e) =>
                          setGeneralSettings({ ...generalSettings, siteName: e.target.value })
                        }
                        className="input"
                      />
                    </div>

                    <div className="form-item">
                      <label className="label">默认抓取间隔</label>
                      <select
                        value={generalSettings.defaultFetchInterval}
                        onChange={(e) =>
                          setGeneralSettings({
                            ...generalSettings,
                            defaultFetchInterval: Number(e.target.value),
                          })
                        }
                        className="select"
                      >
                        <option value={15}>15 分钟</option>
                        <option value={30}>30 分钟</option>
                        <option value={60}>1 小时</option>
                        <option value={180}>3 小时</option>
                        <option value={360}>6 小时</option>
                      </select>
                    </div>

                    <div className="form-item">
                      <label className="label">每页显示事件数</label>
                      <select
                        value={generalSettings.maxEventsPerPage}
                        onChange={(e) =>
                          setGeneralSettings({
                            ...generalSettings,
                            maxEventsPerPage: Number(e.target.value),
                          })
                        }
                        className="select"
                      >
                        <option value={20}>20</option>
                        <option value={50}>50</option>
                        <option value={100}>100</option>
                      </select>
                    </div>
                  </div>
                </div>

                <div className="border-t border-gray-800 pt-6">
                  <h3 className="text-base font-medium text-gray-100 mb-4">功能开关</h3>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between p-3 rounded-lg bg-gray-800/50">
                      <div>
                        <p className="text-sm font-medium text-gray-200">自动标记已读</p>
                        <p className="text-xs text-gray-500 mt-0.5">查看事件详情后自动标记为已读</p>
                      </div>
                      <button
                        onClick={() =>
                          setGeneralSettings({
                            ...generalSettings,
                            autoMarkRead: !generalSettings.autoMarkRead,
                          })
                        }
                        className={cn(
                          'w-10 h-5 rounded-full transition-colors relative',
                          generalSettings.autoMarkRead ? 'bg-primary-500' : 'bg-gray-600'
                        )}
                      >
                        <span
                          className={cn(
                            'absolute top-0.5 w-4 h-4 rounded-full bg-white transition-transform',
                            generalSettings.autoMarkRead ? 'left-5' : 'left-0.5'
                          )}
                        />
                      </button>
                    </div>

                    <div className="flex items-center justify-between p-3 rounded-lg bg-gray-800/50">
                      <div>
                        <p className="text-sm font-medium text-gray-200">启用数据分析</p>
                        <p className="text-xs text-gray-500 mt-0.5">收集使用数据以改进系统</p>
                      </div>
                      <button
                        onClick={() =>
                          setGeneralSettings({
                            ...generalSettings,
                            enableAnalytics: !generalSettings.enableAnalytics,
                          })
                        }
                        className={cn(
                          'w-10 h-5 rounded-full transition-colors relative',
                          generalSettings.enableAnalytics ? 'bg-primary-500' : 'bg-gray-600'
                        )}
                      >
                        <span
                          className={cn(
                            'absolute top-0.5 w-4 h-4 rounded-full bg-white transition-transform',
                            generalSettings.enableAnalytics ? 'left-5' : 'left-0.5'
                          )}
                        />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Notification Settings */}
            {activeTab === 'notifications' && (
              <div className="card-body space-y-6">
                {/* Email */}
                <div>
                  <div className="flex items-center gap-2 mb-4">
                    <Mail className="w-4 h-4 text-primary-400" />
                    <h3 className="text-base font-medium text-gray-100">邮件通知</h3>
                  </div>

                  <div className="space-y-4">
                    <div className="flex items-center justify-between p-3 rounded-lg bg-gray-800/50">
                      <div>
                        <p className="text-sm font-medium text-gray-200">启用邮件通知</p>
                        <p className="text-xs text-gray-500 mt-0.5">发现新事件时发送邮件通知</p>
                      </div>
                      <button
                        onClick={() =>
                          setNotifySettings({
                            ...notifySettings,
                            emailEnabled: !notifySettings.emailEnabled,
                          })
                        }
                        className={cn(
                          'w-10 h-5 rounded-full transition-colors relative',
                          notifySettings.emailEnabled ? 'bg-primary-500' : 'bg-gray-600'
                        )}
                      >
                        <span
                          className={cn(
                            'absolute top-0.5 w-4 h-4 rounded-full bg-white transition-transform',
                            notifySettings.emailEnabled ? 'left-5' : 'left-0.5'
                          )}
                        />
                      </button>
                    </div>

                    {notifySettings.emailEnabled && (
                      <>
                        <div className="form-item">
                          <label className="label">收件邮箱</label>
                          <input
                            type="text"
                            value={notifySettings.emailAddresses}
                            onChange={(e) =>
                              setNotifySettings({
                                ...notifySettings,
                                emailAddresses: e.target.value,
                              })
                            }
                            placeholder="多个邮箱用逗号分隔"
                            className="input"
                          />
                        </div>

                        <div className="form-item">
                          <label className="label">通知级别</label>
                          <div className="flex gap-2">
                            <button
                              onClick={() =>
                                setNotifySettings({
                                  ...notifySettings,
                                  emailOnCritical: !notifySettings.emailOnCritical,
                                })
                              }
                              className={cn(
                                'btn-sm',
                                notifySettings.emailOnCritical
                                  ? 'bg-danger-500/20 border-danger-500/50 text-danger-400'
                                  : 'btn-default'
                              )}
                            >
                              严重
                            </button>
                            <button
                              onClick={() =>
                                setNotifySettings({
                                  ...notifySettings,
                                  emailOnHigh: !notifySettings.emailOnHigh,
                                })
                              }
                              className={cn(
                                'btn-sm',
                                notifySettings.emailOnHigh
                                  ? 'bg-warning-500/20 border-warning-500/50 text-warning-500'
                                  : 'btn-default'
                              )}
                            >
                              高危
                            </button>
                            <button
                              onClick={() =>
                                setNotifySettings({
                                  ...notifySettings,
                                  emailOnMedium: !notifySettings.emailOnMedium,
                                })
                              }
                              className={cn(
                                'btn-sm',
                                notifySettings.emailOnMedium
                                  ? 'bg-yellow-500/20 border-yellow-500/50 text-yellow-500'
                                  : 'btn-default'
                              )}
                            >
                              中危
                            </button>
                          </div>
                        </div>

                        <button onClick={handleTestEmail} className="btn-default">
                          <TestTube className="w-4 h-4" />
                          发送测试邮件
                        </button>
                      </>
                    )}
                  </div>
                </div>

                {/* Webhook */}
                <div className="border-t border-gray-800 pt-6">
                  <div className="flex items-center gap-2 mb-4">
                    <Webhook className="w-4 h-4 text-primary-400" />
                    <h3 className="text-base font-medium text-gray-100">Webhook 通知</h3>
                  </div>

                  <div className="space-y-4">
                    <div className="flex items-center justify-between p-3 rounded-lg bg-gray-800/50">
                      <div>
                        <p className="text-sm font-medium text-gray-200">启用 Webhook</p>
                        <p className="text-xs text-gray-500 mt-0.5">通过 Webhook 推送事件通知</p>
                      </div>
                      <button
                        onClick={() =>
                          setNotifySettings({
                            ...notifySettings,
                            webhookEnabled: !notifySettings.webhookEnabled,
                          })
                        }
                        className={cn(
                          'w-10 h-5 rounded-full transition-colors relative',
                          notifySettings.webhookEnabled ? 'bg-primary-500' : 'bg-gray-600'
                        )}
                      >
                        <span
                          className={cn(
                            'absolute top-0.5 w-4 h-4 rounded-full bg-white transition-transform',
                            notifySettings.webhookEnabled ? 'left-5' : 'left-0.5'
                          )}
                        />
                      </button>
                    </div>

                    {notifySettings.webhookEnabled && (
                      <>
                        <div className="form-item">
                          <label className="label">Webhook URL</label>
                          <input
                            type="url"
                            value={notifySettings.webhookUrl}
                            onChange={(e) =>
                              setNotifySettings({
                                ...notifySettings,
                                webhookUrl: e.target.value,
                              })
                            }
                            placeholder="https://your-server.com/webhook"
                            className="input"
                          />
                        </div>

                        <div className="form-item">
                          <label className="label">Secret（可选）</label>
                          <input
                            type="password"
                            value={notifySettings.webhookSecret}
                            onChange={(e) =>
                              setNotifySettings({
                                ...notifySettings,
                                webhookSecret: e.target.value,
                              })
                            }
                            placeholder="用于签名验证"
                            className="input"
                          />
                        </div>

                        <button onClick={handleTestWebhook} className="btn-default">
                          <TestTube className="w-4 h-4" />
                          测试连接
                        </button>
                      </>
                    )}
                  </div>
                </div>
              </div>
            )}

            {/* API Settings */}
            {activeTab === 'api' && (
              <div className="card-body space-y-6">
                <div>
                  <h3 className="text-base font-medium text-gray-100 mb-2">API 密钥配置</h3>
                  <p className="text-sm text-gray-500 mb-4">
                    配置第三方服务的 API 密钥以启用相关功能
                  </p>

                  <div className="space-y-4">
                    <div className="form-item">
                      <label className="label">GitHub Token</label>
                      <input
                        type="password"
                        value={apiSettings.githubToken}
                        onChange={(e) =>
                          setApiSettings({ ...apiSettings, githubToken: e.target.value })
                        }
                        placeholder="ghp_xxxxxxxxxxxx"
                        className="input font-mono"
                      />
                      <p className="text-xs text-gray-500 mt-1">
                        用于访问 GitHub API，监控仓库安全公告
                      </p>
                    </div>

                    <div className="form-item">
                      <label className="label">NVD API Key</label>
                      <input
                        type="password"
                        value={apiSettings.nvdApiKey}
                        onChange={(e) =>
                          setApiSettings({ ...apiSettings, nvdApiKey: e.target.value })
                        }
                        placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
                        className="input font-mono"
                      />
                      <p className="text-xs text-gray-500 mt-1">
                        用于访问 NVD 漏洞数据库，提高请求频率限制
                      </p>
                    </div>

                    <div className="form-item">
                      <label className="label">OpenAI API Key</label>
                      <input
                        type="password"
                        value={apiSettings.openaiApiKey}
                        onChange={(e) =>
                          setApiSettings({ ...apiSettings, openaiApiKey: e.target.value })
                        }
                        placeholder="sk-xxxxxxxxxxxx"
                        className="input font-mono"
                      />
                      <p className="text-xs text-gray-500 mt-1">
                        用于 AI 分析和报告生成功能
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Database Settings */}
            {activeTab === 'database' && (
              <div className="card-body space-y-6">
                <div>
                  <h3 className="text-base font-medium text-gray-100 mb-4">数据库状态</h3>

                  <div className="grid grid-cols-2 gap-3">
                    <div className="p-3 rounded-lg bg-gray-800/50">
                      <p className="text-xs text-gray-500">数据库类型</p>
                      <p className="text-sm font-medium text-gray-200 mt-1">MySQL 8.0</p>
                    </div>
                    <div className="p-3 rounded-lg bg-gray-800/50">
                      <p className="text-xs text-gray-500">连接状态</p>
                      <p className="text-sm font-medium text-success-500 mt-1 flex items-center gap-1.5">
                        <Check className="w-3.5 h-3.5" />
                        已连接
                      </p>
                    </div>
                    <div className="p-3 rounded-lg bg-gray-800/50">
                      <p className="text-xs text-gray-500">事件总数</p>
                      <p className="text-sm font-medium text-gray-200 mt-1">1,284</p>
                    </div>
                    <div className="p-3 rounded-lg bg-gray-800/50">
                      <p className="text-xs text-gray-500">报告总数</p>
                      <p className="text-sm font-medium text-gray-200 mt-1">156</p>
                    </div>
                  </div>
                </div>

                <div className="border-t border-gray-800 pt-6">
                  <h3 className="text-base font-medium text-gray-100 mb-4">数据操作</h3>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between p-3 rounded-lg bg-gray-800/50">
                      <div>
                        <p className="text-sm font-medium text-gray-200">导出数据</p>
                        <p className="text-xs text-gray-500 mt-0.5">导出所有事件和报告数据</p>
                      </div>
                      <button className="btn-default btn-sm">
                        <Download className="w-3.5 h-3.5" />
                        导出
                      </button>
                    </div>
                    <div className="flex items-center justify-between p-3 rounded-lg bg-gray-800/50">
                      <div>
                        <p className="text-sm font-medium text-gray-200">导入数据</p>
                        <p className="text-xs text-gray-500 mt-0.5">从备份文件恢复数据</p>
                      </div>
                      <button className="btn-default btn-sm">
                        <Upload className="w-3.5 h-3.5" />
                        导入
                      </button>
                    </div>
                    <div className="flex items-center justify-between p-3 rounded-lg bg-gray-800/50">
                      <div>
                        <p className="text-sm font-medium text-gray-200">备份数据库</p>
                        <p className="text-xs text-gray-500 mt-0.5">创建完整数据库备份</p>
                      </div>
                      <button className="btn-default btn-sm">
                        <Shield className="w-3.5 h-3.5" />
                        备份
                      </button>
                    </div>
                  </div>
                </div>

                <div className="border-t border-gray-800 pt-6">
                  <h3 className="text-base font-medium text-danger-400 mb-4">危险操作</h3>
                  <div className="alert alert-danger">
                    <Trash2 className="w-4 h-4 flex-shrink-0" />
                    <div className="flex-1">
                      <p className="text-sm font-medium">清空所有数据</p>
                      <p className="text-xs mt-0.5 opacity-80">此操作不可恢复，请谨慎操作</p>
                    </div>
                    <button className="btn-sm bg-danger-500/20 border-danger-500/50 text-danger-400 hover:bg-danger-500/30">
                      清空数据
                    </button>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
