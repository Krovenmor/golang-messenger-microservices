import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState, type ReactNode } from 'react'
import { tokenStorage } from '../utils/tokenStorage'
import { getOwnProfile } from '../api/profile'
import { socketManager } from '../api/ws'
import { chatStore } from '../api/chatStore'
import { cryptoStore } from '../api/cryptoStore'
import * as keyVault from '../utils/keyVault'
import type { AuthTokens, OwnProfile } from '../types'

interface AuthContextValue {
  isAuthenticated: boolean
  profile: OwnProfile | null
  loading: boolean
  hasProfile: boolean | null
  /** Whether we have the unwrapped private key in memory this session —
   * required before any message can be encrypted/decrypted. A profile can
   * exist (hasProfile === true) while this is still false right after a
   * fresh login, until the user re-enters their key password on /unlock. */
  cryptoReady: boolean
  setTokens: (tokens: AuthTokens) => Promise<void>
  refreshProfile: () => Promise<OwnProfile | null>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(() => !!tokenStorage.get())
  const [profile, setProfile] = useState<OwnProfile | null>(null)
  const [hasProfile, setHasProfile] = useState<boolean | null>(null)
  const [loading, setLoading] = useState(true)
  const [cryptoReady, setCryptoReady] = useState(() => cryptoStore.isUnlocked())

  useEffect(() => cryptoStore.subscribe(setCryptoReady), [])

  const refreshProfile = useCallback(async () => {
    try {
      const data = await getOwnProfile()
      setProfile(data)
      setHasProfile(true)
      if (!cryptoStore.isUnlocked()) {
        const vaultKey = await keyVault.loadKey(data.userId)
        if (vaultKey) cryptoStore.setKeys(vaultKey, data.pubKey)
      }
      return data
    } catch {
      setProfile(null)
      setHasProfile(false)
      return null
    }
  }, [])

  const setTokens = useCallback(
    async (tokens: AuthTokens) => {
      tokenStorage.set(tokens)
      setIsAuthenticated(true)
      // The bootstrap effect only ever runs once, on the initial mount — a
      // login/register happening later in the same tab needs to fetch the
      // profile itself, otherwise `profile` stays null forever and every
      // screen gated on it (the whole chat view) never renders.
      await refreshProfile()
    },
    [refreshProfile]
  )

  const logout = useCallback(() => {
    tokenStorage.clear()
    socketManager.disconnect()
    chatStore.reset()
    cryptoStore.lock()
    setIsAuthenticated(false)
    setProfile(null)
    setHasProfile(null)
  }, [])

  useEffect(() => {
    const onForcedLogout = () => logout()
    window.addEventListener('signalline:logout', onForcedLogout)
    return () => window.removeEventListener('signalline:logout', onForcedLogout)
  }, [logout])

  const bootstrapped = useRef(false)
  useEffect(() => {
    if (bootstrapped.current) return
    bootstrapped.current = true
    const bootstrap = async () => {
      const tokens = tokenStorage.get()
      if (tokens) {
        await refreshProfile()
      }
      setLoading(false)
    }
    bootstrap()
  }, [refreshProfile])

  // The single place that decides whether the socket should be up at all:
  // only once we're authenticated AND a messenger profile actually exists.
  // There's nothing to connect to before profile creation, so we don't try —
  // and this being the one declarative effect (rather than ad-hoc connect()
  // calls scattered through login/register/bootstrap) is also what keeps us
  // from ever having two independent code paths open two sockets.
  useEffect(() => {
    if (isAuthenticated && hasProfile === true) {
      socketManager.connect()
    } else {
      socketManager.disconnect()
    }
  }, [isAuthenticated, hasProfile])

  const value = useMemo(
    () => ({ isAuthenticated, profile, loading, hasProfile, cryptoReady, setTokens, refreshProfile, logout }),
    [isAuthenticated, profile, loading, hasProfile, cryptoReady, setTokens, refreshProfile, logout]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
