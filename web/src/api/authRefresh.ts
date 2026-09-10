import axios from 'axios'
import { tokenStorage } from '../utils/tokenStorage'
import { isJwtExpiringSoon } from '../utils/jwt'
import type { AuthTokens } from '../types'

export const API_BASE_URL: string = import.meta.env.VITE_API_URL || '/api'

export const TOKEN_REFRESHED_EVENT = 'signaline:token-refreshed'
export const SESSION_ENDED_EVENT = 'signaline:logout'

const CROSS_TAB_LOCK_NAME = 'signaline.token-refresh'
const REFRESH_MARGIN_SEC = 60

const SESSION_REJECTED_STATUSES = [400, 401, 403]

function isRefreshTokenRejected(error: unknown): boolean {
  if (!axios.isAxiosError(error)) return false
  const status = error.response?.status
  return status !== undefined && SESSION_REJECTED_STATUSES.includes(status)
}

function endSession() {
  tokenStorage.clear()
  window.dispatchEvent(new CustomEvent(SESSION_ENDED_EVENT))
}

function withCrossTabLock<T>(run: () => Promise<T>): Promise<T> {
  const locks = navigator.locks
  if (!locks) return run()
  return new Promise<T>((resolve, reject) => {
    locks.request(CROSS_TAB_LOCK_NAME, () => run().then(resolve, reject)).catch(reject)
  })
}

async function rotateTokens(rejectedAccessToken?: string): Promise<string | null> {
  const stored = tokenStorage.get()
  if (!stored?.refreshToken) {
    endSession()
    return null
  }

  const rotatedElsewhere = rejectedAccessToken
    ? stored.accessToken !== rejectedAccessToken
    : !isJwtExpiringSoon(stored.accessToken, REFRESH_MARGIN_SEC)
  if (rotatedElsewhere) return stored.accessToken

  try {
    const { data } = await axios.post<AuthTokens>(`${API_BASE_URL}/auth/update`, {
      refreshToken: stored.refreshToken,
    })
    tokenStorage.set(data)
    window.dispatchEvent(new CustomEvent(TOKEN_REFRESHED_EVENT))
    return data.accessToken
  } catch (error) {
    if (isRefreshTokenRejected(error)) endSession()
    return null
  }
}

let refreshInFlight: Promise<string | null> | null = null

export function ensureFreshToken(rejectedAccessToken?: string): Promise<string | null> {
  if (!refreshInFlight) {
    refreshInFlight = withCrossTabLock(() => rotateTokens(rejectedAccessToken)).finally(() => {
      refreshInFlight = null
    })
  }
  return refreshInFlight
}
