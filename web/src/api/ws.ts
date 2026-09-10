import { tokenStorage } from '../utils/tokenStorage'
import { isJwtExpiringSoon } from '../utils/jwt'
import { ensureFreshToken, TOKEN_REFRESHED_EVENT } from './authRefresh'
import type { ProfileStatus, WSEvent, WSRequest } from '../types'

export type WSListener = (event: WSEvent) => void

const WS_BASE_URL: string =
  import.meta.env.VITE_WS_URL || `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/api`

const RECONNECT_BASE_DELAY_MS = 2000
const RECONNECT_MAX_DELAY_MS = 30_000

class SocketManager {
  private socket: WebSocket | null = null
  private listeners = new Set<WSListener>()
  private statusListeners = new Set<(connected: boolean) => void>()
  private reconnectTimer: number | null = null
  private reconnectAttempts = 0
  private manuallyClosed = true
  private isConnected = false

  constructor() {
    window.addEventListener(TOKEN_REFRESHED_EVENT, () => this.reconnectNow())
    window.addEventListener('online', () => this.reconnectNow())
    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'visible') this.reconnectNow()
    })
  }

  connect() {
    this.manuallyClosed = false
    this.reconnectAttempts = 0
    this.clearReconnectTimer()
    if (this.isLiveOrConnecting()) return
    this.open()
  }

  disconnect() {
    this.manuallyClosed = true
    this.clearReconnectTimer()
    this.socket?.close()
    this.socket = null
  }

  private isLiveOrConnecting() {
    const state = this.socket?.readyState
    return state === WebSocket.OPEN || state === WebSocket.CONNECTING
  }

  private clearReconnectTimer() {
    if (this.reconnectTimer !== null) {
      window.clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
  }

  private reconnectNow() {
    if (this.manuallyClosed || this.isLiveOrConnecting()) return
    this.clearReconnectTimer()
    this.reconnectAttempts = 0
    this.socket = null
    this.open()
  }

  private scheduleReconnect() {
    if (this.manuallyClosed) return
    this.clearReconnectTimer()
    const delay = Math.min(RECONNECT_BASE_DELAY_MS * 2 ** this.reconnectAttempts, RECONNECT_MAX_DELAY_MS)
    this.reconnectAttempts += 1
    this.reconnectTimer = window.setTimeout(() => this.open(), delay)
  }

  private async open() {
    const tokens = tokenStorage.get()
    if (!tokens?.accessToken) return

    let accessToken = tokens.accessToken

    if (isJwtExpiringSoon(accessToken)) {
      const fresh = await ensureFreshToken()
      if (!fresh) {
        this.scheduleReconnect()
        return
      }
      accessToken = fresh
      if (this.isLiveOrConnecting()) return
    }

    this.socket = new WebSocket(`${WS_BASE_URL}/ws?token=${encodeURIComponent(accessToken)}`)

    this.socket.onopen = () => {
      this.reconnectAttempts = 0
      this.notifyStatus(true)
    }

    this.socket.onmessage = (ev) => {
      let data: unknown
      try {
        data = JSON.parse(ev.data)
      } catch {
        return
      }

      if (!data || typeof data !== 'object' || !('type' in data) || !('payload' in data)) return
      const event = data as WSEvent

      if (event.type === 'response') {
        if (event.payload.code !== 200) {
          console.warn('[signaline] WS request rejected:', event.payload.msg)
        }
        return
      }

      this.listeners.forEach((cb) => cb(event))
    }

    this.socket.onclose = () => {
      this.notifyStatus(false)
      this.scheduleReconnect()
    }

    this.socket.onerror = () => {
      this.socket?.close()
    }
  }

  private sendRequest(request: WSRequest) {
    if (this.socket?.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(request))
    }
  }

  sendStatus(status: ProfileStatus) {
    this.sendRequest({ req: 1, payload: { newStatus: status } })
  }

  trackUser(userId: string) {
    this.sendRequest({ req: 2, payload: { userId } })
  }

  onEvent(cb: WSListener) {
    this.listeners.add(cb)
    return () => {
      this.listeners.delete(cb)
    }
  }

  onStatus(cb: (connected: boolean) => void) {
    this.statusListeners.add(cb)
    cb(this.isConnected)
    return () => {
      this.statusListeners.delete(cb)
    }
  }

  private notifyStatus(connected: boolean) {
    this.isConnected = connected
    this.statusListeners.forEach((cb) => cb(connected))
  }
}

export const socketManager = new SocketManager()
