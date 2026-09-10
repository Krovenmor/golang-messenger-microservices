import { getChatInfo, getChatHistory, getChatsFull, getMessage } from './messenger'
import { socketManager } from './ws'
import { primeCachedProfile } from './profileCache'
import type { ChatSummary, Message, PublicProfile, WSEvent } from '../types'

const HISTORY_PAGE_SIZE = 50
const MAX_CACHED_MESSAGES_PER_CHAT = 300

function sortChats(list: ChatSummary[]): ChatSummary[] {
  return [...list].sort((a, b) => {
    const at = a.lastMessage?.createdAt ? new Date(a.lastMessage.createdAt).getTime() : 0
    const bt = b.lastMessage?.createdAt ? new Date(b.lastMessage.createdAt).getTime() : 0
    if (bt !== at) return bt - at
    return a.chatId.localeCompare(b.chatId)
  })
}

class ChatStore {
  private chats: ChatSummary[] = []
  private chatListeners = new Set<(chats: ChatSummary[]) => void>()

  private messagesByChat = new Map<string, Message[]>()
  private historyLoaded = new Set<string>()
  private historyExhausted = new Set<string>()
  private loadingMore = new Map<string, Promise<void>>()
  private messageListeners = new Map<string, Set<(messages: Message[]) => void>>()

  private inflightMessageFetches = new Map<string, Promise<Message>>()
  private ownMutations = new Set<string>()
  private pendingSendChats = new Map<string, number>()

  private initPromise: Promise<void> | null = null
  private wsBound = false

  init(): Promise<void> {
    this.bindWs()
    if (this.initPromise) return this.initPromise
    this.initPromise = getChatsFull().then((data) => {
      const byId = new Map<string, ChatSummary>()
      data.forEach((c) => byId.set(c.chatId, c))
      this.chats = sortChats(Array.from(byId.values()))
      this.notifyChats()
    })
    return this.initPromise
  }

  reset() {
    this.chats = []
    this.messagesByChat.clear()
    this.historyLoaded.clear()
    this.historyExhausted.clear()
    this.loadingMore.clear()
    this.inflightMessageFetches.clear()
    this.ownMutations.clear()
    this.pendingSendChats.clear()
    this.initPromise = null
    this.notifyChats()
  }

  getChats(): ChatSummary[] {
    return this.chats
  }

  subscribeChats(cb: (chats: ChatSummary[]) => void): () => void {
    this.chatListeners.add(cb)
    cb(this.chats)
    return () => this.chatListeners.delete(cb)
  }

  getMessages(chatId: string): Message[] {
    return this.messagesByChat.get(chatId) ?? []
  }

  isHistoryLoaded(chatId: string): boolean {
    return this.historyLoaded.has(chatId)
  }

  subscribeMessages(chatId: string, cb: (messages: Message[]) => void): () => void {
    let set = this.messageListeners.get(chatId)
    if (!set) {
      set = new Set()
      this.messageListeners.set(chatId, set)
    }
    set.add(cb)
    cb(this.getMessages(chatId))
    return () => set!.delete(cb)
  }

  async ensureHistory(chatId: string): Promise<void> {
    if (this.historyLoaded.has(chatId)) return
    const data = await getChatHistory(chatId, { q: HISTORY_PAGE_SIZE })
    this.messagesByChat.set(chatId, [...data].reverse())
    this.historyLoaded.add(chatId)
    if (data.length < HISTORY_PAGE_SIZE) {
      this.historyExhausted.add(chatId)
    }
    this.notifyMessages(chatId)
  }

  hasMoreHistory(chatId: string): boolean {
    return this.historyLoaded.has(chatId) && !this.historyExhausted.has(chatId)
  }

  isLoadingMore(chatId: string): boolean {
    return this.loadingMore.has(chatId)
  }

  loadMoreHistory(chatId: string): Promise<void> {
    if (!this.hasMoreHistory(chatId)) return Promise.resolve()
    const existing = this.loadingMore.get(chatId)
    if (existing) return existing

    const promise = (async () => {
      const list = this.messagesByChat.get(chatId) ?? []
      const oldest = list[0]
      if (!oldest) {
        this.historyExhausted.add(chatId)
        return
      }
      const older = await getChatHistory(chatId, { from: oldest.messageId, q: HISTORY_PAGE_SIZE })
      if (older.length === 0) {
        this.historyExhausted.add(chatId)
        return
      }
      const knownIds = new Set(list.map((m) => m.messageId))
      const fresh = [...older].reverse().filter((m) => !knownIds.has(m.messageId))
      if (older.length < HISTORY_PAGE_SIZE) {
        this.historyExhausted.add(chatId)
      }
      if (fresh.length === 0) return
      this.messagesByChat.set(chatId, [...fresh, ...list])
      this.notifyMessages(chatId)
    })().finally(() => {
      this.loadingMore.delete(chatId)
    })

    this.loadingMore.set(chatId, promise)
    return promise
  }

  addLocalChat(chatId: string, own: { userId: string }, peer: PublicProfile) {
    if (this.chats.some((c) => c.chatId === chatId)) return
    primeCachedProfile(peer)
    const now = new Date().toISOString()
    const chat: ChatSummary = {
      chatId: chatId,
      members: [
        { userId: own.userId, joinedAt: now },
        { userId: peer.userId, joinedAt: now },
      ],
      lastMessage: null,
    }
    this.chats = sortChats([chat, ...this.chats])
    this.notifyChats()
  }

  beginOwnSend(chatId: string) {
    this.pendingSendChats.set(chatId, (this.pendingSendChats.get(chatId) ?? 0) + 1)
  }

  endOwnSend(chatId: string) {
    const count = (this.pendingSendChats.get(chatId) ?? 1) - 1
    if (count <= 0) {
      this.pendingSendChats.delete(chatId)
    } else {
      this.pendingSendChats.set(chatId, count)
    }
  }

  applyOwnMessage(chatId: string, message: Message) {
    this.ownMutations.add(message.messageId)
    this.appendMessage(chatId, message)
    this.updatePreview(chatId, message, true)
  }

  applyOwnPatch(chatId: string, messageId: string, patch: Partial<Message>) {
    this.ownMutations.add(messageId)
    const list = this.messagesByChat.get(chatId)
    if (!list) return
    const idx = list.findIndex((m) => m.messageId === messageId)
    if (idx === -1) return
    const patched = { ...list[idx], ...patch }
    const next = [...list]
    next[idx] = patched
    this.messagesByChat.set(chatId, next)
    this.notifyMessages(chatId)

    const chat = this.chats.find((c) => c.chatId === chatId)
    if (chat?.lastMessage?.messageId === messageId) {
      this.patchPreview(chatId, patched)
    }
  }

  private appendMessage(chatId: string, message: Message) {
    const list = this.messagesByChat.get(chatId)
    if (!list) return
    if (list.some((m) => m.messageId === message.messageId)) return
    let next = [...list, message]
    if (next.length > MAX_CACHED_MESSAGES_PER_CHAT) {
      next = next.slice(next.length - MAX_CACHED_MESSAGES_PER_CHAT)
    }
    this.messagesByChat.set(chatId, next)
    this.notifyMessages(chatId)
  }

  private replaceMessage(chatId: string, message: Message) {
    const list = this.messagesByChat.get(chatId)
    if (!list) return
    const idx = list.findIndex((m) => m.messageId === message.messageId)
    if (idx === -1) return
    const next = [...list]
    next[idx] = message
    this.messagesByChat.set(chatId, next)
    this.notifyMessages(chatId)
  }

  private updatePreview(chatId: string, message: Message, isNew: boolean) {
    const chat = this.chats.find((c) => c.chatId === chatId)
    if (!chat) return
    if (isNew || chat.lastMessage?.messageId === message.messageId) {
      this.patchPreview(chatId, message)
    }
  }

  private patchPreview(chatId: string, message: Message) {
    const idx = this.chats.findIndex((c) => c.chatId === chatId)
    if (idx === -1) return
    const next = [...this.chats]
    next[idx] = { ...next[idx], lastMessage: message }
    this.chats = sortChats(next)
    this.notifyChats()
  }

  private notifyChats() {
    this.chatListeners.forEach((cb) => cb(this.chats))
  }

  private notifyMessages(chatId: string) {
    this.messageListeners.get(chatId)?.forEach((cb) => cb(this.getMessages(chatId)))
  }

  private fetchMessageOnce(chatId: string, msgId: string): Promise<Message> {
    const key = `${chatId}:${msgId}`
    let promise = this.inflightMessageFetches.get(key)
    if (!promise) {
      promise = getMessage(chatId, msgId).finally(() => {
        this.inflightMessageFetches.delete(key)
      })
      this.inflightMessageFetches.set(key, promise)
    }
    return promise
  }

  async getOrFetchMessage(chatId: string, messageId: string): Promise<Message | null> {
    const cached = this.messagesByChat.get(chatId)?.find((m) => m.messageId === messageId)
    if (cached) return cached
    try {
      return await this.fetchMessageOnce(chatId, messageId)
    } catch {
      return null
    }
  }

  private patchReaction(chatId: string, msgId: string, emoji: string, userId: string, add: boolean) {
    const list = this.messagesByChat.get(chatId)
    if (!list) return
    const idx = list.findIndex((m) => m.messageId === msgId)
    if (idx === -1) return

    const msg = list[idx]
    const reactions = msg.reactions ? [...msg.reactions] : []
    const groupIdx = reactions.findIndex((r) => r.emoji === emoji)

    if (add) {
      if (groupIdx === -1) {
        reactions.push({ emoji, users: [userId] })
      } else if (!reactions[groupIdx].users.includes(userId)) {
        reactions[groupIdx] = { ...reactions[groupIdx], users: [...reactions[groupIdx].users, userId] }
      } else {
        return
      }
    } else {
      if (groupIdx === -1) return
      const users = reactions[groupIdx].users.filter((u) => u !== userId)
      if (users.length === reactions[groupIdx].users.length) return
      if (users.length === 0) {
        reactions.splice(groupIdx, 1)
      } else {
        reactions[groupIdx] = { ...reactions[groupIdx], users }
      }
    }

    const next = [...list]
    next[idx] = { ...msg, reactions }
    this.messagesByChat.set(chatId, next)
    this.notifyMessages(chatId)
  }

  applyReactionAdded(chatId: string, msgId: string, emoji: string, userId: string) {
    this.patchReaction(chatId, msgId, emoji, userId, true)
  }

  applyReactionRemoved(chatId: string, msgId: string, emoji: string, userId: string) {
    this.patchReaction(chatId, msgId, emoji, userId, false)
  }

  private bindWs() {
    if (this.wsBound) return
    this.wsBound = true

    socketManager.onEvent(async (event: WSEvent) => {
      if (event.type === 'newReaction' || event.type === 'delReaction') {
        const { chatId, msgId, userId, emoji } = event.payload
        this.patchReaction(chatId, msgId, emoji, userId, event.type === 'newReaction')
        return
      }

      if (
        event.type !== 'newChat' &&
        event.type !== 'newMessage' &&
        event.type !== 'messageRedacted' &&
        event.type !== 'messageDeleted'
      ) {
        return
      }

      const chatId = event.payload.chatId
      if (!chatId) return

      if (event.type === 'newChat') {
        if (this.chats.some((c) => c.chatId === chatId)) return
        try {
          const info = await getChatInfo(chatId)
          if (this.chats.some((c) => c.chatId === chatId)) return
          this.chats = sortChats([{ chatId: chatId, members: info.members, lastMessage: null }, ...this.chats])
          this.notifyChats()
        } catch {}
        return
      }

      const msgId = event.payload.msgId
      if (!msgId) return

      if (event.type === 'newMessage' && !this.ownMutations.has(msgId) && this.pendingSendChats.has(chatId)) {
        await new Promise((resolve) => setTimeout(resolve, 250))
      }

      if (this.ownMutations.delete(msgId)) return

      const existingChat = this.chats.find((c) => c.chatId === chatId)
      const isPreviewMessage = event.type === 'newMessage' || existingChat?.lastMessage?.messageId === msgId
      const historyCached = this.historyLoaded.has(chatId)

      if (existingChat && !isPreviewMessage && !historyCached) return

      try {
        const msg = await this.fetchMessageOnce(chatId, msgId)

        if (event.type === 'newMessage') {
          this.appendMessage(chatId, msg)
        } else {
          this.replaceMessage(chatId, msg)
        }

        if (existingChat) {
          this.updatePreview(chatId, msg, event.type === 'newMessage')
        } else {
          const info = await getChatInfo(chatId)
          if (!this.chats.some((c) => c.chatId === chatId)) {
            this.chats = sortChats([{ chatId: chatId, members: info.members, lastMessage: msg }, ...this.chats])
            this.notifyChats()
          }
        }
      } catch {}
    })
  }
}

export const chatStore = new ChatStore()
