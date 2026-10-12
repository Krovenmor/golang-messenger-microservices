import type { ApiError, FieldIssue, FieldIssues } from '../api/errors'

function plural(n: number, one: string, few: string, many: string): string {
  const mod10 = n % 10
  const mod100 = n % 100
  if (mod10 === 1 && mod100 !== 11) return one
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) return few
  return many
}

function symbols(param: string | undefined): string {
  const n = Number(param)
  if (!Number.isFinite(n)) return 'символов'
  return `${n} ${plural(n, 'символ', 'символа', 'символов')}`
}

const fieldSpecificMessages: Record<string, string> = {
  'login.taken': 'Логин уже занят',
  'email.taken': 'Эта почта уже зарегистрирована',
  'code.invalid': 'Неверный или просроченный код',
  'code.len': 'Код — это 6 цифр из письма',
  'confirm.mismatch': 'Пароли не совпадают',
}

export function describeFieldIssue(field: string, issue: FieldIssue): string {
  const specific = fieldSpecificMessages[`${field}.${issue.reason}`]
  if (specific) return specific

  switch (issue.reason) {
    case 'required':
      return 'Обязательное поле'
    case 'notblank':
      return 'Не может состоять из одних пробелов'
    case 'min':
      return `Минимум ${symbols(issue.param)}`
    case 'max':
      return `Максимум ${symbols(issue.param)}`
    case 'len':
      return `Должно быть ровно ${symbols(issue.param)}`
    case 'email':
      return 'Некорректный адрес почты'
    case 'alphanum':
      return 'Только латинские буквы и цифры'
    case 'taken':
      return 'Уже занято'
    case 'invalid':
      return 'Неверное значение'
    default:
      return 'Некорректное значение'
  }
}

export type FieldMessages = Record<string, string>

export function describeFieldIssues(fields: FieldIssues): FieldMessages {
  const messages: FieldMessages = {}
  for (const [field, issue] of Object.entries(fields)) {
    messages[field] = describeFieldIssue(field, issue)
  }
  return messages
}

const serverMessages: Record<string, string> = {
  'invalid login or password': 'Неверный логин или пароль',
  'account is banned': 'Аккаунт временно заблокирован',
  'verification code already sent': 'Код на эту почту уже отправлен — повторите чуть позже',
  'registration not found or expired': 'Время на подтверждение истекло — начните регистрацию заново',
  'invalid refresh token': 'Сессия истекла, войдите заново',
  'internal error': 'Ошибка на сервере, попробуйте позже',
}

export function describeFormError(error: ApiError, fallback: string): string | null {
  if (error.isNetworkError) return 'Сервер недоступен, проверьте подключение'
  if (Object.keys(error.fields).length > 0) return null
  if (error.message && serverMessages[error.message]) return serverMessages[error.message]
  return fallback
}
