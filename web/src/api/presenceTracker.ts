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
      this.trackedUsers.forEach((userId) => socketManager.trackUser(userId))
    })
  }

  getStatus(userId: string): StatusInfo | null {
    return this.statuses.get(userId) ?? null
  }

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
