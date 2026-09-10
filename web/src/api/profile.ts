import { apiClient } from './client'
import type { OwnProfile, PublicProfile } from '../types'

export interface NewProfilePayload {
  name: string
  userName: string
  pubKey: string
  encryptedPrvKey: string
  kdfSalt: string
  keyNonce: string
}

export async function createProfile(payload: NewProfilePayload): Promise<void> {
  await apiClient.post('/profile', payload)
}

export async function getOwnProfile(): Promise<OwnProfile> {
  const { data } = await apiClient.get<OwnProfile>('/profile')
  return data
}

export async function getProfile(target: string): Promise<PublicProfile> {
  const { data } = await apiClient.get<PublicProfile>(`/profile/${target}`)
  return data
}

export async function updateProfileName(name: string): Promise<void> {
  await apiClient.put('/profile/name', { name })
}

export async function updateProfileBio(bio: string): Promise<void> {
  await apiClient.put('/profile/bio', { bio })
}

const BATCH_LIMIT = 50

export async function batchGetProfiles(userIds: string[]): Promise<PublicProfile[]> {
  if (userIds.length === 0) return []
  if (userIds.length <= BATCH_LIMIT) {
    const { data } = await apiClient.post<PublicProfile[]>('/profile/batch', { profiles: userIds })
    return data
  }
  const chunks: string[][] = []
  for (let i = 0; i < userIds.length; i += BATCH_LIMIT) {
    chunks.push(userIds.slice(i, i + BATCH_LIMIT))
  }
  const results = await Promise.all(chunks.map((chunk) => batchGetProfiles(chunk)))
  return results.flat()
}

export async function addProfileAvatar(photoId: string): Promise<void> {
  await apiClient.post(`/profile/avatar/${photoId}`)
}

export async function removeProfileAvatar(photoId: string): Promise<void> {
  await apiClient.delete(`/profile/avatar/${photoId}`)
}
