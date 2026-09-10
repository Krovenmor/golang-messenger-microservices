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
