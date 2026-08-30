/**
 * Decodes a JWT's payload without verifying the signature — that's fine
 * here since this only ever reads our own access token's `exp` claim to
 * decide whether it's worth attempting a connection with; the server is
 * still the one actually validating the token on every real request.
 */
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

/** True if the token is already expired, or will expire within `marginSec`. */
export function isJwtExpiringSoon(token: string, marginSec = 10): boolean {
  const exp = getJwtExpiry(token)
  if (exp === null) return false // can't tell from the token — don't block on it
  return Date.now() >= (exp - marginSec) * 1000
}
