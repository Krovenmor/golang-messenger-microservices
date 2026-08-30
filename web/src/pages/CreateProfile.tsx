import { useEffect, useState, type FormEvent } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { AuthShell } from '../components/AuthShell'
import { Field } from '../components/Field'
import { createProfile } from '../api/profile'
import { getAccount } from '../api/auth'
import { generateProfileKeys, hasRealCrypto } from '../utils/crypto'
import { cryptoStore } from '../api/cryptoStore'
import * as keyVault from '../utils/keyVault'
import { useAuth } from '../context/AuthContext'
import { tokenStorage } from '../utils/tokenStorage'

export default function CreateProfilePage() {
  const location = useLocation()
  // Carried over from registration, so there's no network round-trip in the
  // common case. If this page is reached any other way (reload, or a
  // logged-in account that just doesn't have a profile yet) there's no
  // navigation state to read, so it's fetched from the account endpoint
  // instead — either way, the person is never asked to type a username
  // that's really just their login again.
  const loginFromState = (location.state as { login?: string } | null)?.login ?? null

  const [userName, setUserName] = useState<string | null>(loginFromState)
  const [userNameError, setUserNameError] = useState<string | null>(null)
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [status, setStatus] = useState<'idle' | 'keys' | 'saving'>('idle')
  const { refreshProfile } = useAuth()
  const navigate = useNavigate()

  const hasSession = !!tokenStorage.get()

  useEffect(() => {
    if (userName !== null) return
    let cancelled = false
    getAccount()
      .then((account) => {
        if (!cancelled) setUserName(account.login)
      })
      .catch(() => {
        if (!cancelled) setUserNameError('Не удалось получить логин аккаунта. Обновите страницу.')
      })
    return () => {
      cancelled = true
    }
  }, [userName])

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)

    if (!hasSession) {
      setError('Сессия истекла, войдите заново.')
      return
    }
    if (!userName) {
      setError('Логин аккаунта ещё не загружен, подождите секунду.')
      return
    }

    try {
      setStatus('keys')
      const { material, privateKey } = await generateProfileKeys(password)

      setStatus('saving')
      await createProfile({
        name,
        userName,
        pubKey: material.pubKey,
        encryptedPrvKey: material.encryptedPrvKey,
        kdfSalt: material.kdfSalt,
        keyNonce: material.keyNonce,
      })

      // The key was just generated fresh in this same flow — no need to
      // make the user re-enter the password to unlock it, unlike a
      // returning login (see UnlockProfile.tsx).
      if (privateKey) {
        cryptoStore.setKeys(privateKey, material.pubKey)
      }

      const own = await refreshProfile()
      if (privateKey && own) {
        await keyVault.saveKey(own.userId, privateKey)
      }
      navigate('/', { replace: true })
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Не удалось создать профиль.')
    } finally {
      setStatus('idle')
    }
  }

  const busy = status !== 'idle'

  return (
    <AuthShell
      eyebrow="Последний шаг"
      title="Профиль"
      subtitle="Отображаемое имя видно другим людям. Ключи шифрования создаются прямо в браузере."
    >
      <form onSubmit={handleSubmit}>
        {userNameError && <p className="mb-4 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{userNameError}</p>}

        <Field
          id="name"
          label="Отображаемое имя"
          placeholder="Ада Лавлейс"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
        <Field
          id="keypass"
          label="Пароль от ключа шифрования"
          type="password"
          placeholder="запомните — понадобится при следующем входе"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <p className="mt-2 text-xs text-faint">
          Этот пароль защищает ваш закрытый ключ сквозного шифрования. Он не отправляется на сервер — только
          результат шифрования. Без него сообщения нельзя будет прочитать на новом устройстве, восстановить
          пароль невозможно. На этом устройстве ключ запомнится — пароль спросят снова только после выхода
          из аккаунта или на другом устройстве.
        </p>
        {!hasRealCrypto() && (
          <p className="mt-2 rounded-lg bg-raised px-3 py-2 text-xs text-muted">
            Соединение не защищено (нет HTTPS), поэтому реальные ключи шифрования сгенерировать нельзя —
            будут отправлены случайные значения-заглушки, сообщения не будут по-настоящему зашифрованы.
          </p>
        )}
        {error && <p className="mt-4 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{error}</p>}
        <button
          type="submit"
          disabled={busy || !userName}
          className="mt-6 w-full rounded-lg bg-signal py-2.5 text-sm font-medium text-void transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          {status === 'keys' && 'Генерируем ключи…'}
          {status === 'saving' && 'Сохраняем профиль…'}
          {status === 'idle' && 'Создать профиль'}
        </button>
      </form>
    </AuthShell>
  )
}
