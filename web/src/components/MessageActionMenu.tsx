import { useEffect, useLayoutEffect, useRef, useState, type ReactNode } from 'react'
import { createPortal } from 'react-dom'
import { getSupportedReactions } from '../api/reactions'

const MARGIN = 8

interface Anchor {
  x: number
  y: number
  /** 'point' for a right-click cursor position, 'rect' to hang below a trigger element. */
  mode: 'point' | 'rect'
  rectRight?: number
  rectBottom?: number
}

interface Props {
  anchor: Anchor
  isOwn: boolean
  onClose: () => void
  onReact: (emoji: string) => void
  onReply: () => void
  onEdit: () => void
  onDelete: () => void
}

function ReplyIcon() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M9 14L4 9l5-5" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M4 9h10.5a5.5 5.5 0 015.5 5.5v1" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function EditIcon() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path
        d="M12 20h9M16.5 3.5a2.12 2.12 0 013 3L7 19l-4 1 1-4L16.5 3.5z"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function TrashIcon() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path
        d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2m3 0l-1 14a2 2 0 01-2 2H7a2 2 0 01-2-2L4 6h16z"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function MenuRow({
  icon,
  label,
  onClick,
  danger,
}: {
  icon: ReactNode
  label: string
  onClick: () => void
  danger?: boolean
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`flex w-full items-center gap-2.5 px-3.5 py-2.5 text-left text-sm transition-colors ${
        danger ? 'text-warn hover:bg-warnDim' : 'text-ink hover:bg-signal/10 hover:text-signal'
      }`}
    >
      {icon}
      {label}
    </button>
  )
}

export function MessageActionMenu({ anchor, isOwn, onClose, onReact, onReply, onEdit, onDelete }: Props) {
  const [emojis, setEmojis] = useState<string[]>([])
  const panelRef = useRef<HTMLDivElement>(null)
  const [style, setStyle] = useState<{ left: number; top: number; opacity: number }>({
    left: anchor.x,
    top: anchor.y,
    opacity: 0,
  })

  useEffect(() => {
    getSupportedReactions()
      .then(setEmojis)
      .catch(() => setEmojis([]))
  }, [])

  // Measure the panel once it's actually in the DOM, then clamp it inside
  // the viewport — this is what "floating above everything, never clipped
  // by a scrolling message list" actually requires; you can't know the
  // right position until you know the panel's real size.
  useLayoutEffect(() => {
    const el = panelRef.current
    if (!el) return
    const rect = el.getBoundingClientRect()

    let left = anchor.mode === 'rect' ? (anchor.rectRight ?? anchor.x) - rect.width : anchor.x
    let top = anchor.mode === 'rect' ? (anchor.rectBottom ?? anchor.y) + 6 : anchor.y

    left = Math.min(Math.max(MARGIN, left), window.innerWidth - rect.width - MARGIN)
    top = Math.min(Math.max(MARGIN, top), window.innerHeight - rect.height - MARGIN)

    setStyle({ left, top, opacity: 1 })
  }, [anchor])

  useEffect(() => {
    function handlePointerDown(e: MouseEvent) {
      if (panelRef.current && !panelRef.current.contains(e.target as Node)) onClose()
    }
    function handleKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    // Telegram-style: a menu anchored to a specific message shouldn't
    // linger while the list scrolls out from under it.
    function handleScroll() {
      onClose()
    }
    document.addEventListener('mousedown', handlePointerDown)
    document.addEventListener('keydown', handleKey)
    document.addEventListener('scroll', handleScroll, true)
    return () => {
      document.removeEventListener('mousedown', handlePointerDown)
      document.removeEventListener('keydown', handleKey)
      document.removeEventListener('scroll', handleScroll, true)
    }
  }, [onClose])

  return createPortal(
    <div
      ref={panelRef}
      style={{ left: style.left, top: style.top, opacity: style.opacity }}
      className="fixed z-50 w-60 overflow-hidden rounded-lg border border-signal/30 bg-surface shadow-[0_0_0_1px_rgba(94,234,212,0.08),0_12px_40px_-8px_rgba(94,234,212,0.25),0_8px_24px_-4px_rgba(0,0,0,0.5)] transition-opacity duration-100"
    >
      <div className="h-[2px] w-full bg-gradient-to-r from-signal/0 via-signal to-signal/0" />

      {emojis.length > 0 && (
        <div className="flex flex-wrap gap-0.5 border-b border-line p-2">
          {emojis.slice(0, 16).map((emoji) => (
            <button
              key={emoji}
              type="button"
              onClick={() => {
                onReact(emoji)
                onClose()
              }}
              className="flex h-8 w-8 items-center justify-center rounded-md text-lg leading-none transition-transform hover:scale-110 hover:bg-raised"
            >
              {emoji}
            </button>
          ))}
        </div>
      )}

      <div className="py-1">
        <MenuRow
          icon={<ReplyIcon />}
          label="Ответить"
          onClick={() => {
            onReply()
            onClose()
          }}
        />
        {isOwn && (
          <>
            <MenuRow
              icon={<EditIcon />}
              label="Изменить"
              onClick={() => {
                onEdit()
                onClose()
              }}
            />
            <MenuRow
              icon={<TrashIcon />}
              label="Удалить"
              danger
              onClick={() => {
                onDelete()
                onClose()
              }}
            />
          </>
        )}
      </div>
    </div>,
    document.body
  )
}
