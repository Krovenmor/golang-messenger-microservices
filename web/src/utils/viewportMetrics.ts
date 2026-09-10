const APP_HEIGHT = '--app-height'
const APP_OFFSET_TOP = '--app-offset-top'
const KEYBOARD_SETTLE_MS = 350

type ViewportListener = () => void

const listeners = new Set<ViewportListener>()

export function onViewportChange(listener: ViewportListener): () => void {
  listeners.add(listener)
  return () => {
    listeners.delete(listener)
  }
}

function hasFocusedTextEntry(): boolean {
  const el = document.activeElement as HTMLElement | null
  if (!el) return false
  return el.tagName === 'INPUT' || el.tagName === 'TEXTAREA' || el.isContentEditable
}

function measure() {
  const layoutHeight = document.documentElement.clientHeight
  const visual = window.visualViewport
  if (!visual) return { height: layoutHeight, offsetTop: 0 }

  const canBeCoveredByKeyboard = hasFocusedTextEntry() || visual.scale !== 1
  if (!canBeCoveredByKeyboard) return { height: layoutHeight, offsetTop: 0 }

  return {
    height: Math.round(Math.min(visual.height, layoutHeight)),
    offsetTop: Math.round(visual.offsetTop),
  }
}

export function trackViewportMetrics(): () => void {
  const root = document.documentElement
  let published = { height: -1, offsetTop: -1 }
  let frame = 0
  let settleTimer = 0

  const publish = () => {
    frame = 0
    const next = measure()
    if (next.height === published.height && next.offsetTop === published.offsetTop) return
    published = next
    root.style.setProperty(APP_HEIGHT, `${next.height}px`)
    root.style.setProperty(APP_OFFSET_TOP, `${next.offsetTop}px`)
    listeners.forEach((listener) => listener())
  }

  const schedule = () => {
    if (!frame) frame = requestAnimationFrame(publish)
    window.clearTimeout(settleTimer)
    settleTimer = window.setTimeout(publish, KEYBOARD_SETTLE_MS)
  }

  publish()

  const visual = window.visualViewport
  visual?.addEventListener('resize', schedule)
  visual?.addEventListener('scroll', schedule)
  window.addEventListener('resize', schedule)
  window.addEventListener('orientationchange', schedule)
  window.addEventListener('focusin', schedule)
  window.addEventListener('focusout', schedule)

  return () => {
    if (frame) cancelAnimationFrame(frame)
    window.clearTimeout(settleTimer)
    visual?.removeEventListener('resize', schedule)
    visual?.removeEventListener('scroll', schedule)
    window.removeEventListener('resize', schedule)
    window.removeEventListener('orientationchange', schedule)
    window.removeEventListener('focusin', schedule)
    window.removeEventListener('focusout', schedule)
  }
}
