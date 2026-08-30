import axios from 'axios'
import { tokenStorage } from '../utils/tokenStorage'
import type { AuthTokens } from '../types'

export const API_BASE_URL: string = import.meta.env.VITE_API_URL || '/api'

let isRefreshing = false
let queue: Array<(token: string | null) => void> = []

function flushQueue(token: string | null) {
  queue.forEach((cb) => cb(token))
  queue = []
}

/**
 * Performs (or joins an in-flight) refresh-token exchange, updates storage,
 * and notifies the rest of the app. Shared by the REST 401 interceptor
 * (client.ts) and the WebSocket's proactive pre-connect expiry check
 * (ws.ts) so the two can never race each other into firing two concurrent
 * refreshes — whichever calls first "wins", the other just waits on the
 * same in-flight promise via `queue`.
 *
 * Deliberately lives in its own module with no dependency on client.ts or
 * ws.ts: client.ts needs this to retry a failed REST call, ws.ts needs this
 * to avoid ever reconnecting with a token it already knows is dead — if
 * this lived in either of those files, the other would have to import it,
 * creating an import cycle between them.
 *
 * Returns the fresh access token, or null if the refresh token itself is
 * dead too (session is over — dispatches 'signalline:logout' in that case,
 * which AuthContext listens for). On success, dispatches
 * 'signalline:token-refreshed', which ws.ts listens for to pick up the new
 * token immediately instead of waiting out its own retry backoff.
 */
export async function ensureFreshToken(): Promise<string | null> {
  const tokens = tokenStorage.get()
  if (!tokens?.refreshToken) {
    tokenStorage.clear()
    return null
  }

  if (isRefreshing) {
    return new Promise((resolve) => {
      queue.push(resolve)
    })
  }

  isRefreshing = true
  try {
    const { data } = await axios.post<AuthTokens>(`${API_BASE_URL}/auth/update`, {
      refreshToken: tokens.refreshToken,
    })
    tokenStorage.set(data)
    window.dispatchEvent(new CustomEvent('signalline:token-refreshed'))
    flushQueue(data.accessToken)
    return data.accessToken
  } catch {
    flushQueue(null)
    tokenStorage.clear()
    window.dispatchEvent(new CustomEvent('signalline:logout'))
    return null
  } finally {
    isRefreshing = false
  }
}
