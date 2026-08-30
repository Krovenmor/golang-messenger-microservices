// Profile key management: an ECDH P-256 keypair per user. The public key is
// shared with everyone (needed to encrypt messages TO this user); the
// private key never leaves the browser in plaintext — it's encrypted with a
// key derived from the user's password (PBKDF2) before being sent to the
// server, and only ever lives unwrapped in memory (see cryptoStore.ts).
//
// `crypto.subtle` only exists in secure contexts (HTTPS or localhost) — on
// plain http:// (e.g. hitting the box over the LAN by IP) it's simply
// undefined in every browser, no client-side workaround exists. In that
// case we fall back to random placeholder values so the app doesn't hang or
// crash, but real end-to-end encryption is only possible over https.

function randomBytes(n: number): Uint8Array {
  const arr = new Uint8Array(n)
  if (crypto?.getRandomValues) {
    crypto.getRandomValues(arr)
  } else {
    for (let i = 0; i < n; i++) arr[i] = Math.floor(Math.random() * 256)
  }
  return arr
}

export function bufToBase64(buf: ArrayBuffer | Uint8Array): string {
  const bytes = buf instanceof Uint8Array ? buf : new Uint8Array(buf)
  let binary = ''
  for (let i = 0; i < bytes.length; i++) binary += String.fromCharCode(bytes[i])
  return btoa(binary)
}

export function base64ToBuf(b64: string): Uint8Array {
  const binary = atob(b64)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i)
  return bytes
}

async function deriveWrappingKey(password: string, salt: Uint8Array): Promise<CryptoKey> {
  const enc = new TextEncoder()
  const baseKey = await crypto.subtle.importKey('raw', enc.encode(password), 'PBKDF2', false, [
    'deriveKey',
  ])
  return crypto.subtle.deriveKey(
    { name: 'PBKDF2', salt, iterations: 150_000, hash: 'SHA-256' },
    baseKey,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt', 'decrypt']
  )
}

export function hasRealCrypto(): boolean {
  return !!window.isSecureContext && !!crypto?.subtle
}

export interface ProfileKeyMaterial {
  pubKey: string
  encryptedPrvKey: string
  kdfSalt: string
  keyNonce: string
}

export interface GeneratedProfileKeys {
  material: ProfileKeyMaterial
  privateKey: CryptoKey
}

/** Generates a fresh ECDH P-256 keypair for a new profile, encrypting the
 * private key with a key derived from `password`. Returns both the
 * server-bound material and the live CryptoKey, so the caller can keep the
 * key unlocked in memory for this session without a decrypt round-trip. */
export async function generateProfileKeys(password: string): Promise<GeneratedProfileKeys> {
  if (!hasRealCrypto()) {
    return {
      material: {
        pubKey: bufToBase64(randomBytes(65)),
        encryptedPrvKey: bufToBase64(randomBytes(48)),
        kdfSalt: bufToBase64(randomBytes(16)),
        keyNonce: bufToBase64(randomBytes(12)),
      },
      privateKey: null as unknown as CryptoKey,
    }
  }

  const keyPair = await crypto.subtle.generateKey({ name: 'ECDH', namedCurve: 'P-256' }, true, [
    'deriveKey',
    'deriveBits',
  ])

  // Raw (uncompressed point) format: 65 bytes for P-256 — this is what the
  // rest of the crypto scheme (ephemeral keys, ECDH imports) uses too, so
  // every public key in the system is encoded the same way.
  const pubRaw = await crypto.subtle.exportKey('raw', keyPair.publicKey)
  const prvPkcs8 = await crypto.subtle.exportKey('pkcs8', keyPair.privateKey)

  const salt = randomBytes(16)
  const iv = randomBytes(12)
  const wrappingKey = await deriveWrappingKey(password, salt)
  const encryptedPrv = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, wrappingKey, prvPkcs8)

  // The extractable keyPair.privateKey above only ever existed to produce
  // the pkcs8 bytes we just encrypted for upload — re-import a
  // non-extractable copy from those same bytes for actual use. This is the
  // copy that ends up in cryptoStore and (optionally) IndexedDB; its raw
  // bytes can never be pulled back out via exportKey, from this tab or a
  // future one reading it back from storage.
  const nonExtractablePrivateKey = await crypto.subtle.importKey(
    'pkcs8',
    prvPkcs8,
    { name: 'ECDH', namedCurve: 'P-256' },
    false,
    ['deriveKey', 'deriveBits']
  )

  return {
    material: {
      pubKey: bufToBase64(pubRaw),
      encryptedPrvKey: bufToBase64(encryptedPrv),
      kdfSalt: bufToBase64(salt),
      keyNonce: bufToBase64(iv),
    },
    privateKey: nonExtractablePrivateKey,
  }
}

/** Decrypts a profile's stored private key with the user's password —
 * used when logging back in on a session that doesn't have the key in
 * memory yet. Throws if the password is wrong (AES-GCM auth tag failure). */
export async function unlockPrivateKey(
  password: string,
  encryptedPrvKeyB64: string,
  kdfSaltB64: string,
  keyNonceB64: string
): Promise<CryptoKey> {
  if (!hasRealCrypto()) {
    throw new Error('WebCrypto недоступен (нужен HTTPS или localhost)')
  }
  const salt = base64ToBuf(kdfSaltB64)
  const iv = base64ToBuf(keyNonceB64)
  const wrappingKey = await deriveWrappingKey(password, salt)
  const pkcs8 = await crypto.subtle.decrypt(
    { name: 'AES-GCM', iv },
    wrappingKey,
    base64ToBuf(encryptedPrvKeyB64)
  )
  return crypto.subtle.importKey('pkcs8', pkcs8, { name: 'ECDH', namedCurve: 'P-256' }, false, [
    'deriveKey',
    'deriveBits',
  ])
}
