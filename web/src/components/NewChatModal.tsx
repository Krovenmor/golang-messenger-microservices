import { useState } from 'react'
import { getProfile } from '../api/profile'
import { createChat } from '../api/messenger'
import { Avatar } from './Avatar'
import type { PublicProfile } from '../types'

interface Props {
  onClose: () => void
  onCreated: (chatId: string, peer: PublicProfile) => void
}

export function NewChatModal({ onClose, onCreated }: Props) {
  const [query, setQuery] = useState('')
  const [found, setFound] = useState<PublicProfile | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [searching, setSearching] = useState(false)
  const [creating, setCreating] = useState(false)

  async function handleSearch() {
    if (!query.trim()) return
    setSearching(true)
    setError(null)
    setFound(null)
    try {
      const profile = await getProfile(query.trim())
      setFound(profile)
    } catch {
      setError('Пользователь не найден.')
    } finally {
      setSearching(false)
    }
  }

  async function handleCreate() {
    if (!found) return
    setCreating(true)
    try {
      const { chatId } = await createChat(found.userId)
      onCreated(chatId, found)
    } catch {
      setError('Не удалось создать чат.')
    } finally {
      setCreating(false)
    }
  }

  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/60 px-4" onClick={onClose}>
      <div
        className="w-full max-w-sm animate-rise rounded-2xl border border-line bg-surface p-6 shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        <p className="mb-1 font-mono text-xs uppercase tracking-[0.2em] text-signal">Новый диалог</p>
        <h2 className="mb-4 font-display text-lg font-semibold text-ink">Найти собеседника</h2>

        <div className="flex gap-2">
          <input
            autoFocus
            className="w-full rounded-lg border border-line bg-raised px-3.5 py-2.5 text-sm text-ink placeholder:text-faint outline-none focus:border-signal"
            placeholder="юзернейм или UUID"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          />
          <button
            onClick={handleSearch}
            disabled={searching}
            className="shrink-0 rounded-lg border border-line bg-raised px-4 text-sm text-ink hover:border-signal disabled:opacity-50"
          >
            {searching ? '…' : 'Найти'}
          </button>
        </div>

        {error && <p className="mt-3 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{error}</p>}

        {found && (
          <div className="mt-4 flex items-center gap-3 rounded-xl border border-line bg-raised p-3">
            <Avatar
              avatarId={found.additional?.avatars?.[0]}
              seed={found.userId}
              name={found.name}
              className="h-10 w-10 text-sm"
            />
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-medium text-ink">{found.name}</p>
              <p className="truncate font-mono text-xs text-muted">@{found.userName}</p>
            </div>
            <button
              onClick={handleCreate}
              disabled={creating}
              className="shrink-0 rounded-lg bg-signal px-3 py-1.5 text-xs font-medium text-void hover:opacity-90 disabled:opacity-50"
            >
              {creating ? 'Создаём…' : 'Написать'}
            </button>
          </div>
        )}

        <button onClick={onClose} className="mt-5 w-full text-center text-xs text-faint hover:text-muted">
          Закрыть
        </button>
      </div>
    </div>
  )
}
