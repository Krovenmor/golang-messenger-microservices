import { useEffect, useMemo, useRef, useState } from 'react'
import { EMOJI_CATEGORIES } from '../data/emojis'

const RECENT_KEY = 'signalline.recentEmojis'
const MAX_RECENT = 32

function loadRecent(): string[] {
  try {
    const raw = localStorage.getItem(RECENT_KEY)
    return raw ? JSON.parse(raw) : []
  } catch {
    return []
  }
}

function saveRecent(list: string[]) {
  try {
    localStorage.setItem(RECENT_KEY, JSON.stringify(list))
  } catch {
    // ignore — not critical if it doesn't persist
  }
}

interface Props {
  onSelect: (emoji: string) => void
  onClose: () => void
}

export function EmojiPicker({ onSelect, onClose }: Props) {
  const [query, setQuery] = useState('')
  const [recent, setRecent] = useState<string[]>(() => loadRecent())
  const panelRef = useRef<HTMLDivElement>(null)
  const sectionRefs = useRef<Record<string, HTMLDivElement | null>>({})

  useEffect(() => {
    function handlePointerDown(e: MouseEvent) {
      if (panelRef.current && !panelRef.current.contains(e.target as Node)) {
        onClose()
      }
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

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return null
    const results: string[] = []
    for (const category of EMOJI_CATEGORIES) {
      for (const item of category.emojis) {
        if (item.keywords.some((k) => k.includes(q))) results.push(item.char)
      }
    }
    return results
  }, [query])

  function handlePick(emoji: string) {
    onSelect(emoji)
    setRecent((prev) => {
      const next = [emoji, ...prev.filter((e) => e !== emoji)].slice(0, MAX_RECENT)
      saveRecent(next)
      return next
    })
  }

  function scrollToSection(id: string) {
    sectionRefs.current[id]?.scrollIntoView({ block: 'start' })
  }

  return (
    <div
      ref={panelRef}
      className="absolute bottom-full left-0 z-30 mb-2 flex h-[22rem] w-[19rem] max-w-[calc(100vw-2rem)] animate-rise flex-col overflow-hidden rounded-2xl border border-line bg-surface shadow-2xl shadow-black/40"
    >
      <div className="border-b border-line p-2">
        <input
          autoFocus
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Поиск эмодзи…"
          className="w-full rounded-lg border border-line bg-raised px-3 py-1.5 text-sm text-ink placeholder:text-faint outline-none focus:border-signal"
        />
      </div>

      <div className="flex-1 overflow-y-auto px-2 py-1.5">
        {filtered ? (
          filtered.length > 0 ? (
            <div className="grid grid-cols-7 gap-0.5">
              {filtered.map((emoji, i) => (
                <button
                  key={`${emoji}-${i}`}
                  type="button"
                  onClick={() => handlePick(emoji)}
                  className="flex h-9 w-9 items-center justify-center rounded-lg text-xl leading-none hover:bg-raised"
                >
                  {emoji}
                </button>
              ))}
            </div>
          ) : (
            <p className="mt-10 text-center text-sm text-faint">Ничего не найдено</p>
          )
        ) : (
          <>
            {recent.length > 0 && (
              <div ref={(el) => { sectionRefs.current.recent = el }}>
                <p className="sticky top-0 z-10 bg-surface px-1 py-1 font-mono text-[10px] uppercase tracking-wide text-faint">
                  Недавние
                </p>
                <div className="grid grid-cols-7 gap-0.5 pb-2">
                  {recent.map((emoji, i) => (
                    <button
                      key={`recent-${emoji}-${i}`}
                      type="button"
                      onClick={() => handlePick(emoji)}
                      className="flex h-9 w-9 items-center justify-center rounded-lg text-xl leading-none hover:bg-raised"
                    >
                      {emoji}
                    </button>
                  ))}
                </div>
              </div>
            )}
            {EMOJI_CATEGORIES.map((category) => (
              <div key={category.id} ref={(el) => { sectionRefs.current[category.id] = el }}>
                <p className="sticky top-0 z-10 bg-surface px-1 py-1 font-mono text-[10px] uppercase tracking-wide text-faint">
                  {category.label}
                </p>
                <div className="grid grid-cols-7 gap-0.5 pb-2">
                  {category.emojis.map((item, i) => (
                    <button
                      key={`${category.id}-${i}`}
                      type="button"
                      onClick={() => handlePick(item.char)}
                      title={item.keywords[0]}
                      className="flex h-9 w-9 items-center justify-center rounded-lg text-xl leading-none hover:bg-raised"
                    >
                      {item.char}
                    </button>
                  ))}
                </div>
              </div>
            ))}
          </>
        )}
      </div>

      <div className="flex items-center justify-around border-t border-line px-1 py-1">
        {recent.length > 0 && (
          <button
            type="button"
            onClick={() => scrollToSection('recent')}
            title="Недавние"
            className="rounded-lg p-1.5 text-base leading-none text-muted hover:bg-raised"
          >
            🕘
          </button>
        )}
        {EMOJI_CATEGORIES.map((category) => (
          <button
            key={category.id}
            type="button"
            onClick={() => scrollToSection(category.id)}
            title={category.label}
            className="rounded-lg p-1.5 text-base leading-none hover:bg-raised"
          >
            {category.icon}
          </button>
        ))}
      </div>
    </div>
  )
}
