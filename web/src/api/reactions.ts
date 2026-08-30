import { apiClient } from './client'

let cachedReactions: Promise<string[]> | null = null

/** The supported reaction set is fixed for the session — fetch it once and
 * reuse the same promise for every caller, including concurrent ones. */
export function getSupportedReactions(): Promise<string[]> {
  if (!cachedReactions) {
    cachedReactions = apiClient.get<string[]>('/msg/reactions').then((r) => r.data)
  }
  return cachedReactions
}

export async function addReaction(chatId: string, messageId: string, emoji: string): Promise<void> {
  await apiClient.post(`/msg/chat/${chatId}/message/${messageId}/reaction`, { emoji })
}

export async function removeReaction(chatId: string, messageId: string, emoji: string): Promise<void> {
  await apiClient.delete(`/msg/chat/${chatId}/message/${messageId}/reaction/${encodeURIComponent(emoji)}`)
}
