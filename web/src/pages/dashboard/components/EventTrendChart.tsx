import { useState, useEffect } from 'react'
import ReactECharts from 'echarts-for-react'
import { Loader2 } from 'lucide-react'
import { eventService } from '@/services/event'

interface DayData {
  critical: number
  high: number
  medium: number
  low: number
}

export default function EventTrendChart() {
  const [loading, setLoading] = useState(true)
  const [trendData, setTrendData] = useState<{ dates: string[]; data: DayData[] }>({
    dates: [],
    data: [],
  })

  useEffect(() => {
    const fetchData = async () => {
      try {
        // 获取最近的事件数据
        const res = await eventService.list({ page: 1, size: 200 })

        // 生成最近7天的日期
        const dates: string[] = []
        const dataMap: Record<string, DayData> = {}

        for (let i = 6; i >= 0; i--) {
          const date = new Date()
          date.setDate(date.getDate() - i)
          const dateStr = date.toISOString().split('T')[0]
          const weekDay = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'][date.getDay()]
          dates.push(weekDay)
          dataMap[dateStr] = { critical: 0, high: 0, medium: 0, low: 0 }
        }

        // 按日期统计事件
        res.list.forEach(event => {
          const eventDate = event.event_time?.split('T')[0] || event.created_at?.split('T')[0]
          if (eventDate && dataMap[eventDate]) {
            const severity = event.severity as keyof DayData
            if (severity in dataMap[eventDate]) {
              dataMap[eventDate][severity]++
            }
          }
        })

        const data = Object.values(dataMap)
        setTrendData({ dates, data })
      } catch (error) {
        console.error('获取事件趋势失败:', error)
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center" style={{ height: '300px' }}>
        <Loader2 className="w-6 h-6 animate-spin text-primary-500" />
      </div>
    )
  }

  const option = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#1e293b',
      borderColor: '#334155',
      textStyle: { color: '#f1f5f9' },
    },
    legend: {
      data: ['严重', '高危', '中危', '低危'],
      textStyle: { color: '#94a3b8' },
      top: 0,
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: trendData.dates,
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { color: '#94a3b8' },
    },
    yAxis: {
      type: 'value',
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { color: '#94a3b8' },
      splitLine: { lineStyle: { color: '#1e293b' } },
    },
    series: [
      {
        name: '严重',
        type: 'line',
        smooth: true,
        data: trendData.data.map(d => d.critical),
        itemStyle: { color: '#ef4444' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(239, 68, 68, 0.3)' },
              { offset: 1, color: 'rgba(239, 68, 68, 0)' },
            ],
          },
        },
      },
      {
        name: '高危',
        type: 'line',
        smooth: true,
        data: trendData.data.map(d => d.high),
        itemStyle: { color: '#f97316' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(249, 115, 22, 0.3)' },
              { offset: 1, color: 'rgba(249, 115, 22, 0)' },
            ],
          },
        },
      },
      {
        name: '中危',
        type: 'line',
        smooth: true,
        data: trendData.data.map(d => d.medium),
        itemStyle: { color: '#eab308' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(234, 179, 8, 0.3)' },
              { offset: 1, color: 'rgba(234, 179, 8, 0)' },
            ],
          },
        },
      },
      {
        name: '低危',
        type: 'line',
        smooth: true,
        data: trendData.data.map(d => d.low),
        itemStyle: { color: '#3b82f6' },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
              { offset: 1, color: 'rgba(59, 130, 246, 0)' },
            ],
          },
        },
      },
    ],
  }

  return <ReactECharts option={option} style={{ height: '300px' }} />
}
