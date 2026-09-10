import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState, type ReactNode } from 'react'
import { tokenStorage } from '../utils/tokenStorage'
import { getOwnProfile } from '../api/profile'
import { socketManager } from '../api/ws'
import { SESSION_ENDED_EVENT } from '../api/authRefresh'
import { tokenRefreshScheduler } from '../api/tokenRefreshScheduler'
import { chatStore } from '../api/chatStore'
import { cryptoStore } from '../api/cryptoStore'
import * as keyVault from '../utils/keyVault'
import type { AuthTokens, OwnProfile } from '../types'

interface AuthContextValue {
  isAuthenticated: boolean
  profile: OwnProfile | null
  loading: boolean
  hasProfile: boolean | null
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
      await refreshProfile()
    },
    [refreshProfile]
  )

  const logout = useCallback(() => {
    tokenRefreshScheduler.stop()
    tokenStorage.clear()
    socketManager.disconnect()
    chatStore.reset()
    cryptoStore.lock()
    setIsAuthenticated(false)
    setProfile(null)
    setHasProfile(null)
  }, [])

  useEffect(() => {
    const onSessionEnded = () => logout()
    window.addEventListener(SESSION_ENDED_EVENT, onSessionEnded)
    return () => window.removeEventListener(SESSION_ENDED_EVENT, onSessionEnded)
  }, [logout])

  useEffect(() => {
    if (!isAuthenticated) return
    tokenRefreshScheduler.start()
    return () => tokenRefreshScheduler.stop()
  }, [isAuthenticated])

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
