import { useState } from 'react'
import type { ChatSummary, OwnProfile } from '../types'
import { formatDay } from '../utils/identity'
import { useDecryptedText } from '../utils/useDecryptedText'
import { usePublicProfile } from '../api/profileCache'
import { SignalDot } from './SignalDot'
import { Avatar } from './Avatar'

interface Props {
  chats: ChatSummary[]
  activeChatId: string | null
  ownUserId: string | null
  profile: OwnProfile | null
  connected: boolean
  onSelect: (chatId: string) => void
  onNewChat: () => void
  onLogout: () => void
  onOpenSettings: () => void
}

function peerOf(chat: ChatSummary, ownUserId: string | null) {
  return chat.members.find((m) => m.userId !== ownUserId) ?? chat.members[0] ?? null
}

function ChatRow({
  chat,
  active,
  ownUserId,
  onSelect,
}: {
  chat: ChatSummary
  active: boolean
  ownUserId: string | null
  onSelect: (chatId: string) => void
}) {
  const peer = peerOf(chat, ownUserId)
  const peerProfile = usePublicProfile(peer?.userId ?? null)
  const last = chat.lastMessage
  // Called unconditionally regardless of whether `last` exists/is deleted —
  // the hook itself already handles null/deleted messages internally.
  const decrypted = useDecryptedText(last, ownUserId, chat.chatId)
  const isOwn = !!ownUserId && last?.senderId === ownUserId

  let preview: string
  if (!last) {
    preview = 'нет сообщений'
  } else if (last.deletedAt) {
    preview = 'сообщение удалено'
  } else {
    const prefix = isOwn ? 'Вы: ' : ''
    switch (decrypted.status) {
      case 'ok':
        preview = prefix + decrypted.text
        break
      case 'locked':
        preview = '🔒 сообщение заблокировано'
        break
      case 'error':
        preview = '⚠️ не удалось расшифровать'
        break
      default:
        preview = prefix + '…'
    }
  }

  return (
    <button
      type="button"
      onClick={() => onSelect(chat.chatId)}
      className={`mb-1 flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left transition-colors ${
        active ? 'bg-raised' : 'hover:bg-raised/60'
      }`}
    >
      <Avatar
        avatarId={peerProfile?.additional?.avatars?.[0]}
        seed={peer?.userId ?? chat.chatId}
        name={peerProfile?.name ?? '? ?'}
        className="h-10 w-10 text-sm"
      />
      <div className="min-w-0 flex-1">
        <div className="flex items-center justify-between gap-2">
          <p className="truncate text-sm font-medium text-ink">{peerProfile?.name ?? 'без имени'}</p>
          <span className="shrink-0 font-mono text-[10px] text-faint">
            {formatDay(chat.lastMessage?.createdAt ?? null)}
          </span>
        </div>
        <p className="truncate text-xs text-muted">{preview}</p>
      </div>
    </button>
  )
}

export function Sidebar({
  chats,
  activeChatId,
  ownUserId,
  profile,
  connected,
  onSelect,
  onNewChat,
  onLogout,
  onOpenSettings,
}: Props) {
  const [filter, setFilter] = useState('')

  const normalizedFilter = filter.trim().toLowerCase()
  const filtered = chats.filter((chat) => {
    if (!normalizedFilter) return true
    const peer = peerOf(chat, ownUserId)
    const searchable = [peer?.userId ?? '', chat.chatId].join(' ').toLowerCase()
    return searchable.includes(normalizedFilter)
  })

  return (
    <aside className="flex h-full w-full shrink-0 flex-col border-r border-line bg-surface sm:max-w-[320px]">
      <div className="flex items-center justify-between border-b border-line px-4 py-4">
        <div className="flex items-center gap-2">
          <span className="font-display text-base font-semibold tracking-tight text-ink">Signalline</span>
        </div>
        <SignalDot connected={connected} />
      </div>

      <div className="px-4 py-3">
        <div className="flex gap-2">
          <input
            className="w-full rounded-lg border border-line bg-raised px-3 py-2 text-sm text-ink placeholder:text-faint outline-none focus:border-signal"
            placeholder="Поиск по чатам…"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
          />
          <button
            onClick={onNewChat}
            title="Новый чат"
            className="shrink-0 rounded-lg border border-line bg-raised px-3 text-signal hover:border-signal"
          >
            +
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-2 pb-2">
        {filtered.length === 0 && (
          <p className="px-3 py-6 text-center text-sm text-faint">Пока нет ни одного диалога.</p>
        )}
        {filtered.map((chat) => (
          <ChatRow key={chat.chatId} chat={chat} active={chat.chatId === activeChatId} ownUserId={ownUserId} onSelect={onSelect} />
        ))}
      </div>

      {profile && (
        <div className="flex items-center gap-2 border-t border-line px-4 py-3">
          <button
            type="button"
            onClick={onOpenSettings}
            className="flex min-w-0 flex-1 items-center gap-3 rounded-lg py-1 text-left hover:bg-raised"
          >
            <Avatar
              avatarId={profile.additional?.avatars?.[0]}
              seed={profile.userId}
              name={profile.name}
              className="h-8 w-8 text-xs"
            />
            <div className="min-w-0 flex-1">
              <p className="truncate text-xs font-medium text-ink">{profile.name}</p>
              <p className="truncate font-mono text-[10px] text-faint">@{profile.userName}</p>
            </div>
          </button>
          <button onClick={onLogout} className="shrink-0 font-mono text-[10px] text-faint hover:text-warn">
            выйти
          </button>
        </div>
      )}
      <p className="select-none px-4 py-1 text-center font-mono text-[9px] text-faint/40">v{__APP_VERSION__}</p>
    </aside>
  )
}
