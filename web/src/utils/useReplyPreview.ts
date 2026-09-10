import { useEffect, useState } from 'react'
import type { Message } from '../types'
import { chatStore } from '../api/chatStore'
import { useDecryptedText } from './useDecryptedText'

export function useReplyPreview(chatId: string, replyToId: string | null | undefined, ownUserId: string) {
  const [target, setTarget] = useState<Message | null>(null)

  useEffect(() => {
    if (!replyToId) {
      setTarget(null)
      return
    }
    let cancelled = false
    chatStore.getOrFetchMessage(chatId, replyToId).then((m) => {
      if (!cancelled) setTarget(m)
    })
    return () => {
      cancelled = true
    }
  }, [chatId, replyToId])

  const decrypted = useDecryptedText(target, ownUserId, chatId)
  return { target, decrypted }
}
