import type { AuthTokens } from '../types'

const ACCESS_KEY = 'signalline.accessToken'
const REFRESH_KEY = 'signalline.refreshToken'

export const tokenStorage = {
  get(): AuthTokens | null {
    const accessToken = localStorage.getItem(ACCESS_KEY)
    const refreshToken = localStorage.getItem(REFRESH_KEY)
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
  },
}
