const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

function getApiKey(): string {
  return localStorage.getItem('flagkit_api_key') || ''
}

export function setApiKey(key: string) {
  localStorage.setItem('flagkit_api_key', key)
}

export function getStoredApiKey(): string {
  return getApiKey()
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${getApiKey()}`,
      ...options?.headers,
    },
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || `Request failed: ${res.status}`)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

export interface Flag {
  id: string
  key: string
  name: string
  description: string
  enabled: boolean
  environment: 'production' | 'staging' | 'development'
  rules: Rule[]
  createdAt: string
  updatedAt: string
}

export interface Rule {
  type: 'percentage' | 'allowlist'
  value?: number
  userIds?: string[]
}

export interface AuditLog {
  id: string
  flagId: string | null
  action: 'created' | 'updated' | 'deleted' | 'toggled'
  diff: Record<string, { from: unknown; to: unknown }>
  actor: string
  createdAt: string
}

export interface PaginatedAudit {
  data: AuditLog[]
  total: number
  page: number
}

export const api = {
  listFlags: () => request<Flag[]>('/flags'),
  getFlag: (key: string) => request<Flag>(`/flags/${key}`),
  createFlag: (data: Partial<Flag>) => request<Flag>('/flags', { method: 'POST', body: JSON.stringify(data) }),
  updateFlag: (key: string, data: Partial<Flag>) => request<Flag>(`/flags/${key}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteFlag: (key: string) => request<void>(`/flags/${key}`, { method: 'DELETE' }),
  toggleFlag: (key: string) => request<Flag>(`/flags/${key}/toggle`, { method: 'POST' }),
  listAudit: (params?: { flagId?: string; page?: number; perPage?: number }) => {
    const search = new URLSearchParams()
    if (params?.flagId) search.set('flag_id', params.flagId)
    if (params?.page) search.set('page', String(params.page))
    if (params?.perPage) search.set('per_page', String(params.perPage))
    return request<PaginatedAudit>(`/audit?${search}`)
  },
}
