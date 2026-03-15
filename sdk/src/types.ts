export interface FlagKitOptions {
  apiKey: string
  baseUrl: string
  ttl?: number
}

export interface EvalContext {
  userId?: string
  environment?: string
}

export interface EvalResult {
  enabled: boolean
  reason: string
  flagKey: string
  evaluatedAt: string
}
