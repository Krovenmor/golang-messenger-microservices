// Holds the unwrapped ECDH private key for the current session. Populated
// either from a fresh keypair just generated during profile setup, from
// decrypting the stored encryptedPrvKey with a password (see
// UnlockProfile.tsx), or — most commonly after the first time — from the
// non-extractable copy persisted in IndexedDB (see utils/keyVault.ts), so
// returning users on the same device don't have to re-enter their password
// every session. Cleared on logout and on tab close (it's just a
// module-level variable in memory; the IndexedDB copy is separate and
// deliberately survives logout — see keyVault.ts for why).

class CryptoStore {
  private privateKey: CryptoKey | null = null
  private myPubKey: string | null = null
  private listeners = new Set<(unlocked: boolean) => void>()

  isUnlocked(): boolean {
    return !!this.privateKey
  }

  setKeys(privateKey: CryptoKey, myPubKey: string) {
    this.privateKey = privateKey
    this.myPubKey = myPubKey
    this.listeners.forEach((cb) => cb(true))
  }

  getPrivateKey(): CryptoKey | null {
    return this.privateKey
  }

  getMyPubKey(): string | null {
    return this.myPubKey
  }

  lock() {
    this.privateKey = null
    this.myPubKey = null
    this.listeners.forEach((cb) => cb(false))
  }

  subscribe(cb: (unlocked: boolean) => void): () => void {
    this.listeners.add(cb)
    cb(this.isUnlocked())
    return () => this.listeners.delete(cb)
  }
}

export const cryptoStore = new CryptoStore()
