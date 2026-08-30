import { apiClient } from './client'
import type { MediaFileInfo, MediaStorageInfo } from '../types'

// The static vault serves raw files (unencrypted, public — e.g. avatars)
// from a path that isn't under /api, so it needs its own base rather than
// reusing apiClient's baseURL. Defaults to the same origin as the app,
// which is what nginx.conf's /static/ proxy block expects.
const STATIC_BASE_URL: string = import.meta.env.VITE_STATIC_URL || '/static'

export function getAvatarUrl(photoId: string): string {
  return `${STATIC_BASE_URL}/avatars/${photoId}`
}

export async function getMediaProfileInfo(): Promise<MediaStorageInfo> {
  const { data } = await apiClient.get<MediaStorageInfo>('/media/profile')
  return data
}

export async function getSavedFiles(opts: { from?: string; q?: number } = {}): Promise<MediaFileInfo[]> {
  const { data } = await apiClient.get<MediaFileInfo[]>('/media/profile/files', { params: opts })
  return data
}

/** Uploads an already client-side-resized image (see utils/imageResize.ts —
 * the backend caps avatars at 600x600px, png/jpeg). */
export async function uploadAvatarFile(file: File | Blob): Promise<{ photoId: string }> {
  const form = new FormData()
  form.append('image', file, file instanceof File ? file.name : 'avatar.jpg')
  // Deliberately not setting Content-Type manually — axios/the browser need
  // to add the multipart boundary parameter themselves, which isn't
  // something you can (correctly) do by hand.
  const { data } = await apiClient.post<{ photoId: string }>('/media/public/avatar', form)
  return data
}

/** Deletes the underlying file from media storage (frees quota) — separate
 * from unlinking it from a profile, see profile.ts's removeProfileAvatar. */
export async function deleteAvatarFile(photoId: string): Promise<void> {
  await apiClient.delete(`/media/public/avatar/${photoId}`)
}
