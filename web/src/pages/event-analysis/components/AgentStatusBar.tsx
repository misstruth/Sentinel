const agents = ['采集', '提取', '去重', '评估', '存储']

interface Props {
  currentAgent?: string
}

export default function AgentStatusBar({ currentAgent }: Props) {
  return (
    <div className="flex items-center justify-between bg-slate-800 rounded-lg p-4 mb-4">
      {agents.map((a, i) => {
        const isActive = currentAgent?.includes(a)
        const isPast = currentAgent && agents.indexOf(a) < agents.findIndex(x => currentAgent.includes(x))
        return (
          <div key={a} className="flex items-center">
            <div className={`w-10 h-10 rounded-full flex items-center justify-center text-xs font-medium
              ${isActive ? 'bg-primary-500 text-white' : isPast ? 'bg-green-600 text-white' : 'bg-slate-700 text-slate-400'}`}>
              {a}
            </div>
            {i < 4 && <div className={`w-8 h-0.5 ${isPast ? 'bg-green-600' : 'bg-slate-600'}`} />}
          </div>
        )
      })}
    </div>
  )
}
