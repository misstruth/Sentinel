import { useState, useRef, useEffect } from 'react'
import {
  Send,
  Bot,
  User,
  Loader2,
  Trash2,
  FileText,
  AlertTriangle,
  Activity,
  X,
  Zap,
} from 'lucide-react'
import ReactMarkdown from 'react-markdown'
import { chatService, skillService } from '@/services'
import { Skill } from '@/services/skill'
import { ChatMessage } from '@/types'
import { cn } from '@/utils'
import { useContextStore } from '@/stores/contextStore'

type MessageRole = 'user' | 'assistant'

interface Message {
  id: string
  role: MessageRole
  content: string
  timestamp: Date
  isStreaming?: boolean
}

export default function Chat() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [activeTab, setActiveTab] = useState<'chat' | 'aiops'>('chat')
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)
  const streamingContentRef = useRef('')

  const { currentEventId, currentEventTitle, setContext } = useContextStore()

  // AI Ops 状态
  const [aiOpsType, setAiOpsType] = useState<'log_analysis' | 'anomaly_detection' | 'root_cause'>('log_analysis')
  const [aiOpsData, setAiOpsData] = useState('')
  const [aiOpsResult, setAiOpsResult] = useState<string | null>(null)
  const [aiOpsLoading, setAiOpsLoading] = useState(false)

  // Skills 状态
  const [skills, setSkills] = useState<Skill[]>([])
  const [selectedSkill, setSelectedSkill] = useState<Skill | null>(null)
  const [skillParams, setSkillParams] = useState<Record<string, string>>({})

  // 多Agent模式
  const [multiAgentMode, setMultiAgentMode] = useState(false)

  useEffect(() => {
    skillService.list().then(setSkills).catch(() => {})
  }, [])

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const generateId = () => Math.random().toString(36).substring(2, 15)

  const handleSend = async () => {
    if (!input.trim() || isLoading) return

    // 注入上下文到消息
    let messageContent = input.trim()
    if (currentEventId && currentEventTitle) {
      messageContent = `[当前上下文: 事件#${currentEventId} "${currentEventTitle}"]\n${messageContent}`
    }

    const userMessage: Message = {
      id: generateId(),
      role: 'user',
      content: input.trim(),
      timestamp: new Date(),
    }

    setMessages((prev) => [...prev, userMessage])
    setInput('')
    setIsLoading(true)

    const assistantMessage: Message = {
      id: generateId(),
      role: 'assistant',
      content: '',
      timestamp: new Date(),
      isStreaming: true,
    }

    setMessages((prev) => [...prev, assistantMessage])

    const history: ChatMessage[] = messages.map((m) => ({
      role: m.role,
      content: m.content,
    }))

    streamingContentRef.current = ''
    let updateTimer: ReturnType<typeof setTimeout> | null = null

    const flushContent = () => {
      setMessages((prev) =>
        prev.map((m) =>
          m.id === assistantMessage.id
            ? { ...m, content: streamingContentRef.current }
            : m
        )
      )
    }

    try {
      if (multiAgentMode) {
        // 多Agent模式
        chatService.supervisorChat(
          messageContent,
          (_agent, content) => {
            streamingContentRef.current += content
            flushContent()
          },
          () => {
            setMessages((prev) =>
              prev.map((m) => m.id === assistantMessage.id ? { ...m, isStreaming: false } : m)
            )
            setIsLoading(false)
          }
        )
      } else {
        // 普通模式
        await chatService.chatStream(
        { message: messageContent, history },
        (content) => {
          streamingContentRef.current += content
          if (!updateTimer) {
            updateTimer = setTimeout(() => {
              flushContent()
              updateTimer = null
            }, 50)
          }
        },
        () => {
          if (updateTimer) clearTimeout(updateTimer)
          flushContent()
          setMessages((prev) =>
            prev.map((m) =>
              m.id === assistantMessage.id ? { ...m, isStreaming: false } : m
            )
          )
          setIsLoading(false)
        },
        (error) => {
          setMessages((prev) =>
            prev.map((m) =>
              m.id === assistantMessage.id
                ? { ...m, content: `错误: ${error.message}`, isStreaming: false }
                : m
            )
          )
          setIsLoading(false)
        }
      )
      }
    } catch {
      setIsLoading(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const clearMessages = () => {
    setMessages([])
  }

  const handleAiOps = async () => {
    if (!aiOpsData.trim() || aiOpsLoading) return

    setAiOpsLoading(true)
    setAiOpsResult(null)

    try {
      const result = await chatService.aiOps({
        type: aiOpsType,
        data: aiOpsData,
      })
      setAiOpsResult(result.result)
    } catch (error) {
      setAiOpsResult(`分析失败: ${error instanceof Error ? error.message : '未知错误'}`)
    } finally {
      setAiOpsLoading(false)
    }
  }

  const handleSkillExecute = () => {
    if (!selectedSkill || isLoading) return
    setIsLoading(true)

    const assistantMessage: Message = {
      id: generateId(),
      role: 'assistant',
      content: '',
      timestamp: new Date(),
      isStreaming: true,
    }
    setMessages((prev) => [...prev, assistantMessage])

    let content = ''
    skillService.execute(
      selectedSkill.id,
      skillParams,
      (type, text) => {
        content += (type === 'step' ? `[${text}]\n` : text)
        setMessages((prev) =>
          prev.map((m) => m.id === assistantMessage.id ? { ...m, content } : m)
        )
      },
      () => {
        setMessages((prev) =>
          prev.map((m) => m.id === assistantMessage.id ? { ...m, isStreaming: false } : m)
        )
        setIsLoading(false)
        setSelectedSkill(null)
        setSkillParams({})
      }
    )
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-medium text-gray-100">AI 助手</h1>
          <p className="text-sm text-gray-500 mt-1">智能安全分析与运维助手</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="tabs">
        <button
          onClick={() => setActiveTab('chat')}
          className={cn('tab', activeTab === 'chat' && 'active')}
        >
          <Bot className="w-4 h-4 mr-2" />
          AI 对话
        </button>
        <button
          onClick={() => setActiveTab('aiops')}
          className={cn('tab', activeTab === 'aiops' && 'active')}
        >
          <Activity className="w-4 h-4 mr-2" />
          AI 运维分析
        </button>
      </div>

      {activeTab === 'chat' ? (
        <div className="card flex flex-col" style={{ height: 'calc(100vh - 280px)' }}>
          {/* Chat Messages */}
          <div className="flex-1 overflow-y-auto p-4 space-y-4">
            {messages.length === 0 ? (
              <div className="h-full flex flex-col items-center justify-center text-gray-500">
                <Bot className="w-12 h-12 mb-4 opacity-40" />
                <p className="text-sm">开始与 AI 助手对话</p>
                <p className="text-xs text-gray-600 mt-1">可以询问安全事件、威胁情报、漏洞分析等问题</p>
              </div>
            ) : (
              messages.map((message) => (
                <div
                  key={message.id}
                  className={cn(
                    'flex gap-3',
                    message.role === 'user' ? 'justify-end' : 'justify-start'
                  )}
                >
                  {message.role === 'assistant' && (
                    <div className="w-8 h-8 rounded bg-primary-500/20 flex items-center justify-center flex-shrink-0">
                      <Bot className="w-4 h-4 text-primary-400" />
                    </div>
                  )}
                  <div
                    className={cn(
                      'max-w-[70%] rounded-lg px-4 py-2.5 text-sm',
                      message.role === 'user'
                        ? 'bg-primary-500 text-white'
                        : 'bg-gray-800 text-gray-200'
                    )}
                  >
                    {message.role === 'assistant' ? (
                      <div className="prose prose-invert prose-sm max-w-none">
                        <ReactMarkdown>{message.content}</ReactMarkdown>
                      </div>
                    ) : (
                      <p className="whitespace-pre-wrap">{message.content}</p>
                    )}
                    {message.isStreaming && (
                      <span className="inline-block w-1.5 h-4 bg-primary-400 animate-pulse ml-0.5" />
                    )}
                  </div>
                  {message.role === 'user' && (
                    <div className="w-8 h-8 rounded bg-gray-700 flex items-center justify-center flex-shrink-0">
                      <User className="w-4 h-4 text-gray-400" />
                    </div>
                  )}
                </div>
              ))
            )}
            <div ref={messagesEndRef} />
          </div>

          {/* Input Area */}
          <div className="border-t border-gray-800 p-4">
            {/* 多Agent模式切换 */}
            <div className="mb-2 flex items-center gap-2">
              <button
                onClick={() => setMultiAgentMode(!multiAgentMode)}
                className={cn('btn-sm', multiAgentMode ? 'btn-primary' : 'btn-default')}
              >
                <Bot className="w-3 h-3" />
                {multiAgentMode ? '多Agent模式' : '单Agent模式'}
              </button>
            </div>

            {/* Skills 快捷按钮 */}
            <div className="mb-2 flex gap-2 flex-wrap">
              {skills.map((skill) => (
                <button
                  key={skill.id}
                  onClick={() => { setSelectedSkill(skill); setSkillParams({}) }}
                  className={cn('btn-sm', selectedSkill?.id === skill.id ? 'btn-primary' : 'btn-default')}
                >
                  <Zap className="w-3 h-3" />
                  {skill.name}
                </button>
              ))}
            </div>

            {/* Skill 参数输入 */}
            {selectedSkill && (
              <div className="mb-2 p-3 bg-gray-800 rounded border border-gray-700">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm text-gray-300">{selectedSkill.name}</span>
                  <button onClick={() => setSelectedSkill(null)} className="text-gray-500 hover:text-white">
                    <X className="w-3 h-3" />
                  </button>
                </div>
                <div className="flex gap-2 items-end">
                  {selectedSkill.params.map((p) => (
                    <input
                      key={p.name}
                      type={p.type === 'number' ? 'number' : 'text'}
                      placeholder={p.description}
                      value={skillParams[p.name] || ''}
                      onChange={(e) => setSkillParams({ ...skillParams, [p.name]: e.target.value })}
                      className="input flex-1 text-sm"
                    />
                  ))}
                  <button onClick={handleSkillExecute} disabled={isLoading} className="btn-primary btn-sm">
                    {isLoading ? <Loader2 className="w-3 h-3 animate-spin" /> : '执行'}
                  </button>
                </div>
              </div>
            )}

            {currentEventId && (
              <div className="mb-2 px-3 py-1.5 bg-primary-500/10 border border-primary-500/30 rounded flex items-center justify-between text-xs">
                <span className="text-primary-300">上下文: 事件#{currentEventId} {currentEventTitle}</span>
                <button onClick={() => setContext('', undefined, undefined)} className="text-gray-500 hover:text-white">
                  <X className="w-3 h-3" />
                </button>
              </div>
            )}
            <div className="flex gap-3">
              <textarea
                ref={inputRef}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="输入消息，按 Enter 发送..."
                className="textarea flex-1"
                rows={2}
                disabled={isLoading}
              />
              <div className="flex flex-col gap-2">
                <button
                  onClick={handleSend}
                  disabled={!input.trim() || isLoading}
                  className="btn-primary h-10 w-10 p-0"
                >
                  {isLoading ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <Send className="w-4 h-4" />
                  )}
                </button>
                <button
                  onClick={clearMessages}
                  className="btn-default h-10 w-10 p-0"
                  title="清空对话"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      ) : (
        /* AI Ops */
        <div className="grid grid-cols-2 gap-6" style={{ height: 'calc(100vh - 280px)' }}>
          {/* Left - Input */}
          <div className="card flex flex-col">
            <div className="card-header">输入数据</div>
            <div className="card-body flex-1 flex flex-col gap-4">
              <div className="form-item">
                <label className="label">分析类型</label>
                <div className="flex gap-2">
                  <button
                    onClick={() => setAiOpsType('log_analysis')}
                    className={cn(
                      'btn-sm',
                      aiOpsType === 'log_analysis' ? 'btn-primary' : 'btn-default'
                    )}
                  >
                    <FileText className="w-3.5 h-3.5" />
                    日志分析
                  </button>
                  <button
                    onClick={() => setAiOpsType('anomaly_detection')}
                    className={cn(
                      'btn-sm',
                      aiOpsType === 'anomaly_detection' ? 'btn-primary' : 'btn-default'
                    )}
                  >
                    <AlertTriangle className="w-3.5 h-3.5" />
                    异常检测
                  </button>
                  <button
                    onClick={() => setAiOpsType('root_cause')}
                    className={cn(
                      'btn-sm',
                      aiOpsType === 'root_cause' ? 'btn-primary' : 'btn-default'
                    )}
                  >
                    <Activity className="w-3.5 h-3.5" />
                    根因分析
                  </button>
                </div>
              </div>

              <div className="form-item flex-1 flex flex-col">
                <label className="label">数据内容</label>
                <textarea
                  value={aiOpsData}
                  onChange={(e) => setAiOpsData(e.target.value)}
                  placeholder="粘贴日志、指标数据或错误信息..."
                  className="textarea flex-1 font-mono text-xs"
                />
              </div>

              <button
                onClick={handleAiOps}
                disabled={!aiOpsData.trim() || aiOpsLoading}
                className="btn-primary"
              >
                {aiOpsLoading ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin" />
                    分析中...
                  </>
                ) : (
                  <>
                    <Activity className="w-4 h-4" />
                    开始分析
                  </>
                )}
              </button>
            </div>
          </div>

          {/* Right - Result */}
          <div className="card flex flex-col">
            <div className="card-header">分析结果</div>
            <div className="card-body flex-1 overflow-y-auto">
              {aiOpsResult ? (
                <pre className="whitespace-pre-wrap text-sm text-gray-200 font-mono">
                  {aiOpsResult}
                </pre>
              ) : (
                <div className="h-full flex flex-col items-center justify-center text-gray-500">
                  <Activity className="w-12 h-12 mb-4 opacity-40" />
                  <p className="text-sm">输入数据并点击分析按钮</p>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
