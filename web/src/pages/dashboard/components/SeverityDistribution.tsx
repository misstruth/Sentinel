import { useState, useEffect } from 'react'
import ReactECharts from 'echarts-for-react'
import { Loader2 } from 'lucide-react'
import { eventService } from '@/services/event'

interface SeverityCount {
  critical: number
  high: number
  medium: number
  low: number
  info: number
}

export default function SeverityDistribution() {
  const [loading, setLoading] = useState(true)
  const [counts, setCounts] = useState<SeverityCount>({
    critical: 0,
    high: 0,
    medium: 0,
    low: 0,
    info: 0,
  })

  useEffect(() => {
    const fetchData = async () => {
      try {
        // 获取较多事件来统计分布
        const res = await eventService.list({ page: 1, size: 100 })
        const severityCounts: SeverityCount = {
          critical: 0,
          high: 0,
          medium: 0,
          low: 0,
          info: 0,
        }
        res.list.forEach(event => {
          const severity = event.severity as keyof SeverityCount
          if (severity in severityCounts) {
            severityCounts[severity]++
          }
        })
        setCounts(severityCounts)
      } catch (error) {
        console.error('获取严重级别分布失败:', error)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center py-8">
        <Loader2 className="w-6 h-6 animate-spin text-primary-500" />
      </div>
    )
  }

  const chartData = [
    { value: counts.critical, name: '严重', itemStyle: { color: '#ef4444' } },
    { value: counts.high, name: '高危', itemStyle: { color: '#f97316' } },
    { value: counts.medium, name: '中危', itemStyle: { color: '#eab308' } },
    { value: counts.low, name: '低危', itemStyle: { color: '#3b82f6' } },
    { value: counts.info, name: '信息', itemStyle: { color: '#64748b' } },
  ]

  const option = {
    tooltip: {
      trigger: 'item',
      backgroundColor: '#1e293b',
      borderColor: '#334155',
      textStyle: { color: '#f1f5f9' },
    },
    series: [
      {
        type: 'pie',
        radius: ['50%', '75%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 8,
          borderColor: '#0f172a',
          borderWidth: 2,
        },
        label: {
          show: false,
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 14,
            fontWeight: 'bold',
            color: '#f1f5f9',
          },
        },
        labelLine: {
          show: false,
        },
        data: chartData,
      },
    ],
  }

  const legendData = [
    { label: '严重', value: counts.critical, color: 'bg-red-500' },
    { label: '高危', value: counts.high, color: 'bg-orange-500' },
    { label: '中危', value: counts.medium, color: 'bg-yellow-500' },
    { label: '低危', value: counts.low, color: 'bg-blue-500' },
  ]

  return (
    <div>
      <ReactECharts option={option} style={{ height: '200px' }} />
      <div className="grid grid-cols-2 gap-2 mt-4">
        {legendData.map((item) => (
          <div key={item.label} className="flex items-center gap-2">
            <div className={`w-2 h-2 rounded-full ${item.color}`} />
            <span className="text-sm text-dark-400">{item.label}</span>
            <span className="text-sm font-medium text-dark-200 ml-auto">
              {item.value}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}
