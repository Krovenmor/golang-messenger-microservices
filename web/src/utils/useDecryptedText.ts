import { useEffect, useState } from 'react'
import type { Message } from '../types'
import { cryptoStore } from '../api/cryptoStore'
import { decryptMessage } from './messageCrypto'

const cache = new Map<string, string>()

function cacheKey(message: Message): string {
  return `${message.messageId}:${message.nonce}`
}

export type DecryptState =
  | { status: 'empty' }
  | { status: 'locked' }
  | { status: 'loading' }
  | { status: 'ok'; text: string }
  | { status: 'error' }

export function useDecryptedText(
  message: Message | null | undefined,
  ownUserId: string | null,
  chatId: string
): DecryptState {
  const [state, setState] = useState<DecryptState>(() => computeInitial(message))
  const [cryptoUnlocked, setCryptoUnlocked] = useState(() => cryptoStore.isUnlocked())

  useEffect(() => cryptoStore.subscribe(setCryptoUnlocked), [])

  useEffect(() => {
    setState(computeInitial(message))
    if (!message || message.deletedAt) return

    const cached = cache.get(cacheKey(message))
    if (cached !== undefined) {
      setState({ status: 'ok', text: cached })
      return
    }

    if (!ownUserId) {
      setState({ status: 'loading' })
      return
    }

    const privateKey = cryptoStore.getPrivateKey()
    if (!privateKey) {
      setState({ status: 'locked' })
      return
    }

    let cancelled = false
    const envelope = message.senderId === ownUserId ? message.senderKey : message.receiverKey
    decryptMessage(message.message, envelope, message.nonce, privateKey, chatId)
      .then((text) => {
        cache.set(cacheKey(message), text)
        if (!cancelled) setState({ status: 'ok', text })
      })
      .catch((err) => {
        console.error('[signaline] failed to decrypt message', message.messageId, 'in chat', chatId, err)
        if (!cancelled) setState({ status: 'error' })
      })
    return () => {
      cancelled = true
    }
  }, [message, ownUserId, chatId, cryptoUnlocked])

  return state
}

function computeInitial(message: Message | null | undefined): DecryptState {
  if (!message || message.deletedAt) return { status: 'empty' }
  const cached = cache.get(cacheKey(message))
  if (cached !== undefined) return { status: 'ok', text: cached }
  if (!cryptoStore.isUnlocked()) return { status: 'locked' }
  return { status: 'loading' }
}
