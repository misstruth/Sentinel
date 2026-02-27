import { create } from 'zustand'

interface AgentLog {
  agent: string
  status: 'running' | 'success' | 'error'
  message: string
  timestamp: string
  data?: Record<string, unknown>
}

interface EventStore {
  agentLogs: AgentLog[]
  isProcessing: boolean
  addLog: (log: AgentLog) => void
  clearLogs: () => void
  setProcessing: (v: boolean) => void
}

export const useEventStore = create<EventStore>((set) => ({
  agentLogs: [],
  isProcessing: false,
  addLog: (log) => set((s) => ({ agentLogs: [...s.agentLogs, log] })),
  clearLogs: () => set({ agentLogs: [] }),
  setProcessing: (v) => set({ isProcessing: v }),
}))
