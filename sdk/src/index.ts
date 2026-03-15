import { Cache } from './cache'
import { Client } from './client'
import type { EvalContext, FlagKitOptions } from './types'

export type { FlagKitOptions, EvalContext, EvalResult } from './types'

export class FlagKit {
  private cache: Cache
  private client: Client

  constructor(opts: FlagKitOptions) {
    this.client = new Client(opts.apiKey, opts.baseUrl)
    this.cache = new Cache(opts.ttl ?? 30_000)
  }

  async isEnabled(key: string, ctx?: EvalContext): Promise<boolean> {
    const cacheKey = `${key}:${ctx?.userId ?? ''}:${ctx?.environment ?? ''}`

    const cached = this.cache.get(cacheKey)
    if (cached !== undefined) return cached

    try {
      const result = await this.client.evaluate(key, ctx)
      this.cache.set(cacheKey, result.enabled)
      return result.enabled
    } catch {
      const stale = this.cache.getStale(cacheKey)
      if (stale !== undefined) return stale
      return false
    }
  }
}
