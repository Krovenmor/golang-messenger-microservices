import { useEffect, useRef, useState, type ChangeEvent } from 'react'
import { useAuth } from '../context/AuthContext'
import { Avatar } from './Avatar'
import { getMediaProfileInfo, uploadAvatarFile, deleteAvatarFile } from '../api/media'
import { addProfileAvatar, removeProfileAvatar, updateProfileBio, updateProfileName } from '../api/profile'
import { AvatarBudgetExceededError, UnreadableImageError, prepareAvatarFile } from '../utils/imageResize'
import type { MediaStorageInfo } from '../types'

const BIO_MAX_LENGTH = 200

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
  const [error, setError] = useState<string | null>(null)

  const [nameDraft, setNameDraft] = useState(profile?.name ?? '')
  const [savingName, setSavingName] = useState(false)
  const [nameError, setNameError] = useState<string | null>(null)

  const [bioDraft, setBioDraft] = useState(profile?.additional?.bio ?? '')
  const [savingBio, setSavingBio] = useState(false)
  const [bioError, setBioError] = useState<string | null>(null)

  const fileInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    getMediaProfileInfo()
      .then(setStorage)
      .catch(() => {})
  }, [])

  if (!profile) return null

  const avatars = profile.additional?.avatars ?? []
  const savedBio = profile.additional?.bio ?? ''

  const remainingSpace = storage ? Math.max(0, storage.maxSpace - storage.spaceFilled) : undefined

  async function refreshStorage() {
    try {
      setStorage(await getMediaProfileInfo())
    } catch {}
  }

  async function handleFileChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    e.target.value = ''
    if (!file) return

    if (!file.type.startsWith('image/')) {
      setError('Нужен файл изображения.')
      return
    }

    setError(null)
    setUploading(true)
    try {
      const resized = await prepareAvatarFile(file, remainingSpace)
      const { photoId } = await uploadAvatarFile(resized)
      await addProfileAvatar(photoId)
      await refreshProfile()
      await refreshStorage()
    } catch (err) {
      if (err instanceof AvatarBudgetExceededError) {
        setError(
          `Не хватает места: аватар займёт ${formatBytes(err.encodedBytes)}, ` +
            `а свободно ${formatBytes(err.budgetBytes)}. Удалите старый аватар.`
        )
      } else if (err instanceof UnreadableImageError) {
        setError('Не удалось прочитать изображение — попробуйте PNG или JPEG.')
      } else {
        setError('Не удалось загрузить аватар.')
      }
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

  async function handleSaveName() {
    const trimmed = nameDraft.trim()
    if (!trimmed || trimmed === profile!.name) return
    setSavingName(true)
    setNameError(null)
    try {
      await updateProfileName(trimmed)
      await refreshProfile()
    } catch {
      setNameError('Не удалось сохранить имя.')
    } finally {
      setSavingName(false)
    }
  }

  async function handleSaveBio() {
    const trimmed = bioDraft.trim()
    if (trimmed === savedBio) return
    setSavingBio(true)
    setBioError(null)
    try {
      await updateProfileBio(trimmed)
      await refreshProfile()
    } catch {
      setBioError('Не удалось сохранить описание.')
    } finally {
      setSavingBio(false)
    }
  }

  const usagePct =
    storage && storage.maxSpace > 0 ? Math.min(100, (storage.spaceFilled / storage.maxSpace) * 100) : 0
  const nameChanged = nameDraft.trim() && nameDraft.trim() !== profile.name
  const bioChanged = bioDraft.trim() !== savedBio

  return (
    <div
      className="absolute inset-0 z-20 flex items-start justify-center overflow-y-auto bg-black/60 px-4 py-6 sm:items-center"
      onClick={onClose}
    >
      <div
        className="my-auto w-full max-w-sm animate-rise rounded-2xl border border-line bg-surface p-6 shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        <p className="mb-1 font-mono text-xs uppercase tracking-[0.2em] text-signal">Профиль</p>
        <h2 className="mb-4 font-display text-lg font-semibold text-ink">Настройки</h2>

        <div className="mb-5 flex items-center gap-4">
          <Avatar avatarId={avatars[0]} seed={profile.userId} name={profile.name} className="h-16 w-16 text-xl" />
          <p className="truncate font-mono text-xs text-faint">@{profile.userName}</p>
        </div>

        <label className="mb-1.5 block font-mono text-xs uppercase tracking-wide text-faint" htmlFor="display-name">
          Отображаемое имя
        </label>
        <div className="mb-1 flex gap-2">
          <input
            id="display-name"
            value={nameDraft}
            onChange={(e) => setNameDraft(e.target.value)}
            className="w-full rounded-lg border border-line bg-raised px-3 py-2 text-sm text-ink outline-none transition-colors focus:border-signal"
          />
          {nameChanged && (
            <button
              type="button"
              onClick={handleSaveName}
              disabled={savingName}
              className="shrink-0 rounded-lg bg-signal px-3 text-xs font-medium text-void hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {savingName ? '…' : 'ОК'}
            </button>
          )}
        </div>
        {nameError && <p className="mb-3 text-xs text-warn">{nameError}</p>}

        <label className="mb-1.5 mt-3 block font-mono text-xs uppercase tracking-wide text-faint" htmlFor="bio">
          О себе
        </label>
        <textarea
          id="bio"
          rows={2}
          maxLength={BIO_MAX_LENGTH}
          value={bioDraft}
          onChange={(e) => setBioDraft(e.target.value)}
          placeholder="Пара слов о себе…"
          className="mb-1 w-full resize-none rounded-lg border border-line bg-raised px-3 py-2 text-sm text-ink outline-none transition-colors focus:border-signal"
        />
        <div className="mb-1 flex items-center justify-between">
          <span className="font-mono text-[10px] text-faint">
            {bioDraft.length}/{BIO_MAX_LENGTH}
          </span>
          {bioChanged && (
            <button
              type="button"
              onClick={handleSaveBio}
              disabled={savingBio}
              className="rounded-lg bg-signal px-3 py-1 text-xs font-medium text-void hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {savingBio ? 'Сохраняем…' : 'Сохранить'}
            </button>
          )}
        </div>
        {bioError && <p className="mb-3 text-xs text-warn">{bioError}</p>}

        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={handleFileChange}
          className="hidden"
        />
        <button
          type="button"
          onClick={() => fileInputRef.current?.click()}
          disabled={uploading}
          className="mb-2 mt-3 w-full rounded-lg border border-line bg-raised py-2 text-sm text-ink hover:border-signal disabled:cursor-not-allowed disabled:opacity-50"
        >
          {uploading ? 'Загружаем…' : 'Загрузить аватар'}
        </button>

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
