import { useEffect, useRef, useState, type ChangeEvent } from 'react'
import { useAuth } from '../context/AuthContext'
import { Avatar } from './Avatar'
import { getMediaProfileInfo, uploadAvatarFile, deleteAvatarFile } from '../api/media'
import { addProfileAvatar, removeProfileAvatar, updateProfileBio, updateProfileName } from '../api/profile'
import { prepareAvatarFile } from '../utils/imageResize'
import type { MediaStorageInfo } from '../types'

interface Props {
  onClose: () => void
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} Б`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} КБ`
  return `${(bytes / (1024 * 1024)).toFixed(1)} МБ`
}

export function ProfileSettingsModal({ onClose }: Props) {
  const { profile, refreshProfile } = useAuth()
  const [storage, setStorage] = useState<MediaStorageInfo | null>(null)
  const [uploading, setUploading] = useState(false)
  const [removingId, setRemovingId] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [name, setName] = useState(profile?.name ?? '')
  const [bio, setBio] = useState(profile?.additional?.bio ?? '')
  const fileInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (profile) {
      setName(profile.name)
      setBio(profile.additional?.bio ?? '')
    }
  }, [profile])

  useEffect(() => {
    getMediaProfileInfo()
      .then(setStorage)
      .catch(() => {})
  }, [])

  if (!profile) return null

  const avatars = profile.additional?.avatars ?? []

  async function refreshStorage() {
    try {
      setStorage(await getMediaProfileInfo())
    } catch {
      // non-critical — just leave the last known figure showing
    }
  }

  async function handleFileChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    e.target.value = '' // lets the person pick the same file again later
    if (!file) return

    if (!['image/png', 'image/jpeg'].includes(file.type)) {
      setError('Поддерживаются только PNG и JPEG.')
      return
    }

    setError(null)
    setUploading(true)
    try {
      // Crops to a centered square and downscales to fit the backend's
      // 600x600 limit — the person never has to think about image size.
      const resized = await prepareAvatarFile(file)
      const { photoId } = await uploadAvatarFile(resized)
      await addProfileAvatar(photoId)
      await refreshProfile()
      await refreshStorage()
    } catch {
      setError('Не удалось загрузить аватар.')
    } finally {
      setUploading(false)
    }
  }

  async function handleRemove(photoId: string) {
    setRemovingId(photoId)
    setError(null)
    try {
      await removeProfileAvatar(photoId)
      await deleteAvatarFile(photoId)
      await refreshProfile()
      await refreshStorage()
    } catch {
      setError('Не удалось удалить аватар.')
    } finally {
      setRemovingId(null)
    }
  }

  async function handleSave() {
    const trimmedName = name.trim()
    const nextBio = bio.trim()
    if (!trimmedName) {
      setError('Имя не может быть пустым.')
      return
    }

    setError(null)
    setSaving(true)
    try {
      const hasChangedName = trimmedName !== profile.name
      const hasChangedBio = nextBio !== (profile.additional?.bio ?? '')

      if (hasChangedName) {
        await updateProfileName(trimmedName)
      }
      if (hasChangedBio) {
        await updateProfileBio(nextBio)
      }

      await refreshProfile()
    } catch {
      setError('Не удалось сохранить изменения.')
    } finally {
      setSaving(false)
    }
  }

  const usagePct =
    storage && storage.maxSpace > 0 ? Math.min(100, (storage.spaceFilled / storage.maxSpace) * 100) : 0

  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/60 px-4" onClick={onClose}>
      <div
        className="w-full max-w-sm animate-rise rounded-2xl border border-line bg-surface p-6 shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        <p className="mb-1 font-mono text-xs uppercase tracking-[0.2em] text-signal">Профиль</p>
        <h2 className="mb-4 font-display text-lg font-semibold text-ink">Настройки</h2>

        <div className="mb-5 flex items-center gap-4">
          <Avatar avatarId={avatars[0]} seed={profile.userId} name={profile.name} className="h-16 w-16 text-xl" />
          <div className="min-w-0">
            <p className="truncate text-sm font-medium text-ink">{profile.name}</p>
            <p className="truncate font-mono text-xs text-faint">@{profile.userName}</p>
          </div>
        </div>

        <input
          ref={fileInputRef}
          type="file"
          accept="image/png,image/jpeg"
          onChange={handleFileChange}
          className="hidden"
        />
        <button
          type="button"
          onClick={() => fileInputRef.current?.click()}
          disabled={uploading}
          className="mb-2 w-full rounded-lg border border-line bg-raised py-2 text-sm text-ink hover:border-signal disabled:cursor-not-allowed disabled:opacity-50"
        >
          {uploading ? 'Загружаем…' : 'Загрузить аватар'}
        </button>

        <div className="mb-4 space-y-3">
          <div>
            <label className="mb-1.5 block text-[10px] font-medium uppercase tracking-[0.18em] text-faint">Имя</label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full rounded-lg border border-line bg-raised px-3 py-2 text-sm text-ink outline-none focus:border-signal"
              placeholder="Ваше публичное имя"
            />
          </div>

          <div>
            <label className="mb-1.5 block text-[10px] font-medium uppercase tracking-[0.18em] text-faint">О себе</label>
            <textarea
              value={bio}
              onChange={(e) => setBio(e.target.value)}
              rows={3}
              maxLength={200}
              className="w-full resize-none rounded-lg border border-line bg-raised px-3 py-2 text-sm text-ink outline-none focus:border-signal"
              placeholder="Напишите коротко о себе"
            />
          </div>

          <button
            type="button"
            onClick={handleSave}
            disabled={saving}
            className="w-full rounded-lg bg-signal px-3 py-2 text-sm font-medium text-void hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {saving ? 'Сохраняем…' : 'Сохранить'}
          </button>
        </div>

        {avatars.length > 0 && (
          <>
            <div className="mb-1.5 mt-3 flex flex-wrap gap-2">
              {avatars.map((id) => (
                <div key={id} className="relative">
                  <Avatar avatarId={id} seed={id} name={profile.name} className="h-12 w-12 text-xs" />
                  <button
                    type="button"
                    onClick={() => handleRemove(id)}
                    disabled={removingId === id}
                    aria-label="Удалить аватар"
                    className="absolute -right-1 -top-1 flex h-5 w-5 items-center justify-center rounded-full border border-line bg-void text-[10px] text-faint hover:border-warn hover:text-warn disabled:opacity-50"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
            <p className="mb-4 text-xs text-faint">Первый в списке — тот, что виден остальным.</p>
          </>
        )}

        {error && <p className="mb-4 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{error}</p>}

        {storage && (
          <div className="mb-5">
            <div className="mb-1.5 flex items-center justify-between font-mono text-[10px] text-faint">
              <span>Хранилище</span>
              <span>
                {formatBytes(storage.spaceFilled)} / {formatBytes(storage.maxSpace)}
              </span>
            </div>
            <div className="h-1.5 overflow-hidden rounded-full bg-raised">
              <div className="h-full rounded-full bg-signal transition-all" style={{ width: `${usagePct}%` }} />
            </div>
          </div>
        )}

        <button onClick={onClose} className="w-full text-center text-xs text-faint hover:text-muted">
          Закрыть
        </button>
      </div>
    </div>
  )
}
