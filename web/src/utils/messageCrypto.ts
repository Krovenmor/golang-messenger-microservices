// Per-message encryption. Every message gets its own random AES-256-GCM
// MessageKey. That key is then sealed twice — once for the receiver, once
// for the sender (so the sender can read their own history back) — using a
// fresh ephemeral ECDH keypair per message:
//
//   sharedSecret = ECDH(ephemeralPriv, targetStaticPub)
//   wrapKey      = HKDF-SHA256(sharedSecret)              -> AES-KW key
//   sealed       = ephemeralPub(65B) || AES-KW(MessageKey)(40B)  = 105B
//
// `salt`/`info` for the HKDF step aren't specified by the backend (it never
// looks inside these envelopes) — they only need to be consistent within
// this frontend, so they're fixed constants below.
import { bufToBase64, base64ToBuf } from './crypto'

const HKDF_INFO = new TextEncoder().encode('signalline-message-key-wrap')
const EMPTY_SALT = new Uint8Array(0)

async function importRawPublicKey(raw: Uint8Array): Promise<CryptoKey> {
  return crypto.subtle.importKey('raw', raw, { name: 'ECDH', namedCurve: 'P-256' }, true, [])
}

async function deriveWrapKey(ecdhPrivate: CryptoKey, ecdhPublic: CryptoKey): Promise<CryptoKey> {
  const sharedSecret = await crypto.subtle.deriveBits({ name: 'ECDH', public: ecdhPublic }, ecdhPrivate, 256)
  const hkdfKey = await crypto.subtle.importKey('raw', sharedSecret, 'HKDF', false, ['deriveKey'])
  return crypto.subtle.deriveKey(
    { name: 'HKDF', hash: 'SHA-256', salt: EMPTY_SALT, info: HKDF_INFO },
    hkdfKey,
    { name: 'AES-KW', length: 256 },
    false,
    ['wrapKey', 'unwrapKey']
  )
}

async function sealMessageKey(ephemeralPriv: CryptoKey, ephemeralPubRaw: Uint8Array, targetPub: CryptoKey, messageKeyBytes: Uint8Array): Promise<Uint8Array> {
  const wrapKey = await deriveWrapKey(ephemeralPriv, targetPub)
  const messageKey = await crypto.subtle.importKey('raw', messageKeyBytes, { name: 'AES-GCM' }, true, [
    'encrypt',
    'decrypt',
  ])
  const wrapped = new Uint8Array(await crypto.subtle.wrapKey('raw', messageKey, wrapKey, 'AES-KW'))
  const packed = new Uint8Array(ephemeralPubRaw.length + wrapped.length)
  packed.set(ephemeralPubRaw, 0)
  packed.set(wrapped, ephemeralPubRaw.length)
  return packed
}

export interface EncryptedEnvelope {
  message: string
  senderKey: string
  receiverKey: string
  nonce: string
}

/** Encrypts `plaintext` for a 1:1 chat, producing the full envelope the
 * backend expects. Needs both parties' static public keys (raw, base64).
 * `chatId` is mixed in as AES-GCM additional authenticated data — it isn't
 * secret and isn't stored anywhere in the envelope, but decryption will
 * fail if the ciphertext is ever presented against a different chatId than
 * the one it was encrypted for, which rules out moving/replaying a message
 * from one chat into another. */
export async function encryptMessage(
  plaintext: string,
  myPubKeyB64: string,
  peerPubKeyB64: string,
  chatId: string
): Promise<EncryptedEnvelope> {
  const messageKeyBytes = crypto.getRandomValues(new Uint8Array(32))
  const nonce = crypto.getRandomValues(new Uint8Array(12))
  const aad = new TextEncoder().encode(chatId)

  const messageKey = await crypto.subtle.importKey('raw', messageKeyBytes, { name: 'AES-GCM' }, false, [
    'encrypt',
  ])
  const ciphertext = await crypto.subtle.encrypt(
    { name: 'AES-GCM', iv: nonce, additionalData: aad },
    messageKey,
    new TextEncoder().encode(plaintext)
  )

  const ephemeral = await crypto.subtle.generateKey({ name: 'ECDH', namedCurve: 'P-256' }, true, [
    'deriveKey',
    'deriveBits',
  ])
  const ephemeralPubRaw = new Uint8Array(await crypto.subtle.exportKey('raw', ephemeral.publicKey))

  const myPub = await importRawPublicKey(base64ToBuf(myPubKeyB64))
  const peerPub = await importRawPublicKey(base64ToBuf(peerPubKeyB64))

  const [senderKey, receiverKey] = await Promise.all([
    sealMessageKey(ephemeral.privateKey, ephemeralPubRaw, myPub, messageKeyBytes),
    sealMessageKey(ephemeral.privateKey, ephemeralPubRaw, peerPub, messageKeyBytes),
  ])

  return {
    message: bufToBase64(ciphertext),
    senderKey: bufToBase64(senderKey),
    receiverKey: bufToBase64(receiverKey),
    nonce: bufToBase64(nonce),
  }
}

/** Unseals the MessageKey from a `senderKey`/`receiverKey` envelope (pick
 * whichever matches our role for that message) using our own static
 * private key, then decrypts `message` with it. `chatId` must be the chat
 * the message actually belongs to — see the AAD note on encryptMessage. */
export async function decryptMessage(
  ciphertextB64: string,
  sealedKeyB64: string,
  nonceB64: string,
  myPrivateKey: CryptoKey,
  chatId: string
): Promise<string> {
  const packed = base64ToBuf(sealedKeyB64)
  const ephemeralPubRaw = packed.slice(0, 65)
  const wrapped = packed.slice(65)

  const ephemeralPub = await importRawPublicKey(ephemeralPubRaw)
  const wrapKey = await deriveWrapKey(myPrivateKey, ephemeralPub)
  const messageKey = await crypto.subtle.unwrapKey(
    'raw',
    wrapped,
    wrapKey,
    'AES-KW',
    { name: 'AES-GCM' },
    false,
    ['decrypt']
  )

  const plaintextBuf = await crypto.subtle.decrypt(
    { name: 'AES-GCM', iv: base64ToBuf(nonceB64), additionalData: new TextEncoder().encode(chatId) },
    messageKey,
    base64ToBuf(ciphertextB64)
  )
  return new TextDecoder().decode(plaintextBuf)
}
