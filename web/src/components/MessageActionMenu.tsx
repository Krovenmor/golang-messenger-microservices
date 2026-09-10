import { useEffect, useLayoutEffect, useRef, useState, type ReactNode } from 'react'
import { createPortal } from 'react-dom'
import { getSupportedReactions } from '../api/reactions'

const VIEWPORT_MARGIN = 8
const ANCHOR_GAP = 6
const MAX_REACTIONS = 16

interface Anchor {
  x: number
  y: number
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
  onCopy: () => void
  onEdit: () => void
  onDelete: () => void
}

interface Placement {
  left: number
  top: number
  maxHeight: number
  visible: boolean
}

function clamp(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), Math.max(min, max))
}

function visibleViewport() {
  const visual = window.visualViewport
  return {
    left: visual?.offsetLeft ?? 0,
    top: visual?.offsetTop ?? 0,
    width: visual?.width ?? window.innerWidth,
    height: visual?.height ?? window.innerHeight,
  }
}

function place(anchor: Anchor, contentHeight: number, panelWidth: number): Placement {
  const view = visibleViewport()
  const maxHeight = Math.min(contentHeight, view.height - VIEWPORT_MARGIN * 2)

  const anchorTop = anchor.y
  const anchorBottom = anchor.mode === 'rect' ? anchor.rectBottom ?? anchor.y : anchor.y

  const roomBelow = view.top + view.height - VIEWPORT_MARGIN - (anchorBottom + ANCHOR_GAP)
  const roomAbove = anchorTop - ANCHOR_GAP - (view.top + VIEWPORT_MARGIN)
  const opensDownward = maxHeight <= roomBelow || roomBelow >= roomAbove

  const preferredTop = opensDownward ? anchorBottom + ANCHOR_GAP : anchorTop - ANCHOR_GAP - maxHeight
  const preferredLeft = anchor.mode === 'rect' ? (anchor.rectRight ?? anchor.x) - panelWidth : anchor.x

  return {
    left: clamp(preferredLeft, view.left + VIEWPORT_MARGIN, view.left + view.width - panelWidth - VIEWPORT_MARGIN),
    top: clamp(preferredTop, view.top + VIEWPORT_MARGIN, view.top + view.height - maxHeight - VIEWPORT_MARGIN),
    maxHeight,
    visible: true,
  }
}

function ReplyIcon() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M9 14L4 9l5-5" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M4 9h10.5a5.5 5.5 0 015.5 5.5v1" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function CopyIcon() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <rect x="9" y="9" width="12" height="12" rx="2" />
      <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" strokeLinecap="round" strokeLinejoin="round" />
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
      className={`flex w-full items-center gap-2.5 px-3.5 py-3 text-left text-sm transition-colors sm:py-2.5 ${
        danger ? 'text-warn hover:bg-warnDim' : 'text-ink hover:bg-signal/10 hover:text-signal'
      }`}
    >
      {icon}
      {label}
    </button>
  )
}

export function MessageActionMenu({ anchor, isOwn, onClose, onReact, onReply, onCopy, onEdit, onDelete }: Props) {
  const [emojis, setEmojis] = useState<string[]>([])
  const canCopy = typeof navigator.clipboard?.writeText === 'function'
  const panelRef = useRef<HTMLDivElement>(null)
  const [placement, setPlacement] = useState<Placement>({
    left: anchor.x,
    top: anchor.y,
    maxHeight: 0,
    visible: false,
  })

  useEffect(() => {
    getSupportedReactions()
      .then(setEmojis)
      .catch(() => setEmojis([]))
  }, [])

  useLayoutEffect(() => {
    const reposition = () => {
      const el = panelRef.current
      if (!el) return
      setPlacement(place(anchor, el.scrollHeight, el.offsetWidth))
    }

    reposition()

    const visual = window.visualViewport
    visual?.addEventListener('resize', reposition)
    window.addEventListener('resize', reposition)
    return () => {
      visual?.removeEventListener('resize', reposition)
      window.removeEventListener('resize', reposition)
    }
  }, [anchor, emojis, isOwn, canCopy])

  useEffect(() => {
    function isOutside(target: EventTarget | null) {
      return !!panelRef.current && !panelRef.current.contains(target as Node)
    }
    function handlePointerDown(e: PointerEvent) {
      if (isOutside(e.target)) onClose()
    }
    function handleKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    function handleScroll(e: Event) {
      if (isOutside(e.target)) onClose()
    }
    document.addEventListener('pointerdown', handlePointerDown)
    document.addEventListener('keydown', handleKey)
    document.addEventListener('scroll', handleScroll, true)
    return () => {
      document.removeEventListener('pointerdown', handlePointerDown)
      document.removeEventListener('keydown', handleKey)
      document.removeEventListener('scroll', handleScroll, true)
    }
  }, [onClose])

  return createPortal(
    <div
      ref={panelRef}
      style={{
        left: placement.left,
        top: placement.top,
        maxHeight: placement.maxHeight || undefined,
        opacity: placement.visible ? 1 : 0,
      }}
      className="fixed z-50 w-60 max-w-[calc(100vw_-_1rem)] overflow-y-auto overscroll-contain rounded-lg border border-signal/30 bg-surface shadow-[0_0_0_1px_rgba(94,234,212,0.08),0_12px_40px_-8px_rgba(94,234,212,0.25),0_8px_24px_-4px_rgba(0,0,0,0.5)] transition-opacity duration-100"
    >
      <div className="h-[2px] w-full bg-gradient-to-r from-signal/0 via-signal to-signal/0" />

      {emojis.length > 0 && (
        <div className="flex flex-wrap gap-0.5 border-b border-line p-2">
          {emojis.slice(0, MAX_REACTIONS).map((emoji) => (
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
        {canCopy && (
          <MenuRow
            icon={<CopyIcon />}
            label="Копировать"
            onClick={() => {
              onCopy()
              onClose()
            }}
          />
        )}
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
