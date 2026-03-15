import { describe, it, expect, vi, beforeEach } from 'vitest'
import { FlagKit } from '../src/index'

const mockFetch = vi.fn()
globalThis.fetch = mockFetch

function mockResponse(data: object, ok = true) {
  return { ok, status: ok ? 200 : 500, json: () => Promise.resolve(data) }
}

describe('FlagKit', () => {
  beforeEach(() => {
    mockFetch.mockReset()
  })

  it('returns enabled from API', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ enabled: true, reason: 'rollout', flagKey: 'test', evaluatedAt: '' }))

    const flags = new FlagKit({ apiKey: 'key', baseUrl: 'http://localhost:8080' })
    const result = await flags.isEnabled('test', { userId: 'u1' })
    expect(result).toBe(true)
    expect(mockFetch).toHaveBeenCalledOnce()
  })

  it('returns false from API', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ enabled: false, reason: 'no_match', flagKey: 'test', evaluatedAt: '' }))

    const flags = new FlagKit({ apiKey: 'key', baseUrl: 'http://localhost:8080' })
    const result = await flags.isEnabled('test', { userId: 'u1' })
    expect(result).toBe(false)
  })

  it('caches result on subsequent calls', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ enabled: true, reason: 'rollout', flagKey: 'test', evaluatedAt: '' }))

    const flags = new FlagKit({ apiKey: 'key', baseUrl: 'http://localhost:8080', ttl: 60_000 })
    await flags.isEnabled('test', { userId: 'u1' })
    const result = await flags.isEnabled('test', { userId: 'u1' })
    expect(result).toBe(true)
    expect(mockFetch).toHaveBeenCalledOnce()
  })

  it('returns stale value on API failure', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ enabled: true, reason: 'rollout', flagKey: 'test', evaluatedAt: '' }))

    const flags = new FlagKit({ apiKey: 'key', baseUrl: 'http://localhost:8080', ttl: 1 })
    await flags.isEnabled('test', { userId: 'u1' })

    // Wait for TTL to expire
    await new Promise((r) => setTimeout(r, 5))

    // API fails
    mockFetch.mockRejectedValueOnce(new Error('network error'))
    const result = await flags.isEnabled('test', { userId: 'u1' })
    expect(result).toBe(true) // stale fallback
  })

  it('returns false when no cache and API fails', async () => {
    mockFetch.mockRejectedValueOnce(new Error('network error'))

    const flags = new FlagKit({ apiKey: 'key', baseUrl: 'http://localhost:8080' })
    const result = await flags.isEnabled('test', { userId: 'u1' })
    expect(result).toBe(false) // safe default
  })

  it('uses different cache keys for different contexts', async () => {
    mockFetch
      .mockResolvedValueOnce(mockResponse({ enabled: true, reason: 'rollout', flagKey: 'test', evaluatedAt: '' }))
      .mockResolvedValueOnce(mockResponse({ enabled: false, reason: 'no_match', flagKey: 'test', evaluatedAt: '' }))

    const flags = new FlagKit({ apiKey: 'key', baseUrl: 'http://localhost:8080' })
    const r1 = await flags.isEnabled('test', { userId: 'u1' })
    const r2 = await flags.isEnabled('test', { userId: 'u2' })
    expect(r1).toBe(true)
    expect(r2).toBe(false)
    expect(mockFetch).toHaveBeenCalledTimes(2)
  })

  it('sends correct URL and headers', async () => {
    mockFetch.mockResolvedValueOnce(mockResponse({ enabled: true, reason: 'rollout', flagKey: 'test', evaluatedAt: '' }))

    const flags = new FlagKit({ apiKey: 'my-key', baseUrl: 'http://localhost:8080' })
    await flags.isEnabled('my-flag', { userId: 'u1', environment: 'production' })

    const [url, opts] = mockFetch.mock.calls[0]
    expect(url).toContain('/evaluate/my-flag')
    expect(url).toContain('user_id=u1')
    expect(url).toContain('environment=production')
    expect(opts.headers.Authorization).toBe('Bearer my-key')
  })
})
