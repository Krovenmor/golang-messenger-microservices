import { useEffect, useState } from 'react'
import { batchGetProfiles, getProfile } from './profile'
import type { PublicProfile } from '../types'

const cache = new Map<string, PublicProfile>()

/** Public keys don't change, so once we've looked a profile up there's no
 * reason to ever fetch it again this session. (Display fields like name or
 * avatar could in principle change — there's no push event for that, so a
 * stale avatar/name for someone already cached this session is an accepted
 * tradeoff rather than re-fetching on every render.) */
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

/** Bulk-resolves whichever of the given userIds aren't already cached, via
 * a single POST /profile/batch instead of one GET per user — this is what
 * the chat list uses on load so opening the app with, say, 30 chats costs
 * one request for everyone's name/avatar instead of 30. Already-cached
 * ids are skipped entirely; a fully-cached call makes no request at all. */
export async function primeProfilesFromBatch(userIds: string[]): Promise<void> {
  const missing = [...new Set(userIds)].filter((id) => !cache.has(id))
  if (missing.length === 0) return
  try {
    const profiles = await batchGetProfiles(missing)
    profiles.forEach(primeCachedProfile)
  } catch {
    // Non-fatal — usePublicProfile/getCachedProfile will just fall back to
    // individual lookups for whoever didn't get cached here.
  }
}

/** Reactive wrapper around getCachedProfile, for components that just want
 * to display someone's current name/avatar (e.g. a chat list row) rather
 * than needing the profile for an async operation like encryption. */
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
