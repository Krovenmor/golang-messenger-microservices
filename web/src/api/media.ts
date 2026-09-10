import { apiClient } from './client'
import type { MediaFileInfo, MediaStorageInfo } from '../types'

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

export async function uploadAvatarFile(file: File | Blob): Promise<{ photoId: string }> {
  const form = new FormData()
  form.append('image', file, file instanceof File ? file.name : 'avatar.jpg')
  const { data } = await apiClient.post<{ photoId: string }>('/media/public/avatar', form)
  return data
}

export async function deleteAvatarFile(photoId: string): Promise<void> {
  await apiClient.delete(`/media/public/avatar/${photoId}`)
}
