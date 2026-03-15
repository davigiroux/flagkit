import type { EvalContext, EvalResult } from './types'

export class Client {
  private apiKey: string
  private baseUrl: string

  constructor(apiKey: string, baseUrl: string) {
    this.apiKey = apiKey
    this.baseUrl = baseUrl.replace(/\/$/, '')
  }

  async evaluate(key: string, ctx?: EvalContext): Promise<EvalResult> {
    const params = new URLSearchParams()
    if (ctx?.userId) params.set('user_id', ctx.userId)
    if (ctx?.environment) params.set('environment', ctx.environment)

    const url = `${this.baseUrl}/evaluate/${key}?${params}`
    const res = await fetch(url, {
      headers: { Authorization: `Bearer ${this.apiKey}` },
    })

    if (!res.ok) {
      throw new Error(`FlagKit API error: ${res.status}`)
    }

    return res.json()
  }
}
