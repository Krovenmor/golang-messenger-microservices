import type { AuthTokens } from '../types'

const ACCESS_KEY = 'signaline.accessToken'
const REFRESH_KEY = 'signaline.refreshToken'

const RENAMED_FROM: Record<string, string> = {
  [ACCESS_KEY]: 'signalline.accessToken',
  [REFRESH_KEY]: 'signalline.refreshToken',
}

function readAdoptingPreviousName(key: string): string | null {
  const current = localStorage.getItem(key)
  if (current) return current

  const previous = localStorage.getItem(RENAMED_FROM[key])
  if (!previous) return null

  localStorage.setItem(key, previous)
  localStorage.removeItem(RENAMED_FROM[key])
  return previous
}

export const tokenStorage = {
  get(): AuthTokens | null {
    const accessToken = readAdoptingPreviousName(ACCESS_KEY)
    const refreshToken = readAdoptingPreviousName(REFRESH_KEY)
    if (!accessToken || !refreshToken) return null
    return { accessToken, refreshToken }
  },
  set(tokens: AuthTokens) {
    localStorage.setItem(ACCESS_KEY, tokens.accessToken)
    localStorage.setItem(REFRESH_KEY, tokens.refreshToken)
  },
  clear() {
    localStorage.removeItem(ACCESS_KEY)
    localStorage.removeItem(REFRESH_KEY)
    localStorage.removeItem(RENAMED_FROM[ACCESS_KEY])
    localStorage.removeItem(RENAMED_FROM[REFRESH_KEY])
  },
}
