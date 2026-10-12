import { useEffect, useRef, useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { AuthShell } from '../components/AuthShell'
import { Field } from '../components/Field'
import { confirmRegistration, login as loginRequest, resendCode, startRegistration } from '../api/auth'
import { parseApiError, type ApiError, type FieldIssues } from '../api/errors'
import { useAuth } from '../context/AuthContext'
import { describeFormError } from '../utils/errorMessages'
import { findEmptyFields, useFieldErrors } from '../utils/useFieldErrors'

const CODE_LENGTH = 6
const REGISTRATION_EXPIRED = 'registration not found or expired'
const DETAILS_FIELDS = ['login', 'password', 'email']

type Step = 'details' | 'code'

function hasDetailsIssues(error: ApiError): boolean {
  return DETAILS_FIELDS.some((field) => field in error.fields)
}

export default function RegisterPage() {
  const [step, setStep] = useState<Step>('details')

  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [email, setEmail] = useState('')
  const [code, setCode] = useState('')

  const [registrationId, setRegistrationId] = useState<string | null>(null)
  const [codeExpiresIn, setCodeExpiresIn] = useState<number | null>(null)
  const [cooldown, setCooldown] = useState(0)
  const [resending, setResending] = useState(false)

  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const { fieldErrors, showIssues, clearFieldError, clearAll } = useFieldErrors()
  const { setTokens } = useAuth()
  const navigate = useNavigate()

  const cooldownTimer = useRef<number | null>(null)

  useEffect(() => {
    return () => {
      if (cooldownTimer.current) window.clearInterval(cooldownTimer.current)
    }
  }, [])

  function startCooldown(seconds: number) {
    setCooldown(seconds)
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
  }

  function backToDetails(apiError?: ApiError) {
    setStep('details')
    setCode('')
    setRegistrationId(null)
    if (apiError) {
      showIssues(apiError.fields)
      setError(describeFormError(apiError, 'Не удалось зарегистрироваться. Проверьте данные аккаунта.'))
    } else {
      clearAll()
      setError(null)
    }
  }

  function handleCodeStepError(apiError: ApiError, fallback: string) {
    if (apiError.message === REGISTRATION_EXPIRED || hasDetailsIssues(apiError)) {
      backToDetails(apiError)
      return
    }
    showIssues(apiError.fields)
    setError(describeFormError(apiError, fallback))
  }

  function validateDetails(): FieldIssues {
    const issues = findEmptyFields({ login, password, confirm, email })
    if (!issues.password && !issues.confirm && password !== confirm) {
      issues.confirm = { reason: 'mismatch' }
    }
    return issues
  }

  async function handleDetailsSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)

    const issues = validateDetails()
    showIssues(issues)
    if (Object.keys(issues).length > 0) return

    setSubmitting(true)
    try {
      const res = await startRegistration(login, password, email.trim())
      setRegistrationId(res.registrationId)
      setCodeExpiresIn(res.expiresIn)
      startCooldown(res.retryAfter)
      setStep('code')
    } catch (err) {
      const apiError = parseApiError(err)
      showIssues(apiError.fields)
      setError(describeFormError(apiError, 'Не удалось начать регистрацию. Проверьте данные аккаунта.'))
    } finally {
      setSubmitting(false)
    }
  }

  async function handleResend() {
    if (!registrationId || resending || cooldown > 0) return
    setError(null)
    clearFieldError('code')
    setResending(true)
    try {
      const res = await resendCode(registrationId)
      setCodeExpiresIn(res.expiresIn)
      startCooldown(res.retryAfter)
    } catch (err) {
      handleCodeStepError(parseApiError(err), 'Не удалось отправить код ещё раз.')
    } finally {
      setResending(false)
    }
  }

  async function handleCodeSubmit(e: FormEvent) {
    e.preventDefault()
    if (!registrationId) return
    setError(null)

    if (code.length !== CODE_LENGTH) {
      showIssues({ code: { reason: 'len', param: String(CODE_LENGTH) } })
      return
    }
    clearAll()

    setSubmitting(true)
    try {
      await confirmRegistration(registrationId, code)
    } catch (err) {
      handleCodeStepError(parseApiError(err), 'Не удалось подтвердить регистрацию.')
      setSubmitting(false)
      return
    }

    try {
      const tokens = await loginRequest(login, password)
      await setTokens(tokens)
      navigate('/setup', { replace: true, state: { login } })
    } catch (err) {
      setError(describeFormError(parseApiError(err), 'Аккаунт создан, но войти не удалось. Попробуйте на странице входа.'))
    } finally {
      setSubmitting(false)
    }
  }

  const footer = (
    <span>
      Уже есть аккаунт?{' '}
      <Link to="/login" className="text-signal hover:underline">
        Войти
      </Link>
    </span>
  )

  if (step === 'code') {
    return (
      <AuthShell
        eyebrow="Регистрация · шаг 2"
        title="Подтвердите почту"
        subtitle={`Мы отправили код на ${email.trim()}.${codeExpiresIn != null ? ` Он действителен ${Math.round(codeExpiresIn / 60)} мин.` : ''}`}
        footer={footer}
      >
        <form onSubmit={handleCodeSubmit} noValidate>
          <Field
            id="code"
            label="Код из письма"
            value={code}
            onChange={(e) => {
              setCode(e.target.value.replace(/\D/g, '').slice(0, CODE_LENGTH))
              clearFieldError('code')
            }}
            autoComplete="one-time-code"
            inputMode="numeric"
            maxLength={CODE_LENGTH}
            placeholder="000000"
            autoFocus
            error={fieldErrors.code}
            action={
              <button
                type="button"
                onClick={handleResend}
                disabled={resending || cooldown > 0}
                className="shrink-0 whitespace-nowrap rounded-lg border border-line bg-raised px-3 text-xs text-ink hover:border-signal disabled:cursor-not-allowed disabled:opacity-50"
              >
                {resending ? '…' : cooldown > 0 ? `${cooldown}с` : 'ещё раз'}
              </button>
            }
          />

          {error && <p className="mt-4 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{error}</p>}
          <button
            type="submit"
            disabled={submitting}
            className="mt-6 w-full rounded-lg bg-signal py-2.5 text-sm font-medium text-void transition-opacity hover:opacity-90 disabled:opacity-50"
          >
            {submitting ? 'Проверяем…' : 'Подтвердить'}
          </button>
          <button
            type="button"
            onClick={() => backToDetails()}
            disabled={submitting}
            className="mt-3 w-full text-xs text-faint hover:text-ink disabled:opacity-50"
          >
            Изменить данные
          </button>
        </form>
      </AuthShell>
    )
  }

  return (
    <AuthShell
      eyebrow="Регистрация · шаг 1"
      title="Создать аккаунт"
      subtitle="Понадобится подтвердить почту, а после — ещё один шаг: публичный профиль, прежде чем начать переписку."
      footer={footer}
    >
      <form onSubmit={handleDetailsSubmit} noValidate>
        <Field
          id="login"
          label="Логин"
          value={login}
          onChange={(e) => {
            setLogin(e.target.value)
            clearFieldError('login')
          }}
          autoComplete="username"
          error={fieldErrors.login}
        />
        <Field
          id="password"
          label="Пароль"
          type="password"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value)
            clearFieldError('password')
            clearFieldError('confirm')
          }}
          autoComplete="new-password"
          error={fieldErrors.password}
        />
        <Field
          id="confirm"
          label="Повторите пароль"
          type="password"
          value={confirm}
          onChange={(e) => {
            setConfirm(e.target.value)
            clearFieldError('confirm')
          }}
          autoComplete="new-password"
          error={fieldErrors.confirm}
        />
        <Field
          id="email"
          label="Почта"
          type="email"
          value={email}
          onChange={(e) => {
            setEmail(e.target.value)
            clearFieldError('email')
          }}
          autoComplete="email"
          error={fieldErrors.email}
        />

        {error && <p className="mt-4 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{error}</p>}
        <button
          type="submit"
          disabled={submitting}
          className="mt-6 w-full rounded-lg bg-signal py-2.5 text-sm font-medium text-void transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          {submitting ? 'Отправляем код…' : 'Получить код'}
        </button>
      </form>
    </AuthShell>
  )
}
