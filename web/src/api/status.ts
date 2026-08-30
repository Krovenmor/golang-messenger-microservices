import { apiClient } from './client'
import type { StatusInfo } from '../types'

export async function getStatus(profileId: string): Promise<StatusInfo> {
  const { data } = await apiClient.get<StatusInfo>(`/status/${profileId}`)
  return data
}
