import { tokenStorage } from '../utils/tokenStorage'
import { isJwtExpiringSoon } from '../utils/jwt'
import { ensureFreshToken } from './authRefresh'
import type { ProfileStatus, WSEvent, WSRequest } from '../types'

export type WSListener = (event: WSEvent) => void

const WS_BASE_URL: string =
  import.meta.env.VITE_WS_URL || `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/api`

class SocketManager {
  private socket: WebSocket | null = null
  private listeners = new Set<WSListener>()
  private statusListeners = new Set<(connected: boolean) => void>()
  private reconnectTimer: number | null = null
  private manuallyClosed = true
  private isConnected = false

  constructor() {
    // A REST call elsewhere in the app may have refreshed the token (via
    // the 401 interceptor) — if we're currently down and retrying, don't
    // wait out the rest of the backoff timer, grab it now.
    window.addEventListener('signalline:token-refreshed', () => this.reconnectWithFreshToken())
  }

  /**
   * Starts the connection using whatever access token is currently in
   * storage. Safe to call repeatedly and from multiple places (e.g. React
   * effects that run twice under StrictMode) — if a socket is already open
   * or in the middle of connecting, this is a no-op instead of spinning up
   * a second live connection.
   */
  connect() {
    this.manuallyClosed = false
    if (this.reconnectTimer) window.clearTimeout(this.reconnectTimer)
    if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
      return
    }
    this.open()
  }

  /**
   * Called after the REST client silently refreshes the access token. Only
   * actually does anything if the socket isn't currently alive (or actively
   * connecting) — a live WS session doesn't care that the REST access token
   * expired, since that token only ever mattered for the initial handshake.
   * Forcing a reconnect on a perfectly good connection would just be a
   * pointless disruption (and, done carelessly, briefly look like two
   * connections). This exists purely so a socket that's already down and
   * retrying with a stale token doesn't have to wait out its own backoff
   * timer before picking up the new one.
   */
  reconnectWithFreshToken() {
    if (this.manuallyClosed) return
    const state = this.socket?.readyState
    if (state === WebSocket.OPEN || state === WebSocket.CONNECTING) return
    if (this.reconnectTimer) window.clearTimeout(this.reconnectTimer)
    this.socket = null
    this.open()
  }

  private async open() {
    const tokens = tokenStorage.get()
    if (!tokens?.accessToken) return

    let accessToken = tokens.accessToken

    // The bug this closes: if nothing else in the app happens to make a
    // REST call around the time the access token expires, there's nothing
    // to trigger the usual 401-driven refresh — left alone, a socket that
    // drops for any reason would just retry forever with the same dead
    // token every 2s. Checking the token's own exp claim here means the
    // retry loop can refresh itself, with no dependency on unrelated REST
    // traffic happening to occur.
    if (isJwtExpiringSoon(accessToken)) {
      const fresh = await ensureFreshToken()
      if (!fresh) {
        // Refresh token is dead too — the session is genuinely over.
        // ensureFreshToken already dispatched signalline:logout (which
        // tears this connection down via disconnect()); don't schedule
        // another retry that can only ever fail the same way.
        this.manuallyClosed = true
        if (this.reconnectTimer) window.clearTimeout(this.reconnectTimer)
        return
      }
      accessToken = fresh
      // A concurrent path (e.g. the token-refreshed listener above) may
      // have already opened a new connection with this same fresh token
      // while we were awaiting — don't duplicate it.
      if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
        return
      }
    }

    this.socket = new WebSocket(`${WS_BASE_URL}/ws?token=${encodeURIComponent(accessToken)}`)

    this.socket.onopen = () => {
      this.notifyStatus(true)
    }

    this.socket.onmessage = (ev) => {
      let data: unknown
      try {
        data = JSON.parse(ev.data)
      } catch {
        return // malformed frame — ignore
      }

      if (!data || typeof data !== 'object' || !('type' in data) || !('payload' in data)) return
      const event = data as WSEvent

      // `response` is the ack for a request we sent ourselves (post status /
      // track user) — it's not something general event listeners care
      // about, just log if the server rejected it.
      if (event.type === 'response') {
        if (event.payload.code !== 200) {
          console.warn('[signalline] WS request rejected:', event.payload.msg)
        }
        return
      }

      this.listeners.forEach((cb) => cb(event))
    }

    this.socket.onclose = () => {
      this.notifyStatus(false)
      if (!this.manuallyClosed) {
        // Re-reads the token from storage on every attempt, so if a REST
        // call refreshed it in the meantime, the retry uses the new one
        // instead of looping forever on an expired token.
        this.reconnectTimer = window.setTimeout(() => this.open(), 2000)
      }
    }

    this.socket.onerror = () => {
      this.socket?.close()
    }
  }

  disconnect() {
    this.manuallyClosed = true
    if (this.reconnectTimer) window.clearTimeout(this.reconnectTimer)
    this.socket?.close()
    this.socket = null
  }

  /** Sends a request frame, if the socket is actually open. */
  private sendRequest(request: WSRequest) {
    if (this.socket?.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(request))
    }
  }

  /** Reports our own presence status — see ProfileStatus in types. Online
   * is assigned automatically on connect and Offline on disconnect, so
   * this is only ever used to announce Away/Typing (and Online again once
   * either of those ends). */
  sendStatus(status: ProfileStatus) {
    this.sendRequest({ req: 1, payload: { newStatus: status } })
  }

  /** Subscribes this connection to live status push events for `userId`
   * (delivered as `type: "status"` events to onEvent listeners). */
  trackUser(userId: string) {
    this.sendRequest({ req: 2, payload: { userId } })
  }

  onEvent(cb: WSListener) {
    this.listeners.add(cb)
    return () => this.listeners.delete(cb)
  }

  onStatus(cb: (connected: boolean) => void) {
    this.statusListeners.add(cb)
    // Sync the subscriber immediately — otherwise a component that mounts
    // after the socket already opened (e.g. after a hard reload, where the
    // connection happens before React finishes bootstrapping auth) would
    // never learn about that and would show "reconnecting" forever even
    // though the socket is fine.
    cb(this.isConnected)
    return () => this.statusListeners.delete(cb)
  }

  private notifyStatus(connected: boolean) {
    this.isConnected = connected
    this.statusListeners.forEach((cb) => cb(connected))
  }
}

export const socketManager = new SocketManager()
