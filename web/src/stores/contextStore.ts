import { create } from 'zustand'

interface ContextState {
  currentPage: string
  currentEventId: number | null
  currentEventTitle: string
  setContext: (page: string, eventId?: number, title?: string) => void
}

export const useContextStore = create<ContextState>((set) => ({
  currentPage: '',
  currentEventId: null,
  currentEventTitle: '',
  setContext: (page, eventId, title) => set({
    currentPage: page,
    currentEventId: eventId ?? null,
    currentEventTitle: title ?? '',
  }),
}))
