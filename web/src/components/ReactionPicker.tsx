import { useEffect, useRef, useState } from 'react'
import { getSupportedReactions } from '../api/reactions'

interface Props {
  onSelect: (emoji: string) => void
  onClose: () => void
}

export function ReactionPicker({ onSelect, onClose }: Props) {
  const [emojis, setEmojis] = useState<string[]>([])
  const panelRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    getSupportedReactions()
      .then(setEmojis)
      .catch(() => setEmojis([]))
  }, [])

  useEffect(() => {
    function handlePointerDown(e: MouseEvent) {
      if (panelRef.current && !panelRef.current.contains(e.target as Node)) onClose()
    }
    function handleKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('mousedown', handlePointerDown)
    document.addEventListener('keydown', handleKey)
    return () => {
      document.removeEventListener('mousedown', handlePointerDown)
      document.removeEventListener('keydown', handleKey)
    }
  }, [onClose])

  return (
    <div
      ref={panelRef}
      className="absolute bottom-full z-30 mb-2 flex max-h-48 w-56 flex-wrap content-start gap-0.5 overflow-y-auto rounded-2xl border border-line bg-surface p-2 shadow-2xl shadow-black/40 animate-rise"
    >
      {emojis.length === 0 && <p className="w-full py-4 text-center text-xs text-faint">Загрузка…</p>}
      {emojis.map((emoji) => (
        <button
          key={emoji}
          type="button"
          onClick={() => onSelect(emoji)}
          className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-xl leading-none hover:bg-raised"
        >
          {emoji}
        </button>
      ))}
    </div>
  )
}
