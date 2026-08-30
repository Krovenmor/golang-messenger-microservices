import { apiClient } from './client'
import type { ChatInfo, ChatSummary, Message } from '../types'
import type { EncryptedEnvelope } from '../utils/messageCrypto'

export async function getChatIds(): Promise<string[]> {
  const { data } = await apiClient.get<string[]>('/msg/chats')
  return data
}

export async function getChatsFull(): Promise<ChatSummary[]> {
  const { data } = await apiClient.get<ChatSummary[]>('/msg/chats/full')
  return data
}

export async function createChat(userId: string): Promise<{ chatId: string }> {
  const { data } = await apiClient.post<{ chatId: string }>('/msg/chat', { userId })
  return data
}

export async function postMessage(
  chatId: string,
  envelope: EncryptedEnvelope,
  replyToId?: string | null
): Promise<{ messageId: string }> {
  const { data } = await apiClient.post<{ messageId: string }>(`/msg/chat/${chatId}`, {
    ...envelope,
    ...(replyToId ? { replyToId } : {}),
  })
  return data
}

export async function getChatHistory(
  chatId: string,
  opts: { from?: string; q?: number } = {}
): Promise<Message[]> {
  const { data } = await apiClient.get<Message[]>(`/msg/chat/${chatId}`, {
    params: opts,
  })
  return data
}

export async function getChatInfo(chatId: string): Promise<ChatInfo> {
  const { data } = await apiClient.get<ChatInfo>(`/msg/chat/${chatId}/info`)
  return data
}

export async function getMessage(chatId: string, messageId: string): Promise<Message> {
  const { data } = await apiClient.get<Message>(`/msg/chat/${chatId}/message/${messageId}`)
  return data
}

export async function editMessage(chatId: string, messageId: string, envelope: EncryptedEnvelope): Promise<void> {
  await apiClient.put(`/msg/chat/${chatId}/message/${messageId}`, envelope)
}

export async function deleteMessage(chatId: string, messageId: string): Promise<void> {
  await apiClient.delete(`/msg/chat/${chatId}/message/${messageId}`)
}
