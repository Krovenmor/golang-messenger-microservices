import { useEffect, useState } from 'react'
import { batchGetProfiles, getProfile } from './profile'
import type { PublicProfile } from '../types'

const cache = new Map<string, PublicProfile>()

export async function getCachedProfile(userId: string): Promise<PublicProfile> {
  const hit = cache.get(userId)
  if (hit) return hit
  const profile = await getProfile(userId)
  cache.set(userId, profile)
  return profile
}

export function primeCachedProfile(profile: PublicProfile) {
  cache.set(profile.userId, profile)
}

export function getCachedProfileSync(userId: string): PublicProfile | null {
  return cache.get(userId) ?? null
}

export async function primeProfilesFromBatch(userIds: string[]): Promise<void> {
  const missing = [...new Set(userIds)].filter((id) => !cache.has(id))
  if (missing.length === 0) return
  try {
    const profiles = await batchGetProfiles(missing)
    profiles.forEach(primeCachedProfile)
  } catch {}
}

export function usePublicProfile(userId: string | null): PublicProfile | null {
  const [profile, setProfile] = useState<PublicProfile | null>(() => (userId ? cache.get(userId) ?? null : null))

  useEffect(() => {
    if (!userId) {
      setProfile(null)
      return
    }
    let cancelled = false
    getCachedProfile(userId)
      .then((p) => {
        if (!cancelled) setProfile(p)
      })
      .catch(() => {})
    return () => {
      cancelled = true
    }
  }, [userId])

  return profile
}
