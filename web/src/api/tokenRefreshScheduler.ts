import { tokenStorage } from '../utils/tokenStorage'
import { getJwtExpiry, isJwtExpiringSoon } from '../utils/jwt'
import { ensureFreshToken, TOKEN_REFRESHED_EVENT } from './authRefresh'

const REFRESH_MARGIN_MS = 60_000
const MIN_DELAY_MS = 1_000
const MAX_DELAY_MS = 10 * 60_000
const TRANSIENT_FAILURE_RETRY_MS = 15_000

function millisecondsUntilRefreshDue(accessToken: string): number {
  const expiry = getJwtExpiry(accessToken)
  if (expiry === null) return TRANSIENT_FAILURE_RETRY_MS
  const due = expiry * 1000 - REFRESH_MARGIN_MS - Date.now()
  return Math.min(Math.max(due, MIN_DELAY_MS), MAX_DELAY_MS)
}

class TokenRefreshScheduler {
  private timer: number | null = null
  private active = false

  start() {
    if (this.active) return
    this.active = true
    window.addEventListener(TOKEN_REFRESHED_EVENT, this.reschedule)
    window.addEventListener('storage', this.reschedule)
    window.addEventListener('online', this.refreshIfDue)
    document.addEventListener('visibilitychange', this.refreshIfDue)
    this.reschedule()
  }

  stop() {
    this.active = false
    this.clearTimer()
    window.removeEventListener(TOKEN_REFRESHED_EVENT, this.reschedule)
    window.removeEventListener('storage', this.reschedule)
    window.removeEventListener('online', this.refreshIfDue)
    document.removeEventListener('visibilitychange', this.refreshIfDue)
  }

  private clearTimer() {
    if (this.timer !== null) {
      window.clearTimeout(this.timer)
      this.timer = null
    }
  }

  private reschedule = () => {
    this.clearTimer()
    if (!this.active) return
    const tokens = tokenStorage.get()
    if (!tokens) return
    this.timer = window.setTimeout(this.refreshNow, millisecondsUntilRefreshDue(tokens.accessToken))
  }

  private refreshNow = async () => {
    if (!this.active) return
    const token = await ensureFreshToken()
    if (!this.active) return
    if (token) {
      this.reschedule()
      return
    }
    this.clearTimer()
    if (tokenStorage.get()) {
      this.timer = window.setTimeout(this.refreshNow, TRANSIENT_FAILURE_RETRY_MS)
    }
  }

  private refreshIfDue = () => {
    if (!this.active || document.visibilityState === 'hidden') return
    const tokens = tokenStorage.get()
    if (!tokens) return
    if (isJwtExpiringSoon(tokens.accessToken, REFRESH_MARGIN_MS / 1000)) {
      this.refreshNow()
    } else {
      this.reschedule()
    }
  }
}

export const tokenRefreshScheduler = new TokenRefreshScheduler()
