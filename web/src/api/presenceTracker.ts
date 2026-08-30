import { socketManager } from './ws'
import { getStatus } from './status'
import type { StatusInfo } from '../types'

class PresenceTracker {
  private statuses = new Map<string, StatusInfo>()
  private trackedUsers = new Set<string>()
  private listeners = new Map<string, Set<(status: StatusInfo | null) => void>>()
  private wsBound = false

  private bindWs() {
    if (this.wsBound) return
    this.wsBound = true

    socketManager.onEvent((event) => {
      if (event.type !== 'status') return
      const { userId, newStatus, eventTime } = event.payload
      this.statuses.set(userId, { status: newStatus, lastSeen: eventTime })
      this.notify(userId)
    })

    socketManager.onStatus((connected) => {
      if (!connected) return
      // The server's "who am I tracking" list lives with the connection —
      // after a reconnect we have to ask again for everyone we cared about.
      this.trackedUsers.forEach((userId) => socketManager.trackUser(userId))
    })
  }

  getStatus(userId: string): StatusInfo | null {
    return this.statuses.get(userId) ?? null
  }

  /**
   * Starts live-tracking a user's status. The first subscriber for a given
   * userId triggers a one-time `GET /api/status/{id}` (so the UI has a
   * baseline immediately, not just "unknown" until their next status
   * change) plus a `req:2` track subscription over the socket; everyone
   * after that just rides the same cached value and push updates — no
   * per-open-chat polling.
   */
  subscribe(userId: string, cb: (status: StatusInfo | null) => void): () => void {
    this.bindWs()

    let set = this.listeners.get(userId)
    if (!set) {
      set = new Set()
      this.listeners.set(userId, set)
    }
    set.add(cb)
    cb(this.getStatus(userId))

    if (!this.trackedUsers.has(userId)) {
      this.trackedUsers.add(userId)
      socketManager.trackUser(userId)
      getStatus(userId)
        .then((status) => {
          this.statuses.set(userId, status)
          this.notify(userId)
        })
        .catch(() => {})
    }

    return () => {
      set!.delete(cb)
    }
  }

  private notify(userId: string) {
    this.listeners.get(userId)?.forEach((cb) => cb(this.getStatus(userId)))
  }
}

export const presenceTracker = new PresenceTracker()
