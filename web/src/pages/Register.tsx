import { useEffect, useRef, useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { AuthShell } from '../components/AuthShell'
import { Field } from '../components/Field'
import { register, sendCode } from '../api/auth'
import { login as loginRequest } from '../api/auth'
import { useAuth } from '../context/AuthContext'

export default function RegisterPage() {
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [email, setEmail] = useState('')
  const [code, setCode] = useState('')

  const [codeSent, setCodeSent] = useState(false)
  const [sendingCode, setSendingCode] = useState(false)
  const [cooldown, setCooldown] = useState(0)
  const [codeExpiresIn, setCodeExpiresIn] = useState<number | null>(null)
  const [codeError, setCodeError] = useState<string | null>(null)

  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const { setTokens } = useAuth()
  const navigate = useNavigate()

  const cooldownTimer = useRef<number | null>(null)

  useEffect(() => {
    return () => {
      if (cooldownTimer.current) window.clearInterval(cooldownTimer.current)
    }
  }, [])

  async function handleSendCode() {
    if (!email.trim() || sendingCode || cooldown > 0) return
    setCodeError(null)
    setSendingCode(true)
    try {
      const res = await sendCode(email.trim())
      setCodeSent(true)
      setCodeExpiresIn(res.expiresIn)
      setCooldown(res.retryAfter)
      if (cooldownTimer.current) window.clearInterval(cooldownTimer.current)
      cooldownTimer.current = window.setInterval(() => {
        setCooldown((s) => {
          if (s <= 1) {
            if (cooldownTimer.current) window.clearInterval(cooldownTimer.current)
            return 0
          }
          return s - 1
        })
      }, 1000)
    } catch (err: any) {
      setCodeError(err?.response?.data?.message || 'Не удалось отправить код. Проверьте адрес почты.')
    } finally {
      setSendingCode(false)
    }
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)

    if (password !== confirm) {
      setError('Пароли не совпадают.')
      return
    }
    if (password.length < 8) {
      setError('Пароль должен быть не короче 8 символов.')
      return
    }
    if (!/^\d{6}$/.test(code.trim())) {
      setError('Код — это 6 цифр из письма.')
      return
    }

    setSubmitting(true)
    try {
      await register(login, password, email.trim(), code.trim())
      const tokens = await loginRequest(login, password)
      await setTokens(tokens)
      navigate('/setup', { replace: true, state: { login } })
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Не удалось зарегистрироваться. Проверьте код и данные аккаунта.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <AuthShell
      eyebrow="Регистрация"
      title="Создать аккаунт"
      subtitle="Понадобится подтвердить почту, а после — ещё один шаг: публичный профиль, прежде чем начать переписку."
      footer={
        <span>
          Уже есть аккаунт?{' '}
          <Link to="/login" className="text-signal hover:underline">
            Войти
          </Link>
        </span>
      }
    >
      <form onSubmit={handleSubmit}>
        <Field
          id="login"
          label="Логин"
          value={login}
          onChange={(e) => setLogin(e.target.value)}
          autoComplete="username"
          required
        />
        <Field
          id="password"
          label="Пароль"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          autoComplete="new-password"
          required
        />
        <Field
          id="confirm"
          label="Повторите пароль"
          type="password"
          value={confirm}
          onChange={(e) => setConfirm(e.target.value)}
          autoComplete="new-password"
          required
        />

        <label className="mb-4 block last:mb-0" htmlFor="email">
          <span className="mb-1.5 block font-mono text-xs uppercase tracking-wide text-faint">Почта</span>
          <div className="flex gap-2">
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              autoComplete="email"
              required
              className="w-full rounded-lg border border-line bg-raised px-3.5 py-2.5 text-sm text-ink placeholder:text-faint outline-none transition-colors focus:border-signal"
            />
            <button
              type="button"
              onClick={handleSendCode}
              disabled={!email.trim() || sendingCode || cooldown > 0}
              className="shrink-0 whitespace-nowrap rounded-lg border border-line bg-raised px-3 text-xs text-ink hover:border-signal disabled:cursor-not-allowed disabled:opacity-50"
            >
              {sendingCode ? '…' : cooldown > 0 ? `${cooldown}с` : codeSent ? 'ещё раз' : 'получить код'}
            </button>
          </div>
        </label>

        {codeSent && (
          <p className="-mt-3 mb-4 text-xs text-signal">
            Код отправлен на {email.trim()} — проверьте почту.
            {codeExpiresIn != null && ` Действителен ${Math.round(codeExpiresIn / 60)} мин.`}
          </p>
        )}
        {codeError && <p className="mt-2 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{codeError}</p>}

        <Field
          id="code"
          label="Код из письма"
          value={code}
          onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
          autoComplete="one-time-code"
          inputMode="numeric"
          pattern="\d{6}"
          maxLength={6}
          placeholder="000000"
          disabled={!codeSent}
          required
        />

        {error && <p className="mt-4 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{error}</p>}
        <button
          type="submit"
          disabled={submitting}
          className="mt-6 w-full rounded-lg bg-signal py-2.5 text-sm font-medium text-void transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          {submitting ? 'Создаём…' : 'Зарегистрироваться'}
        </button>
      </form>
    </AuthShell>
  )
}
