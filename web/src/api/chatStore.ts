// Single source of truth for chats + message history, shared by the
// sidebar and the open chat window. Exists specifically so that:
//
//  1. Every WebSocket event triggers AT MOST one GET to the backend for the
//     one message it references — never once per UI consumer that cares
//     about it. `inflightMessageFetches` dedupes even truly concurrent
//     calls. Events for our own actions (send/edit/delete) never fetch at
//     all — the backend broadcasts chat events to every subscriber
//     including the sender, so our own echo is guaranteed, and by then we
//     already have the exact result from the POST/PUT/DELETE response.
//  2. Chat history is cached in memory per chat. Opening a chat fetches its
//     history once; switching away and back re-reads the cache, no network
//     call. Cache is capped per chat to bound memory in a long-lived tab.
import { getChatInfo, getChatHistory, getChatsFull, getMessage } from './messenger'
import { socketManager } from './ws'
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
  // Message IDs from our own send/edit/delete — every one of these is
  // guaranteed to echo back as a WS event (the backend broadcasts to every
  // subscriber of the chat, sender included), so we mark it here and
  // consume the marker the moment that echo shows up, skipping the fetch
  // entirely since we already know the exact result from our own POST/PUT/
  // DELETE response.
  private ownMutations = new Set<string>()
  // chatId -> count of our own sends still awaiting their POST response.
  // Used only to close the race with the WS push, handled in bindWs below.
  private pendingSendChats = new Map<string, number>()

  private initPromise: Promise<void> | null = null
  private wsBound = false

  /** Fetches the chat list exactly once (idempotent) — the only bulk call
   * in the whole app. Everything after this is incremental. */
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

  /** Call on logout so a fresh login in the same tab doesn't inherit the
   * previous account's cached chats/messages. */
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

  /** Loads history once per chat. Subsequent calls for an already-loaded
   * chat resolve immediately with no network request. */
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

  /** Whether scrolling up in this chat could still reveal older messages
   * we haven't fetched yet. False until the initial load has happened, and
   * false again once we've confirmed we hit the beginning of the chat. */
  hasMoreHistory(chatId: string): boolean {
    return this.historyLoaded.has(chatId) && !this.historyExhausted.has(chatId)
  }

  isLoadingMore(chatId: string): boolean {
    return this.loadingMore.has(chatId)
  }

  /** Fetches the next page of older messages (scrolled-up-into-the-past),
   * using the oldest cached message as the pivot. Only ever one request in
   * flight per chat — a second call while one's already running just waits
   * on the same promise instead of firing another GET. */
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

  /** Adds a chat we just created ourselves from data already on hand — no
   * network round-trip. */
  addLocalChat(chatId: string, own: { userId: string; name: string }, peer: PublicProfile) {
    if (this.chats.some((c) => c.chatId === chatId)) return
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

  /** Call right before `POST /msg/chat/{id}` and again (in a `finally`)
   * once it settles — lets the WS handler recognize "there's a send from us
   * in flight for this chat" even before we know the message's real ID. */
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

  /** Optimistic local update for a message we just sent. */
  applyOwnMessage(chatId: string, message: Message) {
    this.ownMutations.add(message.messageId)
    this.appendMessage(chatId, message)
    this.updatePreview(chatId, message, true)
  }

  /** Optimistic local update for a message we just edited/deleted. */
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
    if (!list) return // history not cached for this chat — nothing to append to
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

  /** Shares one in-flight (or just-completed) GET across every caller that
   * asks for the same chatId+msgId at roughly the same time. */
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

  /** For resolving a reply-to reference: checks the already-loaded history
   * first (the common case — you can only reply to a message you can see),
   * falling back to a single deduped GET for anything older that isn't
   * cached yet (e.g. a reply to a message from before the loaded window). */
  async getOrFetchMessage(chatId: string, messageId: string): Promise<Message | null> {
    const cached = this.messagesByChat.get(chatId)?.find((m) => m.messageId === messageId)
    if (cached) return cached
    try {
      return await this.fetchMessageOnce(chatId, messageId)
    } catch {
      return null
    }
  }

  /** Adds/removes a single user's reaction on a message, purely locally —
   * no network call. Used both for our own optimistic toggle and for
   * newReaction/delReaction WS events (which carry the full result already,
   * so there's nothing to fetch either way). Idempotent: adding a userId
   * that's already in the group, or removing one that isn't, is a no-op —
   * so calling this twice for the same logical change (e.g. our own
   * optimistic update followed by its WS echo) is harmless. */
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
        return // status / response frames aren't ours to handle here
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
        } catch {
          // ignore
        }
        return
      }

      const msgId = event.payload.msgId
      if (!msgId) return

      // The WS push for a message we just sent can (and often does) arrive
      // before the POST response completes — they're two independent round
      // trips from the same server action, with no ordering guarantee. If
      // there's a send from us in flight for this exact chat and we don't
      // yet recognize this message, give our own request a brief chance to
      // land and register it before assuming this is someone else's
      // message and hitting the network for it.
      if (event.type === 'newMessage' && !this.ownMutations.has(msgId) && this.pendingSendChats.has(chatId)) {
        await new Promise((resolve) => setTimeout(resolve, 250))
      }

      // We caused this ourselves — we already have the exact result from
      // our own POST/PUT/DELETE response, so there's nothing to fetch.
      if (this.ownMutations.delete(msgId)) return

      const existingChat = this.chats.find((c) => c.chatId === chatId)
      const isPreviewMessage = event.type === 'newMessage' || existingChat?.lastMessage?.messageId === msgId
      const historyCached = this.historyLoaded.has(chatId)

      // Nobody (sidebar preview, open chat window) would use this fetch —
      // skip it. E.g. an edit to a message buried in history of a chat
      // that isn't currently open.
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
          // Event for a chat we don't know about locally yet.
          const info = await getChatInfo(chatId)
          if (!this.chats.some((c) => c.chatId === chatId)) {
            this.chats = sortChats([{ chatId: chatId, members: info.members, lastMessage: msg }, ...this.chats])
            this.notifyChats()
          }
        }
      } catch {
        // message may already be gone — ignore
      }
    })
  }
}

export const chatStore = new ChatStore()
