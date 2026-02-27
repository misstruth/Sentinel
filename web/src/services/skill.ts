import api from './api'
import { ApiResponse } from '@/types'

export interface Skill {
  id: string
  name: string
  description: string
  category: string
  params: { name: string; type: string; description: string; required: boolean }[]
}

export const skillService = {
  list: async () => {
    const res = await api.get<ApiResponse<{ skills: Skill[] }>>('/skill/v1/list')
    return res.data.data.skills
  },

  execute: (
    skillId: string,
    params: Record<string, unknown>,
    onMessage: (type: string, content: string) => void,
    onDone: () => void
  ) => {
    fetch('/api/skill/v1/execute', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ skill_id: skillId, params }),
    }).then(async res => {
      const reader = res.body?.getReader()
      if (!reader) { onDone(); return }
      const decoder = new TextDecoder()
      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        const text = decoder.decode(value)
        for (const line of text.split('\n')) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6)
            if (data === '[DONE]') { onDone(); return }
            try {
              const { type, content } = JSON.parse(data)
              onMessage(type, content)
            } catch { /* ignore */ }
          }
        }
      }
      onDone()
    }).catch(() => onDone())
  },
}
