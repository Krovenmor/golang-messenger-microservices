import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { getAvatarUrl } from '../api/media'
import { fingerprintColor, formatExactDateTime, initials } from '../utils/identity'
import { isStatusActive, statusLabel } from '../utils/statusLabel'
import type { PublicProfile, StatusInfo } from '../types'

interface Props {
  profile: PublicProfile
  status: StatusInfo | null
  onClose: () => void
}

function AvatarPreview({ avatarId, profile }: { avatarId: string | null; profile: PublicProfile }) {
  const [errored, setErrored] = useState(false)

  useEffect(() => {
    setErrored(false)
  }, [avatarId])

  if (!avatarId || errored) {
    return (
      <div
        className="flex aspect-square w-full items-center justify-center rounded-2xl font-display text-6xl font-semibold text-void"
        style={{ backgroundColor: fingerprintColor(profile.userId) }}
      >
        {initials(profile.name)}
      </div>
    )
  }

  return (
    <img
      src={getAvatarUrl(avatarId)}
      alt={profile.name}
      onError={() => setErrored(true)}
      className="aspect-square w-full rounded-2xl bg-raised object-cover"
    />
  )
}

function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-baseline justify-between gap-3 border-t border-line py-2.5">
      <span className="shrink-0 text-xs text-faint">{label}</span>
      <span className="min-w-0 text-right text-xs text-muted">{value}</span>
    </div>
  )
}

export function PeerProfileModal({ profile, status, onClose }: Props) {
  const avatars = profile.additional?.avatars ?? []
  const [shownAvatar, setShownAvatar] = useState<string | null>(avatars[0] ?? null)

  useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', handleKey)
    return () => document.removeEventListener('keydown', handleKey)
  }, [onClose])

  const bio = profile.additional?.bio?.trim()
  const lastSeenAt = status?.lastSeen ? new Date(status.lastSeen * 1000) : null
  const isOnlineNow = isStatusActive(status)

  return createPortal(
    <div
      onClick={onClose}
      className="fixed left-0 right-0 top-[var(--app-offset-top,0px)] z-40 flex h-[var(--app-height,100vh)] items-start justify-center overflow-y-auto bg-black/70 px-4 py-6 sm:items-center"
    >
      <div
        onClick={(e) => e.stopPropagation()}
        className="my-auto w-full max-w-sm animate-rise rounded-2xl border border-line bg-surface p-5 shadow-2xl shadow-black/50"
      >
        <div className="mb-4 flex items-start justify-between gap-3">
          <p className="font-mono text-xs uppercase tracking-[0.2em] text-signal">Профиль</p>
          <button
            type="button"
            onClick={onClose}
            aria-label="Закрыть"
            className="-mr-1 -mt-1 shrink-0 rounded-md p-1 text-faint hover:bg-raised hover:text-ink"
          >
            ✕
          </button>
        </div>

        <AvatarPreview avatarId={shownAvatar} profile={profile} />

        {avatars.length > 1 && (
          <div className="mt-2 flex gap-2 overflow-x-auto pb-1">
            {avatars.map((id) => (
              <button
                key={id}
                type="button"
                onClick={() => setShownAvatar(id)}
                aria-label="Показать аватар"
                className={`h-12 w-12 shrink-0 overflow-hidden rounded-lg border transition-colors ${
                  id === shownAvatar ? 'border-signal' : 'border-line hover:border-signal/50'
                }`}
              >
                <img src={getAvatarUrl(id)} alt="" className="h-full w-full object-cover" />
              </button>
            ))}
          </div>
        )}

        <div className="mt-4">
          <p className="break-words font-display text-xl font-semibold text-ink">{profile.name}</p>
          <p className="mt-0.5 break-all font-mono text-sm text-signal">@{profile.userName}</p>
          <p className={`mt-1.5 text-xs ${isStatusActive(status) ? 'text-signal' : 'text-faint'}`}>
            {statusLabel(status) || 'статус неизвестен'}
          </p>
        </div>

        {bio ? (
          <p className="mt-4 whitespace-pre-wrap break-words rounded-lg bg-raised px-3 py-2.5 text-sm text-muted">
            {bio}
          </p>
        ) : (
          <p className="mt-4 rounded-lg border border-dashed border-line px-3 py-2.5 text-sm italic text-faint">
            Описание не заполнено
          </p>
        )}

        <div className="mt-4">
          <InfoRow
            label={isOnlineNow ? 'В сети с' : 'Был(а) в сети'}
            value={lastSeenAt ? formatExactDateTime(lastSeenAt) : 'неизвестно'}
          />
          <InfoRow label="Регистрация" value={formatExactDateTime(new Date(profile.createdAt)) || 'неизвестно'} />
        </div>
      </div>
    </div>,
    document.body
  )
}
