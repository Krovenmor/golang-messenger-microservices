// Persists the unwrapped ECDH private key in IndexedDB, keyed by userId, so
// a returning user on the SAME device doesn't have to re-enter their key
// password every session. This is safe specifically because the key stored
// here is always imported non-extractable (see crypto.ts) — IndexedDB can
// hold a CryptoKey object directly via structured clone, and a
// non-extractable key can still be used for deriveBits/deriveKey after
// being read back, but its raw bytes can never be pulled out via
// exportKey, from this tab or a future one.
//
// This is device-local storage by design (not synced anywhere), matching
// how every real E2E messenger handles "stay logged in" — losing the
// device means losing the key, same as this app's `/unlock` flow already
// implies losing it means re-entering the password.

const DB_NAME = 'signalline-keyvault'
const DB_VERSION = 1
const STORE_NAME = 'privateKeys'

function isSupported(): boolean {
  return typeof indexedDB !== 'undefined'
}

function openDb(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(DB_NAME, DB_VERSION)
    req.onupgradeneeded = () => {
      if (!req.result.objectStoreNames.contains(STORE_NAME)) {
        req.result.createObjectStore(STORE_NAME)
      }
    }
    req.onsuccess = () => resolve(req.result)
    req.onerror = () => reject(req.error)
  })
}

export async function saveKey(userId: string, key: CryptoKey): Promise<void> {
  if (!isSupported()) return
  try {
    const db = await openDb()
    await new Promise<void>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, 'readwrite')
      tx.objectStore(STORE_NAME).put(key, userId)
      tx.oncomplete = () => resolve()
      tx.onerror = () => reject(tx.error)
    })
    db.close()
  } catch (err) {
    // Not fatal — worst case the user just gets asked for their password
    // again next time, same as before this feature existed.
    console.warn('[signalline] failed to persist key to IndexedDB:', err)
  }
}

export async function loadKey(userId: string): Promise<CryptoKey | null> {
  if (!isSupported()) return null
  try {
    const db = await openDb()
    const key = await new Promise<CryptoKey | null>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, 'readonly')
      const req = tx.objectStore(STORE_NAME).get(userId)
      req.onsuccess = () => resolve((req.result as CryptoKey | undefined) ?? null)
      req.onerror = () => reject(req.error)
    })
    db.close()
    return key
  } catch (err) {
    console.warn('[signalline] failed to read key from IndexedDB:', err)
    return null
  }
}

/** "Forget this device" — clears the persisted key for one user. Not called
 * automatically on logout (logging out and back in on the same device
 * should still auto-unlock), only meant for an explicit user action. */
export async function clearKey(userId: string): Promise<void> {
  if (!isSupported()) return
  try {
    const db = await openDb()
    await new Promise<void>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, 'readwrite')
      tx.objectStore(STORE_NAME).delete(userId)
      tx.oncomplete = () => resolve()
      tx.onerror = () => reject(tx.error)
    })
    db.close()
  } catch {
    // ignore
  }
}
