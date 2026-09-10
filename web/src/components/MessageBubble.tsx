import { useEffect, useRef, useState, type MouseEvent, type TouchEvent } from 'react'
import type { Message } from '../types'
import { fingerprintColor, formatTime } from '../utils/identity'
import { useDecryptedText } from '../utils/useDecryptedText'
import { useReplyPreview } from '../utils/useReplyPreview'
import { MessageActionMenu } from './MessageActionMenu'

interface Props {
  message: Message
  chatId: string
  isOwn: boolean
  ownUserId: string
  peerName: string
  onEdit: (messageId: string, text: string) => Promise<void>
  onDelete: (messageId: string) => Promise<void>
  onReply: (message: Message) => void
  onToggleReaction: (messageId: string, emoji: string) => void
}

type MenuAnchor = { x: number; y: number; mode: 'point' | 'rect'; rectRight?: number; rectBottom?: number }

const LONG_PRESS_MS = 450
const LONG_PRESS_MOVE_TOLERANCE_PX = 10
const SWIPE_AXIS_LOCK_PX = 12
const SWIPE_TRIGGER_PX = 56
const SWIPE_MAX_PX = 72
const SWIPE_RETURN_MS = 160

function ReplyIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M9 14L4 9l5-5" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M4 9h10.5a5.5 5.5 0 015.5 5.5v1" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function ReplyQuote({
  chatId,
  replyToId,
  ownUserId,
  peerName,
}: {
  chatId: string
  replyToId: string
  ownUserId: string
  peerName: string
}) {
  const { target, decrypted } = useReplyPreview(chatId, replyToId, ownUserId)

  const authorLabel = !target ? '' : target.senderId === ownUserId ? 'Вы' : peerName
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

export function MessageBubble({ message, chatId, isOwn, ownUserId, peerName, onEdit, onDelete, onReply, onToggleReaction }: Props) {
  const decrypted = useDecryptedText(message, ownUserId, chatId)
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState('')
  const [busy, setBusy] = useState(false)
  const [menuAnchor, setMenuAnchor] = useState<MenuAnchor | null>(null)
  const triggerRef = useRef<HTMLButtonElement>(null)
  const longPressTimer = useRef<number | null>(null)
  const longPressOrigin = useRef<{ x: number; y: number } | null>(null)
  const longPressOpenedMenu = useRef(false)
  const bubbleRef = useRef<HTMLDivElement>(null)
  const replyHintRef = useRef<HTMLSpanElement>(null)
  const swipeOrigin = useRef<{ x: number; y: number } | null>(null)
  const swipeAxis = useRef<'undecided' | 'horizontal' | 'vertical'>('undecided')
  const swipeOffset = useRef(0)
  const color = fingerprintColor(message.senderId)
  const swipeDirection = isOwn ? -1 : 1

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
    if (decrypted.status !== 'ok' || editing) return
    e.preventDefault()
    cancelLongPress()
    if (longPressOpenedMenu.current) return
    setMenuAnchor({ x: e.clientX, y: e.clientY, mode: 'point' })
  }

  function paintSwipe(offset: number) {
    const bubble = bubbleRef.current
    if (bubble) bubble.style.transform = offset === 0 ? '' : `translateX(${offset}px)`
    const hint = replyHintRef.current
    if (hint) hint.style.opacity = `${Math.min(1, Math.abs(offset) / SWIPE_TRIGGER_PX)}`
  }

  function resetSwipe(animated: boolean) {
    const bubble = bubbleRef.current
    if (bubble) bubble.style.transition = animated ? `transform ${SWIPE_RETURN_MS}ms ease-out` : 'none'
    paintSwipe(0)
    swipeOffset.current = 0
    swipeOrigin.current = null
    swipeAxis.current = 'undecided'
  }

  function cancelLongPress() {
    if (longPressTimer.current !== null) {
      window.clearTimeout(longPressTimer.current)
      longPressTimer.current = null
    }
    longPressOrigin.current = null
  }

  function handleTouchStart(e: TouchEvent<HTMLDivElement>) {
    if (decrypted.status !== 'ok' || editing) return
    const touch = e.touches[0]
    if (!touch) return

    cancelLongPress()
    longPressOpenedMenu.current = false
    longPressOrigin.current = { x: touch.clientX, y: touch.clientY }

    swipeOrigin.current = { x: touch.clientX, y: touch.clientY }
    swipeAxis.current = 'undecided'
    swipeOffset.current = 0
    if (bubbleRef.current) bubbleRef.current.style.transition = 'none'

    const { clientX, clientY } = touch
    longPressTimer.current = window.setTimeout(() => {
      longPressTimer.current = null
      longPressOpenedMenu.current = true
      navigator.vibrate?.(10)
      setMenuAnchor({ x: clientX, y: clientY, mode: 'point' })
    }, LONG_PRESS_MS)
  }

  function handleTouchMove(e: TouchEvent<HTMLDivElement>) {
    const touch = e.touches[0]
    if (!touch) return

    const pressOrigin = longPressOrigin.current
    if (
      pressOrigin &&
      Math.hypot(touch.clientX - pressOrigin.x, touch.clientY - pressOrigin.y) > LONG_PRESS_MOVE_TOLERANCE_PX
    ) {
      cancelLongPress()
    }

    const origin = swipeOrigin.current
    if (!origin) return

    const dx = touch.clientX - origin.x
    const dy = touch.clientY - origin.y

    if (swipeAxis.current === 'undecided') {
      if (Math.abs(dx) < SWIPE_AXIS_LOCK_PX && Math.abs(dy) < SWIPE_AXIS_LOCK_PX) return
      swipeAxis.current = Math.abs(dx) > Math.abs(dy) ? 'horizontal' : 'vertical'
    }
    if (swipeAxis.current !== 'horizontal') return

    const travel = Math.max(0, Math.min(dx * swipeDirection, SWIPE_MAX_PX))
    swipeOffset.current = travel * swipeDirection
    paintSwipe(swipeOffset.current)
  }

  function handleTouchEnd(e: TouchEvent<HTMLDivElement>) {
    cancelLongPress()

    const swipedFarEnough = Math.abs(swipeOffset.current) >= SWIPE_TRIGGER_PX
    resetSwipe(true)
    if (swipedFarEnough) {
      navigator.vibrate?.(10)
      onReply(message)
    }

    if (!longPressOpenedMenu.current) return
    longPressOpenedMenu.current = false
    e.preventDefault()
  }

  function handleTouchCancel() {
    cancelLongPress()
    resetSwipe(true)
  }

  function copyText() {
    if (decrypted.status !== 'ok') return
    navigator.clipboard.writeText(decrypted.text).catch(() => {})
  }

  return (
    <div className={`flex ${isOwn ? 'justify-end' : 'justify-start'}`}>
      <div className="relative my-1 flex max-w-[85%] sm:max-w-[70%] flex-col gap-1">
        <span
          ref={replyHintRef}
          aria-hidden
          style={{ opacity: 0, [isOwn ? 'right' : 'left']: 0 }}
          className="pointer-events-none absolute top-0 flex h-11 items-center text-signal"
        >
          <ReplyIcon />
        </span>
        <div
          ref={bubbleRef}
          onContextMenu={openMenuAtCursor}
          onTouchStart={handleTouchStart}
          onTouchMove={handleTouchMove}
          onTouchEnd={handleTouchEnd}
          onTouchCancel={handleTouchCancel}
          className={`group relative touch-pan-y rounded-xl border-l-2 px-4 py-2.5 text-sm leading-relaxed shadow-sm [@media(hover:none)]:select-none [@media(hover:none)]:[-webkit-touch-callout:none] ${
            isOwn ? 'bg-raised text-ink' : 'bg-surface text-ink'
          }`}
          style={{ borderLeftColor: color }}
        >
          {message.replyToId && (
            <ReplyQuote chatId={chatId} replyToId={message.replyToId} ownUserId={ownUserId} peerName={peerName} />
          )}

          {editing ? (
            <div className="flex flex-col gap-2">
              <textarea
                autoFocus
                className="w-full select-text resize-none rounded-lg border border-line bg-void px-2.5 py-1.5 text-sm text-ink outline-none focus:border-signal"
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
              className="absolute -top-2.5 right-2 hidden h-6 w-6 items-center justify-center rounded-md border border-line bg-raised text-muted opacity-0 transition-opacity hover:border-signal/50 hover:text-signal [@media(hover:hover)]:group-hover:flex [@media(hover:hover)]:group-hover:opacity-100"
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
          onCopy={copyText}
          onReply={() => onReply(message)}
          onEdit={() => setEditing(true)}
          onDelete={() => onDelete(message.messageId)}
        />
      )}
    </div>
  )
}
