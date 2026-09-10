import { useEffect, useLayoutEffect, useRef, useState, type KeyboardEvent } from 'react'
import type { ChatMember, Message, StatusInfo } from '../types'
import { MessageBubble } from './MessageBubble'
import { EmojiPicker } from './EmojiPicker'
import { PeerProfileModal } from './PeerProfileModal'
import { deleteMessage, editMessage, postMessage } from '../api/messenger'
import { addReaction, removeReaction } from '../api/reactions'
import { presenceTracker } from '../api/presenceTracker'
import { chatStore } from '../api/chatStore'
import { getCachedProfile, usePublicProfile } from '../api/profileCache'
import { encryptMessage } from '../utils/messageCrypto'
import { presence } from '../utils/presence'
import { Avatar } from './Avatar'
import { isStatusActive, statusLabel } from '../utils/statusLabel'
import { useDecryptedText } from '../utils/useDecryptedText'
import { onViewportChange } from '../utils/viewportMetrics'

const LOAD_MORE_SCROLL_THRESHOLD = 80

interface Props {
  chatId: string
  ownUserId: string
  ownPubKey: string
  peer: ChatMember | null
  onBack?: () => void
}

export function ChatWindow({ chatId, ownUserId, ownPubKey, peer, onBack }: Props) {
  const [messages, setMessages] = useState<Message[]>(() => chatStore.getMessages(chatId))
  const [draft, setDraft] = useState('')
  const [sending, setSending] = useState(false)
  const [loading, setLoading] = useState(() => !chatStore.isHistoryLoaded(chatId))
  const [loadingMore, setLoadingMore] = useState(false)
  const [peerStatus, setPeerStatus] = useState<StatusInfo | null>(null)
  const [showEmojiPicker, setShowEmojiPicker] = useState(false)
  const [sendError, setSendError] = useState<string | null>(null)
  const [replyingTo, setReplyingTo] = useState<Message | null>(null)
  const [showPeerProfile, setShowPeerProfile] = useState(false)

  const peerProfile = usePublicProfile(peer?.userId ?? null)

  const scrollRef = useRef<HTMLDivElement>(null)
  const bottomRef = useRef<HTMLDivElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const isPrependingRef = useRef(false)
  const loadingMoreRef = useRef(false)
  const readyRef = useRef(false)
  const bottomDistanceRef = useRef(0)

  const replyPreview = useDecryptedText(replyingTo, ownUserId, chatId)

  useEffect(() => {
    readyRef.current = false
    setLoading(!chatStore.isHistoryLoaded(chatId))
    const unsub = chatStore.subscribeMessages(chatId, setMessages)
    chatStore.ensureHistory(chatId).finally(() => setLoading(false))
    return unsub
  }, [chatId])

  useEffect(() => {
    setReplyingTo(null)
  }, [chatId])

  useLayoutEffect(() => {
    if (loading) return
    const el = scrollRef.current
    if (!el) return

    if (!readyRef.current) {
      el.scrollTop = el.scrollHeight
      bottomDistanceRef.current = 0
      readyRef.current = true
      return
    }

    if (!isPrependingRef.current) {
      bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
    }
  }, [loading, messages])

  useEffect(() => {
    if (!peer) {
      setPeerStatus(null)
      return
    }
    return presenceTracker.subscribe(peer.userId, setPeerStatus)
  }, [peer])

  useLayoutEffect(() => {
    const el = textareaRef.current
    if (!el) return
    el.style.height = 'auto'
    el.style.height = `${Math.min(el.scrollHeight, 160)}px`
  }, [draft])

  useEffect(
    () =>
      onViewportChange(() => {
        const el = scrollRef.current
        if (!el) return
        el.scrollTop = Math.max(0, el.scrollHeight - el.clientHeight - bottomDistanceRef.current)
      }),
    []
  )

  async function handleScroll() {
    const scrolled = scrollRef.current
    if (scrolled) {
      bottomDistanceRef.current = Math.max(0, scrolled.scrollHeight - scrolled.scrollTop - scrolled.clientHeight)
    }

    if (!readyRef.current || loadingMoreRef.current) return
    const el = scrollRef.current
    if (!el) return
    if (el.scrollTop > LOAD_MORE_SCROLL_THRESHOLD) return
    if (!chatStore.hasMoreHistory(chatId)) return

    loadingMoreRef.current = true
    setLoadingMore(true)
    isPrependingRef.current = true
    const prevScrollHeight = el.scrollHeight
    const prevScrollTop = el.scrollTop

    await chatStore.loadMoreHistory(chatId)

    requestAnimationFrame(() => {
      const node = scrollRef.current
      if (node) {
        node.scrollTop = node.scrollHeight - prevScrollHeight + prevScrollTop
      }
      isPrependingRef.current = false
      loadingMoreRef.current = false
      setLoadingMore(false)
    })
  }

  async function handleSend() {
    const text = draft.trim()
    if (!text || !peer) return
    const replyToId = replyingTo?.messageId ?? null
    setSending(true)
    setSendError(null)
    setDraft('')
    setReplyingTo(null)
    chatStore.beginOwnSend(chatId)
    try {
      const peerProfile = await getCachedProfile(peer.userId)
      const envelope = await encryptMessage(text, ownPubKey, peerProfile.pubKey, chatId)
      const { messageId } = await postMessage(chatId, envelope, replyToId)
      chatStore.applyOwnMessage(chatId, {
        messageId,
        senderId: ownUserId,
        message: envelope.message,
        senderKey: envelope.senderKey,
        receiverKey: envelope.receiverKey,
        nonce: envelope.nonce,
        createdAt: new Date().toISOString(),
        redactedAt: null,
        deletedAt: null,
        replyToId,
      })
    } catch {
      setSendError('Не удалось зашифровать или отправить сообщение.')
      setDraft(text)
    } finally {
      chatStore.endOwnSend(chatId)
      setSending(false)
    }
  }

  async function handleEdit(messageId: string, text: string) {
    if (!peer) return
    const peerProfile = await getCachedProfile(peer.userId)
    const envelope = await encryptMessage(text, ownPubKey, peerProfile.pubKey, chatId)
    await editMessage(chatId, messageId, envelope)
    chatStore.applyOwnPatch(chatId, messageId, {
      message: envelope.message,
      senderKey: envelope.senderKey,
      receiverKey: envelope.receiverKey,
      nonce: envelope.nonce,
      redactedAt: new Date().toISOString(),
    })
  }

  async function handleDelete(messageId: string) {
    await deleteMessage(chatId, messageId)
    chatStore.applyOwnPatch(chatId, messageId, { deletedAt: new Date().toISOString(), message: '' })
    if (replyingTo?.messageId === messageId) setReplyingTo(null)
  }

  function handleReply(message: Message) {
    setReplyingTo(message)
    textareaRef.current?.focus()
  }

  async function handleToggleReaction(messageId: string, emoji: string) {
    const target = messages.find((m) => m.messageId === messageId)
    const alreadyReacted = target?.reactions?.find((r) => r.emoji === emoji)?.users.includes(ownUserId) ?? false

    if (alreadyReacted) {
      chatStore.applyReactionRemoved(chatId, messageId, emoji, ownUserId)
    } else {
      chatStore.applyReactionAdded(chatId, messageId, emoji, ownUserId)
    }

    try {
      if (alreadyReacted) {
        await removeReaction(chatId, messageId, emoji)
      } else {
        await addReaction(chatId, messageId, emoji)
      }
    } catch {
      if (alreadyReacted) {
        chatStore.applyReactionAdded(chatId, messageId, emoji, ownUserId)
      } else {
        chatStore.applyReactionRemoved(chatId, messageId, emoji, ownUserId)
      }
    }
  }

  function insertEmoji(emoji: string) {
    const el = textareaRef.current
    const start = el?.selectionStart ?? draft.length
    const end = el?.selectionEnd ?? draft.length
    const next = draft.slice(0, start) + emoji + draft.slice(end)
    setDraft(next)
    presence.notifyTyping()
    requestAnimationFrame(() => {
      const node = textareaRef.current
      if (!node) return
      const pos = start + emoji.length
      node.focus()
      node.setSelectionRange(pos, pos)
    })
  }

  function handleKeyDown(e: KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
    if (e.key === 'Escape' && replyingTo) {
      setReplyingTo(null)
    }
  }

  const replyAuthor = replyingTo && (replyingTo.senderId === ownUserId ? 'Вы' : peerProfile?.name ?? 'Собеседник')
  const replyText =
    replyPreview.status === 'ok'
      ? replyPreview.text
      : replyPreview.status === 'locked'
        ? '🔒 сообщение заблокировано'
        : replyPreview.status === 'error'
          ? 'не удалось расшифровать'
          : '…'

  return (
    <div className="flex h-full min-h-0 flex-1 flex-col overflow-hidden bg-void">
      {showPeerProfile && peerProfile && (
        <PeerProfileModal profile={peerProfile} status={peerStatus} onClose={() => setShowPeerProfile(false)} />
      )}

      <div className="flex shrink-0 items-center gap-3 border-b border-line px-4 pb-4 pt-[calc(1rem_+_env(safe-area-inset-top))] sm:px-6">
        {onBack && (
          <button
            type="button"
            onClick={onBack}
            aria-label="Назад к чатам"
            className="-ml-1 shrink-0 rounded-lg p-1.5 text-muted hover:bg-raised hover:text-ink sm:hidden"
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M15 18l-6-6 6-6" strokeLinecap="round" strokeLinejoin="round" />
            </svg>
          </button>
        )}
        <button
          type="button"
          onClick={() => setShowPeerProfile(true)}
          disabled={!peerProfile}
          aria-label="Профиль собеседника"
          className="-mx-2 flex min-w-0 flex-1 items-center gap-3 rounded-lg px-2 py-1 text-left transition-colors enabled:hover:bg-raised disabled:cursor-default"
        >
          <Avatar
            avatarId={peerProfile?.additional?.avatars?.[0]}
            seed={peer?.userId ?? chatId}
            name={peerProfile?.name ?? '? ?'}
            className="h-9 w-9 text-xs"
          />
          <div className="min-w-0">
            <p className="truncate text-sm font-medium text-ink">{peerProfile?.name ?? 'Диалог'}</p>
            {peerStatus && (
              <p className={`truncate font-mono text-xs ${isStatusActive(peerStatus) ? 'text-signal' : 'text-faint'}`}>
                {statusLabel(peerStatus)}
              </p>
            )}
          </div>
        </button>
      </div>

      <div
        ref={scrollRef}
        onScroll={handleScroll}
        style={{ overflowAnchor: 'none' }}
        className="min-h-0 flex-1 overflow-y-auto overflow-x-hidden px-4 py-4 sm:px-6"
      >
        {loading && <p className="text-center text-sm text-faint">Загружаем историю…</p>}
        {!loading && loadingMore && (
          <p className="pb-2 text-center text-xs text-faint">Загружаем более ранние сообщения…</p>
        )}
        {!loading && messages.length === 0 && (
          <p className="mt-10 text-center text-sm text-faint">Начните переписку — здесь пока пусто.</p>
        )}
        {!loading &&
          messages.map((m) => (
            <MessageBubble
              key={m.messageId}
              message={m}
              chatId={chatId}
              isOwn={m.senderId === ownUserId}
              ownUserId={ownUserId}
              peerName={peerProfile?.name ?? 'Собеседник'}
              onEdit={handleEdit}
              onDelete={handleDelete}
              onReply={handleReply}
              onToggleReaction={handleToggleReaction}
            />
          ))}
        <div ref={bottomRef} />
      </div>

      <div className="relative shrink-0 border-t border-line px-3 pt-2.5 pb-[calc(0.625rem_+_env(safe-area-inset-bottom))] sm:px-4">
        {showEmojiPicker && (
          <EmojiPicker onSelect={insertEmoji} onClose={() => setShowEmojiPicker(false)} />
        )}
        {sendError && <p className="mb-2 text-xs text-warn">{sendError}</p>}

        {replyingTo && (
          <div className="mb-2 flex items-start gap-2 rounded-lg border border-line bg-surface px-3 py-2">
            <div className="min-w-0 flex-1 border-l-2 border-signal pl-2">
              <p className="font-mono text-[10px] text-signal">{replyAuthor}</p>
              <p className="truncate text-xs text-muted">{replyText}</p>
            </div>
            <button
              type="button"
              onClick={() => setReplyingTo(null)}
              aria-label="Отменить ответ"
              className="shrink-0 rounded-md p-1 text-faint hover:bg-raised hover:text-ink"
            >
              ✕
            </button>
          </div>
        )}

        <div className="flex items-center gap-2 overflow-hidden rounded-xl border border-line bg-surface p-1.5 transition-colors focus-within:border-signal">
          <button
            type="button"
            onClick={() => setShowEmojiPicker((v) => !v)}
            aria-label="Эмодзи"
            className={`flex h-[38px] w-[38px] shrink-0 items-center justify-center rounded-lg text-lg leading-none hover:bg-raised ${
              showEmojiPicker ? 'bg-raised' : ''
            }`}
          >
            🙂
          </button>
          <textarea
            ref={textareaRef}
            rows={1}
            value={draft}
            onChange={(e) => {
              setDraft(e.target.value)
              presence.notifyTyping()
            }}
            onKeyDown={handleKeyDown}
            placeholder="Написать сообщение…"
            className="max-h-40 min-h-[38px] min-w-0 flex-1 appearance-none resize-none border-none bg-transparent px-1 py-2 text-sm leading-relaxed text-ink placeholder:text-faint outline-none ring-0 focus:outline-none focus:ring-0"
          />
          <button
            onClick={handleSend}
            disabled={sending || !draft.trim()}
            className="flex h-[38px] shrink-0 items-center justify-center rounded-lg bg-signal px-4 text-sm font-medium text-void transition-opacity hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-40"
          >
            Отправить
          </button>
        </div>
      </div>
    </div>
  )
}
