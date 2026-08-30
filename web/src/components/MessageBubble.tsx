import { useEffect, useRef, useState, type MouseEvent } from 'react'
import type { ChatMember, Message } from '../types'
import { fingerprintColor, formatTime } from '../utils/identity'
import { useDecryptedText } from '../utils/useDecryptedText'
import { useReplyPreview } from '../utils/useReplyPreview'
import { usePublicProfile } from '../api/profileCache'
import { MessageActionMenu } from './MessageActionMenu'

interface Props {
  message: Message
  chatId: string
  isOwn: boolean
  ownUserId: string
  peer: ChatMember | null
  onEdit: (messageId: string, text: string) => Promise<void>
  onDelete: (messageId: string) => Promise<void>
  onReply: (message: Message) => void
  onToggleReaction: (messageId: string, emoji: string) => void
}

type MenuAnchor = { x: number; y: number; mode: 'point' | 'rect'; rectRight?: number; rectBottom?: number }

function ReplyQuote({ chatId, replyToId, ownUserId, peer }: { chatId: string; replyToId: string; ownUserId: string; peer: ChatMember | null }) {
  const { target, decrypted } = useReplyPreview(chatId, replyToId, ownUserId)
  const peerProfile = usePublicProfile(peer?.userId ?? null)

  const authorLabel = !target ? '' : target.senderId === ownUserId ? 'Вы' : peerProfile?.name ?? 'Собеседник'
  const text =
    decrypted.status === 'ok'
      ? decrypted.text
      : decrypted.status === 'locked'
        ? '🔒 сообщение заблокировано'
        : decrypted.status === 'error'
          ? 'не удалось расшифровать'
          : '…'

  return (
    <div className="mb-1.5 rounded-md border-l-2 border-signal/50 bg-void/40 px-2 py-1">
      {target && <p className="font-mono text-[10px] text-signal/80">{authorLabel}</p>}
      <p className="truncate text-xs text-muted">{target ? text : 'сообщение недоступно'}</p>
    </div>
  )
}

export function MessageBubble({ message, chatId, isOwn, ownUserId, peer, onEdit, onDelete, onReply, onToggleReaction }: Props) {
  const decrypted = useDecryptedText(message, ownUserId, chatId)
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState('')
  const [busy, setBusy] = useState(false)
  const [menuAnchor, setMenuAnchor] = useState<MenuAnchor | null>(null)
  const triggerRef = useRef<HTMLButtonElement>(null)
  const color = fingerprintColor(message.senderId)

  // Seed the edit textarea from the decrypted text once it's available —
  // there's nothing sensible to edit before that (the raw field is ciphertext).
  useEffect(() => {
    if (decrypted.status === 'ok') setDraft(decrypted.text)
  }, [decrypted])

  if (message.deletedAt) {
    return (
      <div className={`flex ${isOwn ? 'justify-end' : 'justify-start'}`}>
        <div className="my-1 max-w-[85%] sm:max-w-[70%] rounded-xl border border-dashed border-line px-4 py-2 text-xs italic text-faint">
          сообщение удалено
        </div>
      </div>
    )
  }

  async function submitEdit() {
    if (decrypted.status !== 'ok' || !draft.trim() || draft === decrypted.text) {
      setEditing(false)
      return
    }
    setBusy(true)
    try {
      await onEdit(message.messageId, draft.trim())
      setEditing(false)
    } finally {
      setBusy(false)
    }
  }

  function openMenuFromTrigger() {
    const rect = triggerRef.current?.getBoundingClientRect()
    if (!rect) return
    setMenuAnchor({ x: rect.left, y: rect.top, mode: 'rect', rectRight: rect.right, rectBottom: rect.bottom })
  }

  function openMenuAtCursor(e: MouseEvent) {
    if (decrypted.status !== 'ok') return
    e.preventDefault()
    setMenuAnchor({ x: e.clientX, y: e.clientY, mode: 'point' })
  }

  return (
    <div className={`group flex ${isOwn ? 'justify-end' : 'justify-start'}`}>
      <div className="my-1 flex max-w-[85%] sm:max-w-[70%] flex-col gap-1">
        <div
          onContextMenu={openMenuAtCursor}
          className={`relative rounded-xl border-l-2 px-4 py-2.5 text-sm leading-relaxed shadow-sm ${
            isOwn ? 'bg-raised text-ink' : 'bg-surface text-ink'
          }`}
          style={{ borderLeftColor: color }}
        >
          {message.replyToId && (
            <ReplyQuote chatId={chatId} replyToId={message.replyToId} ownUserId={ownUserId} peer={peer} />
          )}

          {editing ? (
            <div className="flex flex-col gap-2">
              <textarea
                autoFocus
                className="w-full resize-none rounded-lg border border-line bg-void px-2.5 py-1.5 text-sm text-ink outline-none focus:border-signal"
                rows={2}
                value={draft}
                onChange={(e) => setDraft(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault()
                    submitEdit()
                  }
                  if (e.key === 'Escape') setEditing(false)
                }}
              />
              <div className="flex justify-end gap-2 text-xs">
                <button className="text-faint hover:text-ink" onClick={() => setEditing(false)}>
                  отмена
                </button>
                <button className="text-signal hover:underline" disabled={busy} onClick={submitEdit}>
                  сохранить
                </button>
              </div>
            </div>
          ) : (
            <p className="whitespace-pre-wrap break-words">
              {decrypted.status === 'ok' && decrypted.text}
              {decrypted.status === 'loading' && <span className="opacity-50">…</span>}
              {decrypted.status === 'locked' && <span className="italic text-faint">🔒 расшифровка заблокирована</span>}
              {decrypted.status === 'error' && <span className="italic text-warn">не удалось расшифровать</span>}
            </p>
          )}

          {!editing && decrypted.status === 'ok' && (
            <button
              ref={triggerRef}
              type="button"
              onClick={openMenuFromTrigger}
              aria-label="Действия с сообщением"
              className="absolute -top-2.5 right-2 hidden h-6 w-6 items-center justify-center rounded-md border border-line bg-raised text-muted opacity-0 transition-opacity hover:border-signal/50 hover:text-signal group-hover:flex group-hover:opacity-100"
            >
              <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor">
                <circle cx="5" cy="12" r="2" />
                <circle cx="12" cy="12" r="2" />
                <circle cx="19" cy="12" r="2" />
              </svg>
            </button>
          )}
        </div>

        {!!message.reactions?.length && (
          <div className={`flex flex-wrap gap-1 ${isOwn ? 'justify-end' : 'justify-start'}`}>
            {message.reactions.map((r) => {
              const mine = r.users.includes(ownUserId)
              return (
                <button
                  key={r.emoji}
                  type="button"
                  onClick={() => onToggleReaction(message.messageId, r.emoji)}
                  className={`flex items-center gap-1 rounded-full border px-2 py-0.5 text-xs transition-colors ${
                    mine
                      ? 'border-signal bg-signal/10 text-signal'
                      : 'border-line bg-raised text-muted hover:border-signal/50'
                  }`}
                >
                  <span>{r.emoji}</span>
                  <span className="font-mono text-[10px]">{r.users.length}</span>
                </button>
              )
            })}
          </div>
        )}
        <span
          style={{ fontFamily: 'system-ui, -apple-system, "Segoe UI", Roboto, Arial, sans-serif' }}
          className={`tabular-nums text-[10px] text-faint ${isOwn ? 'text-right' : 'text-left'}`}
        >
          {formatTime(message.createdAt)}
          {message.redactedAt && ' · изменено'}
        </span>
      </div>

      {menuAnchor && (
        <MessageActionMenu
          anchor={menuAnchor}
          isOwn={isOwn}
          onClose={() => setMenuAnchor(null)}
          onReact={(emoji) => onToggleReaction(message.messageId, emoji)}
          onReply={() => onReply(message)}
          onEdit={() => setEditing(true)}
          onDelete={() => onDelete(message.messageId)}
        />
      )}
    </div>
  )
}
