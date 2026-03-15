interface CacheEntry {
  value: boolean
  expiresAt: number
}

export class Cache {
  private store = new Map<string, CacheEntry>()
  private ttl: number

  constructor(ttl: number) {
    this.ttl = ttl
  }

  get(key: string): boolean | undefined {
    const entry = this.store.get(key)
    if (!entry) return undefined
    if (Date.now() > entry.expiresAt) return undefined
    return entry.value
  }

  getStale(key: string): boolean | undefined {
    const entry = this.store.get(key)
    if (!entry) return undefined
    return entry.value
  }

  set(key: string, value: boolean): void {
    this.store.set(key, { value, expiresAt: Date.now() + this.ttl })
  }
}
