import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { AuthShell } from '../components/AuthShell'
import { Field } from '../components/Field'
import { useAuth } from '../context/AuthContext'
import { login as loginRequest } from '../api/auth'

export default function LoginPage() {
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const { setTokens } = useAuth()
  const navigate = useNavigate()

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setSubmitting(true)
    try {
      const tokens = await loginRequest(login, password)
      await setTokens(tokens)
      navigate('/', { replace: true })
    } catch (err: any) {
      setError(err?.response?.data?.message || 'Не удалось войти. Проверьте логин и пароль.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <AuthShell
      eyebrow="Вход"
      title="С возвращением"
      subtitle="Введите данные учётной записи, чтобы продолжить переписку."
      footer={
        <span>
          Нет аккаунта?{' '}
          <Link to="/register" className="text-signal hover:underline">
            Зарегистрироваться
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
          autoComplete="current-password"
          required
        />
        {error && <p className="mt-4 rounded-lg bg-warnDim px-3 py-2 text-sm text-warn">{error}</p>}
        <button
          type="submit"
          disabled={submitting}
          className="mt-6 w-full rounded-lg bg-signal py-2.5 text-sm font-medium text-void transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          {submitting ? 'Входим…' : 'Войти'}
        </button>
      </form>
    </AuthShell>
  )
}
