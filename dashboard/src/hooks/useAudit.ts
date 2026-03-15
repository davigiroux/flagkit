import { useQuery } from '@tanstack/react-query'
import { api } from '../api/client'

export function useAuditLogs(params?: { flagId?: string; page?: number; perPage?: number }) {
  return useQuery({
    queryKey: ['audit', params],
    queryFn: () => api.listAudit(params),
  })
}
