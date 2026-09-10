const DB_NAME = 'signaline-keyvault'
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
    console.warn('[signaline] failed to persist key to IndexedDB:', err)
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
    console.warn('[signaline] failed to read key from IndexedDB:', err)
    return null
  }
}

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
  } catch {}
}
