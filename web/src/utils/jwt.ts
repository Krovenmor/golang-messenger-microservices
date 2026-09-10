export function getJwtExpiry(token: string): number | null {
  try {
    const payload = token.split('.')[1]
    if (!payload) return null
    const base64 = payload.replace(/-/g, '+').replace(/_/g, '/')
    const json = atob(base64)
    const data = JSON.parse(json)
    return typeof data.exp === 'number' ? data.exp : null
  } catch {
    return null
  }
}

export function isJwtExpiringSoon(token: string, marginSec = 10): boolean {
  const exp = getJwtExpiry(token)
  if (exp === null) return false
  return Date.now() >= (exp - marginSec) * 1000
}
