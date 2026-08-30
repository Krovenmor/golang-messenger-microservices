import { useEffect, useState } from 'react'
import type { Message } from '../types'
import { chatStore } from '../api/chatStore'
import { useDecryptedText } from './useDecryptedText'

/** For a message that has `replyToId` set, resolves the message it's
 * replying to (from the already-loaded history, or a single deduped GET if
 * it's older than what's cached) and decrypts it for a quoted preview. */
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
