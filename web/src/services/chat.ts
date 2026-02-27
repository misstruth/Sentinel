import api from './api'
import {
  ChatMessage,
  ChatRequest,
  ChatResponse,
  AIOpsRequest,
  AIOpsResponse,
  ApiResponse,
} from '@/types'

// 生成会话ID
const getSessionId = () => {
  let sessionId = localStorage.getItem('chat_session_id')
  if (!sessionId) {
    sessionId = 'session-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9)
    localStorage.setItem('chat_session_id', sessionId)
  }
  return sessionId
}

export const chatService = {
  // Supervisor 多Agent对话
  supervisorChat: (
    query: string,
    onMessage: (agent: string, content: string) => void,
    onDone: () => void
  ) => {
    fetch('/api/chat/v1/supervisor', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query }),
    }).then(async res => {
      const reader = res.body?.getReader()
      if (!reader) { onDone(); return }
      const decoder = new TextDecoder()
      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        for (const line of decoder.decode(value).split('\n')) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6)
            if (data === '[DONE]') { onDone(); return }
            try {
              const { agent, content } = JSON.parse(data)
              onMessage(agent, content)
            } catch {}
          }
        }
      }
      onDone()
    }).catch(() => onDone())
  },

  // 普通对话
  async chat(data: ChatRequest): Promise<ChatResponse> {
    const res = await api.post<ApiResponse<{ answer: string }>>('/chat', {
      id: getSessionId(),
      question: data.message,
    })
    return { reply: res.data.data.answer }
  },

  // 流式对话
  async chatStream(
    data: { message: string; history?: ChatMessage[] },
    onMessage: (content: string) => void,
    onDone: () => void,
    onError: (error: Error) => void
  ): Promise<void> {
    try {
      const response = await fetch('/api/chat_stream', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token') || ''}`,
        },
        body: JSON.stringify({
          id: getSessionId(),
          question: data.message,
        }),
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const reader = response.body?.getReader()
      if (!reader) {
        throw new Error('No reader available')
      }

      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) {
          onDone()
          break
        }

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6)
            if (data === '[DONE]') {
              onDone()
              return
            }
            try {
              const parsed = JSON.parse(data)
              if (parsed.content) {
                onMessage(parsed.content)
              }
            } catch {
              // 非 JSON 数据，直接作为内容
              onMessage(data)
            }
          }
        }
      }
    } catch (error) {
      onError(error instanceof Error ? error : new Error('Unknown error'))
    }
  },

  // AI 运维分析
  async aiOps(data: AIOpsRequest): Promise<AIOpsResponse> {
    const res = await api.post<ApiResponse<AIOpsResponse>>('/ai_ops', data)
    return res.data.data
  },

  // 文件上传
  async uploadFile(file: File): Promise<{ file_id: string; filename: string }> {
    const formData = new FormData()
    formData.append('file', file)
    const res = await api.post<ApiResponse<{ file_id: string; filename: string }>>('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return res.data.data
  },

  // 清除会话
  clearSession() {
    localStorage.removeItem('chat_session_id')
  },
}
