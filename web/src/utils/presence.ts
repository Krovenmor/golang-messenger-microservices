import { socketManager } from '../api/ws'
import { ProfileStatus } from '../types'

// Timings are ours to pick (nothing in the API dictates them):
const AWAY_AFTER_MS = 60_000 // no activity for a minute -> Away
const TYPING_HOLD_MS = 2_000 // 2s of continuous typing before announcing Typing
const TYPING_STOP_MS = 3_000 // no keystroke for 3s -> stop announcing Typing
const IDLE_CHECK_INTERVAL_MS = 5_000

class Presence {
  private tracking = false
  private lastActivityAt = Date.now()
  private isAway = false

  private typingBurstStart: number | null = null
  private isTyping = false
  private typingStopTimer: number | null = null
  private idleCheckTimer: number | null = null

  private handleActivity = () => {
    this.lastActivityAt = Date.now()
    if (this.isAway) {
      this.isAway = false
      socketManager.sendStatus(ProfileStatus.Online)
    }
  }

  private checkIdle = () => {
    if (!this.isAway && Date.now() - this.lastActivityAt >= AWAY_AFTER_MS) {
      this.isAway = true
      socketManager.sendStatus(ProfileStatus.Away)
    }
  }

  /** Starts watching for activity — only meaningful while the socket is
   * actually connected, so this is driven by socketManager's own status
   * events rather than called directly from UI code. */
  start() {
    if (this.tracking) return
    this.tracking = true
    this.lastActivityAt = Date.now()
    this.isAway = false

    const opts: AddEventListenerOptions = { passive: true }
    window.addEventListener('mousemove', this.handleActivity, opts)
    window.addEventListener('mousedown', this.handleActivity, opts)
    window.addEventListener('keydown', this.handleActivity, opts)
    window.addEventListener('touchstart', this.handleActivity, opts)
    window.addEventListener('scroll', this.handleActivity, opts)
    this.idleCheckTimer = window.setInterval(this.checkIdle, IDLE_CHECK_INTERVAL_MS)
  }

  stop() {
    if (!this.tracking) return
    this.tracking = false
    window.removeEventListener('mousemove', this.handleActivity)
    window.removeEventListener('mousedown', this.handleActivity)
    window.removeEventListener('keydown', this.handleActivity)
    window.removeEventListener('touchstart', this.handleActivity)
    window.removeEventListener('scroll', this.handleActivity)
    if (this.idleCheckTimer) window.clearInterval(this.idleCheckTimer)
    if (this.typingStopTimer) window.clearTimeout(this.typingStopTimer)
    this.idleCheckTimer = null
    this.typingStopTimer = null
    this.typingBurstStart = null
    this.isTyping = false
  }

  /** Call on every keystroke in the message composer. Only actually
   * announces Typing once the user has been typing continuously for
   * TYPING_HOLD_MS, so a single keystroke doesn't spam a status change. */
  notifyTyping() {
    if (!this.tracking) return
    this.handleActivity() // typing counts as activity for the away-timer too

    const now = Date.now()
    if (this.typingBurstStart === null) {
      this.typingBurstStart = now
    } else if (!this.isTyping && now - this.typingBurstStart >= TYPING_HOLD_MS) {
      this.isTyping = true
      socketManager.sendStatus(ProfileStatus.Typing)
    }

    if (this.typingStopTimer) window.clearTimeout(this.typingStopTimer)
    this.typingStopTimer = window.setTimeout(() => {
      this.typingBurstStart = null
      if (this.isTyping) {
        this.isTyping = false
        socketManager.sendStatus(ProfileStatus.Online)
      }
    }, TYPING_STOP_MS)
  }
}

export const presence = new Presence()

// Presence only makes sense while there's a live connection to report over.
socketManager.onStatus((connected) => {
  if (connected) {
    presence.start()
  } else {
    presence.stop()
  }
})
