import type { WSEvent } from '../types'

export type ToastKind = 'info' | 'success' | 'warn'

export interface ToastItem {
  id: number
  title: string
  message?: string
  kind: ToastKind
}

const subscribers = new Set<(items: ToastItem[]) => void>()
const queue: ToastItem[] = []
let nextId = 1
let activeChatId: string | null = null

function emit() {
  subscribers.forEach((cb) => cb([...queue]))
}

export function setActiveChatId(chatId: string | null) {
  activeChatId = chatId
}

export function pushToast(title: string, message?: string, kind: ToastKind = 'info') {
  const item: ToastItem = { id: nextId++, title, message, kind }
  queue.push(item)
  while (queue.length > 4) queue.shift()
  emit()

  const timeout = window.setTimeout(() => {
    const idx = queue.findIndex((toast) => toast.id === item.id)
    if (idx !== -1) {
      queue.splice(idx, 1)
      emit()
    }
  }, 4200)

  if (typeof timeout === 'number') {
    // keep the timer reference alive for the browser runtime; no extra work needed here.
  }
}

export function subscribeToasts(cb: (items: ToastItem[]) => void) {
  subscribers.add(cb)
  cb([...queue])
  return () => subscribers.delete(cb)
}

export async function ensureNotificationPermission(): Promise<NotificationPermission | 'unsupported'> {
  if (!('Notification' in window)) return 'unsupported'

  if (Notification.permission === 'granted') return 'granted'
  if (Notification.permission === 'denied') return 'denied'

  return Notification.requestPermission()
}

export function browserNotify(title: string, body?: string) {
  if (!('Notification' in window)) return

  if (Notification.permission === 'default') {
    Notification.requestPermission().catch(() => {})
    return
  }

  if (Notification.permission !== 'granted') return

  const notification = new Notification(title, {
    body,
    tag: title,
    silent: false,
  })

  notification.onclick = () => {
    window.focus()
    notification.close()
  }
}

export function handleSocketToast(event: WSEvent, ownUserId: string | null, peerNameForChat?: (chatId: string) => string) {
  const isCurrentChatEvent = event.type === 'newMessage' || event.type === 'newReaction' || event.type === 'delReaction'
  if (isCurrentChatEvent && activeChatId === event.payload.chatId) return

  switch (event.type) {
    case 'newChat': {
      pushToast('Новый чат', 'Открыт новый диалог', 'success')
      browserNotify('Новый чат', 'Открыт новый диалог')
      return
    }
    case 'newMessage': {
      if (event.payload.userId === ownUserId) return
      const peerName = peerNameForChat?.(event.payload.chatId) ?? 'Собеседник'
      pushToast(peerName, 'Новое сообщение', 'info')
      browserNotify('Новое сообщение', peerName)
      return
    }
    case 'newReaction': {
      if (event.payload.userId === ownUserId) return
      const peerName = peerNameForChat?.(event.payload.chatId) ?? 'Собеседник'
      pushToast(`${peerName} добавил реакцию`, event.payload.emoji, 'success')
      browserNotify('Реакция', `${peerName}: ${event.payload.emoji}`)
      return
    }
    case 'delReaction': {
      if (event.payload.userId === ownUserId) return
      const peerName = peerNameForChat?.(event.payload.chatId) ?? 'Собеседник'
      pushToast(`${peerName} убрал реакцию`, event.payload.emoji, 'warn')
      browserNotify('Реакция', `${peerName}: ${event.payload.emoji}`)
      return
    }
    default:
      return
  }
}
