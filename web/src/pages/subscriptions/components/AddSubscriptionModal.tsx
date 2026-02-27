import { useState } from 'react'
import { X, Github, Rss, Database, Loader2, Globe, Webhook, Shield, Bug, AlertTriangle } from 'lucide-react'
import { subscriptionService } from '@/services/subscription'
import toast from 'react-hot-toast'

import type { SourceType } from '@/types'

interface AddSubscriptionModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: () => void
}

const sourceTypes = [
  {
    value: 'github_repo',
    label: 'GitHub 仓库',
    icon: Github,
    description: '监控 GitHub 仓库的 Release 和安全公告',
    placeholder: 'https://github.com/owner/repo',
  },
  {
    value: 'rss',
    label: 'RSS 订阅',
    icon: Rss,
    description: '订阅安全博客、新闻等 RSS 源',
    placeholder: 'https://example.com/feed.xml',
  },
  {
    value: 'nvd',
    label: 'NVD 漏洞库',
    icon: Database,
    description: '订阅 NVD 国家漏洞数据库',
    placeholder: '关键词过滤（可选）',
  },
  {
    value: 'cve',
    label: 'CVE 数据',
    icon: Bug,
    description: '订阅 CVE 漏洞数据',
    placeholder: '关键词过滤（可选）',
  },
  {
    value: 'threat_intel',
    label: '威胁情报',
    icon: Shield,
    description: '订阅威胁情报源',
    placeholder: 'https://threatfeed.example.com/api',
  },
  {
    value: 'vendor_advisory',
    label: '厂商公告',
    icon: Globe,
    description: '订阅厂商安全公告',
    placeholder: 'https://vendor.com/security/advisories',
  },
  {
    value: 'attack_activity',
    label: '攻击活动',
    icon: AlertTriangle,
    description: '监控攻击活动情报',
    placeholder: 'https://attack.example.com/feed',
  },
  {
    value: 'webhook',
    label: 'Webhook',
    icon: Webhook,
    description: '接收外部系统推送的安全事件',
    placeholder: '自动生成 Webhook URL',
  },
]

// 将分钟转换为 cron 表达式
function minutesToCron(minutes: number): string {
  if (minutes < 60) {
    return `*/${minutes} * * * *`
  }
  const hours = Math.floor(minutes / 60)
  if (hours < 24) {
    return `0 */${hours} * * *`
  }
  return '0 0 * * *'
}

export default function AddSubscriptionModal({
  isOpen,
  onClose,
  onSuccess,
}: AddSubscriptionModalProps) {
  const [step, setStep] = useState(1)
  const [sourceType, setSourceType] = useState<SourceType | ''>('')
  const [formData, setFormData] = useState({
    name: '',
    url: '',
    fetch_interval: 60,
    keywords: '',
    auth_type: 'none',
    auth_config: '',
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

  const selectedSource = sourceTypes.find((s) => s.value === sourceType)

  const handleSubmit = async () => {
    setIsSubmitting(true)
    try {
      const config: Record<string, unknown> = {}
      if (formData.keywords) {
        config.keywords = formData.keywords.split(',').map(k => k.trim())
      }
      if (formData.auth_type !== 'none') {
        config.auth_type = formData.auth_type
        config.auth_config = formData.auth_config
      }

      await subscriptionService.create({
        name: formData.name,
        description: '',
        source_type: sourceType as SourceType,
        source_url: formData.url,
        cron_expr: minutesToCron(formData.fetch_interval),
        config: Object.keys(config).length > 0 ? JSON.stringify(config) : '',
      })
      toast.success('订阅创建成功')
      onSuccess?.()
      handleClose()
    } catch (error) {
      toast.error('创建订阅失败')
      console.error(error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleClose = () => {
    setStep(1)
    setSourceType('')
    setFormData({ name: '', url: '', fetch_interval: 60, keywords: '', auth_type: 'none', auth_config: '' })
    onClose()
  }

  if (!isOpen) return null

  return (
    <>
      {/* Backdrop */}
      <div className="modal-overlay" onClick={handleClose}>
        {/* Modal */}
        <div className="modal w-full max-w-lg" onClick={(e) => e.stopPropagation()}>
          {/* Header */}
          <div className="modal-header">
            <div>
              <h2 className="modal-title">添加订阅</h2>
              <p className="text-sm text-gray-500 mt-0.5">
                {step === 1 ? '选择数据源类型' : '配置订阅信息'}
              </p>
            </div>
            <button
              onClick={handleClose}
              className="p-1.5 rounded text-gray-400 hover:text-gray-200 hover:bg-gray-800 transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          <div className="modal-body max-h-[60vh]">
            {/* Step 1: Select Source Type */}
            {step === 1 && (
              <div className="grid grid-cols-2 gap-3">
                {sourceTypes.map((source) => (
                  <button
                    key={source.value}
                    onClick={() => {
                      setSourceType(source.value as SourceType)
                      setStep(2)
                    }}
                    className="flex flex-col items-center gap-2 p-4 rounded-lg border border-gray-700 hover:border-primary-500 hover:bg-gray-800/50 transition-colors text-center"
                  >
                    <div className="w-10 h-10 rounded bg-gray-800 flex items-center justify-center">
                      <source.icon className="w-5 h-5 text-gray-400" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-gray-200">{source.label}</p>
                      <p className="text-xs text-gray-500 mt-0.5 line-clamp-2">{source.description}</p>
                    </div>
                  </button>
                ))}
              </div>
            )}

            {/* Step 2: Configure */}
            {step === 2 && selectedSource && (
              <div className="space-y-4">
                {/* Source Type Badge */}
                <div className="flex items-center gap-3 p-3 rounded-lg bg-gray-800">
                  <selectedSource.icon className="w-5 h-5 text-primary-400" />
                  <span className="text-sm text-gray-200">{selectedSource.label}</span>
                  <button
                    onClick={() => setStep(1)}
                    className="ml-auto text-xs text-primary-400 hover:text-primary-300"
                  >
                    更换类型
                  </button>
                </div>

                {/* Name */}
                <div className="form-item">
                  <label className="label">订阅名称 <span className="text-danger-500">*</span></label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="输入订阅名称"
                    className="input"
                  />
                </div>

                {/* URL */}
                {sourceType !== 'webhook' && (
                  <div className="form-item">
                    <label className="label">
                      {sourceType === 'nvd' || sourceType === 'cve' ? '关键词' : '订阅地址'}
                      {sourceType !== 'nvd' && sourceType !== 'cve' && <span className="text-danger-500">*</span>}
                    </label>
                    <input
                      type="text"
                      value={formData.url}
                      onChange={(e) => setFormData({ ...formData, url: e.target.value })}
                      placeholder={selectedSource.placeholder}
                      className="input"
                    />
                  </div>
                )}

                {/* Keywords */}
                <div className="form-item">
                  <label className="label">关键词过滤</label>
                  <input
                    type="text"
                    value={formData.keywords}
                    onChange={(e) => setFormData({ ...formData, keywords: e.target.value })}
                    placeholder="多个关键词用逗号分隔"
                    className="input"
                  />
                  <p className="text-xs text-gray-500 mt-1">只抓取包含这些关键词的事件</p>
                </div>

                {/* Fetch Interval */}
                <div className="form-item">
                  <label className="label">抓取间隔</label>
                  <select
                    value={formData.fetch_interval}
                    onChange={(e) => setFormData({ ...formData, fetch_interval: Number(e.target.value) })}
                    className="select"
                  >
                    <option value={15}>15 分钟</option>
                    <option value={30}>30 分钟</option>
                    <option value={60}>1 小时</option>
                    <option value={180}>3 小时</option>
                    <option value={360}>6 小时</option>
                    <option value={720}>12 小时</option>
                    <option value={1440}>24 小时</option>
                  </select>
                </div>

                {/* Auth Type */}
                <div className="form-item">
                  <label className="label">认证方式</label>
                  <select
                    value={formData.auth_type}
                    onChange={(e) => setFormData({ ...formData, auth_type: e.target.value })}
                    className="select"
                  >
                    <option value="none">无需认证</option>
                    <option value="api_key">API Key</option>
                    <option value="basic">Basic Auth</option>
                    <option value="bearer">Bearer Token</option>
                    <option value="oauth2">OAuth 2.0</option>
                  </select>
                </div>

                {/* Auth Config */}
                {formData.auth_type !== 'none' && (
                  <div className="form-item">
                    <label className="label">
                      {formData.auth_type === 'api_key' && 'API Key'}
                      {formData.auth_type === 'basic' && '用户名:密码'}
                      {formData.auth_type === 'bearer' && 'Token'}
                      {formData.auth_type === 'oauth2' && 'OAuth 配置 (JSON)'}
                    </label>
                    <input
                      type={formData.auth_type === 'basic' ? 'password' : 'text'}
                      value={formData.auth_config}
                      onChange={(e) => setFormData({ ...formData, auth_config: e.target.value })}
                      placeholder="输入认证信息"
                      className="input"
                    />
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Footer */}
          {step === 2 && (
            <div className="modal-footer">
              <button onClick={handleClose} className="btn-default">
                取消
              </button>
              <button
                onClick={handleSubmit}
                disabled={!formData.name || (sourceType !== 'webhook' && sourceType !== 'nvd' && sourceType !== 'cve' && !formData.url) || isSubmitting}
                className="btn-primary"
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin" />
                    创建中...
                  </>
                ) : (
                  '创建订阅'
                )}
              </button>
            </div>
          )}
        </div>
      </div>
    </>
  )
}
