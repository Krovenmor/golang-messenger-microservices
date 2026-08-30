import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { AuthShell } from '../components/AuthShell'
import { Field } from '../components/Field'
import { unlockPrivateKey, hasRealCrypto } from '../utils/crypto'
import { cryptoStore } from '../api/cryptoStore'
import * as keyVault from '../utils/keyVault'
import { useAuth } from '../context/AuthContext'

export default function UnlockProfilePage() {
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const { profile, logout } = useAuth()
  const navigate = useNavigate()

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)

    if (!profile) {
      setError('Профиль ещё не загружен, попробуйте снова через секунду.')
      return
    }

    if (!hasRealCrypto()) {
      setError('Нужен HTTPS (или localhost) для расшифровки ключа в браузере.')
      return
    }

    setSubmitting(true)
    try {
      const privateKey = await unlockPrivateKey(password, profile.encryptedPrvKey, profile.kdfSalt, profile.keyNonce)
      cryptoStore.setKeys(privateKey, profile.pubKey)
      // Remember it on this device (non-extractable — see keyVault.ts) so
      // this screen doesn't show up again next time, only after logging
      // out and back in, or on a different device/browser.
      await keyVault.saveKey(profile.userId, privateKey)
      navigate('/', { replace: true })
    } catch {
      setError('Неверный пароль — не удалось расшифровать ключ.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <AuthShell
      eyebrow="Разблокировка"
      title="Введите пароль от ключа"
      subtitle="Он нужен, чтобы расшифровать вашу переписку в этом сеансе — сервер его не хранит и восстановить не может. Это устройство запомнит ключ, так что вводить пароль снова не придётся."
      footer={
        <button onClick={logout} className="text-signal hover:underline">
          Выйти из другого аккаунта
        </button>
      }
    >
      <form onSubmit={handleSubmit}>
        <Field
          id="unlock-password"
          label="Пароль от ключа шифрования"
          type="password"
          autoFocus
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        {error && <p className="mt-4 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{error}</p>}
        <button
          type="submit"
          disabled={submitting}
          className="mt-6 w-full rounded-lg bg-signal py-2.5 text-sm font-medium text-void transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          {submitting ? 'Расшифровываем…' : 'Разблокировать'}
        </button>
      </form>
    </AuthShell>
  )
}
