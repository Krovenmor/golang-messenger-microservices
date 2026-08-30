// Every UUID gets a stable hue derived from its own bytes, so each person's
// messages carry a consistent "fingerprint" color across the whole app —
// no lookups, no palette table, just the id itself.

export function fingerprintHue(id: string): number {
  let hash = 0
  for (let i = 0; i < id.length; i++) {
    hash = (hash * 31 + id.charCodeAt(i)) >>> 0
  }
  return hash % 360
}

export function fingerprintColor(id: string, opts: { saturation?: number; lightness?: number } = {}): string {
  const { saturation = 65, lightness = 62 } = opts
  return `hsl(${fingerprintHue(id)}, ${saturation}%, ${lightness}%)`
}

export function initials(name: string): string {
  const parts = name.trim().split(/\s+/)
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase()
  return (parts[0][0] + parts[1][0]).toUpperCase()
}

export function formatTime(iso: string | null): string {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })
}

export function formatDay(iso: string | null): string {
  if (!iso) return ''
  const d = new Date(iso)
  const today = new Date()
  const isToday = d.toDateString() === today.toDateString()
  if (isToday) return 'сегодня'
  return d.toLocaleDateString(undefined, { day: '2-digit', month: 'short' })
}
